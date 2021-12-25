package stores

import (
	"backendServer/app/api/models"
	customErrors "backendServer/pkg/errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createCheckListItemMockDB() (*CheckListItemStore, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	return &CheckListItemStore{db: gdb}, mock, nil
}

func TestCreateCheckListItem(t *testing.T) {
	t.Parallel()

	repo, mock, err := createCheckListItemMockDB()
	if err != nil {
		t.Fatalf("cant create mockDB: %s", err)
	}

	checkListItem := new(models.CheckListItem)
	if err := faker.FakeData(checkListItem); err != nil {
		t.Error(err)
	}

	// success
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "check_list_items" (.+) RETURNING`).WithArgs(
		checkListItem.CHLID,
		checkListItem.Text,
		checkListItem.Status,
		checkListItem.CHLIID,
	).WillReturnRows(sqlmock.NewRows([]string{"1"}))
	mock.ExpectCommit()

	err = repo.Create(checkListItem)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "check_list_items" (.+) RETURNING`).WithArgs(
		checkListItem.CHLID,
		checkListItem.Text,
		checkListItem.Status,
		checkListItem.CHLIID,
	).WillReturnError(customErrors.ErrInternal)
	mock.ExpectRollback()

	err = repo.Create(checkListItem)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteCheckListItem(t *testing.T) {
	t.Parallel()

	repo, mock, err := createCheckListItemMockDB()
	if err != nil {
		t.Fatalf("cant create mockDB: %s", err)
	}

	testCHLIID := uint(1)

	// success
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "check_list_items" WHERE "check_list_items"."chli_id"`).WithArgs(
		testCHLIID,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.Delete(testCHLIID)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)

	// error
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "check_list_items" WHERE "check_list_items"."chli_id"`).WithArgs(
		testCHLIID,
	).WillReturnError(customErrors.ErrInternal)
	mock.ExpectRollback()

	err = repo.Delete(testCHLIID)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
