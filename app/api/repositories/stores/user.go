package stores

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/hasher"
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aofei/cameron"

	"github.com/streadway/amqp"

	_ "golang.org/x/image/bmp"

	"github.com/google/uuid"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"

	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserStore struct {
	db                *gorm.DB
	avatarPath        string
	defaultAvatarName string
	channel           *amqp.Channel
	queueName         string
}

func CreateUserRepository(db *gorm.DB, avatarPath, defaultAvatarName string, channel *amqp.Channel, queueName string) repositories.UserRepository {
	return &UserStore{db: db, avatarPath: avatarPath, defaultAvatarName: defaultAvatarName, channel: channel, queueName: queueName}
}

func (userStore *UserStore) Create(user *models.User) (err error) {
	isUserExist, err := userStore.IsUserExist(user)
	if isUserExist {
		return
	}

	isEmailUsed, err := userStore.IsEmailUsed(user)
	if isEmailUsed {
		return
	}

	user.HashedPassword, err = hasher.HashPassword(user.Password)
	if err != nil {
		return
	}

	isCustomAvatarCreated := false

	fileNameID := uuid.NewString()
	fileName := strings.Join([]string{userStore.avatarPath, "/", fileNameID, ".webp"}, "")

	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)

	if err == nil {
		err = webp.Encode(out, cameron.Identicon([]byte(fileNameID), 540, 60), options)
	}

	if err == nil {
		fileName = strings.Replace(fileName, "/backend", "", -1)
		user.Avatar = fileName
		isCustomAvatarCreated = true
	}

	if !isCustomAvatarCreated {
		user.Avatar = strings.Join([]string{userStore.avatarPath, "/", userStore.defaultAvatarName}, "")
		user.Avatar = strings.Replace(user.Avatar, "/backend", "", -1)
	}

	err = userStore.db.Create(user).Error
	if err != nil {
		return
	}

	publicData, err := userStore.GetPublicData(user.UID)
	if err != nil {
		return
	}

	body, err := json.Marshal(publicData)
	if err != nil {
		return
	}

	err = userStore.channel.Publish("", userStore.queueName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Transient,
		ContentType:  "text/plain",
		Body:         body,
	})
	return
}

func (userStore *UserStore) Update(user *models.User) (err error) {
	oldUser, err := userStore.GetByID(user.UID)
	if err != nil {
		return
	}

	if user.Login != "" && user.Login != oldUser.Login {
		var isNewLoginExist bool
		isNewLoginExist, err = userStore.IsUserExist(user)
		if isNewLoginExist {
			return
		}
		oldUser.Login = user.Login
	}

	if user.Email != "" && user.Email != oldUser.Email {
		var isNewEmailUsed bool
		isNewEmailUsed, err = userStore.IsEmailUsed(user)
		if isNewEmailUsed {
			return
		}
		oldUser.Email = user.Email
	}

	if user.Password != "" && !hasher.IsPasswordsEqual(user.Password, oldUser.HashedPassword) {
		oldUser.HashedPassword, err = hasher.HashPassword(user.Password)
		if err != nil {
			return
		}
	}

	if user.Description != "" && user.Description != oldUser.Description {
		oldUser.Description = user.Description
	}

	return userStore.db.Save(oldUser).Error
}

func (userStore *UserStore) UpdateAvatar(user *models.User, avatar *multipart.FileHeader) (err error) {
	oldUser, err := userStore.GetByID(user.UID)
	if err != nil {
		return
	}

	if user.Avatar != "" {
		fileNameID := uuid.NewString()
		fileName := strings.Join([]string{userStore.avatarPath, "/", fileNameID, ".webp"}, "")

		in, err := avatar.Open()
		if err != nil {
			return err
		}
		defer func(in multipart.File) {
			_ = in.Close()
		}(in)

		img, _, err := image.Decode(in)
		if err != nil {
			return err
		}

		out, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer func(out *os.File) {
			_ = out.Close()
		}(out)

		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
		if err != nil {
			return err
		}

		err = webp.Encode(out, img, options)
		if err != nil {
			return err
		}

		defaultAvatar := strings.Join([]string{userStore.avatarPath, "/", userStore.defaultAvatarName}, "")
		defaultAvatar = strings.Replace(defaultAvatar, "/backend", "", -1)
		if oldUser.Avatar != "" && oldUser.Avatar != defaultAvatar {
			err = os.Remove(oldUser.Avatar)
			if err != nil {
				oldUser.Avatar = strings.Replace(oldUser.Avatar, "/backend", "", -1)
				return err
			}
		}

		fileName = strings.Replace(fileName, "/backend", "", -1)
		user.Avatar = fileName
		oldUser.Avatar = fileName
	}

	return userStore.db.Save(oldUser).Error
}

func (userStore *UserStore) GetByLogin(login string) (*models.User, error) {
	user := new(models.User)
	if res := userStore.db.Where("login = ?", login).Find(user); res.RowsAffected == 0 {
		return nil, customErrors.ErrUserNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (userStore *UserStore) GetByID(uid uint) (*models.User, error) {
	user := new(models.User)
	if res := userStore.db.Find(user, uid); res.RowsAffected == 0 {
		return nil, customErrors.ErrUserNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (userStore *UserStore) FindAllByLogin(text string, amount int) (users *[]models.PublicUserInfo, err error) {
	users = new([]models.PublicUserInfo)
	text = strings.Join([]string{"%", text, "%"}, "")
	err = userStore.db.Model(&models.User{}).Where("LOWER(login) LIKE ?", strings.ToLower(text)).Limit(amount).Find(users).Error
	return
}

func (userStore *UserStore) FindBoardMembersByLogin(bid uint, text string, amount int) (users *[]models.PublicUserInfo, err error) {
	users = new([]models.PublicUserInfo)
	text = strings.Join([]string{"%", text, "%"}, "")

	err = userStore.db.Raw("? UNION ?",
		userStore.db.Table("users").
			Joins("LEFT OUTER JOIN users_teams ON users_teams.user_uid = users.uid").
			Joins("LEFT OUTER JOIN teams ON users_teams.team_t_id = teams.t_id").
			Joins("JOIN boards ON teams.t_id = boards.t_id").
			Where("boards.b_id = ? AND LOWER(users.login) LIKE ?", bid, strings.ToLower(text)).
			Select("users.uid, users.login, users.avatar"),
		userStore.db.Table("users").
			Joins("LEFT OUTER JOIN users_boards ON users_boards.user_uid = users.uid").
			Joins("LEFT OUTER JOIN boards ON users_boards.board_b_id = boards.b_id").
			Where("boards.b_id = ? AND LOWER(users.login) LIKE ?", bid, strings.ToLower(text)).
			Select("users.uid, users.login, users.avatar"),
	).Limit(amount).Find(users).Error

	return
}

func (userStore *UserStore) FindBoardInvitedMembersByLogin(bid uint, text string, amount int) (users *[]models.PublicUserInfo, err error) {
	var uids []uint
	userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_teams ON users_teams.user_uid = users.uid").
		Joins("LEFT OUTER JOIN teams ON users_teams.team_t_id = teams.t_id").
		Joins("JOIN boards ON teams.t_id = boards.t_id").
		Where("boards.b_id = ?", bid).Pluck("users.uid", &uids)

	users = new([]models.PublicUserInfo)
	text = strings.Join([]string{"%", text, "%"}, "")

	err = userStore.db.Model(&models.User{}).Not(uids).Where("LOWER(login) LIKE ?", strings.ToLower(text)).Limit(amount).Find(users).Error
	//
	//err = userStore.db.Table("users").
	//	Joins("LEFT OUTER JOIN users_boards ON users_boards.user_uid = users.uid").
	//	Joins("LEFT OUTER JOIN boards ON users_boards.board_b_id = boards.b_id").
	//	Where("boards.b_id = ? AND LOWER(users.login) LIKE ?", bid, strings.ToLower(text)).
	//	Select("users.uid, users.login, users.avatar").Limit(amount).Find(users).Error

	return
}

func (userStore *UserStore) GetUserTeams(uid uint) (teams *[]models.Team, err error) {
	teams = new([]models.Team)
	err = userStore.db.Model(&models.User{UID: uid}).Association("Teams").Find(teams)
	return
}

func (userStore *UserStore) GetUserToggledBoards(uid uint) (boards *[]models.Board, err error) {
	boards = new([]models.Board)
	err = userStore.db.Model(&models.User{UID: uid}).Association("Boards").Find(boards)
	return
}

func (userStore *UserStore) AddUserToTeam(uid, tid uint) (err error) {
	user, err := userStore.GetByID(uid)
	if err != nil {
		return
	}
	if isMember, _ := userStore.IsUserInTeam(uid, tid); isMember {

		boards := new([]models.Board)
		err = userStore.db.Model(&models.Team{TID: tid}).Association("Boards").Find(boards)
		if err != nil {
			return
		}

		for _, board := range *boards {
			cards := new([]models.Card)
			err = userStore.db.Model(&models.Board{BID: board.BID}).Association("Cards").Find(cards)
			if err != nil {
				return
			}

			for _, card := range *cards {
				if isAssigned, _ := userStore.IsCardAssigned(uid, card.CID); isAssigned {
					err = userStore.AddUserToCard(uid, card.CID)
					if err != nil {
						return
					}
				}
			}
		}

		err = userStore.db.Model(&models.Team{TID: tid}).Association("Users").Delete(user)
	} else {
		err = userStore.db.Model(&models.Team{TID: tid}).Association("Users").Append(user)
	}
	return
}

func (userStore *UserStore) AddUserToBoard(uid, bid uint) (err error) {
	user, err := userStore.GetByID(uid)
	if err != nil {
		return
	}
	if isAccessed, _ := userStore.IsBoardAccessed(uid, bid); isAccessed {
		cards := new([]models.Card)
		err = userStore.db.Model(&models.Board{BID: bid}).Association("Cards").Find(cards)
		if err != nil {
			return
		}

		for _, card := range *cards {
			if isAssigned, _ := userStore.IsCardAssigned(uid, card.CID); isAssigned {
				err = userStore.AddUserToCard(uid, card.CID)
				if err != nil {
					return
				}
			}
		}

		err = userStore.db.Model(&models.Board{BID: bid}).Association("Users").Delete(user)
	} else {
		err = userStore.db.Model(&models.Board{BID: bid}).Association("Users").Append(user)
	}
	return
}

func (userStore *UserStore) AddUserToCard(uid, cid uint) (err error) {
	user, err := userStore.GetByID(uid)
	if err != nil {
		return
	}

	if isAccessed, err := userStore.IsCardAccessed(uid, cid); !isAccessed {
		return err
	}

	if isAssigned, _ := userStore.IsCardAssigned(uid, cid); isAssigned {
		err = userStore.db.Model(&models.Card{CID: cid}).Association("Users").Delete(user)
	} else {
		err = userStore.db.Model(&models.Card{CID: cid}).Association("Users").Append(user)
	}
	return
}

func (userStore *UserStore) GetPublicData(uid uint) (user *models.PublicUserInfo, err error) {
	user = new(models.PublicUserInfo)
	err = userStore.db.Model(&models.User{}).Find(user, uid).Error
	return
}

func (userStore *UserStore) IsUserExist(user *models.User) (bool, error) {
	if res := userStore.db.Select("login").Where("login = ?", user.Login).Find(user); res.RowsAffected == 0 {
		return false, nil
	} else if res.Error != nil {
		return true, res.Error
	}
	return true, customErrors.ErrUserAlreadyCreated
}

func (userStore *UserStore) IsEmailUsed(user *models.User) (bool, error) {
	if res := userStore.db.Select("email").Where("email = ?", user.Email).Find(user); res.RowsAffected == 0 {
		return false, nil
	} else if res.Error != nil {
		return true, res.Error
	}
	return true, customErrors.ErrEmailAlreadyUsed
}

func (userStore *UserStore) IsUserInTeam(uid uint, tid uint) (isMember bool, err error) {
	user := new(models.User)
	if err = userStore.db.Model(&models.Team{TID: tid}).Association("Users").Find(user, uid); err != nil {
		return false, err
	} else if user.UID == 0 {
		return false, nil
	}
	return true, nil
}

func (userStore *UserStore) IsBoardAccessed(uid uint, bid uint) (isAccessed bool, err error) {
	result := userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_teams ON users_teams.user_uid = users.uid").
		Joins("LEFT OUTER JOIN teams ON users_teams.team_t_id = teams.t_id").
		Joins("JOIN boards ON teams.t_id = boards.t_id").
		Where("users.uid = ? AND boards.b_id = ?", uid, bid).
		Select("teams.t_id").Find(&models.Team{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
		return
	}

	result = userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_boards ON users_boards.user_uid = users.uid").
		Joins("LEFT OUTER JOIN boards ON users_boards.board_b_id = boards.b_id").
		Where("users.uid = ? AND boards.b_id = ?", uid, bid).
		Select("boards.b_id").Find(&models.Board{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
	} else {
		err = customErrors.ErrNoAccess
	}
	return
}

func (userStore *UserStore) IsCardListAccessed(uid uint, clid uint) (isAccessed bool, err error) {
	result := userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_teams ON users_teams.user_uid = users.uid").
		Joins("LEFT OUTER JOIN teams ON users_teams.team_t_id = teams.t_id").
		Joins("JOIN boards ON teams.t_id = boards.t_id").
		Joins("JOIN card_lists ON card_lists.b_id = boards.b_id").
		Where("users.uid = ? AND card_lists.cl_id = ?", uid, clid).
		Select("teams.t_id").Find(&models.Team{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
		return
	}

	result = userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_boards ON users_boards.user_uid = users.uid").
		Joins("LEFT OUTER JOIN boards ON users_boards.board_b_id = boards.b_id").
		Joins("JOIN card_lists ON card_lists.b_id = boards.b_id").
		Where("users.uid = ? AND card_lists.cl_id = ?", uid, clid).
		Select("boards.b_id").Find(&models.Board{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
	} else {
		err = customErrors.ErrNoAccess
	}
	return
}

func (userStore *UserStore) IsCardAccessed(uid uint, cid uint) (isAccessed bool, err error) {
	result := userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_teams ON users_teams.user_uid = users.uid").
		Joins("LEFT OUTER JOIN teams ON users_teams.team_t_id = teams.t_id").
		Joins("JOIN boards ON teams.t_id = boards.t_id").
		Joins("JOIN cards ON cards.b_id = boards.b_id").
		Where("users.uid = ? AND cards.c_id = ?", uid, cid).
		Select("teams.t_id").Find(&models.Team{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
		return
	}

	result = userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_boards ON users_boards.user_uid = users.uid").
		Joins("LEFT OUTER JOIN boards ON users_boards.board_b_id = boards.b_id").
		Joins("JOIN cards ON cards.b_id = boards.b_id").
		Where("users.uid = ? AND cards.c_id = ?", uid, cid).
		Select("boards.b_id").Find(&models.Board{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
	} else {
		err = customErrors.ErrNoAccess
	}
	return
}

func (userStore *UserStore) IsCardAssigned(uid uint, cid uint) (isAssigned bool, err error) {
	result := userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_cards ON users_cards.user_uid = users.uid").
		Joins("LEFT OUTER JOIN cards ON users_cards.card_c_id = cards.c_id").
		Where("users.uid = ? AND cards.c_id = ?", uid, cid).
		Select("cards.c_id").Find(&models.Card{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAssigned = true
	} else {
		err = customErrors.ErrNoAccess
	}
	return
}

func (userStore *UserStore) IsCommentAccessed(uid uint, cmid uint) (isAccessed bool, err error) {
	result := userStore.db.Where("cm_id = ? AND uid = ?", cmid, uid).Find(&models.Comment{})

	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
	} else {
		err = customErrors.ErrNoAccess
	}
	return
}

func (userStore *UserStore) IsCheckListAccessed(uid uint, chlid uint) (isAccessed bool, err error) {
	result := userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_teams ON users_teams.user_uid = users.uid").
		Joins("LEFT OUTER JOIN teams ON users_teams.team_t_id = teams.t_id").
		Joins("JOIN boards ON teams.t_id = boards.t_id").
		Joins("JOIN cards ON cards.b_id = boards.b_id").
		Joins("JOIN check_lists ON check_lists.c_id = cards.c_id").
		Where("users.uid = ? AND check_lists.chl_id = ?", uid, chlid).
		Select("teams.t_id").Find(&models.Team{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
		return
	}

	result = userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_boards ON users_boards.user_uid = users.uid").
		Joins("LEFT OUTER JOIN boards ON users_boards.board_b_id = boards.b_id").
		Joins("JOIN cards ON cards.b_id = boards.b_id").
		Joins("JOIN check_lists ON check_lists.c_id = cards.c_id").
		Where("users.uid = ? AND check_lists.chl_id = ?", uid, chlid).
		Select("boards.b_id").Find(&models.Board{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
	} else {
		err = customErrors.ErrNoAccess
	}
	return
}

func (userStore *UserStore) IsCheckListItemAccessed(uid uint, chliid uint) (isAccessed bool, err error) {
	result := userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_teams ON users_teams.user_uid = users.uid").
		Joins("LEFT OUTER JOIN teams ON users_teams.team_t_id = teams.t_id").
		Joins("JOIN boards ON teams.t_id = boards.t_id").
		Joins("JOIN cards ON cards.b_id = boards.b_id").
		Joins("JOIN check_lists ON check_lists.c_id = cards.c_id").
		Joins("JOIN check_list_items ON check_list_items.chl_id = check_lists.chl_id").
		Where("users.uid = ? AND check_list_items.chli_id = ?", uid, chliid).
		Select("teams.t_id").Find(&models.Team{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
		return
	}

	result = userStore.db.Table("users").
		Joins("LEFT OUTER JOIN users_boards ON users_boards.user_uid = users.uid").
		Joins("LEFT OUTER JOIN boards ON users_boards.board_b_id = boards.b_id").
		Joins("JOIN cards ON cards.b_id = boards.b_id").
		Joins("JOIN check_lists ON check_lists.c_id = cards.c_id").
		Joins("JOIN check_list_items ON check_list_items.chl_id = check_lists.chl_id").
		Where("users.uid = ? AND check_list_items.chli_id = ?", uid, chliid).
		Select("boards.b_id").Find(&models.Board{})
	err = result.Error
	if err != nil {
		return
	}

	if result.RowsAffected > 0 {
		isAccessed = true
	} else {
		err = customErrors.ErrNoAccess
	}
	return
}
