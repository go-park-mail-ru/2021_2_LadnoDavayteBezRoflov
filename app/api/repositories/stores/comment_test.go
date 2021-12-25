package stores

import (
	"backendServer/app/api/models"
	customErrors "backendServer/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createCommentMockDB() (*CommentStore, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	return &CommentStore{db: gdb}, mock, nil
}

func TestCreateComment(t *testing.T) {
	t.Parallel()

	repo, mock, err := createCommentMockDB()
	if err != nil {
		t.Fatalf("cant create mockDB: %s", err)
	}

	comment := new(models.Comment)
	if err := faker.FakeData(comment); err != nil {
		t.Error(err)
	}

	// success
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "comments" (.+) RETURNING`).WithArgs(
		comment.CID,
		comment.UID,
		comment.Text,
		comment.Date,
		comment.CMID,
	).WillReturnRows(sqlmock.NewRows([]string{"1"}))
	mock.ExpectCommit()

	err = repo.Create(comment)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "comments" (.+) RETURNING`).WithArgs(
		comment.CID,
		comment.UID,
		comment.Text,
		comment.Date,
		comment.CMID,
	).WillReturnError(customErrors.ErrInternal)
	mock.ExpectRollback()

	err = repo.Create(comment)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

//
//func TestUpdateComment(t *testing.T) {
//	t.Parallel()
//
//	repo, mock, err := createCommentMockDB()
//	if err != nil {
//		t.Fatalf("cant create mockDB: %s", err)
//	}
//
//	comment := new(models.Comment)
//	if err := faker.FakeData(comment); err != nil {
//		t.Error(err)
//	}
//	testNewText := "some new text"
//
//	// success
//	mock.ExpectBegin()
//	mock.ExpectQuery(`SELECT * FROM "comments" WHERE cm_id = $1`).
//		WithArgs(comment.CMID).WillReturnRows(sqlmock.NewRows([]string{"cm_id", "c_id", "uid", "text", "date"}).
//		AddRow(comment.CMID, comment.CID, comment.UID, comment.Text, comment.Date))
//	mock.ExpectExec(`UPDATE "comments" text = $1 WHERE cm_id = $2`).WithArgs(testNewText, comment.CMID).WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectCommit()
//
//	err = repo.Update(comment)
//	assert.NoError(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.NoError(t, err)
//
//	//// error
//	//mock.ExpectBegin()
//	//mock.ExpectQuery(`INSERT INTO "comments" (.+) RETURNING`).WithArgs(
//	//	comment.CID,
//	//	comment.UID,
//	//	comment.Text,
//	//	comment.Date,
//	//	comment.CMID,
//	//).WillReturnError(customErrors.ErrInternal)
//	//mock.ExpectRollback()
//	//
//	//err = repo.Create(comment)
//	//assert.Error(t, err)
//	//
//	//err = mock.ExpectationsWereMet()
//	//assert.NoError(t, err)
//}

func TestDeleteComment(t *testing.T) {
	t.Parallel()

	repo, mock, err := createCommentMockDB()
	if err != nil {
		t.Fatalf("cant create mockDB: %s", err)
	}

	testCMID := uint(1)

	// success
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "comments" WHERE "comments"."cm_id"`).WithArgs(
		testCMID,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Delete(testCMID)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// error
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "comments" WHERE "comments"."cm_id"`).WithArgs(
		testCMID,
	).WillReturnError(customErrors.ErrInternal)
	mock.ExpectRollback()

	err = repo.Delete(testCMID)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
