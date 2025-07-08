package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/e-commerce/api"
	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/gapi"
	"github.com/e-commerce/pb"
	"github.com/e-commerce/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// load ENV variables
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln("cannot load config file:", err)
	}

	// conn to database
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalln("cannot connect to the db:", err)
	}

	// run db migrations
	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(connPool)

	// runGinServer(config, store)
	runGrpcServer(config, store)
}

func runDBMigration(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalln("cannot create new migration instance:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalln("failed to run migrate up:", err)
	}

	fmt.Println("db migrated successfully")
}

// run grpc Server for grpc requests
func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create grpc server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEcommerceServer(grpcServer, server)
	reflection.Register(grpcServer) // client to know which funcs available on server

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalln("cannot create listener:", err)
	}

	log.Printf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln("cannot start grpc server:", err)
	}

}

// run Gin Server for HTTP requests
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create gin server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatalln("cannot start server:", err)
	}
}
