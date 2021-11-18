package stores

//
//func createMockDB() (*UserStore, sqlmock.Sqlmock, error) {
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
//	return &UserStore{db: gdb}, mock, nil
//}
//
//func TestCreateUserRepository(t *testing.T) {
//	t.Parallel()
//
//	db, _, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("cant create mockDB: %s", err)
//	}
//	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
//	if err != nil {
//		_ = db.Close()
//		t.Fatalf("cant create mockDB: %s", err)
//	}
//
//	expectedUserRepo := &UserStore{db: gdb}
//
//	require.Equal(t, expectedUserRepo, CreateUserRepository(gdb, "", ""))
//}
//
//func createMockQueryInsertUser(mock *sqlmock.Sqlmock, user *models.User, isSuccessful bool, err error) {
//	(*mock).ExpectBegin()
//	expectedQuery := (*mock).ExpectQuery(`INSERT INTO "users" (.+) RETURNING`).WithArgs(
//		user.Login,
//		user.Email,
//		sqlmock.AnyArg(),
//		sqlmock.AnyArg(),
//		sqlmock.AnyArg(),
//	)
//	if isSuccessful {
//		expectedQuery.WillReturnRows(sqlmock.NewRows([]string{"uid"}).AddRow(1))
//		(*mock).ExpectCommit()
//	} else {
//		expectedQuery.WillReturnError(err)
//		(*mock).ExpectRollback()
//	}
//}
//
//func TestUserRepositoryCreateSuccess(t *testing.T) {
//	t.Parallel()
//
//	repo, mock, err := createMockDB()
//	if err != nil {
//		t.Fatalf("cant create mockDB: %s", err)
//	}
//
//	newUser := &models.User{}
//	if err := faker.FakeData(newUser); err != nil {
//		t.Error(err)
//	}
//
//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "login" FROM "users" WHERE login = $1`)).WithArgs(
//		newUser.Login,
//	).WillReturnRows(sqlmock.NewRows([]string{"login"}))
//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "email" FROM "users" WHERE email = $1`)).WithArgs(
//		newUser.Email,
//	).WillReturnRows(sqlmock.NewRows([]string{"email"}))
//
//	createMockQueryInsertUser(&mock, newUser, true, nil)
//	err = repo.Create(newUser)
//	require.NoError(t, err)
//	require.NotEqual(t, 0, newUser.UID)
//
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
//
//func TestUserRepositoryCreateFail(t *testing.T) {
//	//t.Parallel()
//	//
//	//repo, mock, err := createMockDB()
//	//if err != nil {
//	//	t.Fatalf("cant create mockDB: %s", err)
//	//}
//	//
//	//existingUser := &models.User{}
//	//if err := faker.FakeData(existingUser); err != nil {
//	//	t.Error(err)
//	//}
//	//
//	//createMockQueryInsertUser(&mock, existingUser, true, nil)
//	//
//	//if err := repo.Create(existingUser); err != nil {
//	//	t.Error(err)
//	//}
//	//if err := mock.ExpectationsWereMet(); err != nil {
//	//	t.Errorf("there were unfulfilled expectations: %s", err)
//	//}
//	//
//	//existingUser.UID = 0
//	//mock.ExpectQuery(regexp.QuoteMeta(`SELECT "login" FROM "users" WHERE login = $1`)).WithArgs(
//	//	existingUser.Login,
//	//).WillReturnRows(sqlmock.NewRows([]string{"login"}).AddRow(existingUser.Login))
//	//createMockQueryInsertUser(&mock, existingUser, false, customErrors.ErrUserAlreadyCreated)
//	//
//	//errUserIsExist := repo.Create(existingUser)
//	//require.ErrorIs(t, customErrors.ErrUserAlreadyCreated, errUserIsExist)
//	//if err := mock.ExpectationsWereMet(); err != nil {
//	//	t.Errorf("there were unfulfilled expectations: %s", err)
//	//}
//	//
//	//newUser := &models.User{}
//	//if err := faker.FakeData(newUser); err != nil {
//	//	t.Error(err)
//	//}
//	//newUser.UID = 0
//	//newUser.Email = existingUser.Email
//	//
//	//createMockQueryInsertUser(&mock, newUser, false, customErrors.ErrEmailAlreadyUsed)
//	//errEmailIsUsed := repo.Create(newUser)
//	//require.ErrorIs(t, customErrors.ErrEmailAlreadyUsed, errEmailIsUsed)
//	//if err := mock.ExpectationsWereMet(); err != nil {
//	//	t.Errorf("there were unfulfilled expectations: %s", err)
//	//}
//}
//
//func TestUserRepositoryUpdateSuccess(t *testing.T) {
//}
//
//func TestUserRepositoryUpdateFail(t *testing.T) {
//}
//
//func TestUserRepositoryGetByLoginSuccess(t *testing.T) {
//}
//
//func TestUserRepositoryGetByLoginFail(t *testing.T) {
//}
//
//func TestUserRepositoryGetByIDSuccess(t *testing.T) {
//}
//
//func TestUserRepositoryGetByIDFail(t *testing.T) {
//}
//
//func TestUserRepositoryGetUserTeamsSuccess(t *testing.T) {
//}
//
//func TestUserRepositoryGetUserTeamsFail(t *testing.T) {
//}
//
//func TestUserRepositoryAddUserToTeamSuccess(t *testing.T) {
//}
//
//func TestUserRepositoryAddUserToTeamFail(t *testing.T) {
//}
//
//func TestUserRepositoryIsUserExist(t *testing.T) {
//	t.Parallel()
//
//	repo, mock, err := createMockDB()
//	if err != nil {
//		t.Fatalf("cant create mockDB: %s", err)
//	}
//
//	user := &models.User{}
//	if err := faker.FakeData(user); err != nil {
//		t.Error(err)
//	}
//
//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "login" FROM "users" WHERE login = $1`)).WithArgs(
//		user.Login,
//	).WillReturnRows(sqlmock.NewRows([]string{"login"}))
//	exist, err := repo.IsUserExist(user)
//	require.NoError(t, err)
//	require.False(t, exist)
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//
//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "login" FROM "users" WHERE login = $1`)).WithArgs(
//		user.Login,
//	).WillReturnRows(sqlmock.NewRows([]string{"login"}).AddRow(user.Login))
//
//	exist, err = repo.IsUserExist(user)
//	require.Error(t, err)
//	require.Equal(t, customErrors.ErrUserAlreadyCreated, err)
//	require.True(t, exist)
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
//
//func TestUserRepositoryIsEmailUsed(t *testing.T) {
//	t.Parallel()
//
//	repo, mock, err := createMockDB()
//	if err != nil {
//		t.Fatalf("cant create mockDB: %s", err)
//	}
//
//	user := &models.User{}
//	if err := faker.FakeData(user); err != nil {
//		t.Error(err)
//	}
//
//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "email" FROM "users" WHERE email = $1`)).WithArgs(
//		user.Email,
//	).WillReturnRows(sqlmock.NewRows([]string{"email"}))
//	used, err := repo.IsEmailUsed(user)
//	require.NoError(t, err)
//	require.False(t, used)
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//
//	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "email" FROM "users" WHERE email = $1`)).WithArgs(
//		user.Email,
//	).WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow(user.Email))
//
//	used, err = repo.IsEmailUsed(user)
//	require.Error(t, err)
//	require.Equal(t, customErrors.ErrEmailAlreadyUsed, err)
//	require.True(t, used)
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
