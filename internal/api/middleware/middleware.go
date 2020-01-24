package middleware

import (
	"context"
	"time"

	m "github.com/grpc-ecosystem/go-grpc-middleware"
	lg "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type LogData struct {
	Log  lg.FieldLogger
	Name string
}

func (l *LogData) WithServerUnaryInterceptors() grpc.ServerOption {
	return grpc.UnaryInterceptor(m.ChainUnaryServer(l.LoggerInterceptor))
}

// Logger unary interceptor function to handle logging RPC call
func (l *LogData) LoggerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	// Calls the handler
	h, err := handler(ctx, req)
	// Logging
	l.Log.Printf("%s call: %s   [data: %v] [time: %s] [ret: %v]\n", l.Name, info.FullMethod, req, time.Since(start), err)

	return h, err
}
