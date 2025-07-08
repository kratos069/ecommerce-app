package gapi

import (
	"fmt"

	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/pb"
	"github.com/e-commerce/token"
	"github.com/e-commerce/util"
)

// servers gRPC requests for the insta-app
type Server struct {
	pb.UnimplementedEcommerceServer // for forward compatibility
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	// taskDistributor worker.TaskDistributor
}

// Creates gRPC server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		// taskDistributor: taskDistributor,
	}

	return server, nil
}
