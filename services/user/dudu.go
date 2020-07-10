package main

import (
	"context"
	"fmt"
	model "grpc_crud/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {}
type UserData struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name"`
	Email string `bson:"email"`
	Address string `bson:"address"`
}

func (u *UserServiceServer) Createuser(ctx context.Context, req *model.CreateUserReq) (*model.CreateUserRes, error) {
	usr := req.GetUser()
	data := UserData{
		Name: usr.Name,
		Email: usr.Email,
		Address: usr.Address,
	}

	result, err := userdb.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	oid := result.InsertedID.(primitive.ObjectID)
	usr.Id = oid.Hex()

	return &model.CreateUserRes{User: usr}, nil
}

func (u *UserServiceServer) UpdateUser(ctx context.Context, req *model.UpdateUserReq) (*model.UpdateUserRes, error) {
	usr := req.GetUser()
	oid, err := primitive.ObjectIDFromHex(usr.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Could not convert the supplied user id to a MongoDB ObjectId: %v", err))
	}

	update := bson.M{
		"name": usr.GetName(),
		"email": usr.GetEmail(),
		"address": usr.GetAddress(),
	}

	filter := bson.M{"_id": oid}
	result := userdb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))
	decoded := UserData{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Could not find user with supplied ID: %v", err))
	}

	return &model.UpdateUserRes{
		User: &model.User{
			Id: decoded.Id.Hex(),
			Address: decoded.Address,
			Email: decoded.Email,
			Name: decoded.Name,
		},
	}, nil
}

func (u *UserServiceServer) DeleteUser(ctx context.Context, req *model.DeleteUserReq) (*model.DeleteUserRes, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	_, err = userdb.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete user with id %s: %v", req.GetId(), err))
	}

	return &model.DeleteUserRes{
		Success: true,
	}, nil
}

func (u *UserServiceServer) GetUser(ctx context.Context, req *model.GetUserReq) (*model.GetUserRes, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}

	result := userdb.FindOne(ctx, bson.M{"_id": oid})
	data := &UserData{}
	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find user with Object Id %s: %v", req.GetId(), err))
	}

	return &model.GetUserRes{
		User: &model.User{
			Id: oid.Hex(),
			Name: data.Name,
			Email: data.Email,
			Address: data.Address,
		},
	}, nil
}

func (u *UserServiceServer) ListUsers(req *model.ListUsersReq, server model.UserService_ListUsersServer) error {
	data := &UserData{}
	cursor, err := userdb.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		err := cursor.Decode(data)
		if err != nil {
			return status.Error(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}

		_ = server.Send(&model.ListUsersRes{
			User: &model.User{
				Id:      data.Id.Hex(),
				Name:    data.Name,
				Email:   data.Email,
				Address: data.Address,
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}

	return nil
}
