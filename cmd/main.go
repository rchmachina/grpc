package main

import (
	"log"
	"net"

	"github.com/joho/godotenv"
	"github.com/rchmachina/grpc/cmd/config/database"
	mw "github.com/rchmachina/grpc/cmd/config/auth"
	"github.com/rchmachina/grpc/cmd/services"
	userPb "github.com/rchmachina/grpc/dto/authpb"
	"google.golang.org/grpc"
	streamPb "github.com/rchmachina/grpc/dto/streamingService"
	

)

const (
	port = ":55001"
)

func main() {
	var err error
	//
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error: failed to load the .env file")
	}

	log.Println("Hello there!")
	log.Println("Connected to database")
	netListen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err, " failed to listening port: ", port)
	}
	db := database.DatabaseConnection()
	//grpcServer := grpc.NewServer()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(mw.AuthInterceptor([]string{
			 // Add methods that require authentication
			"/auth.AuthService/TestingMw",
			"/StreamingService.StreamingService/BidirectionalStreaming",
		})),
		grpc.StreamInterceptor(mw.AuthStreamInterceptor([]string{
			// Methods that require authentication for streaming calls
			"/StreamingService.StreamingService/BidirectionalStreaming",
		})),
	)
	authService := services.AuthService{Db: db}
	streamService := services.StreamingService{Db:db}
	userPb.RegisterAuthServiceServer(grpcServer, &authService)
	streamPb.RegisterStreamingServiceServer(grpcServer, &streamService)


	log.Printf("server started at %v", netListen.Addr())
	if err := grpcServer.Serve(netListen); err != nil {
		log.Fatal(" failed to serve: ", err.Error())

	}

}
