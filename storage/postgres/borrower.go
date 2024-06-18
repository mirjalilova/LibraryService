package postgres

import (
	"database/sql"
	"fmt"
	pb "library_service/genproto/library_service"

	_ "github.com/lib/pq"
	"time"
)

type BorrowerRepo struct {
	db *sql.DB
}

func NewBorrowerRepo(db *sql.DB) *BorrowerRepo {
	return &BorrowerRepo{
		db: db,
	}
}

func (b *BorrowerRepo) Create(req *pb.BorrowerCreateReq) (*pb.BorrowerRes, error) {
	res := &pb.BorrowerRes{
		User: &pb.User{},
		Book: &pb.BookRes{
			Author: &pb.Author{},
            Genre:  &pb.Genre{},
		},
	}

    query := `INSERT INTO borrowers (user_id, book_id) VALUES ($1, $2) RETURNING id`

    row := b.db.QueryRow(query, req.UserId, req.BookId)

    err := row.Scan(&res.Id)
    if err!= nil {
        return nil, fmt.Errorf("Error scaning id: %v", err)
    }

    res, err = b.Get(&pb.GetByIdReq{Id: res.Id})

    return res, err
}

func (b *BorrowerRepo) Get(req *pb.GetByIdReq) (*pb.BorrowerRes, error) {
	res := &pb.BorrowerRes{
		User: &pb.User{},
		Book: &pb.BookRes{
			Author: &pb.Author{},
            Genre:  &pb.Genre{},
		},
	}

    query := `SELECT 
				b.id, 
				bk.id,
				bk.title,
                a.id,
                a.name,
                a.biography,
                g.id,
				g.name,
                bk.summary
				b.borrow_date,
				b.return_date
			FROM borrowers b
			JOIN books bk on b.book_id=bk.id
			JOIN authors a on bk.author_id=a.id
			JOIN genres g on bk.genre_id=g.id
			WHERE id = $1`

	
	var bDate, rDate time.Time
    row := b.db.QueryRow(query, req.Id)
    err := row.Scan(
        &res.Id,
        &res.Book.Id,
		&res.Book.Title,
        &res.Book.Author.Id,
        &res.Book.Author.Name,
        &res.Book.Author.Biography,
        &res.Book.Genre.Id,
		&res.Book.Genre.Name,
		&res.Book.Summary,
        &bDate,
		&rDate,
    )

	res.BorrowDate = bDate.Format("2006-01-02 15:04:05")
	res.ReturnDate = rDate.Format("2006-01-02 15:04:05")

    if err!= nil {
        return nil, fmt.Errorf("can't get borrower: %w", err)
    }

    return res, nil
}

func (b *BorrowerRepo) GetAll(req *pb.BorrowerGetAllReq) (*pb.BorrowerGetAllRes, error) {
	res := &pb.BorrowerGetAllRes{
        Borrowers: []*pb.BorrowerRes{},
    }

    query := `SELECT 
                b.id, 
                bk.id,
                bk.title,
                a.id,
                a.name,
                a.biography,
                g.id,
                g.name,
                bk.summary
                b.borrow_date,
                b.return_date
            FROM borrowers b
            JOIN books bk on b.book_id=bk.id
            JOIN authors a on bk.author_id=a.id
            JOIN genres g on bk.genre_id=g.id
            WHERE b.deleted_at=0`

    rows, err := b.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error while getting borrowers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		borrower := &pb.BorrowerRes{
			User: &pb.User{},
			Book: &pb.BookRes{
				Author: &pb.Author{},
				Genre:  &pb.Genre{},
			},
		}

		var bDate, rDate time.Time
		err := rows.Scan(
			&borrower.Id,
			&borrower.Book.Id,
			&borrower.Book.Title,
			&borrower.Book.Author.Id,
			&borrower.Book.Author.Name,
			&borrower.Book.Author.Biography,
			&borrower.Book.Genre.Id,
			&borrower.Book.Genre.Name,
			&borrower.Book.Summary,
			&bDate,
			&rDate,
		)

		borrower.BorrowDate = bDate.Format("2006-01-02 15:04:05")
		borrower.ReturnDate = rDate.Format("2006-01-02 15:04:05")

		if err!= nil {
			return nil, fmt.Errorf("can't get borrower: %w", err)
		}

		res.Borrowers = append(res.Borrowers, borrower)
	}

	return res, nil
}

func (b *BorrowerRepo) Update(req *pb.BorrowerUpdateReq) (*pb.BorrowerRes, error) {
	res := &pb.BorrowerRes{
		User: &pb.User{},
        Book: &pb.BookRes{
            Author: &pb.Author{},
            Genre:  &pb.Genre{},
        },
    }

	query := `UPDATE borrowers SET user_id=$1, book_id=$2, borrow_date=$3, return_date=$4, updated_at=now() WHERE id=$5 RETURNING id`

	var id string

	row := b.db.QueryRow(query, req.UpdateBorrower.UserId, req.UpdateBorrower.BookId, req.UpdateBorrower.BorrowDate, req.UpdateBorrower.ReturnDate, req.Id.Id)
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("Error scaning id: %v", err)
    }
	
	res, err = b.Get(&pb.GetByIdReq{Id: id})
	if err!= nil {
        return nil, fmt.Errorf("can't get borrower: %w", err)
    }

	return res, nil
}

func (b *BorrowerRepo) Delete(req *pb.GetByIdReq) (*pb.Void, error) {
	res := &pb.Void{}

    query := `UPDATE borrowers SET deleted_at=EXTRACT(EPOACH FROM NOW()) WHERE id=$1`

	b.db.Exec(query, req.Id)

    return res, nil
}

func (b *BorrowerRepo) GetOverdueBooks(req *pb.Void) (*pb.BorrowerGetAllRes, error) {
	res := &pb.BorrowerGetAllRes{
        Borrowers: []*pb.BorrowerRes{},
    }

    query := `SELECT 
                b.id, 
                bk.id,
                bk.title,
                a.id,
                a.name,
                a.biography,
                g.id,
                g.name,
                bk.summary
                b.borrow_date,
                b.return_date
            FROM borrowers b
            JOIN books bk on b.book_id=bk.id
            JOIN authors a on bk.author_id=a.id
            JOIN genres g on bk.genre_id=g.id
            WHERE b.deleted_at=0 AND b.return_date < now()`

    rows, err := b.db.Query(query)
    if err!= nil {
        return nil, fmt.Errorf("Error while getting borrowers: %v", err)
    }

	defer rows.Close()

	for rows.Next() {
		borrower := &pb.BorrowerRes{
			User: &pb.User{},
            Book: &pb.BookRes{
                Author: &pb.Author{},
                Genre:  &pb.Genre{},
            },
		}

		var bDate, rDate time.Time
		err := rows.Scan(
			&borrower.Id,
			&borrower.Book.Id,
			&borrower.Book.Title,
			&borrower.Book.Author.Id,
			&borrower.Book.Author.Name,
			&borrower.Book.Author.Biography,
			&borrower.Book.Genre.Id,
			&borrower.Book.Genre.Name,
			&borrower.Book.Summary,
			&bDate,
			&rDate,
		)

		borrower.BorrowDate = bDate.Format("2006-01-02 15:04:05")
		borrower.ReturnDate = rDate.Format("2006-01-02 15:04:05")

		if err!= nil {
			return nil, fmt.Errorf("can't get borrower: %w", err)
		}

		res.Borrowers = append(res.Borrowers, borrower)
	}

	return res, nil
}