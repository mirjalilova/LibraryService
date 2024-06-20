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

	old_borrower, err := s.storage.BorrowerS.Get(&pb.GetByIdReq{
		Id: req.Id.Id,
	})
	if err != nil {
		slog.Error("can't get borrower: %v", err)
		return nil, fmt.Errorf("can't get borrower: %w", err)
	}

	if req.UpdateBorrower.UserId == "string" {
		req.UpdateBorrower.UserId = old_borrower.User.Id
	}
	if req.UpdateBorrower.BookId == "string" {
		req.UpdateBorrower.BookId = old_borrower.Book.Id
	}
	if req.UpdateBorrower.BorrowDate == "string" {
		req.UpdateBorrower.BorrowDate = old_borrower.BorrowDate
	}
	if req.UpdateBorrower.ReturnDate == "string" {
		req.UpdateBorrower.ReturnDate = old_borrower.ReturnDate
	}

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
	if err != nil {
		slog.Error("can't get overdue books: %v", err)
		return nil, fmt.Errorf("can't get overdue books: %w", err)
	}
	return borrowers, nil
}

func (s *BorrowerService) GetBorrowedBooks(ctx context.Context, req *pb.BorrowedBooksReq) (*pb.BorrowedBooksRes, error) {
	borrowers, err := s.storage.BorrowerS.GetBorrowedBooks(req)
	if err != nil {
		slog.Error("can't get borrowed books: %v", err)
		return nil, fmt.Errorf("can't get borrowed books: %w", err)
	}
	return borrowers, nil
}

func (s *BorrowerService) GetBorrowingHistory(ctx context.Context, req *pb.BorrowedBooksReq) (*pb.BorrowedBooksRes, error) {
	borrowers, err := s.storage.BorrowerS.GetBorrowingHistory(req)
	if err != nil {
		slog.Error("can't get borrowing history: %v", err)
		return nil, fmt.Errorf("can't get borrowing history: %w", err)
	}
	return borrowers, nil
}
