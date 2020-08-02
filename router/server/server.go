package main

import (
	"context"
	"flag"
	"fmt"
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

func main() {
	flag.Int64Var(&baseWay, "b", 1, "")
	flag.StringVar(&addr, "H", ":8080", "")
	flag.BoolVar(&useTLS, "tls", false, "")
	flag.Parse()

	startService()
}

func isGrpcRequest(r *http.Request) bool {
	return r.ProtoMajor == 2 && r.Method == "PRI" && r.RequestURI == "*"
}

func startService() {

	healthCheckMux := http.NewServeMux()
	healthCheckMux.HandleFunc("/", hcHandler)
	healthCheckMux.HandleFunc("/healthz", hcHandler)
	healthCheckMux.HandleFunc("/_ah/health", hcHandler)

	// httpMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	if isGrpcRequest(r) {
	// 		s.ServeHTTP(w, r)
	// 		return
	// 	}
	// 	healthCheckMux.ServeHTTP(w, r)
	// })

	// gRPC サーバ
	s := grpc.NewServer()
	sv := &server{baseWay: baseWay}
	router.RegisterRouterServer(s, sv)

	if useTLS {
		// TLS なら http を介して動く
		err := http.ListenAndServeTLS(addr, "./cert/cert.pem", "./cert/key.pem", s)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
	} else {
		// 非TLS の場合、http を介して動作しない
		// err := http.ListenAndServe(addr, s)
		// if err != nil {
		// 	log.Fatalf("failed to listen: %v", err)
		// }

		// 非TLSの場合、Lister直なら動いた
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
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
func (sv *server) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (sv *server) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

func hcHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
