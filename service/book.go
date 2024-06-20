package service

import (
	"context"
	"fmt"
	pb "library_service/genproto/library_service"
	st "library_service/storage/postgres"

	"golang.org/x/exp/slog"
)

type BookService struct {
	storage st.Storage
	pb.UnimplementedBookServiceServer
}

func NewBookService(storage *st.Storage) *BookService {
	return &BookService{
		storage: *storage,
	}
}

func (s *BookService) Create(ctx context.Context, req *pb.BookCreateReq) (*pb.BookRes, error) {
	book, err := s.storage.BookS.Create(req)
	if err != nil {
		slog.Error("can't create book: %v", err)
		return nil, fmt.Errorf("can't create book: %w", err)
	}
	return book, nil
}

func (s *BookService) Get(ctx context.Context, req *pb.GetByIdReq) (*pb.BookRes, error) {
	book, err := s.storage.BookS.Get(req)
	if err != nil {
		slog.Error("can't get book: %v", err)
		return nil, fmt.Errorf("can't get book: %w", err)
	}
	return book, nil
}

func (s *BookService) GetAll(ctx context.Context, req *pb.BookGetAllReq) (*pb.BookGetAllRes, error) {
	books, err := s.storage.BookS.GetAll(req)
	if err != nil {
		slog.Error("can't get all books: %v", err)
		return nil, fmt.Errorf("can't get all books: %w", err)
	}
	return books, nil
}

func (s *BookService) Update(ctx context.Context, req *pb.BookUpdateReq) (*pb.BookRes, error) {

	old_book, err := s.storage.BookS.Get(&pb.GetByIdReq{
		Id: req.Id.Id,
	})
	if err != nil {
		slog.Error("can't get book: %v", err)
		return nil, fmt.Errorf("can't get book: %w", err)
	}

	if req.UpdateBook.Title == "string" {
		req.UpdateBook.Title = old_book.Title
	}
	if req.UpdateBook.Summary == "string" {
		req.UpdateBook.Summary = old_book.Summary
	}
	if req.UpdateBook.GenreId == "string" {
		req.UpdateBook.GenreId = old_book.Genre.Id
	}
	if req.UpdateBook.AuthorId == "string" {
		req.UpdateBook.AuthorId = old_book.Author.Id
	}

	book, err := s.storage.BookS.Update(req)
	if err != nil {
		slog.Error("can't update book: %v", err)
		return nil, fmt.Errorf("can't update book: %w", err)
	}
	return book, nil
}

func (s *BookService) Delete(ctx context.Context, req *pb.GetByIdReq) (*pb.Void, error) {
	_, err := s.storage.BookS.Delete(req)
	if err != nil {
		slog.Error("can't delete book: %v", err)
		return nil, fmt.Errorf("can't delete book: %w", err)
	}
	return &pb.Void{}, nil
}

func (s *BookService) Search(ctx context.Context, req *pb.BookSearchReq) (*pb.BookGetAllRes, error) {
	books, err := s.storage.BookS.Search(req)
	if err != nil {
		slog.Error("can't search book: %v", err)
		return nil, fmt.Errorf("can't search book: %w", err)
	}
	return books, nil
}
