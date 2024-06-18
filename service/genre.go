package service

import (
    pb "library_service/genproto/library_service"
    st "library_service/storage/postgres"
    "golang.org/x/exp/slog"
	"context"
	"fmt"
)

type GenreService struct {
    storage st.Storage
    pb.UnimplementedGenreServiceServer
}

func NewGenreService(storage *st.Storage) *GenreService {
	return &GenreService{
        storage: *storage,
    }
}

func (s *GenreService) Create(ctx context.Context, req *pb.GenreCreateReq) (*pb.GenreRes, error) {
	genre, err := s.storage.GenreS.Create(req)
    if err != nil {
        slog.Error("can't create genre: %v", err)
        return nil, fmt.Errorf("can't create genre: %w", err)
    }
    return genre, nil
}

func (s *GenreService) Get(ctx context.Context, req *pb.GetByIdReq) (*pb.GenreRes, error) {
	genre, err := s.storage.GenreS.Get(req)
    if err != nil {
        slog.Error("can't get genre: %v", err)
        return nil, fmt.Errorf("can't get genre: %w", err)
    }
    return genre, nil
}

func (s *GenreService) GetAll(ctx context.Context, req *pb.GenreGetAllReq) (*pb.GenreGetAllRes, error) {
	genres, err := s.storage.GenreS.GetAll(req)
    if err != nil {
        slog.Error("can't get all genres: %v", err)
        return nil, fmt.Errorf("can't get all genres: %w", err)
    }
    return genres, nil
}

func (s *GenreService) Update(ctx context.Context, req *pb.GenreUpdateReq) (*pb.GenreRes, error) {
	genre, err := s.storage.GenreS.Update(req)
    if err != nil {
        slog.Error("can't update genre: %v", err)
        return nil, fmt.Errorf("can't update genre: %w", err)
    }
    return genre, nil
}

func (s *GenreService) Delete(ctx context.Context, req *pb.GetByIdReq) (*pb.Void, error) {
	_, err := s.storage.GenreS.Delete(req)
    if err != nil {
        slog.Error("can't delete genre: %v", err)
        return nil, fmt.Errorf("can't delete genre: %w", err)
    }
    return &pb.Void{}, nil
}