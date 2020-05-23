package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/74th/try-envoy/router"
)

type server struct {
	baseWay int64
}

func main() {
	var baseWay int64
	var addr string
	flag.Int64Var(&baseWay, "b", 1, "")
	flag.StringVar(&addr, "H", "0.0.0.0:8080", "")
	flag.Parse()

	startService(baseWay, addr)
}

func isGrpcRequest(r *http.Request) bool {
	return r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc")
}

func startService(baseWay int64, addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sv := &server{
		baseWay: baseWay,
	}
	http.HandleFunc("/", hchandler)
	http.HandleFunc("/_ah/health", hchandler)

	muxHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isGrpcRequest(r) {
			s.ServeHTTP(w, r)
			return
		}
		http.DefaultServeMux.ServeHTTP(w, r)
	})

	router.RegisterRouterServer(s, sv)
	grpc_health_v1.RegisterHealthServer(s, sv)

	if err := http.Serve(lis, h2c.NewHandler(
		muxHandler,
		&http2.Server{},
	)); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (sv *server) Recommend(ctx context.Context, req *router.RouteRequest) (*router.RouteResponse, error) {
	log.Printf("req: %v\n", req.Gps[0].Timestamp)
	ways := make([]int64, 3)
	for i := range ways {
		ways[i] = rand.Int63n(100) * sv.baseWay
	}

	res := new(router.RouteResponse)
	res.Ways = ways
	return res, nil
}

// copy of https://github.com/GoogleCloudPlatform/grpc-gke-nlb-tutorial/blob/master/echo-grpc/health/health.go
func (sv *server) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (sv *server) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

func hchandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
