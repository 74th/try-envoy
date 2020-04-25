package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/74th/try-envoy/router"

	"google.golang.org/grpc"
)

var addr string
var base int64

func init() {
	flag.StringVar(&addr, "H", "localhost:50000", "")
	flag.Int64Var(&base, "b", 1, "")
	flag.Parse()
}

func main() {

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
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
