package tlstcp

import (
	"context"
	"fmt"
	"net"
	"time"

	"crypto/tls"
	"crypto/x509"
	"encoding/binary"

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/netext"

	//"github.com/loadimpact/k6/lib/netext/httpext"
	"github.com/loadimpact/k6/js/modules"
	"github.com/loadimpact/k6/lib/types"
	//"github.com/loadimpact/k6/stats"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/tlstcp".
func init() {
	modules.Register("k6/x/tlstcp", New())
}

var ErrTcpInInitContext = common.NewInitContextError("using tlstcp in the init context is not supported")

type TlsTcp struct{}

func New() *TlsTcp {
	return &TlsTcp{}
}

//Connect https://stackoverflow.com/q/65451276
//, args ...goja.Value
func (*TlsTcp) Connect(ctx context.Context, network string, addr string, root_ca_pem string) (*tls.Conn, error) {
	state := lib.GetState(ctx)
	if state == nil {
		return nil, ErrTcpInInitContext
	}
	// First, create the set of root certificates. For this example we only
	// have one. It's also possible to omit this in order to use the
	// default root set of the current operating system.
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(root_ca_pem))
	if !ok {
		fmt.Println("failed to parse root cert pool")
		return nil, common.NewInitContextError("failed to parse root cert pool")
	}
	//
	var tlsConfig *tls.Config
	if state.TLSConfig != nil {
		tlsConfig = state.TLSConfig.Clone()
	}
	tlsConfig.RootCAs = roots
	tlsConfig.MinVersion = tls.VersionTLS13
	tlsConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		//tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	}
	tlsConfig.CurvePreferences = []tls.CurveID{tls.X25519}
	//Todo Move this to caller or use state only
	// tlsConfig := &tls.Config{
	// 	MinVersion: tls.VersionTLS13,
	// 	CipherSuites: []uint16{
	//     tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	//     //tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	// },
	// 	CurvePreferences:   []tls.CurveID{tls.X25519},
	// 	RootCAs:            roots,
	// 	InsecureSkipVerify: true,
	// }
	dialer := netext.NewDialer(net.Dialer{
		Timeout:   2 * time.Second,
		KeepAlive: 10 * time.Second,
	}, netext.NewResolver(net.LookupIP, 0, types.DNSfirst, types.DNSpreferIPv4))

	conon, err := dialer.DialContext(ctx, network, addr)
	if err != nil {
		fmt.Println("failed to connect root cert pool")
		return nil, common.NewInitContextError("failed to connect root cert pool")
	}

	//tags := state.CloneTags()
	// Override any global tags with request-specific ones.
	// for k, v := range preq.Tags {
	// 	tags[k] = v
	// }

	// Only set the name system tag if the user didn't explicitly set it beforehand,
	// and the Name was generated from a tagged template string (via http.url).
	// if _, ok := tags["name"]; !ok && state.Options.SystemTags.Has(stats.TagName) &&
	// 	preq.URL.Name != "" && preq.URL.Name != preq.URL.Clean() {
	// 	tags["name"] = preq.URL.Name
	// }

	// Check rate limit *after* we've prepared a request; no need to wait with that part.
	if rpsLimit := state.RPSLimit; rpsLimit != nil {
		if err := rpsLimit.Wait(ctx); err != nil {
			return nil, err
		}
	}

	conn := tls.Client(conon, tlsConfig)

	//tracerTransport := httpext.newTransport(ctx, state, tags)

	return conn, nil
}

//CloseConn the socket connection
func (*TlsTcp) CloseConn(ctx context.Context, requester *tls.Conn, data string) {
	requester.Close()
}

//Send bytes to a socket
func (*TlsTcp) Send(ctx context.Context, requester *tls.Conn, data []byte) {

	bufHeader := make([]byte, 4)
	fmt.Printf("sending %d bytes: %v\n ", len(data), data)

	binary.LittleEndian.PutUint32(bufHeader, 13) //len(data)

	singlePacket := make([]byte, len(bufHeader)+len(data))

	//add on header the bytes
	copy(singlePacket[:], bufHeader[:])
	copy(singlePacket[len(bufHeader):], data[:])
	requester.Write(singlePacket)
}

//Receive bytes from socket
func (*TlsTcp) Receive(ctx context.Context, requester *tls.Conn, data []byte) []byte {

	var arr []byte
	len, error := requester.Read(arr)
	if error != nil {
		fmt.Println("failed to read")
		return arr
	}
	fmt.Printf("lenght received %d\n", len)
	return arr
}

//custom root chain https://gist.github.com/StevenACoffman/844ad083a2b2998c67c6ff7b56710851
// type Client struct {
//     conn *tls.Conn
// }
// // XClient represents the Client constructor (i.e. `new redis.Client()`) and
// // returns a new Redis client object.
// func (r *TlsTcp) XClient(ctxPtr *context.Context, conf *tls.Config) interface{} {
//     rt := common.GetRuntime(*ctxPtr)
//     //{client: conn}
//     return common.Bind(rt, &Client{conn: conn}, ctxPtr)
// }
