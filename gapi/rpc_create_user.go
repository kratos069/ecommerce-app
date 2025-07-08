package gapi

import (
	"context"

	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/pb"
	"github.com/e-commerce/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context,
	req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to hash password")

	}

	arg := db.CreateUserParams{
		Username:     req.GetUsername(),
		PasswordHash: hashedPassword,
		Email:        req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"failed to create user in db")

	}

	resp := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return resp, nil
}
