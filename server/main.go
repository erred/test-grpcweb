package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"

	hw "seankhliao.com/grpctest/helloworld"
)

var serverOption grpc.ServerOption

func init() {
	// // Read cert and key file
	// BackendCert, _ := ioutil.ReadFile("./backend.cert")
	// BackendKey, _ := ioutil.ReadFile("./backend.key")
	//
	// // Generate Certificate struct
	// cert, err := tls.X509KeyPair(BackendCert, BackendKey)
	// if err != nil {
	// 	log.Fatalf("failed to parse certificate: %v", err)
	// }
	//
	// // Create credentials
	// creds := credentials.NewServerTLSFromCert(&cert)
	// serverOption = grpc.Creds(creds)
}

func main() {
	// svr := grpc.NewServer(serverOption)
	svr := grpc.NewServer()
	hw.RegisterGreeterServer(svr, &Server{})
	lis, _ := net.Listen("tcp", ":8080")
	svr.Serve(lis)
	// upstreamGRPCserver := &SVR{svr, "grpc-web -> grpc"}
	// proxyGRPCserver := &SVR{grpcweb.New(upstreamGRPCserver), "client -> grpc-web"}

	// h2server := &http2.Server{}
	// h2chandler := h2c.NewHandler(proxyGRPCserver, h2server)

	// http.Handle("/helloworld.Greeter/", h2chandler)
	// http.Handle("/helloworld.Greeter/", upstreamGRPCserver)
	// http.Handle("/helloworld.Greeter/", proxyGRPCserver)
	// http.Handle("/helloworld.Greeter/", svr)
	// http.Handle("/", http.FileServer(http.Dir("./web")))
	// http.ListenAndServe(":8080", nil)
	// log.Fatal(http.ListenAndServeTLS(":8080", "./backend.cert", "./backend.key", nil))
}

// ============================= Interceptor
//			client -> grpc-web -> interceptor -> grpc

type ic struct {
	http.ResponseWriter
	buf []byte
}

func (i *ic) Write(p []byte) (n int, err error) {
	n, err = i.ResponseWriter.Write(p)
	i.buf = append(i.buf, p[:n]...)
	return
}
func (i *ic) Flush()                   { i.ResponseWriter.(http.Flusher).Flush() }
func (i *ic) CloseNotify() <-chan bool { return i.ResponseWriter.(http.CloseNotifier).CloseNotify() }

type SVR struct {
	http.Handler
	name string
}

func (s *SVR) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// b, _ := httputil.DumpRequest(r, true)
	// fmt.Printf("\n%s dump request:\n%#v\n", s.name, string(b))

	ww := &ic{w, nil}
	s.Handler.ServeHTTP(ww, r)

	// res, _ := http.ReadResponse(bufio.NewReader(bytes.NewBuffer(ww.buf)), r)
	// b, _ := httputil.DumpResponse(res, true)
	fmt.Printf("\n%s dump response:\n%#v\n", s.name, string(ww.buf))
}

// ============================= GRPC server

type Server struct{}

func (s *Server) SayHello(ctx context.Context, r *hw.HelloRequest) (*hw.HelloReply, error) {
	return &hw.HelloReply{
		Message: "hello " + r.GetName(),
	}, nil
}
func (s *Server) SayRepeatHello(r *hw.RepeatHelloRequest, svr hw.Greeter_SayRepeatHelloServer) error {
	for i := 0; i < int(r.GetCount()); i++ {
		svr.Send(&hw.HelloReply{
			Message: "hello " + r.GetName() + " " + strconv.Itoa(i) + "\n",
		})
	}
	return nil
}
func (s *Server) SayHelloAfterDelay(ctx context.Context, r *hw.HelloRequest) (*hw.HelloReply, error) {
	time.Sleep(5 * time.Second)
	return &hw.HelloReply{
		Message: "hello " + r.GetName(),
	}, nil
}
