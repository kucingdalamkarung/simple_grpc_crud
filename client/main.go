package main

import (
	"context"
	"fmt"
	model "grpc_crud/proto"
	"io"
	"log"

	"google.golang.org/grpc"
)

func deleteUser() {
	requestOptions := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:9000", requestOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	client := model.NewUserServiceClient(conn)
	res, err := client.DeleteUser(context.TODO(), &model.DeleteUserReq{Id: "5f08667523fd5cf9c9562708"})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(res.Success)
}

func updateUser() {
	newData := &model.UpdateUserReq{
		User: &model.User{
			Id: "5f0867e923fd5cf9c9562709",
			Name: "alya",
			Email: "fiqrikhoirul@yahoo.com",
			Address: "Bandung barat",
		},
	}

	reqOption := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:9000", reqOption)
	if err != nil {
		log.Fatal(err.Error())
	}

	client := model.NewUserServiceClient(conn)
	res, err := client.UpdateUser(context.Background(), newData)
	fmt.Println(res)
}

func createUser() {
	user := &model.CreateUserReq{
		User: &model.User{
			Email: "fiqrikhoirul@yahoo.com",
			Name: "fiqri",
			Address: "Bandung barat",
		},
	}

	requestOptions := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:9000", requestOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	client := model.NewUserServiceClient(conn)
	res, err := client.Createuser(context.TODO(), user)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(res.User)
}

func getUser() {
	reqOptions := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:9000", reqOptions)
	if err != nil {
		log.Fatal(err.Error())
	}

	client := model.NewUserServiceClient(conn)
	res, err := client.GetUser(context.Background(), &model.GetUserReq{Id: "5f0867e923fd5cf9c9562709"})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(res)
}

func listUser() {
	reqOptions := grpc.WithInsecure()
	conn, _ := grpc.Dial("localhost:9000", reqOptions)
	client := model.NewUserServiceClient(conn)
	stream, _ := client.ListUsers(context.Background(), &model.ListUsersReq{})

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(res.GetUser())
	}
}

func main() {
	listUser()
}
