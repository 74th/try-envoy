package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/74th/try-envoy/router"
)

var baseWay int64
var addr string
var useTLS bool

type server struct {
	baseWay int64
}

type healthCheck struct{}

func main() {
	flag.Int64Var(&baseWay, "b", 1, "")
	flag.StringVar(&addr, "H", ":8080", "")
	flag.BoolVar(&useTLS, "tls", false, "")
	flag.Parse()

	startService()
}

func startService() {

	// gRPC Server
	s := grpc.NewServer()

	// Application
	sv := &server{baseWay: baseWay}
	router.RegisterRouterServer(s, sv)

	// Health Check
	hc := &healthCheck{}
	grpc_health_v1.RegisterHealthServer(s, hc)

	if useTLS {
		// TLS
		err := http.ListenAndServeTLS(addr, "./cert/cert.pem", "./cert/key.pem", s)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
	} else {
		// if not using TLS, it cannot use http
		// err := http.ListenAndServe(addr, s)
		// if err != nil {
		// 	log.Fatalf("failed to listen: %v", err)
		// }

		// not TLS
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
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
func (hc *healthCheck) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (hc *healthCheck) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}
