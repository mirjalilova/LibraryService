package postgres

import (
	"database/sql"
	"fmt"
	pb "library_service/genproto/library_service"

	_ "github.com/lib/pq"
)

type AuthorRepo struct {
	db *sql.DB
}

func NewAuthorRepo(db *sql.DB) *AuthorRepo {
	return &AuthorRepo{
		db: db,
	}
}

func (a *AuthorRepo) Create(req *pb.AuthorCreateReq) (*pb.AuthorRes, error) {
	res := &pb.AuthorRes{}

	query := `INSERT INTO authors (name, biography) VALUES ($1, $2) RETURNING id`

	row := a.db.QueryRow(query, req.Name, req.Biography)

	err := row.Scan(&res.Id)
	if err != nil {
		return nil, fmt.Errorf("Error scaning id: %v", err)
	}

	res, err = a.Get(&pb.GetByIdReq{Id: res.Id})
	if err != nil {
		return nil, fmt.Errorf("can't get author: %w", err)
	}

	return res, err
}

func (a *AuthorRepo) Get(req *pb.GetByIdReq) (*pb.AuthorRes, error) {

	res := &pb.AuthorRes{}

	query := `SELECT id, name, biography FROM authors WHERE id = $1`

	row := a.db.QueryRow(query, req.Id)
	err := row.Scan(
		&res.Id,
		&res.Name,
		&res.Biography,
	)

	if err != nil {
		return nil, fmt.Errorf("can't get author: %w", err)
	}

	return res, nil
}

func (a *AuthorRepo) GetAll(req *pb.AuthorGetAllReq) (*pb.AuthorGetAllRes, error) {
	res := &pb.AuthorGetAllRes{
		Authors: []*pb.AuthorRes{},
	}

	query := `SELECT id, name, biography FROM authors WHERE deleted_at=0`

	var args []interface{}

	if req.Name != "" {
		args = append(args, "%"+req.Name+"%")
		query += fmt.Sprintf(" AND name ILIKE $%d", len(args))
	}

	var defaultLimit int32
	err := a.db.QueryRow("SELECT COUNT(1) FROM authors WHERE deleted_at=0").Scan(&defaultLimit)
	if err != nil {
		return nil, err
	}
	if req.Filter.Limit == 0 {
		req.Filter.Limit = defaultLimit
	}

	args = append(args, req.Filter.Limit, req.Filter.Offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := a.db.Query(query, args...)

	if err != nil {
		return nil, fmt.Errorf("Error while getting authors: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		author := &pb.AuthorRes{}

		err := rows.Scan(
			&author.Id,
			&author.Name,
			&author.Biography,
		)

		if err != nil {
			return nil, fmt.Errorf("Error while scaning authors : %v", err)
		}

		res.Authors = append(res.Authors, author)
	}

	return res, nil
}

func (a *AuthorRepo) Update(req *pb.AuthorUpdateReq) (*pb.AuthorRes, error) {
	res := &pb.AuthorRes{}

	query := `UPDATE authors SET name = $1, biography = $2, updated_at = now() WHERE id = $3 RETURNING id`

	row := a.db.QueryRow(query, req.UpdateAuthor.Name, req.UpdateAuthor.Biography, req.Id.Id)

	err := row.Scan(&res.Id)
	if err != nil {
		return nil, fmt.Errorf("Error scaning id: %v", err)
	}

	res, err = a.Get(&pb.GetByIdReq{Id: res.Id})

	return res, err
}

func (a *AuthorRepo) Delete(req *pb.GetByIdReq) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE authors SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1 RETURNING id`

	_, err := a.db.Exec(query, req.Id)

	if err != nil {
		return nil, fmt.Errorf("can't delete author: %v", err)
	}

	return res, nil
}
