package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/microservices/email/repository"
	"backendServer/app/microservices/email/usecase"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailUseCaseImpl struct {
	emailRepository repository.EmailRepository
	mailUsername    string
}

func CreateEmailUseCase(emailRepository repository.EmailRepository, mailUsername string) usecase.EmailUseCase {
	return &EmailUseCaseImpl{emailRepository: emailRepository, mailUsername: mailUsername}
}

func (emailUseCase *EmailUseCaseImpl) SendFirstLetter(userInfo *models.PublicUserInfo) (emailLetter *gomail.Message) {
	emailLetter = gomail.NewMessage()
	emailLetter.SetHeader("From", emailUseCase.mailUsername)
	emailLetter.SetHeader("To", userInfo.Email)
	emailLetter.SetHeader("Subject", "Добро пожаловать в Brrrello!")
	emailLetter.SetBody("text/plain", "Рады вас видеть у себя на сайте!")
	return
}

func (emailUseCase *EmailUseCaseImpl) SendNotifications() (emailLetters *[]gomail.Message, err error) {
	emailLetters = new([]gomail.Message)

	userDeadlines := make(map[models.PublicUserInfo]map[string][]models.Card)

	// Получаем список всех карточек
	cards, err := emailUseCase.emailRepository.GetAllCards()
	if err != nil {
		return
	}

	// В каждой карточке проверяем дедлайн
	for _, card := range *cards {
		if card.Deadline == "" {
			continue
		}

		parsedDeadline, err := time.Parse("2006-01-02T15:04", card.Deadline)
		if err != nil {
			continue
		}

		timeNow := time.Now()
		timeNow = timeNow.Add(3 * time.Hour)
		tomorrow := timeNow.Add(time.Duration(24) * time.Hour)
		tooLateDeadline := timeNow.Add(-time.Duration(3*24) * time.Hour)

		if parsedDeadline.Before(tomorrow) && parsedDeadline.After(tooLateDeadline) {
			assignees, err := emailUseCase.emailRepository.GetAssignedUsers(card.CID)
			if err != nil {
				continue
			}

			boardTitle, err := emailUseCase.emailRepository.FindBoardTitleByID(card.BID)
			if err != nil {
				continue
			}

			for _, assignee := range *assignees {
				if _, userAlreadyFound := userDeadlines[assignee]; !userAlreadyFound {
					userDeadlines[assignee] = make(map[string][]models.Card)
				}

				if cards, smthAlreadyFound := userDeadlines[assignee][boardTitle]; smthAlreadyFound {
					userDeadlines[assignee][boardTitle] = append(cards, card)
				} else {
					userDeadlines[assignee][boardTitle] = []models.Card{card}
				}
			}
		}
	}

	// Для каждого пользователя в map формируем письмо с дедлайнами
	for user, boards := range userDeadlines {
		emailLetter := gomail.NewMessage()
		emailLetter.SetHeader("From", emailUseCase.mailUsername)
		emailLetter.SetHeader("To", user.Email)
		emailLetter.SetHeader("Subject", "Уведомление о истечении сроков карточек")
		body := "Добрый день!\nУ вас истекают или уже истекли сроки следующих карточек.\n"

		for boardTitle, cards := range boards {
			boardInfo := strings.Join([]string{"Карточки на доске \"", boardTitle, "\":"}, "")
			cardsInfo := boardInfo

			for _, card := range cards {
				parsedDeadline, err := time.Parse("2006-01-02T15:04", card.Deadline)
				if err != nil {
					continue
				}

				timeNow := time.Now()
				timeNow = timeNow.Add(3 * time.Hour)
				tomorrow := timeNow.Add(time.Duration(24) * time.Hour)

				var cardInfo string

				if parsedDeadline.Before(tomorrow) && parsedDeadline.After(timeNow) {
					cardInfo = strings.Join([]string{"карточка \"", card.Title, "\"", " истекает в течение дня;"}, "")
				} else {
					cardInfo = strings.Join([]string{"карточка \"", card.Title, "\"", " истекла менее 3-х дней назад;"}, "")
				}

				cardsInfo = strings.Join([]string{cardsInfo, cardInfo}, "\n\t- ")
			}

			body = strings.Join([]string{body, cardsInfo}, "\n")
		}

		emailLetter.SetBody("text/plain", body)
		*emailLetters = append(*emailLetters, *emailLetter)
	}

	return
}
