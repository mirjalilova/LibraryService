package postgres

import (
	"database/sql"
	"fmt"
	pb "library_service/genproto/library_service"

	_ "github.com/lib/pq"
)

type GenreRepo struct {
	db *sql.DB
}

func NewGenreRepo(db *sql.DB) *GenreRepo {
	return &GenreRepo{
		db: db,
	}
}

func (g *GenreRepo) Create(req *pb.GenreCreateReq) (*pb.GenreRes, error) {
	res := &pb.GenreRes{}

	query := `INSERT INTO genres (name) VALUES ($1) RETURNING id`

	row := g.db.QueryRow(query, req.Name)

	var id string

	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("Error scaning id: %v", err)
	}

	res, err = g.Get(&pb.GetByIdReq{Id: id})
	if err != nil {
		return nil, fmt.Errorf("can't get genre: %w", err)
	}

	return res, nil
}

func (g *GenreRepo) Get(req *pb.GetByIdReq) (*pb.GenreRes, error) {
	res := &pb.GenreRes{}

	query := `SELECT id, name FROM genres WHERE id = $1 AND deleted_at=0`

	row := g.db.QueryRow(query, req.Id)

	err := row.Scan(&res.Id, &res.Name)
	if err != nil {
		return nil, fmt.Errorf("can't get genre: %w", err)
	}

	return res, nil
}

func (g *GenreRepo) GetAll(req *pb.GenreGetAllReq) (*pb.GenreGetAllRes, error) {
	res := &pb.GenreGetAllRes{
		Genres: []*pb.GenreRes{},
	}

	query := `SELECT id, name FROM genres WHERE deleted_at=0`

	var args []interface{}

	if req.Name != "" {
		args = append(args, "%"+req.Name+"%")
		query += fmt.Sprintf(" AND name ILIKE $%d", len(args))
	}

	var defaultLimit int32
	err := g.db.QueryRow("SELECT COUNT(1) FROM genres WHERE deleted_at=0").Scan(&defaultLimit)
	if err != nil {
		return nil, err
	}
	if req.Filter.Limit == 0 {
		req.Filter.Limit = defaultLimit
	}

	args = append(args, req.Filter.Limit, req.Filter.Offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := g.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Error while getting genres: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		genre := &pb.GenreRes{}

		err := rows.Scan(&genre.Id, &genre.Name)
		if err != nil {
			return nil, fmt.Errorf("can't get genre: %w", err)
		}

		res.Genres = append(res.Genres, genre)
	}

	return res, nil
}

func (g *GenreRepo) Update(req *pb.GenreUpdateReq) (*pb.GenreRes, error) {
	res := &pb.GenreRes{}

	query := `UPDATE genres SET name=$1, updated_at=now() WHERE id=$2 RETURNING id`

	row := g.db.QueryRow(query, req.UpdateGenre.Name, req.Id.Id)

	var id string

	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("Error scaning id: %v", err)
	}

	res, err = g.Get(&pb.GetByIdReq{Id: id})
	if err != nil {
		return nil, fmt.Errorf("can't get genre: %w", err)
	}

	return res, nil
}

func (g *GenreRepo) Delete(req *pb.GetByIdReq) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE genres SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1 RETURNING id`

	_, err := g.db.Exec(query, req.Id)

	if err != nil {
		return nil, fmt.Errorf("can't delete genre: %v", err)
	}

	return res, nil
}
