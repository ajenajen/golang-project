package main

import (
	"client/services"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	creds := insecure.NewCredentials() // add insecure ไปก่อน ก่อนทำ tls

	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close() //close grpc Dial

	calculatorClient := services.NewCalculatorClient(cc)
	calculatorService := services.NewCalculatorService(calculatorClient)

	// err = calculatorService.Hello("Test")
	// err = calculatorService.Fibonacci(3)
	// err = calculatorService.Average(1, 12, 123, 1234, 12345)
	err = calculatorService.Sum(1, 2, 3, 4, 5)
	if err != nil {
		log.Fatal(err)
	}
}
