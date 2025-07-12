package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/e-commerce/api"
	db "github.com/e-commerce/db/sqlc"
	_ "github.com/e-commerce/doc/statik"
	"github.com/e-commerce/gapi"
	"github.com/e-commerce/pb"
	"github.com/e-commerce/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	go runGatewayServer(config, store)
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

// run both grpc + http requests
func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalln("cannot create grpc server:", err)
	}

	// for having snake case for http response
	jsonOption := runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterEcommerceHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalln("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// serve static resources from backend (swagger folder)
	// fs := http.FileServer(http.Dir("./doc/swagger"))

	// server static resources as binary of Go's backend server
	staticFS, err := fs.New()
	if err != nil {
		log.Fatalln("cannot create statik fs:", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/",
		http.FileServer(staticFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatalln("cannot create listener:", err)
	}

	log.Printf("start HTTP gateway server at %s",
		listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalln("cannot start HTTP gateway server:", err)
	}

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
