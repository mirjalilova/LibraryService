package postgres

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	pb "library_service/genproto/library_service"
)

func setupBookRepo(t *testing.T) *BookRepo {
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

		// Initialize any necessary tables or data
		err = initializeDatabase(db)
		if err != nil {
			t.Fatalf("failed to initialize database: %v", err)
		}
	}

	repo := NewBookRepo(testDB)
	return repo
}

func teardownBookRepo() {
	if testDB != nil {
		testDB.Close()
		testDB = nil
	}
}

func initializeDatabase(db *sql.DB) error {
	// Initialize your database schema or insert necessary data here
	// Example:
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id UUID PRIMARY KEY,
			title TEXT NOT NULL,
			author_id UUID NOT NULL,
			genre_id UUID NOT NULL,
			summary TEXT
		);

		-- Additional table creation or data insertion as needed
	`)
	if err != nil {
		return err
	}

	return nil
}

func TestBookRepo(t *testing.T) {
	t.Run("TestCreateBook", TestCreateBook)
	t.Run("TestGetBook", TestGetBook)
	t.Run("TestGetAllBooks", TestGetAllBooks)
	t.Run("TestUpdateBook", TestUpdateBook)
	t.Run("TestDeleteBook", TestDeleteBook)
	t.Run("TestSearchBooks", TestSearchBooks)
}

func TestCreateBook(t *testing.T) {
	repo := setupBookRepo(t)
	defer teardownBookRepo()

	// Create a book
	req := &pb.BookCreateReq{
		Title:    "Sample Book",
		AuthorId: "b3929343-7b1f-4572-be55-8e6cf1f7f9e1", // Replace with valid author ID
		GenreId:  "3fa85f64-5717-4562-b3fc-2c963f66afa6", // Replace with valid genre ID
		Summary:  "A sample book summary",
	}
	res, err := repo.Create(req)
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	assert.NotEmpty(t, res.Id, "book ID should not be empty")
	assert.Equal(t, req.Title, res.Title)
}

func TestGetBook(t *testing.T) {
	repo := setupBookRepo(t)
	defer teardownBookRepo()

	// First, create a book to retrieve later
	createReq := &pb.BookCreateReq{
		Title:    "Sample Book",
		AuthorId: "b3929343-7b1f-4572-be55-8e6cf1f7f9e1", // Replace with valid author ID
		GenreId:  "3fa85f64-5717-4562-b3fc-2c963f66afa6", // Replace with valid genre ID
		Summary:  "A sample book summary",
	}
	createRes, err := repo.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	req := &pb.GetByIdReq{Id: createRes.Id}
	res, err := repo.Get(req)
	if err != nil {
		t.Fatalf("failed to get book: %v", err)
	}

	assert.Equal(t, createRes.Id, res.Id)
	assert.Equal(t, createReq.Title, res.Title)
}

func TestGetAllBooks(t *testing.T) {
	repo := setupBookRepo(t)
	defer teardownBookRepo()

	// Create some books
	books := []pb.BookCreateReq{
		{Title: "Book One", AuthorId: "b3929343-7b1f-4572-be55-8e6cf1f7f9e1", GenreId: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Summary: "Summary One"},
		{Title: "Book Two", AuthorId: "b3929343-7b1f-4572-be55-8e6cf1f7f9e1", GenreId: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Summary: "Summary Two"},
	}
	for _, book := range books {
		_, err := repo.Create(&book)
		if err != nil {
			t.Fatalf("failed to create book: %v", err)
		}
	}

	req := &pb.BookGetAllReq{
		Filter: &pb.Filter{
			Limit:  10,
			Offset: 0,
		},
	}
	res, err := repo.GetAll(req)
	if err != nil {
		t.Fatalf("failed to get books: %v", err)
	}

	assert.GreaterOrEqual(t, len(res.Books), 2, "expected at least two books in the result")
}

func TestUpdateBook(t *testing.T) {
	repo := setupBookRepo(t)
	defer teardownBookRepo()

	// First, create a book to update later
	createReq := &pb.BookCreateReq{
		Title:    "Initial Book",
		AuthorId: "b3929343-7b1f-4572-be55-8e6cf1f7f9e1", // Replace with valid author ID
		GenreId:  "3fa85f64-5717-4562-b3fc-2c963f66afa6", // Replace with valid genre ID
		Summary:  "Initial summary",
	}
	createRes, err := repo.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	updateReq := &pb.BookUpdateReq{
		Id: &pb.GetByIdReq{Id: createRes.Id},
		UpdateBook: &pb.BookCreateReq{
			Title:    "Updated Book",
			AuthorId: "b3929343-7b1f-4572-be55-8e6cf1f7f9e1", // Replace with valid author ID
			GenreId:  "3fa85f64-5717-4562-b3fc-2c963f66afa6", // Replace with valid genre ID
			Summary:  "Updated summary",
		},
	}
	res, err := repo.Update(updateReq)
	if err != nil {
		t.Fatalf("failed to update book: %v", err)
	}

	assert.Equal(t, updateReq.UpdateBook.Title, res.Title)
}

func TestDeleteBook(t *testing.T) {
	repo := setupBookRepo(t)
	defer teardownBookRepo()

	// First, create a book to delete later
	createReq := &pb.BookCreateReq{
		Title:    "Delete Me",
		AuthorId: "b3929343-7b1f-4572-be55-8e6cf1f7f9e1", // Replace with valid author ID
		GenreId:  "3fa85f64-5717-4562-b3fc-2c963f66afa6", // Replace with valid genre ID
		Summary:  "To be deleted",
	}
	createRes, err := repo.Create(createReq)
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	req := &pb.GetByIdReq{Id: createRes.Id}
	_, err = repo.Delete(req)
	if err != nil {
		t.Fatalf("failed to delete book: %v", err)
	}

	// Verify that the book was deleted successfully
	getRes, err := repo.Get(req)
	assert.Error(t, err, "expected error when trying to fetch deleted book")
	assert.Nil(t, getRes, "expected nil result when fetching deleted book")
}

func TestSearchBooks(t *testing.T) {
	repo := setupBookRepo(t)
	defer teardownBookRepo()

	// Create some books
	books := []pb.BookCreateReq{
		{Title: "Sample Book", AuthorId: "b3929343-7b1f-4572-be55-8e6cf1f7f9e1", GenreId: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Summary: "Summary One"},
		{Title: "Another Book", AuthorId: "b3929343-7b1f-4572-be55-8e6cf1f7f9e1", GenreId: "3fa85f64-5717-4562-b3fc-2c963f66afa6", Summary: "Summary Two"},
	}
	for _, book := range books {
		_, err := repo.Create(&book)
		if err != nil {
			t.Fatalf("failed to create book: %v", err)
		}
	}

	req := &pb.BookSearchReq{Title: "Sample"}
	res, err := repo.Search(req)
	if err != nil {
		t.Fatalf("failed to search books: %v", err)
	}

	assert.GreaterOrEqual(t, len(res.Books), 1, "expected at least one book in the search result")
}
