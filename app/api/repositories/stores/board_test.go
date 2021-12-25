package stores

//func createBoardMockDB() (*BoardStore, sqlmock.Sqlmock, error) {
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
//	return &BoardStore{db: gdb}, mock, nil
//}

//func TestCreateBoard(t *testing.T) {
//	t.Parallel()
//
//	repo, mock, err := createBoardMockDB()
//	if err != nil {
//		t.Fatalf("cant create mockDB: %s", err)
//	}
//
//	board := new(models.Board)
//	if err := faker.FakeData(board); err != nil {
//		t.Error(err)
//	}
//	board.Users = []models.User{}
//	board.CardLists = []models.CardList{}
//	board.Cards = []models.Card{}
//	board.Tags = []models.Tag{}
//
//	// success
//	mock.ExpectBegin()
//	mock.ExpectQuery(`INSERT INTO "boards" (.+) RETURNING`).WithArgs(
//		board.TID,
//		board.Title,
//		board.Description,
//		board.AccessPath,
//		board.BID,
//	).WillReturnRows(sqlmock.NewRows([]string{"1"}))
//	mock.ExpectCommit()
//
//	err = repo.Create(board)
//	assert.NoError(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.NoError(t, err)
//
//	// error
//	mock.ExpectBegin()
//	mock.ExpectQuery(`INSERT INTO "boards" (.+) RETURNING`).WithArgs(
//		board.TID,
//		board.Title,
//		board.Description,
//		board.AccessPath,
//		board.BID,
//	).WillReturnError(customErrors.ErrInternal)
//	mock.ExpectRollback()
//
//	err = repo.Create(board)
//	assert.Error(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.NoError(t, err)
//}
//
//func TestDeleteBoard(t *testing.T) {
//	t.Parallel()
//
//	repo, mock, err := createBoardMockDB()
//	if err != nil {
//		t.Fatalf("cant create mockDB: %s", err)
//	}
//
//	testBID := uint(1)
//
//	// success
//	mock.ExpectBegin()
//	mock.ExpectExec(`DELETE FROM "boards" WHERE "boards"."b_id"`).WithArgs(
//		testBID,
//	).WillReturnResult(sqlmock.NewResult(1, 1))
//	mock.ExpectCommit()
//
//	err = repo.Delete(testBID)
//	assert.NoError(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.NoError(t, err)
//
//	// error
//	mock.ExpectBegin()
//	mock.ExpectExec(`DELETE FROM "boards" WHERE "boards"."b_id"`).WithArgs(
//		testBID,
//	).WillReturnError(customErrors.ErrInternal)
//	mock.ExpectRollback()
//
//	err = repo.Delete(testBID)
//	assert.Error(t, err)
//
//	err = mock.ExpectationsWereMet()
//	assert.NoError(t, err)
//}
