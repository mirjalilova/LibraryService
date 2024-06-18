package service

import (
    pb "library_service/genproto/library_service"
    st "library_service/storage/postgres"
    "golang.org/x/exp/slog"
    "context"
    "fmt"
)

type AuthorService struct {
	storage st.Storage
	pb.UnimplementedAuthorServiceServer
}

func NewAuthorService(storage *st.Storage) *AuthorService {
	return &AuthorService{
		storage: *storage,
	}
}

func (s *AuthorService) Create(ctx context.Context, req *pb.AuthorCreateReq) (*pb.AuthorRes, error) {
    author, err := s.storage.AuthorS.Create(req)
    if err != nil {
        slog.Error("can't create author: %v", err)
        return nil, fmt.Errorf("can't create author: %w", err)
    }
    return author, nil
}

func (s *AuthorService) Get(ctx context.Context, req *pb.GetByIdReq) (*pb.AuthorRes, error) {
    author, err := s.storage.AuthorS.Get(req)
    if err != nil {
        slog.Error("can't get author: %v", err)
        return nil, fmt.Errorf("can't get author: %w", err)
    }
    return author, nil
}

func (s *AuthorService) GetAll(ctx context.Context, req *pb.AuthorGetAllReq) (*pb.AuthorGetAllRes, error) {
    authors, err := s.storage.AuthorS.GetAll(req)
    if err != nil {
        slog.Error("can't get all authors: %v", err)
        return nil, fmt.Errorf("can't get all authors: %w", err)
    }
    return authors, nil
}


func (s *AuthorService) Update(ctx context.Context, req *pb.AuthorUpdateReq) (*pb.AuthorRes, error) {
    author, err := s.storage.AuthorS.Update(req)
    if err != nil {
        slog.Error("can't update author: %v", err)
        return nil, fmt.Errorf("can't update author: %w", err)
    }
    return author, nil
}

func (s *AuthorService) Delete(ctx context.Context, req *pb.GetByIdReq) (*pb.Void, error) {
    _, err := s.storage.AuthorS.Delete(req)
    if err != nil {
        slog.Error("can't delete author: %v", err)
        return nil, fmt.Errorf("can't delete author: %w", err)
    }
    return &pb.Void{}, nil
}


