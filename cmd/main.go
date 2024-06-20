package main

import (
	"golang.org/x/exp/slog"
	"net"

	cf "library_service/config"

	pb "library_service/genproto/library_service"
	service "library_service/service"
	"library_service/storage/postgres"

	"google.golang.org/grpc"
)

func main() {
	config := cf.Load()
	db, err := postgres.Connect(config)
	if err != nil {
		slog.Error("can't connect to db: %v", err)
		return
	}

	listener, err := net.Listen("tcp", config.LIBRARY_PORT)
	if err != nil {
		slog.Error("can't listen: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterAuthorServiceServer(s, service.NewAuthorService(db))
	pb.RegisterGenreServiceServer(s, service.NewGenreService(db))
	pb.RegisterBookServiceServer(s, service.NewBookService(db))
	pb.RegisterBorrowerServiceServer(s, service.NewBorrowerService(db))

	slog.Info("server started port", config.LIBRARY_PORT)
	if err := s.Serve(listener); err != nil {
		slog.Error("can't serve: %v", err)
		return
	}
}
