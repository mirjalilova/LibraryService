package storage

import pb "library_service/genproto/library_service"

type StorageI interface {
	Author() AuthorI
	Book() BookI
	Borrower() BorrowerI
	Genre() GenreI
}

type AuthorI interface {
	Create(*pb.AuthorCreateReq) (*pb.AuthorRes, error)
	Get(*pb.GetByIdReq) (*pb.AuthorRes, error)
	GetAll(*pb.AuthorGetAllReq) (*pb.AuthorGetAllRes, error)
	Update(*pb.AuthorUpdateReq) (*pb.AuthorRes, error)
	Delete(*pb.GetByIdReq) (*pb.Void, error)
}

type BookI interface {
	Create(*pb.BookCreateReq) (*pb.BookRes, error)
	Get(*pb.GetByIdReq) (*pb.BookRes, error)
	GetAll(*pb.BookGetAllReq) (*pb.BookGetAllRes, error)
	Update(*pb.BookUpdateReq) (*pb.BookRes, error)
	Delete(*pb.GetByIdReq) (*pb.Void, error)
	Search(*pb.BookSearchReq) (*pb.BookGetAllRes, error)
}

type BorrowerI interface {
	Create(*pb.BorrowerCreateReq) (*pb.BorrowerRes, error)
	Get(*pb.GetByIdReq) (*pb.BorrowerRes, error)
	GetAll(*pb.BorrowerGetAllReq) (*pb.BorrowerGetAllRes, error)
	Update(*pb.BorrowerUpdateReq) (*pb.BorrowerRes, error)
	Delete(*pb.GetByIdReq) (*pb.Void, error)
	GetOverdueBooks(*pb.Void) (*pb.BorrowerGetAllRes, error)
	GetBorrowedBooks(*pb.BorrowedBooksReq) (*pb.BorrowedBooksRes, error)
	GetBorrowingHistory(*pb.BorrowedBooksReq) (*pb.BorrowedBooksRes, error)
}

type GenreI interface {
	Create(*pb.GenreCreateReq) (*pb.GenreRes, error)
	Get(*pb.GetByIdReq) (*pb.GenreRes, error)
	GetAll(*pb.GenreGetAllReq) (*pb.GenreGetAllRes, error)
	Update(*pb.GenreUpdateReq) (*pb.GenreRes, error)
	Delete(*pb.GetByIdReq) (*pb.Void, error)
}