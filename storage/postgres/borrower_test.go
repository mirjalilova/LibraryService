package postgres

//
//import (
//	"context"
//	"database/sql"
//	"testing"
//
//	g "library_service/genproto"
//
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/assert"
//)
//
//func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
//	db, mock, err := sqlmock.New()
//	assert.NoError(t, err)
//	return db, mock, func() { db.Close() }
//}
//
//func TestCreateBorrower(t *testing.T) {
//	db, mock, cleanup := setupMockDB(t)
//	defer cleanup()
//
//	repo := NewBorrowerRepo(db)
//
//	req := &g.{
//		Borrower: &g.Borrower{
//			UserId:     uuid.NewString(),
//			BookId:     uuid.NewString(),
//			BorrowDate: "2024-01-01",
//			ReturnDate: "2024-02-01",
//		},
//	}
//
//	mock.ExpectExec("INSERT INTO borrowers \\(id, user_id, book_id, borrow_date, return_date\\)").
//		WithArgs(sqlmock.AnyArg(), req.Borrower.UserId, req.Borrower.BookId, req.Borrower.BorrowDate, req.Borrower.ReturnDate).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//
//	res, err := repo.Create(context.Background(), req)
//	assert.NoError(t, err)
//	assert.Equal(t, "Borrower created successfully", res.Message)
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestGetBorrower(t *testing.T) {
//	db, mock, cleanup := setupMockDB(t)
//	defer cleanup()
//
//	repo := NewBorrowerRepo(db)
//	borrowerId := uuid.New().String()
//
//	req := &g.GetBorrowerRequest{Id: borrowerId}
//
//	rows := sqlmock.NewRows([]string{
//		"id", "user_id", "borrow_date", "return_date",
//		"k.id", "k.title", "k.summary",
//		"a.id", "a.name", "a.biography",
//		"g.id", "g.name",
//	}).
//		AddRow(
//			borrowerId, "user-id", "2024-01-01", "2024-02-01",
//			"book-id", "Book Title", "Book Summary",
//			"author-id", "Author Name", "Author Biography",
//			"genre-id", "Genre Name",
//		)
//
//	mock.ExpectQuery("SELECT b.id, b.user_id, b.borrow_date, b.return_date, k.id, k.title, k.summary, a.id, a.name, a.biography, g.id, g.name FROM borrowers b").
//		WithArgs(borrowerId).
//		WillReturnRows(rows)
//
//	res, err := repo.Get(context.Background(), req)
//	assert.NoError(t, err)
//	assert.Equal(t, borrowerId, res.Borrower.Id)
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestDeleteBorrower(t *testing.T) {
//	db, mock, cleanup := setupMockDB(t)
//	defer cleanup()
//
//	repo := NewBorrowerRepo(db)
//	borrowerId := uuid.New().String()
//
//	req := &g.DeleteBorrowerRequest{Id: borrowerId}
//
//	mock.ExpectExec("DELETE FROM borrowers WHERE id = \\$1").
//		WithArgs(borrowerId).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//
//	res, err := repo.Delete(context.Background(), req)
//	assert.NoError(t, err)
//	assert.Equal(t, "borrower deleted successfully", res.Message)
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestGetAllBorrowers(t *testing.T) {
//	db, mock, cleanup := setupMockDB(t)
//	defer cleanup()
//
//	repo := NewBorrowerRepo(db)
//
//	req := &g.GetAllBorrowerRequest{UserId: "user-id"}
//
//	rows := sqlmock.NewRows([]string{
//		"id", "user_id", "borrow_date", "return_date",
//		"k.id", "k.title", "k.summary",
//		"a.id", "a.name", "a.biography",
//		"g.id", "g.name",
//	}).
//		AddRow(
//			"borrower-id", "user-id", "2024-01-01", "2024-02-01",
//			"book-id", "Book Title", "Book Summary",
//			"author-id", "Author Name", "Author Biography",
//			"genre-id", "Genre Name",
//		)
//
//	mock.ExpectQuery("SELECT b.id, b.user_id, b.borrow_date, b.return_date, k.id, k.title, k.summary, a.id, a.name, a.biography, g.id, g.name FROM borrowers b").
//		WillReturnRows(rows)
//
//	res, err := repo.GetAll(context.Background(), req)
//	assert.NoError(t, err)
//	assert.Len(t, res.Borrowers, 1)
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestGetActiveBorrowers(t *testing.T) {
//	db, mock, cleanup := setupMockDB(t)
//	defer cleanup()
//
//	repo := NewBorrowerRepo(db)
//
//	req := &g.GetByUserIdActiveRequest{UserId: "user-id"}
//
//	rows := sqlmock.NewRows([]string{
//		"id", "user_id", "bk.id", "bk.title", "bk.summary",
//		"a.id", "a.name", "a.biography",
//		"g.id", "g.name", "borrow_date", "return_date",
//	}).
//		AddRow(
//			"borrower-id", "user-id", "book-id", "Book Title", "Book Summary",
//			"author-id", "Author Name", "Author Biography",
//			"genre-id", "Genre Name", "2024-01-01", "2024-02-01",
//		)
//
//	mock.ExpectQuery("SELECT b.id, b.user_id, bk.id, bk.title, bk.summary, a.id, a.name, a.biography, g.id, g.name, b.borrow_date, b.return_date FROM borrowers b").
//		WithArgs("user-id").
//		WillReturnRows(rows)
//
//	res, err := repo.GetActiveBorrowers(context.Background(), req)
//	assert.NoError(t, err)
//	assert.Len(t, res.Borrowers, 1)
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestGetInActiveBorrowers(t *testing.T) {
//	db, mock, cleanup := setupMockDB(t)
//	defer cleanup()
//
//	repo := NewBorrowerRepo(db)
//
//	req := &g.GetByUserIdInActiveRequest{UserId: "user-id"}
//
//	rows := sqlmock.NewRows([]string{
//		"id", "user_id", "bk.id", "bk.title", "bk.summary",
//		"a.id", "a.name", "a.biography",
//		"g.id", "g.name", "borrow_date", "return_date",
//	}).
//		AddRow(
//			"borrower-id", "user-id", "book-id", "Book Title", "Book Summary",
//			"author-id", "Author Name", "Author Biography",
//			"genre-id", "Genre Name", "2024-01-01", "2024-02-01",
//		)
//
//	mock.ExpectQuery("SELECT b.id, b.user_id, bk.id, bk.title, bk.summary, a.id, a.name, a.biography, g.id, g.name, b.borrow_date, b.return_date FROM borrowers b").
//		WithArgs("user-id").
//		WillReturnRows(rows)
//
//	res, err := repo.GetInActiveBorrowers(context.Background(), req)
//	assert.NoError(t, err)
//	assert.Len(t, res.Borrowers, 1)
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
