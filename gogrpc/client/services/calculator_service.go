package services

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CalculatorService interface {
	Hello(name string) error
	Fibonacci(n uint32) error
	Average(numbers ...float64) error
	Sum(numbers ...int32) error
}

type calculatorService struct {
	calculatorClient CalculatorClient
}

func NewCalculatorService(calculatorClient CalculatorClient) CalculatorService {

	return calculatorService{calculatorClient}
}

// reciever function
func (base calculatorService) Hello(name string) error {
	req := HelloRequest{
		Name:        name,
		CreatedDate: timestamppb.Now(),
	}

	res, err := base.calculatorClient.Hello(context.Background(), &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service : Hello\n")
	fmt.Printf("Request : %v\n", req.Name)
	fmt.Printf("Respone : %v\n", res.Result)
	return nil
}

func (base calculatorService) Fibonacci(n uint32) error {
	req := FibonacciRequest{
		N: n,
	}

	// จัดการเรื่อง timeout ถ้า stream ที่วิ่งมามันนานเกิน
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := base.calculatorClient.Fibonacci(ctx, &req)
	if err != nil {
		return err
	}

	fmt.Printf("Service : Fibonacci\n")
	fmt.Printf("Request : %v\n", req.N)
	for {
		res, err := stream.Recv() // stream receiver
		if err == io.EOF {        // เมื่อ stream หมด มันจะส่ง err io.EOF
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Response: %v\n", res.Result)
	}

	return nil
}

func (base calculatorService) Average(numbers ...float64) error {
	stream, err := base.calculatorClient.Average(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Service : Average\n")
	for _, number := range numbers {
		req := AverageRequest{
			Number: number,
		}
		stream.Send(&req)
		fmt.Printf("Request : %v\n", req.Number)
		time.Sleep(time.Second) // buffer send req each 1 second
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Printf("Response: %v\n", res.Result)
	return nil
}

func (base calculatorService) Sum(numbers ...int32) error {
	stream, err := base.calculatorClient.Sum(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Service : Sum\n")

	//request routine
	go func() {
		for _, number := range numbers {
			req := SumRequest{
				Number: number,
			}
			stream.Send(&req)
			fmt.Printf("Request: %v\n", req.Number)
			time.Sleep(time.Second * 1)
		}
		stream.CloseSend()
	}()

	// กำหนดจุดจบของ routine
	done := make(chan bool)  //channel
	errs := make(chan error) //channel

	//response routine
	go func() {
		// เราไม่รู้ว่า request จะหยุดส่งมาเมื่อไร เราจึงวนลูปไม่รู้จบไว้
		for {
			res, err := stream.Recv()
			if err == io.EOF { // เมื่อ stream หมด มันจะส่ง err io.EOF
				break
			}
			if err != nil {
				errs <- err
			}
			fmt.Printf("Response: %v\n", res.Result)
		}
		done <- true // ส่งค่าให้ channel ว่า true
	}()

	select {
	case <-done:
		return nil
	case err := <-errs: //case ที่เป็น err
		return err
	}

}
