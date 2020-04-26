package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net"

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

func startService(baseWay int64, addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sv := &server{
		baseWay: baseWay,
	}
	router.RegisterRouterServer(s, sv)
	grpc_health_v1.RegisterHealthServer(s, sv)

	if err := s.Serve(lis); err != nil {
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
