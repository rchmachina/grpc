package chatgrpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	streamPb "github.com/rchmachina/grpc/dto/streamingService"
	"google.golang.org/grpc"
)

func Chat() {
	conn, err := grpc.Dial("localhost:55001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := streamPb.NewStreamingServiceClient(conn)
	stream, err := client.BidirectionalStreaming(context.Background())
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	waitc := make(chan struct{})

	// Goroutine to send messages to the server
	go func() {
		for i := 0; i < 5; i++ {
			if err := stream.Send(&streamPb.StreamRequest{
				Message: fmt.Sprintf("Message %d from client", i),
			}); err != nil {
				log.Fatalf("Failed to send message: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend() // Close the send stream
	}()

	// Goroutine to receive messages from the server
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Failed to receive message: %v", err)
				break
			}
			log.Printf("Received from server: %s", res.Reply)
		}
		close(waitc)
	}()

	<-waitc
}