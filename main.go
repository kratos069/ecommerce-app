package main

import (
	"context"
	"net"
	"net/http"
	"os"

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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	// load ENV variables
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config file:")
	}

	// pretty logs in dev (instead of json)
	if config.Environment == "development" {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// conn to database
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to the db")
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
		log.Fatal().Msg("cannot create new migration instance:")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("failed to run migrate up:")
	}

	log.Info().Msg("db migrated successfully")
}

// run both grpc + http requests
func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create grpc server:")
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
		log.Fatal().Msg("cannot register handler server:")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// serve static resources from backend (swagger folder)
	// fs := http.FileServer(http.Dir("./doc/swagger"))

	// server static resources as binary of Go's backend server
	staticFS, err := fs.New()
	if err != nil {
		log.Fatal().Msg("cannot create statik fs:")
	}

	swaggerHandler := http.StripPrefix("/swagger/",
		http.FileServer(staticFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener:")
	}

	log.Info().Msgf("start HTTP gateway server at %s",
		listener.Addr().String())

	handler := gapi.HttpLogger(mux) // logs of http requests
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot start HTTP gateway server:")
	}

}

// run grpc Server for grpc requests
func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create grpc server:")
	}

	// logs for server
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)

	// 10 mb max streaming RPC msg size
	grpcServer := grpc.NewServer(grpcLogger, grpc.MaxRecvMsgSize(10<<20))
	pb.RegisterEcommerceServer(grpcServer, server)
	reflection.Register(grpcServer) // client to know which funcs available on server

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener:")
	}

	log.Info().Msgf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start grpc server:")
	}

}

// run Gin Server for HTTP requests
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create gin server:")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server:")
	}
}
