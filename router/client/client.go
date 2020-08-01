package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/74th/try-envoy/router"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr string
var base int64
var tls bool

func init() {
	flag.StringVar(&addr, "H", "localhost:50000", "")
	flag.BoolVar(&tls, "tls", false, "")
	flag.Int64Var(&base, "b", 1, "")
	flag.Parse()
}

func main() {

	log.Printf("connect to %s", addr)

	var conn *grpc.ClientConn
	if tls {
		creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
		if err != nil {
			log.Fatal(err.Error())
		}
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
	} else {
		var err error
		conn, err = grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
	}
	defer conn.Close()
	c := router.NewRouterClient(conn)
	for i := 0; i < 100; i++ {
		request(c)
		time.Sleep(time.Millisecond * 500)
	}
}

func request(c router.RouterClient) {
	ctx := context.Background()
	req := new(router.RouteRequest)
	req.Gps = make([]*router.GPSPoint, 3)
	for i := range req.Gps {
		req.Gps[i] = &router.GPSPoint{
			Timestamp: rand.Int63n(100) * base,
			Latitude:  rand.Float64(),
			Lognitude: rand.Float64(),
		}
	}
	res, err := c.Recommend(ctx, req)
	if err != nil {
		log.Printf("err: %s\n", err.Error())
	} else {
		log.Printf("res: %v\n", res)
	}
}
