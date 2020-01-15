package gateway

import (
	"context"
	"fmt"
	"net/http"

	lg "github.com/sirupsen/logrus"
	bk "test_grpc/api/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// Data for HTTP REST gateway
type proxyREST struct {
	log   lg.FieldLogger
	Host  string
	Port  string
	GPort string
	Conn  string
	GConn string
}

// New creates new HTTP gateway struct
func New(ghost string, gport string, port string, log lg.FieldLogger) *proxyREST {
	str := fmt.Sprintf(":%s", port)
	gstr := fmt.Sprintf("%s:%s", ghost, gport)
	return &proxyREST{
		log:   log,
		Host:  ghost,
		Port:  port,
		GPort: gport,
		Conn:  str,
		GConn: gstr,
	}
}

// Start REST HTTP gateway
func (s *proxyREST) Start() error {
	s.log.Println("HTTP gateway init.")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := bk.RegisterLibraryHandlerFromEndpoint(ctx, mux, s.GConn, opts)
	if err != nil {
		return err
	}

	s.log.Println("HTTP gateway start.")
	return http.ListenAndServe(s.Conn, mux)
}
