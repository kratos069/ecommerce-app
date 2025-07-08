package gapi

import (
	"context"

	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/pb"
	"github.com/e-commerce/util"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context,
	req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound,
				"email doesnot exist")

		}

		return nil, status.Errorf(codes.Internal,
			"unable to get user from db")

	}

	err = util.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied,
			"wrong password")
	}

	// creating access token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.UserID,
		user.Role,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"unable to create access token")
	}

	// creating refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		user.UserID,
		user.Role,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"unable to create refresh token")
	}

	// create session
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     accessPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"unable to create session")
	}

	resp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return resp, nil
}
