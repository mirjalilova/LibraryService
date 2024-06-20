package postgres

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	pb "library_service/genproto/library_service"
)

var (
	authorId = "b3929343-7b1f-4572-be55-8e6cf1f7f9e1"
	genreId  = "61678c2e-5526-4940-afed-44ce65bd4a04"
	userId   = "89843524-4aa4-4ade-936c-f0a7975c48bc"
	bookId   = ""
)

// Helper function to create a new test AuthorRepo with a real database connection
func NewTestAuthorRepo(t *testing.T) *AuthorRepo {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		"postgres",
		"feruza1727",
		"localhost",
		"5432",
		"mock")

	db, err := sql.Open("postgres", connString)
	if err != nil {
		t.Fatal("failed to open connection", err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal("failed to ping database", err)
	}
	return NewAuthorRepo(db)
}

func TestCreateAuthor(t *testing.T) {
	repo := NewTestAuthorRepo(t)
	defer repo.db.Close()

	req := &pb.AuthorCreateReq{
		Name:      "John Doe",
		Biography: "Some biography",
	}
	res, err := repo.Create(req)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the author was created successfully
	assert.NotEmpty(t, res.Id, "author ID should not be empty")
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, req.Biography, res.Biography)
}

func TestGetAuthor(t *testing.T) {
	repo := NewTestAuthorRepo(t)
	defer repo.db.Close()

	// First, create an author to retrieve later
	createReq := &pb.AuthorCreateReq{
		Name:      "Jane Doe",
		Biography: "Another biography",
	}
	createRes, err := repo.Create(createReq)
	if err != nil {
		t.Fatal(err)
	}

	req := &pb.GetByIdReq{Id: createRes.Id}
	res, err := repo.Get(req)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the retrieved author matches the created one
	assert.Equal(t, createRes.Id, res.Id)
	assert.Equal(t, createRes.Name, res.Name)
	assert.Equal(t, createRes.Biography, res.Biography)
}

func TestGetAllAuthors(t *testing.T) {
	repo := NewTestAuthorRepo(t)
	defer repo.db.Close()

	// Create some authors
	authors := []pb.AuthorCreateReq{
		{Name: "Author One", Biography: "Bio One"},
		{Name: "Author Two", Biography: "Bio Two"},
	}
	for _, author := range authors {
		_, err := repo.Create(&author)
		if err != nil {
			t.Fatal(err)
		}
	}

	req := &pb.AuthorGetAllReq{
		Filter: &pb.Filter{
			Limit:  10,
			Offset: 0,
		},
	}
	res, err := repo.GetAll(req)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that we have at least two authors in the result
	assert.GreaterOrEqual(t, len(res.Authors), 2)
}

func TestUpdateAuthor(t *testing.T) {
	repo := NewTestAuthorRepo(t)
	defer repo.db.Close()

	// First, create an author to update later
	createReq := &pb.AuthorCreateReq{
		Name:      "Initial Name",
		Biography: "Initial biography",
	}
	createRes, err := repo.Create(createReq)
	if err != nil {
		t.Fatal(err)
	}

	updateReq := &pb.AuthorUpdateReq{
		Id: &pb.GetByIdReq{Id: createRes.Id},
		UpdateAuthor: &pb.AuthorCreateReq{
			Name:      "Updated Name",
			Biography: "Updated biography",
		},
	}
	res, err := repo.Update(updateReq)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the author was updated successfully
	assert.Equal(t, updateReq.UpdateAuthor.Name, res.Name)
	assert.Equal(t, updateReq.UpdateAuthor.Biography, res.Biography)
}

func TestDeleteAuthor(t *testing.T) {
	repo := NewTestAuthorRepo(t)
	defer repo.db.Close()

	// First, create an author to delete later
	createReq := &pb.AuthorCreateReq{
		Name:      "Delete Me",
		Biography: "To be deleted",
	}
	createRes, err := repo.Create(createReq)
	if err != nil {
		t.Fatal(err)
	}

	req := &pb.GetByIdReq{Id: createRes.Id}
	_, err = repo.Delete(req)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the author was deleted successfully
	getRes, err := repo.Get(req)
	assert.NoError(t, err)
	assert.NotNil(t, getRes)
}
