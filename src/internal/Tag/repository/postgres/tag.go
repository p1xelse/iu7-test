package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"timetracker/internal/Tag/repository"
	"timetracker/models"
)

type Tag struct {
	ID     uint64 `gorm:"column:id"`
	UserID uint64 `gorm:"column:user_id"`
	Name   string `gorm:"column:name"`
	About  string `gorm:"column:about"`
	Color  string `gorm:"column:color"`
}

type TagEntryRelation struct {
	TagID   uint64 `gorm:"column:tag_id"`
	EntryID uint64 `gorm:"column:entry_id"`
}

func (TagEntryRelation) TableName() string {
	return "tag_entry"
}

func (Tag) TableName() string {
	return "tag"
}

func toPostgresTag(t *models.Tag) *Tag {
	return &Tag{
		ID:     t.ID,
		UserID: t.UserID,
		Name:   t.Name,
		About:  t.About,
		Color:  t.Color,
	}
}

func toModelTag(t *Tag) *models.Tag {
	return &models.Tag{
		ID:     t.ID,
		UserID: t.UserID,
		Name:   t.Name,
		About:  t.About,
		Color:  t.Color,
	}
}

func toModelTags(tags []*Tag) []*models.Tag {
	out := make([]*models.Tag, len(tags))

	for i, b := range tags {
		out[i] = toModelTag(b)
	}

	return out
}

type tagRepository struct {
	db *gorm.DB
}

func (tr tagRepository) CreateTag(t *models.Tag) error {
	postgresTag := toPostgresTag(t)

	tx := tr.db.Create(postgresTag)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table tag)")
	}

	t.ID = postgresTag.ID
	return nil
}

func (tr tagRepository) UpdateTag(t *models.Tag) error {
	postgresTag := toPostgresTag(t)

	tx := tr.db.Omit("id").Updates(postgresTag)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table tag)")
	}

	return nil
}

func (tr tagRepository) GetTag(id uint64) (*models.Tag, error) {
	var tag Tag

	tx := tr.db.Where("id = ?", id).Take(&tag)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table tag)")
	}

	return toModelTag(&tag), nil
}

func (tr tagRepository) DeleteTag(id uint64) error {
	tx := tr.db.Delete(&Tag{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table tag)")
	}

	return nil
}

func (tr tagRepository) GetUserTags(ustrID uint64) ([]*models.Tag, error) {
	entries := make([]*Tag, 0, 10)

	tx := tr.db.Where(&Tag{UserID: ustrID}).Find(&entries)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table tag)")
	}

	return toModelTags(entries), nil
}

func (tr tagRepository) GetEntryTags(entryID uint64) ([]*models.Tag, error) {
	tagEntryRels := make([]*TagEntryRelation, 0, 10)
	tx := tr.db.Where(&TagEntryRelation{EntryID: entryID}).Find(&tagEntryRels)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table tag)")
	}

	tags := make([]*Tag, 0, 10)

	for idx := range tagEntryRels {
		var tag Tag
		tx := tr.db.Where(&Tag{ID: tagEntryRels[idx].TagID}).Take(&tag)

		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, models.ErrNotFound
		} else if tx.Error != nil {
			return nil, errors.Wrap(tx.Error, "database error (table tag)")
		}

		tags = append(tags, &tag)
	}

	return toModelTags(tags), nil
}

func (tr tagRepository) CreateEntryTags(entryID uint64, tagList []models.Tag) error {
	tagEntryRels := make([]*TagEntryRelation, 0, len(tagList))

	for idx := range tagList {
		tagEntryRels = append(tagEntryRels, &TagEntryRelation{EntryID: entryID, TagID: tagList[idx].ID})
	}

	tx := tr.db.Create(tagEntryRels)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table tag)")
	}

	return nil
}

func (tr tagRepository) UpdateEntryTags(entryID uint64, tagList []models.Tag) error {
	err := tr.DeleteEntryTags(entryID)

	if err != nil {
		return err
	}

	err = tr.CreateEntryTags(entryID, tagList)

	return err
}

func (tr tagRepository) DeleteEntryTags(entryID uint64) error {
	tx := tr.db.Where("entry_id = ?", entryID).Delete(&TagEntryRelation{})

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table tags)")
	}

	return nil
}

func NewTagRepository(db *gorm.DB) repository.RepositoryI {
	return &tagRepository{
		db: db,
	}
}
