package proxy

import (
	"context"
	"fmt"
	"net/http"

	bk "test_grpc/api/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// Data for handlers
type proxyREST struct {
	Host  string
	Port  string
	GPort string
	Conn  string
	GConn string
}

// New creates new proxy struct
func New(ghost string, gport string, port string) *proxyREST {
	str := fmt.Sprintf(":%s", port)
	gstr := fmt.Sprintf("%s:%s", ghost, gport)
	return &proxyREST{
		Host:  ghost,
		Port:  port,
		GPort: gport,
		Conn:  str,
		GConn: gstr,
	}
}

// Start REST proxy
func (s *proxyREST) Start() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := bk.RegisterLibraryHandlerFromEndpoint(ctx, mux, s.GConn, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(s.Conn, mux)
}
