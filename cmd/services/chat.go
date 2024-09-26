package services

import (
	"io"
	"log"

	streamPb "github.com/rchmachina/grpc/dto/streamingService"
	"gorm.io/gorm"
)

type StreamingService struct {
	streamPb.UnimplementedStreamingServiceServer
	Db *gorm.DB
}

func (a *StreamingService) BidirectionalStreaming(stream streamPb.StreamingService_BidirectionalStreamingServer) error {
	log.Println("Bidirectional stream started")
	for {
		// Receive message from the client
		req, err := stream.Recv()
		if err == io.EOF {
			// Client has finished sending messages
			return nil
		}
		if err != nil {
			log.Printf("Error receiving stream: %v", err)
			return err
		}

		log.Printf("Received message from client: %s", req.Message)

		// Process the message and send a response
		response := &streamPb.StreamResponse{
			Reply: "Server received: " + req.Message,
		}
		err = stream.Send(response)
		if err != nil {
			log.Printf("Error sending response: %v", err)
			return err
		}
	}
}
