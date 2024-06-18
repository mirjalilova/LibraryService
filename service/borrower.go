package service

import (
	"context"
	"fmt"
	pb "library_service/genproto/library_service"
	st "library_service/storage/postgres"

	"golang.org/x/exp/slog"
)

type BorrowerService struct {
	storage st.Storage
	pb.UnimplementedBorrowerServiceServer
}

func NewBorrowerService(storage *st.Storage) *BorrowerService {
	return &BorrowerService{
		storage: *storage,
	}
}

func (s *BorrowerService) Create(ctx context.Context, req *pb.BorrowerCreateReq) (*pb.BorrowerRes, error) {
	borrower, err := s.storage.BorrowerS.Create(req)
	if err != nil {
		slog.Error("can't create borrower: %v", err)
		return nil, fmt.Errorf("can't create borrower: %w", err)
	}
	return borrower, nil
}

func (s *BorrowerService) Get(ctx context.Context, req *pb.GetByIdReq) (*pb.BorrowerRes, error) {
	borrower, err := s.storage.BorrowerS.Get(req)
	if err != nil {
		slog.Error("can't get borrower: %v", err)
		return nil, fmt.Errorf("can't get borrower: %w", err)
	}
	return borrower, nil
}

func (s *BorrowerService) GetAll(ctx context.Context, req *pb.BorrowerGetAllReq) (*pb.BorrowerGetAllRes, error) {

	borrowers, err := s.storage.BorrowerS.GetAll(req)
	if err != nil {
		slog.Error("can't get all borrowers: %v", err)
		return nil, fmt.Errorf("can't get all borrowers: %w", err)
	}

	return borrowers, nil
}

func (s *BorrowerService) Update(ctx context.Context, req *pb.BorrowerUpdateReq) (*pb.BorrowerRes, error) {
	borrower, err := s.storage.BorrowerS.Update(req)
	if err != nil {
		slog.Error("can't update borrower: %v", err)
		return nil, fmt.Errorf("can't update borrower: %w", err)
	}
	return borrower, nil
}

func (s *BorrowerService) Delete(ctx context.Context, req *pb.GetByIdReq) (*pb.Void, error) {
	_, err := s.storage.BorrowerS.Delete(req)
	if err != nil {
		slog.Error("can't delete borrower: %v", err)
		return nil, fmt.Errorf("can't delete borrower: %w", err)
	}
	return &pb.Void{}, nil
}

func (s *BorrowerService) GetOverdueBooks(ctx context.Context, req *pb.Void) (*pb.BorrowerGetAllRes, error) {
	borrowers, err := s.storage.BorrowerS.GetOverdueBooks(req)
    if err!= nil {
        slog.Error("can't get overdue books: %v", err)
        return nil, fmt.Errorf("can't get overdue books: %w", err)
    }
    return borrowers, nil
}

