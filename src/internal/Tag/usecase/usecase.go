package usecase

import (
	"github.com/pkg/errors"
	tagRep "timetracker/internal/Tag/repository"
	"timetracker/models"
)

type UsecaseI interface {
	CreateTag(t *models.Tag) error
	UpdateTag(t *models.Tag) error
	GetTag(id uint64) (*models.Tag, error)
	DeleteTag(id uint64, userID uint64) error
	GetUserTags(userID uint64) ([]*models.Tag, error)
}

type usecase struct {
	tagRepository tagRep.RepositoryI
}

func (u *usecase) CreateTag(t *models.Tag) error {
	err := u.tagRepository.CreateTag(t)

	if err != nil {
		return errors.Wrap(err, "Error in func Tag.Usecase.CreateTag")
	}

	return nil
}

func (u *usecase) UpdateTag(t *models.Tag) error {
	_, err := u.tagRepository.GetTag(t.ID)

	if err != nil {
		return errors.Wrap(err, "Error in func Tag.Usecase.CreateTag")
	}

	err = u.tagRepository.UpdateTag(t)

	if err != nil {
		return errors.Wrap(err, "Error in func Tag.Usecase.CreateTag")
	}

	return nil
}

func (u *usecase) GetTag(id uint64) (*models.Tag, error) {
	resTag, err := u.tagRepository.GetTag(id)

	if err != nil {
		return nil, errors.Wrap(err, "Tag.usecase.GetTag error while get Tag info")
	}

	return resTag, nil
}

func (u *usecase) DeleteTag(id uint64, userID uint64) error {
	existedTag, err := u.tagRepository.GetTag(id)
	if err != nil {
		return err
	}

	if existedTag == nil {
		return errors.New("Tag not found") //TODO models error
	}

	if existedTag.UserID != userID {
		return errors.New("Permission denied")
	}

	err = u.tagRepository.DeleteTag(id)

	if err != nil {
		return errors.Wrap(err, "Tag.repository delete error")
	}

	return nil
}

func (u *usecase) GetUserTags(userID uint64) ([]*models.Tag, error) {
	entries, err := u.tagRepository.GetUserTags(userID)

	if err != nil {
		return nil, errors.Wrap(err, "Error in func Tag.Usecase.GetUserPosts")
	}

	return entries, nil
}

func New(tRep tagRep.RepositoryI) UsecaseI {
	return &usecase{
		tagRepository: tRep,
	}
}
