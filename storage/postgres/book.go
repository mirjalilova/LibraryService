package postgres

import (
	"database/sql"
	"fmt"
	pb "library_service/genproto/library_service"
	"strings"

	_ "github.com/lib/pq"
)

type BookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) *BookRepo {
	return &BookRepo{
		db: db,
	}
}

func (b *BookRepo) Create(req *pb.BookCreateReq) (*pb.BookRes, error) {
	res := &pb.BookRes{
		Author: &pb.Author{},
		Genre:  &pb.Genre{},
	}

	query := `INSERT INTO books (
				title, 
				author_id, 
				genre_id, 
				summary
			) VALUES ($1, $2, $3, $4) 
			RETURNING 
				id`

	row := b.db.QueryRow(query, req.Title, req.AuthorId, req.GenreId, req.Summary)

	var id string
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("Error scaning id: %v", err)
	}

	fmt.Println(id, res)

	res, err = b.Get(&pb.GetByIdReq{Id: id})
	if err!= nil {
        return nil, fmt.Errorf("can't get book: %w", err)
    }

	return res, nil
}

func (b *BookRepo) Get(req *pb.GetByIdReq) (*pb.BookRes, error) {
	res := &pb.BookRes{
		Author: &pb.Author{},
		Genre:  &pb.Genre{},
	}

	query := `SELECT
				b.id, 
				b.title, 
				a.id, 
				a.name, 
				a.biography, 
				g.id, 
				g.name, 
				b.summary 
			FROM books b 
			JOIN authors a on b.author_id = a.id 
			JOIN genres g on b.genre_id = g.id 
			WHERE b.id = $1 AND b.deleted_at=0`
	row := b.db.QueryRow(query, req.Id)

	err := row.Scan(
		&res.Id,
		&res.Title,
		&res.Author.Id,
		&res.Author.Name,
		&res.Author.Biography,
		&res.Genre.Id,
		&res.Genre.Name,
		&res.Summary,
	)

	if err != nil {
		return nil, fmt.Errorf("can't get book: %w", err)
	}

	return res, nil
}

func (b *BookRepo) GetAll(req *pb.BookGetAllReq) (*pb.BookGetAllRes, error) {
	res := &pb.BookGetAllRes{
		Books: []*pb.BookRes{},
	}

	query := `SELECT 
                b.id, 
                b.title, 
                a.id, 
                a.name, 
                a.biography, 
                g.id, 
                g.name, 
                b.summary 
            FROM books b 
            JOIN authors a on b.author_id = a.id 
            JOIN genres g on b.genre_id = g.id 
            WHERE b.deleted_at=0`

	var args []interface{}
	var conditions []string

	if req.AuthorId != "" {
		args = append(args, req.AuthorId)
		conditions = append(conditions, fmt.Sprintf("author_id = $%d", len(args)))
	}
	if req.GenreId != "" {
		args = append(args, req.GenreId)
		conditions = append(conditions, fmt.Sprintf("genre_id = $%d", len(args)))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var defaultLimit int32
	err := b.db.QueryRow("SELECT COUNT(1) FROM books WHERE deleted_at=0").Scan(&defaultLimit)
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
		return nil, fmt.Errorf("Error query: %v", err)
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
			return nil, fmt.Errorf("Error while scaning books : %v", err)
		}

		res.Books = append(res.Books, book)
	}

	return res, nil
}

func (b *BookRepo) Update(req *pb.BookUpdateReq) (*pb.BookRes, error) {
	res := &pb.BookRes{
		Author: &pb.Author{},
		Genre:  &pb.Genre{},
	}

	query := `UPDATE books SET 
                title = $1, 
                author_id = $2, 
                genre_id = $3, 
                summary = $4,
				updated_at = now()
            WHERE id = $5 
            RETURNING 
                id`

	row := b.db.QueryRow(query, req.UpdateBook.Title, req.UpdateBook.AuthorId, req.UpdateBook.GenreId, req.UpdateBook.Summary, req.Id.Id)

	var id string
	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("Error scaning id: %v", err)
	}

	res, err = b.Get(&pb.GetByIdReq{Id: id})

	return res, nil
}

func (b *BookRepo) Delete(req *pb.GetByIdReq) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE books SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1`

	_, err := b.db.Exec(query, req.Id)

	if err != nil {
		return nil, fmt.Errorf("can't delete book: %w", err)
	}

	return res, nil
}

func (b *BookRepo) Search(req *pb.BookSearchReq) (*pb.BookGetAllRes, error) {

	res := &pb.BookGetAllRes{}

	query := `SELECT 
				b.id, 
				b.title, 
				a.id, 
				a.name, 
				a.biography, 
				g.id, 
				g.name, 
				b.summary 
			FROM books b 
			JOIN authors a on b.author_id = a.id 
			JOIN genres g on b.genre_id = g.id 
			WHERE b.deleted_at=0`

	var args []interface{}
	var conditions []string

	if req.Title != "" {
		args = append(args, "%"+req.Title+"%")
		query += fmt.Sprintf(" AND b.title ILIKE $%d", len(args))
	}
	if req.Author != "" {
		args = append(args, "%"+req.Author+"%")
		query += fmt.Sprintf(" OR a.name ILIKE $%d", len(args))
		
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	fmt.Println(query)
	rows, err := b.db.Query(query, args...)

	if err != nil {
		return nil, fmt.Errorf("Error while searching books: %v", err)
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
			return nil, fmt.Errorf("Error while scaning books : %v", err)
		}

		res.Books = append(res.Books, book)
	}

	return res, nil
}
