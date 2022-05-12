package main

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Channel(stream HelloService_ChannelServer) error {

		_, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("== io.EOF")
				return nil
			}
			fmt.Println("== " + err.Error())
			return err
		}
	for {
		reply := &String {Value: "Hello," + time.Now().String()}
		 err = stream.Send(reply)
		 if err != nil {
		 	return err
		 }
		 time.Sleep(time.Second * 2)
	}
}

func main() {
	grpcServer := grpc.NewServer()
	RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}