package stores

//
//func createCheckListMockDB() (*CheckListStore, sqlmock.Sqlmock, error) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		return nil, nil, err
//	}
//	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
//	if err != nil {
//		_ = db.Close()
//		return nil, nil, err
//	}
//
//	return &CheckListStore{db: gdb}, mock, nil
//}
//
//func TestCreateCheckList(t *testing.T) {
//	t.Parallel()
//
//	repo, mock, err := createCheckListMockDB()
//	if err != nil {
//		t.Fatalf("cant create mockDB: %s", err)
//	}
//
//	checkList := new(models.CheckList)
//	if err := faker.FakeData(checkList); err != nil {
//		t.Error(err)
//	}
//	checkList.CheckListItems = []models.CheckListItem{}
//
//	// success
//	mock.ExpectBegin()
//	mock.ExpectQuery(`INSERT INTO "check_lists" (.+) RETURNING`).WithArgs(
//		checkList.CID,
//		checkList.Title,
//		checkList.CHLID,
//	).WillReturnRows(sqlmock.NewRows([]string{"1"}))
//	mock.ExpectCommit()
//
//	err = repo.Create(checkList)
//	assert.NoError(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.NoError(t, err)
//
//	// error
//	mock.ExpectBegin()
//	mock.ExpectQuery(`INSERT INTO "check_lists" (.+) RETURNING`).WithArgs(
//		checkList.CID,
//		checkList.Title,
//		checkList.CHLID,
//	).WillReturnError(customErrors.ErrInternal)
//	mock.ExpectRollback()
//
//	err = repo.Create(checkList)
//	assert.Error(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.NoError(t, err)
//}
//
//func TestDeleteCheckList(t *testing.T) {
//	t.Parallel()
//
//	repo, mock, err := createCheckListMockDB()
//	if err != nil {
//		t.Fatalf("cant create mockDB: %s", err)
//	}
//
//	testCHLID := uint(1)
//
//	// success
//	mock.ExpectBegin()
//	mock.ExpectExec(`DELETE FROM "check_lists" WHERE "check_lists"."chl_id"`).WithArgs(
//		testCHLID,
//	).WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectCommit()
//
//	err = repo.Delete(testCHLID)
//	assert.NoError(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.NoError(t, err)
//
//	// error
//	mock.ExpectBegin()
//	mock.ExpectExec(`DELETE FROM "check_lists" WHERE "check_lists"."chl_id"`).WithArgs(
//		testCHLID,
//	).WillReturnError(customErrors.ErrInternal)
//	mock.ExpectRollback()
//
//	err = repo.Delete(testCHLID)
//	assert.Error(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.NoError(t, err)
//}
