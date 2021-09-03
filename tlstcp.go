package tlstcp

import (
	"context"
	"fmt"
	"time"

	"crypto/tls"
	"crypto/x509"
	"encoding/binary"

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/metrics"
	"github.com/loadimpact/k6/stats"

	//"github.com/loadimpact/k6/lib/netext/httpext"
	"github.com/loadimpact/k6/js/modules"
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

	//
	tags := state.CloneTags()

	//argument parsing goes around here from GOJA VM
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
	// },
	// 	CurvePreferences:   []tls.CurveID{tls.X25519},
	// 	RootCAs:            roots,
	// }

	start := time.Now()
	conon, err := state.Dialer.DialContext(ctx, network, addr)
	connectionEnd := time.Now()
	connectionDuration := stats.D(connectionEnd.Sub(start))

	sampleTags := stats.IntoSampleTags(&tags)
	stats.PushIfNotDone(ctx, state.Samples, stats.ConnectedSamples{
		Samples: []stats.Sample{
			{Metric: metrics.HTTPReqs, Time: start, Tags: sampleTags, Value: 1},
			{Metric: metrics.HTTPReqConnecting, Time: start, Tags: sampleTags, Value: connectionDuration},
		},
		Tags: sampleTags,
		Time: start,
	})

	if err != nil {
		fmt.Println(fmt.Sprintf("errorfound %s", err))
		return nil, common.NewInitContextError(fmt.Sprintf("Dial error %s", err))

	}

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
	//fmt.Printf("request size without headers %d bytes: %v\n ", len(data), data)

	binary.LittleEndian.PutUint32(bufHeader, uint32(len(data)))

	singlePacket := make([]byte, len(bufHeader)+len(data))

	//add on header the bytes
	copy(singlePacket[:], bufHeader[:])
	copy(singlePacket[len(bufHeader):], data[:])

	_, error := requester.Write(singlePacket)
	//fmt.Printf("write ret code: %d \n ", cd)

	if error != nil {
		fmt.Println("send failure")
		fmt.Printf("send failure %s", error.Error())
	}
}

//Receive bytes from socket
func (*TlsTcp) Receive(ctx context.Context, requester *tls.Conn, data []byte) []byte {

	var arr []byte
	_, error := requester.Read(arr)
	if error != nil {
		fmt.Println("failed to read")
		return arr
	}
	//fmt.Printf("lenght received %d\n", len(arr))
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
