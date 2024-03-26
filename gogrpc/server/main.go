package main

import (
	"fmt"
	"log"
	"net"
	"server/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	s := grpc.NewServer()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	services.RegisterCalculatorServer(s, services.NewCalculatorServer()) // register service เพื่อมาบริการงานที่ทำไว้
	reflection.Register(s)

	fmt.Println("gRPC server listening on port 50051")
	err = s.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
