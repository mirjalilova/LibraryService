package postgres

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	pb "library_service/genproto/library_service"
)

var testDB *sql.DB

func setup(t *testing.T) *GenreRepo {
	if testDB == nil {
		connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
			"postgres",
			"feruza1727",
			"localhost",
			"5432",
			"mock") // Replace with your actual database details
		db, err := sql.Open("postgres", connString)
		if err != nil {
			t.Fatalf("failed to open database connection: %v", err)
		}
		testDB = db
	}

	repo := NewGenreRepo(testDB)
	return repo
}

func teardown() {
	if testDB != nil {
		testDB.Close()
		testDB = nil
	}
}

func TestCreateGenre(t *testing.T) {
	repo := setup(t)
	defer teardown()

	req := &pb.GenreCreateReq{
		Name: "Fantasy",
	}
	res, err := repo.Create(req)
	if err != nil {
		t.Fatalf("failed to create genre: %v", err)
	}

	assert.NotEmpty(t, res.Id, "genre ID should not be empty")
	assert.Equal(t, req.Name, res.Name)
}

func TestGetGenre(t *testing.T) {
	repo := setup(t)
	defer teardown()

	// First, create a genre to retrieve later
	createReq := &pb.GenreCreateReq{
		Name: "Sci-Fi",
	}
	createRes, err := repo.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create genre: %v", err)
	}

	req := &pb.GetByIdReq{Id: createRes.Id}
	res, err := repo.Get(req)
	if err != nil {
		t.Fatalf("failed to get genre: %v", err)
	}

	assert.Equal(t, createRes.Id, res.Id)
	assert.Equal(t, createReq.Name, res.Name)
}

func TestGetAllGenres(t *testing.T) {
	repo := setup(t)
	defer teardown()

	// Create some genres
	genres := []pb.GenreCreateReq{
		{Name: "Adventure"},
		{Name: "Mystery"},
	}
	for _, genre := range genres {
		_, err := repo.Create(&genre)
		if err != nil {
			t.Fatalf("failed to create genre: %v", err)
		}
	}

	req := &pb.GenreGetAllReq{
		Filter: &pb.Filter{
			Limit:  10,
			Offset: 0,
		},
	}
	res, err := repo.GetAll(req)
	if err != nil {
		t.Fatalf("failed to get genres: %v", err)
	}

	assert.GreaterOrEqual(t, len(res.Genres), 2, "expected at least two genres in the result")
}

func TestUpdateGenre(t *testing.T) {
	repo := setup(t)
	defer teardown()

	// First, create a genre to update later
	createReq := &pb.GenreCreateReq{
		Name: "Romance",
	}
	createRes, err := repo.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create genre: %v", err)
	}

	updateReq := &pb.GenreUpdateReq{
		Id: &pb.GetByIdReq{Id: createRes.Id},
		UpdateGenre: &pb.GenreCreateReq{
			Name: "Updated Romance",
		},
	}
	res, err := repo.Update(updateReq)
	if err != nil {
		t.Fatalf("failed to update genre: %v", err)
	}

	assert.Equal(t, updateReq.UpdateGenre.Name, res.Name)
}

func TestDeleteGenre(t *testing.T) {
	repo := setup(t)
	defer teardown()

	// First, create a genre to delete later
	createReq := &pb.GenreCreateReq{
		Name: "To Be Deleted",
	}
	createRes, err := repo.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create genre: %v", err)
	}

	req := &pb.GetByIdReq{Id: createRes.Id}
	_, err = repo.Delete(req)
	if err != nil {
		t.Fatalf("failed to delete genre: %v", err)
	}

	// Verify that the genre was deleted successfully
	getRes, err := repo.Get(req)
	assert.Error(t, err, "expected error when trying to fetch deleted genre")
	assert.Nil(t, getRes, "expected nil result when fetching deleted genre")
}
