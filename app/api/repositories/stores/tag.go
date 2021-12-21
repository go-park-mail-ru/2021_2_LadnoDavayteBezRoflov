package stores

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	customErrors "backendServer/pkg/errors"

	"gorm.io/gorm"
)

type TagStore struct {
	db *gorm.DB
}

func CreateTagRepository(db *gorm.DB) repositories.TagRepository {
	return &TagStore{db: db}
}

func (tagStore *TagStore) Create(tag *models.Tag) (err error) {
	tag.ColorID = tag.Color.CLRID
	return tagStore.db.Create(tag).Error
}

func (tagStore *TagStore) Update(tag *models.Tag) (err error) {
	oldTag, err := tagStore.GetByID(tag.TGID)
	if err != nil {
		return
	}

	if tag.Title != "" && tag.Title != oldTag.Title {
		oldTag.Title = tag.Title
	}

	if tag.Color.CLRID != 0 && tag.Color.CLRID != oldTag.Color.CLRID {
		tag.Color = models.AvailableColors[tag.Color.CLRID-1]
		oldTag.Color = tag.Color
		tag.ColorID = oldTag.Color.CLRID
		oldTag.ColorID = tag.ColorID
	}

	return tagStore.db.Save(oldTag).Error
}

func (tagStore *TagStore) Delete(tgid uint) (err error) {
	tag, err := tagStore.GetByID(tgid)
	if err != nil {
		return err
	}

	cards := new([]models.Card)
	err = tagStore.db.Model(&models.Board{BID: tag.BID}).Association("Cards").Find(cards)
	if err != nil {
		return
	}

	for _, card := range *cards {
		if isAssigned, _ := tagStore.IsCardAssigned(tgid, card.CID); isAssigned {
			err = tagStore.AddTagToCard(tgid, card.CID)
			if err != nil {
				return
			}
		}
	}

	return tagStore.db.Delete(&models.Tag{}, tgid).Error
}

func (tagStore *TagStore) GetByID(tgid uint) (*models.Tag, error) {
	tag := new(models.Tag)
	if res := tagStore.db.Find(tag, tgid); res.RowsAffected == 0 {
		return nil, customErrors.ErrCommentNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}
	tag.Color = models.AvailableColors[tag.ColorID-1]
	return tag, nil
}

func (tagStore *TagStore) AddTagToCard(tgid, cid uint) (err error) {
	tag, err := tagStore.GetByID(tgid)
	if err != nil {
		return
	}

	if isAssigned, _ := tagStore.IsCardAssigned(tgid, cid); isAssigned {
		err = tagStore.db.Model(&models.Card{CID: cid}).Association("Tags").Delete(tag)
	} else {
		err = tagStore.db.Model(&models.Card{CID: cid}).Association("Tags").Append(tag)
	}
	return
}

func (tagStore *TagStore) IsCardAssigned(tgid uint, cid uint) (isAssigned bool, err error) {
	result := tagStore.db.Table("tags").
		Joins("LEFT OUTER JOIN tags_cards ON tags_cards.tag_tg_id = tags.tg_id").
		Joins("LEFT OUTER JOIN cards ON tags_cards.card_c_id = cards.c_id").
		Where("tags.tg_id = ? AND cards.c_id = ?", tgid, cid).
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
