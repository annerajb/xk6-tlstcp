package tlstcp

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"encoding/base64"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/metrics"
	"go.k6.io/k6/stats"

	//"go.k6.io/k6/lib/netext/httpext"
	"go.k6.io/k6/js/modules"
	//"go.k6.io/k6/stats"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/tlstcp".
func init() {
	modules.Register("k6/x/tlstcp", New())
}

var ErrTcpInInitContext = common.NewInitContextError("using tlstcp in the init context is not supported")

type IPLevel4Response struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
	Body   string `json:"body"`
}

type TlsTcp struct{}

func New() *TlsTcp {
	return &TlsTcp{}
}

//Connect https://stackoverflow.com/q/65451276
//, args ...goja.Value
// func (*TlsTcp) Connect(ctx context.Context, network string, addr string, root_ca_pem string) (*tls.Conn, error) {
// 	state := lib.GetState(ctx)
// 	if state == nil {
// 		return nil, ErrTcpInInitContext
// 	}

// 	//
// 	tags := state.CloneTags()

// 	//argument parsing goes around here from GOJA VM
// 	// First, create the set of root certificates. For this example we only
// 	// have one. It's also possible to omit this in order to use the
// 	// default root set of the current operating system.
// 	roots := x509.NewCertPool()
// 	ok := roots.AppendCertsFromPEM([]byte(root_ca_pem))
// 	if !ok {
// 		fmt.Println("failed to parse root cert pool")
// 		return nil, common.NewInitContextError("failed to parse root cert pool")
// 	}
// 	//
// 	var tlsConfig *tls.Config
// 	if state.TLSConfig != nil {
// 		tlsConfig = state.TLSConfig.Clone()
// 	}
// 	tlsConfig.RootCAs = roots
// 	tlsConfig.MinVersion = tls.VersionTLS13
// 	tlsConfig.InsecureSkipVerify = true
// 	tlsConfig.CipherSuites = []uint16{
// 		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
// 		//tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
// 	}
// 	tlsConfig.CurvePreferences = []tls.CurveID{tls.X25519}

// 	//Todo Move this to caller or use state only
// 	// tlsConfig := &tls.Config{
// 	// 	MinVersion: tls.VersionTLS13,

// 	start := time.Now()
// 	conon, err := state.Dialer.DialContext(ctx, network, addr)
// 	connectionEnd := time.Now()
// 	connectionDuration := stats.D(connectionEnd.Sub(start))

// 	sampleTags := stats.IntoSampleTags(&tags)

// 	if err != nil {
// 		fmt.Printf("errorfound %s\n", err)
// 		return nil, common.NewInitContextError(fmt.Sprintf("Dial error %s", err))
// 	}

// 	// Only set the name system tag if the user didn't explicitly set it beforehand,
// 	// and the Name was generated from a tagged template string (via http.url).
// 	// if _, ok := tags["name"]; !ok && state.Options.SystemTags.Has(stats.TagName) &&
// 	// 	preq.URL.Name != "" && preq.URL.Name != preq.URL.Clean() {
// 	// 	tags["name"] = preq.URL.Name
// 	// }

// 	// Check rate limit *after* we've prepared a request; no need to wait with that part.
// 	if rpsLimit := state.RPSLimit; rpsLimit != nil {
// 		if err := rpsLimit.Wait(ctx); err != nil {
// 			return nil, err
// 		}
// 	}

// 	conn := tls.Client(conon, tlsConfig)
// 	//tracerTransport := httpext.newTransport(ctx, state, tags)

// 	startHS := time.Now()
// 	err = conn.Handshake()
// 	endHS := time.Now()
// 	hsDuration := stats.D(endHS.Sub(startHS))

// 	if err != nil {
// 		fmt.Printf("handshake failure: %s", err.Error())
// 		return nil, err
// 	}

// 	stats.PushIfNotDone(ctx, state.Samples, stats.ConnectedSamples{
// 		Samples: []stats.Sample{
// 			{Metric: metrics.HTTPReqs, Time: start, Tags: sampleTags, Value: 1},
// 			{Metric: metrics.HTTPReqConnecting, Time: start, Tags: sampleTags, Value: connectionDuration},
// 			{Metric: metrics.HTTPReqTLSHandshaking, Time: startHS, Tags: sampleTags, Value: hsDuration},
// 		},
// 		Tags: sampleTags,
// 		Time: start,
// 	})
// 	return conn, nil
// }

//CloseConn the socket connection
// func (*TlsTcp) CloseConn(ctx context.Context, requester *tls.Conn, data string) {
// 	requester.Close()
// }

//Send bytes to a socket
// func (*TlsTcp) Send(ctx context.Context, requester *tls.Conn, datastr string) {
// 	//fmt.Println("sending bytes")
// 	data, err := base64.StdEncoding.DecodeString(datastr)
// 	if err != nil {
// 		log.Fatalf("Some error occured during base64 decode. Error %s", err.Error())
// 	}

// 	bufHeader := make([]byte, 4)
// 	//fmt.Printf("request size without headers %d bytes: %v\n ", len(data), data)

// 	binary.LittleEndian.PutUint32(bufHeader, uint32(len(data)))

// 	singlePacket := make([]byte, len(bufHeader)+len(data))

// 	//add on header the bytes
// 	copy(singlePacket[:], bufHeader[:])
// 	copy(singlePacket[len(bufHeader):], data[:])

// 	sendStart := time.Now()
// 	_, error := requester.Write(singlePacket)
// 	sendEnd := time.Now()
// 	sendDuration := stats.D(sendEnd.Sub(sendStart))
// 	state := lib.GetState(ctx)
// 	stats.PushIfNotDone(ctx, state.Samples, stats.ConnectedSamples{
// 		Samples: []stats.Sample{
// 			{Metric: metrics.HTTPReqSending, Time: sendStart, Value: sendDuration},
// 		},
// 		Time: sendStart,
// 	})

// 	if error != nil {
// 		fmt.Println("send failure")
// 		fmt.Printf("send failure: %s\n", error.Error())
// 		return
// 	}
// }

// //Receive bytes from socket
// func (*TlsTcp) Receive(ctx context.Context, requester *tls.Conn) []byte {
// 	//fmt.Println("receiving bytes")

// 	header_replay := make([]byte, 4)

// 	recvHeaderStart := time.Now()
// 	_, error := requester.Read(header_replay)
// 	recvHeaderEnd := time.Now()
// 	recvHeaderDuration := stats.D(recvHeaderEnd.Sub(recvHeaderStart))

// 	if error != nil {
// 		fmt.Printf("failed to read header: %s\n", error.Error())
// 		return header_replay
// 	}
// 	expectedLengh := binary.LittleEndian.Uint32(header_replay)

// 	//fmt.Printf("lenght received %d\n", expectedLengh)

// 	receivedContent := make([]byte, expectedLengh)

// 	recvStart := time.Now()
// 	_, error = requester.Read(receivedContent)
// 	recvEnd := time.Now()
// 	recvDuration := stats.D(recvEnd.Sub(recvStart)) + recvHeaderDuration
// 	//todo error check state value
// 	state := lib.GetState(ctx)
// 	stats.PushIfNotDone(ctx, state.Samples, stats.ConnectedSamples{
// 		Samples: []stats.Sample{
// 			{Metric: metrics.HTTPReqReceiving, Time: recvStart, Value: recvDuration},
// 		},
// 		Time: recvStart,
// 	})

// 	if error != nil {
// 		fmt.Printf("failed to read header: %s\n", error.Error())
// 		return receivedContent
// 	}
// 	return receivedContent
// }

//Connect https://stackoverflow.com/q/65451276
//, args ...goja.Value
func (*TlsTcp) Connect(ctx context.Context, network string, addr string, root_ca_pem string) (net.Conn, error) {
	state := lib.GetState(ctx)
	if state == nil {
		return nil, ErrTcpInInitContext
	}

	//
	tags := state.CloneTags()

	//argument parsing goes around here from GOJA VM

	//Todo Move this to caller or use state only
	// tlsConfig := &tls.Config{
	// 	MinVersion: tls.VersionTLS13,

	start := time.Now()
	conon, err := state.Dialer.DialContext(ctx, network, addr)
	connectionEnd := time.Now()
	connectionDuration := stats.D(connectionEnd.Sub(start))

	sampleTags := stats.IntoSampleTags(&tags)

	if err != nil {
		fmt.Printf("error found %s\n", err)
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

	// s, err := net.ResolveUDPAddr(network, addr)
	// if err != nil {
	// 	log.Fatalf("dns resolution fail: %s", err.Error())
	// }
	//tracerTransport := httpext.newTransport(ctx, state, tags)

	//startHS := time.Now()
	//endHS := time.Now()
	//hsDuration := stats.D(endHS.Sub(startHS))

	if err != nil {
		fmt.Printf("handshake failure: %s", err.Error())
		return nil, err
	}

	stats.PushIfNotDone(ctx, state.Samples, stats.ConnectedSamples{
		Samples: []stats.Sample{
			{Metric: metrics.HTTPReqs, Time: start, Tags: sampleTags, Value: 1},
			{Metric: metrics.HTTPReqConnecting, Time: start, Tags: sampleTags, Value: connectionDuration},
			//{Metric: metrics.HTTPReqTLSHandshaking, Time: startHS, Tags: sampleTags, Value: hsDuration},
		},
		Tags: sampleTags,
		Time: start,
	})
	return conon, nil
}
func (*TlsTcp) CloseConn(ctx context.Context, requester net.Conn, data string) {
	requester.Close()
}

//Send bytes to a socket
func (*TlsTcp) Send(ctx context.Context, requester net.Conn, datastr string) *IPLevel4Response {
	data, err := base64.StdEncoding.DecodeString(datastr)
	if err != nil {
		log.Fatalf("Some error occured during base64 decode. Error %s", err.Error())
	}
	requester.SetDeadline(time.Now().Add(60 * time.Second))
	sendStart := time.Now()
	_, err = requester.Write(data)
	sendEnd := time.Now()
	sendDuration := stats.D(sendEnd.Sub(sendStart))

	state := lib.GetState(ctx)
	stats.PushIfNotDone(ctx, state.Samples, stats.ConnectedSamples{
		Samples: []stats.Sample{
			{Metric: metrics.HTTPReqSending, Time: sendStart, Value: sendDuration},
		},
		Time: sendStart,
	})
	lvl4Response := IPLevel4Response{
		Status: 200,
	}
	if err != nil {
		fmt.Println("send failure")
		fmt.Printf("send failure: %s\n", err.Error())
		lvl4Response.Status = 300
	}
	return &lvl4Response
}

//Receive bytes from socket
func (*TlsTcp) Receive(ctx context.Context, requester net.Conn) *goja.ArrayBuffer {
	rt := common.GetRuntime(ctx)

	header_replay := make([]byte, 148)

	requester.SetDeadline(time.Now().Add(5 * time.Second))
	recvStart := time.Now()
	n, error := requester.Read(header_replay)
	recvHeaderEnd := time.Now()
	recvHeaderDuration := stats.D(recvHeaderEnd.Sub(recvStart))

	//todo error check state value
	state := lib.GetState(ctx)
	stats.PushIfNotDone(ctx, state.Samples, stats.ConnectedSamples{
		Samples: []stats.Sample{
			{Metric: metrics.HTTPReqReceiving, Time: recvStart, Value: recvHeaderDuration},
		},
		Time: recvStart,
	})
	// lvl4Response := IPLevel4Response{
	// 	Status: 200,
	// }
	if error != nil {
		//fmt.Printf("failed to read header: %s\n", error.Error())
		common.Throw(common.GetRuntime(ctx), error)
		return nil
	}
	ab := rt.NewArrayBuffer(header_replay[:n])
	//arb := rt.ToValue(&ab)
	//lvl4Response.Body
	//&lvl4Response
	return &ab
}
