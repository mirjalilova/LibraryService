package postgres

import (
	"database/sql"
	"fmt"
	"library_service/client"
	pb "library_service/genproto/library_service"

	"time"

	_ "github.com/lib/pq"
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
	if err != nil {
		return nil, fmt.Errorf("Error scaning id: %v", err)
	}

	res, err = b.Get(&pb.GetByIdReq{Id: res.Id})
	if err != nil {
		return nil, fmt.Errorf("can't get borrower: %w", err)
	}

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
				b.id AS borrower_id,
				bk.id AS book_id,
				bk.title AS book_title,
				a.id AS author_id,
				a.name AS author_name,
				a.biography AS author_biography,
				g.id AS genre_id,
				g.name AS genre_name,
				bk.summary AS book_summary,
				b.borrow_date,
				b.return_date
			FROM borrowers b
			JOIN books bk ON b.book_id = bk.id
			JOIN authors a ON bk.author_id = a.id
			JOIN genres g ON bk.genre_id = g.id
			WHERE b.id = $1`

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

	query = `SELECT user_id from borrowers WHERE id = $1`
	row = b.db.QueryRow(query, req.Id)

	var userId string
	err = row.Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("Error while getting user_id : %v", err)
	}

	us, err := client.GetUser(userId)
	if err != nil {
		return nil, fmt.Errorf("Error while getting user : %v", err)
	}

	res.User.Id = us.ID
	res.User.Username = us.Username
	res.User.Email = us.Email

	res.BorrowDate = bDate.Format("2006-01-02 15:04:05")
	res.ReturnDate = rDate.Format("2006-01-02 15:04:05")

	if err != nil {
		return nil, fmt.Errorf("can't get borrower: %w", err)
	}

	return res, nil
}

func (b *BorrowerRepo) GetAll(req *pb.BorrowerGetAllReq) (*pb.BorrowerGetAllRes, error) {
	res := &pb.BorrowerGetAllRes{
		Borrowers: []*pb.BorrowerRes{},
	}

	query := `SELECT
				b.id AS borrower_id,
				bk.id AS book_id,
				bk.title AS book_title,
				a.id AS author_id,
				a.name AS author_name,
				a.biography AS author_biography,
				g.id AS genre_id,
				g.name AS genre_name,
				bk.summary AS book_summary,
				b.borrow_date,
				b.return_date
			FROM borrowers b
			JOIN books bk ON b.book_id = bk.id
			JOIN authors a ON bk.author_id = a.id
			JOIN genres g ON bk.genre_id = g.id
			WHERE b.deleted_at = 0
			`

	var args []interface{}

	if req.UserId != "" {
		args = append(args, req.UserId)
		query += fmt.Sprintf(" AND b.user_id = $%d", len(args))
	}
	if req.BookId != "" {
		args = append(args, req.BookId)
		query += fmt.Sprintf(" AND b.book_id = $%d", len(args))
	}
	if req.BorrowDate != "" {
		args = append(args, req.BorrowDate)
		query += fmt.Sprintf(" AND b.borrow_date > $%d", len(args))
	}
	if req.ReturnDate != "" {
		args = append(args, req.ReturnDate)
		query += fmt.Sprintf(" AND b.return_date > $%d", len(args))
	}

	var defaultLimit int32
	err := b.db.QueryRow("SELECT COUNT(1) FROM borrowers WHERE deleted_at=0").Scan(&defaultLimit)
	if err != nil {
		return nil, err
	}
	if req.Filter.Limit == 0 {
		req.Filter.Limit = defaultLimit
	}

	args = append(args, req.Filter.Limit, req.Filter.Offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := b.db.Query(query, args...)
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

		query = `SELECT user_id from borrowers WHERE id = $1`
		row := b.db.QueryRow(query, borrower.Id)

		var userId string
		err = row.Scan(&userId)
		if err != nil {
			return nil, fmt.Errorf("Error while getting user id : %v", err)
		}

		us, err := client.GetUser(userId)
		if err != nil {
			return nil, fmt.Errorf("Error while getting user : %v", err)
		}

		borrower.User.Id = us.ID
		borrower.User.Username = us.Username
		borrower.User.Email = us.Email

		borrower.BorrowDate = bDate.Format("2006-01-02 15:04:05")
		borrower.ReturnDate = rDate.Format("2006-01-02 15:04:05")

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
	if err != nil {
		return nil, fmt.Errorf("can't get borrower: %w", err)
	}

	return res, nil
}

func (b *BorrowerRepo) Delete(req *pb.GetByIdReq) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE borrowers SET deleted_at=EXTRACT(EPOCH FROM NOW()) WHERE id=$1`

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
            bk.summary,
            b.borrow_date,
            b.return_date
          FROM borrowers b
          JOIN books bk ON b.book_id = bk.id
          JOIN authors a ON bk.author_id = a.id
          JOIN genres g ON bk.genre_id = g.id
          WHERE b.deleted_at = 0 AND b.return_date < NOW()`

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

		if err != nil {
			return nil, fmt.Errorf("can't get borrower: %w", err)
		}

		res.Borrowers = append(res.Borrowers, borrower)
	}

	return res, nil
}

func (b *BorrowerRepo) GetBorrowedBooks(req *pb.BorrowedBooksReq) (*pb.BorrowedBooksRes, error) {
	res := &pb.BorrowedBooksRes{
		Books: []*pb.BookRes{},
	}

	query := `SELECT
                bk.id,
                bk.title,
                a.id,
                a.name,
                a.biography,
                g.id,
                g.name,
                bk.summary
            FROM borrowers b 
			JOIN books bk on b.book_id=bk.id
			JOIN authors a on bk.author_id=a.id
			JOIN genres g on bk.genre_id=g.id
			WHERE b.user_id=$1 AND b.deleted_at=0 AND b.return_date > now()`

	rows, err := b.db.Query(query, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("Error while getting borrowers: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		book := &pb.BookRes{
			Author: &pb.Author{},
			Genre:  &pb.Genre{},
		}

		err := rows.Scan(
			&book.Id,
			&book.Title,
			&book.Author.Id,
			&book.Author.Name,
			&book.Author.Biography,
			&book.Genre.Id,
			&book.Genre.Name,
			&book.Summary,
		)

		if err != nil {
			return nil, fmt.Errorf("can't get borrower: %w", err)
		}

		res.Books = append(res.Books, book)
	}

	return res, err
}

func (b *BorrowerRepo) GetBorrowingHistory(req *pb.BorrowedBooksReq) (*pb.BorrowedBooksRes, error) {
	res := &pb.BorrowedBooksRes{
		Books: []*pb.BookRes{},
	}

	query := `SELECT
				bk.id,
				bk.title,
				a.id,
				a.name,
				a.biography,
				g.id,
				g.name,
				bk.summary
			FROM borrowers b 
			JOIN books bk on b.book_id=bk.id
			JOIN authors a on bk.author_id=a.id
			JOIN genres g on bk.genre_id=g.id
			WHERE b.user_id=$1 AND b.deleted_at=0`

	rows, err := b.db.Query(query, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("Error while getting borrowers: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		book := &pb.BookRes{
			Author: &pb.Author{},
			Genre:  &pb.Genre{},
		}

		err := rows.Scan(
			&book.Id,
			&book.Title,
			&book.Author.Id,
			&book.Author.Name,
			&book.Author.Biography,
			&book.Genre.Id,
			&book.Genre.Name,
			&book.Summary,
		)

		if err != nil {
			return nil, fmt.Errorf("can't get borrower: %w", err)
		}

		res.Books = append(res.Books, book)
	}

	return res, err
}
