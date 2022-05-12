package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"io"
	"log"
	"time"

	gbackoff "google.golang.org/grpc/backoff"
)

const (
	grpcInitialWindowSize     = 1 << 21
	grpcInitialConnWindowSize = 1 << 30
	grpcMaxCallRecvMsgSize    = 1 << 26
)

func main() {
	//conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	conn, err := grpc.DialContext(context.Background(), "localhost:1234", grpc.WithInsecure(),
		grpc.WithInitialWindowSize(grpcInitialWindowSize),
		grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallRecvMsgSize)),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: gbackoff.Config{
				BaseDelay:  time.Second,
				Multiplier: 1.1,
				Jitter:     0.1,
				MaxDelay:   3 * time.Second,
			},
			MinConnectTimeout: 3 * time.Second,
		}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := NewHelloServiceClient(conn)
	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if err := stream.Send(&String{Value: "hi"}); err != nil {
		log.Fatal(err)
	}
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("== io.EOF")
				return
			}
			fmt.Println("== " + err.Error())
			log.Fatal(err)
		}
		fmt.Println(reply.GetValue())
	}

}
