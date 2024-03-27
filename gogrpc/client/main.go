package main

import (
	"client/services"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	var cc *grpc.ClientConn
	var err error
	var creds credentials.TransportCredentials

	host := flag.String("host", "localhost:50051", "gRPC server host")
	tls := flag.Bool("tls", false, "use a secure TLS connection")
	flag.Parse()

	if *tls {
		certFile := "../tls/ca.crt"
		creds, err = credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatal(err)
		}

	} else {
		// add insecure ไปก่อน ก่อนทำ tls
		creds = insecure.NewCredentials()
	}

	cc, err = grpc.Dial(*host, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close() //close grpc Dial

	calculatorClient := services.NewCalculatorClient(cc)
	calculatorService := services.NewCalculatorService(calculatorClient)

	err = calculatorService.Hello("JaneJane")
	// err = calculatorService.Fibonacci(3)
	// err = calculatorService.Average(1, 12, 123, 1234, 12345)
	// err = calculatorService.Sum(1, 2, 3, 4, 5)

	if err != nil {
		// มันสามารถแยก err ให้ได้ ว่ามาจาก grpc หรือจากที่อื่น
		if grpcErr, ok := status.FromError(err); ok { //ถ้าใช่ error นั้นมาจาก grpc ค่าของ ok = true และได้ grpcError มา
			// และเช็คต่อว่า ถ้า ok ทำต่อ
			log.Printf("[%v] %v", grpcErr.Code(), grpcErr.Message())
		} else {
			log.Fatal(err)
		}
	}
}
