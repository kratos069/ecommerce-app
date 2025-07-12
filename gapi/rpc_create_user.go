package gapi

import (
	"context"

	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/pb"
	"github.com/e-commerce/util"
	"github.com/e-commerce/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context,
	req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

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

func validateCreateUserRequest(req *pb.CreateUserRequest) (
	violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
