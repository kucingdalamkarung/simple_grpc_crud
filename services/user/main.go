package main

import (
	"context"
	"fmt"
	model "grpc_crud/proto"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var db *mongo.Client
var userdb *mongo.Collection

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("Unable to listen on port :9000")
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	srv := &UserServiceServer{}

	model.RegisterUserServiceServer(s, srv)

	fmt.Println("Connecting to MongoDB")
	mongoCtx := context.Background()
	db, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(os.Getenv("MONGO_HOST")))
	if err != nil {
		log.Fatal("Unable to connect to MongoDB")
	}
	fmt.Println("Connecting successfully")
	userdb = db.Database("go_learn").Collection("users")

	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :9000")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("\nStopping the server...")
	s.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	db.Disconnect(mongoCtx)
	fmt.Println("Done.")
}
