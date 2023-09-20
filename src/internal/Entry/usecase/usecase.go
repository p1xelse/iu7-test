package usecase

import (
	"time"
	entryRep "timetracker/internal/Entry/repository"
	tagRep "timetracker/internal/Tag/repository"
	userRep "timetracker/internal/User/repository"
	"timetracker/models"

	"github.com/pkg/errors"
)

type UsecaseI interface {
	CreateEntry(e *models.Entry) error
	UpdateEntry(e *models.Entry) error
	GetEntry(id uint64) (*models.Entry, error)
	DeleteEntry(id uint64, userID uint64) error
	GetUserEntries(userID uint64) ([]*models.Entry, error)
	GetUserEntriesForDay(userID uint64, date time.Time) ([]*models.Entry, error)
}

type usecase struct {
	entryRepository entryRep.RepositoryI
	tagRepository   tagRep.RepositoryI
	userRepository  userRep.RepositoryI
}

func New(eRep entryRep.RepositoryI, tRep tagRep.RepositoryI, uRep userRep.RepositoryI) UsecaseI {
	return &usecase{
		entryRepository: eRep,
		tagRepository:   tRep,
		userRepository:  uRep,
	}
}

func (u *usecase) CreateEntry(e *models.Entry) error {
	err := u.entryRepository.CreateEntry(e)

	if err != nil {
		return errors.Wrap(err, "Error in func entry.Usecase.CreateEntry")
	}

	if e.TagList != nil && len(e.TagList) != 0 {
		err = u.tagRepository.CreateEntryTags(e.ID, e.TagList)

		if err != nil {
			return errors.Wrap(err, "Error in func entry.Usecase.CreateEntry")
		}
	}

	return nil
}

func (u *usecase) UpdateEntry(e *models.Entry) error {
	existedEntry, err := u.entryRepository.GetEntry(e.ID)

	if err != nil {
		return errors.Wrap(err, "Error in func goal.Usecase.Update.UpdateEntry")
	}

	if *existedEntry.UserID != *e.UserID {
		return models.ErrPermissionDenied
	}

	err = u.entryRepository.UpdateEntry(e)

	if err != nil {
		return errors.Wrap(err, "Error in func entry.Usecase.UpdateEntry")
	}

	if e.TagList != nil && len(e.TagList) != 0 {
		err = u.tagRepository.UpdateEntryTags(e.ID, e.TagList)

		if err != nil {
			return errors.Wrap(err, "Error in func entry.Usecase.UpdateEntry")
		}
	}

	if err != nil {
		return errors.Wrap(err, "Error in func entry.Usecase.UpdateEntry")
	}

	return nil
}

func (u *usecase) addAdditionalFieldsToEntry(entry *models.Entry) error {
	err := u.addTagsToEntry(entry)

	if err != nil {
		return errors.Wrap(err, "error while get tags")
	}

	return nil
}

func (u *usecase) addTagsToEntry(entry *models.Entry) error {
	tags, err := u.tagRepository.GetEntryTags(entry.ID)

	if err != nil {
		return errors.Wrap(err, "Error in func addPostAttachmentsAuthors")
	}

	entry.TagList = make([]models.Tag, 0, 10)

	for _, att := range tags {
		entry.TagList = append(entry.TagList, *att)
	}

	return nil
}

func (u *usecase) GetEntry(id uint64) (*models.Entry, error) {
	resEntry, err := u.entryRepository.GetEntry(id)

	if err != nil {
		return nil, errors.Wrap(err, "entry.usecase.GetEntry error while get entry info")
	}

	err = u.addAdditionalFieldsToEntry(resEntry)

	if err != nil {
		return nil, errors.Wrap(err, "entry.usecase.GetEntry error while get additional info")
	}

	return resEntry, nil
}

func (u *usecase) DeleteEntry(id uint64, userId uint64) error {
	existedEntry, err := u.entryRepository.GetEntry(id)
	if err != nil {
		return err
	}

	if existedEntry == nil {
		return models.ErrNotFound
	}

	if *existedEntry.UserID != userId {
		return models.ErrPermissionDenied
	}

	err = u.entryRepository.DeleteEntry(id) // TODO потом откатывать транзакцию если теги удалились, придумать как

	if err != nil {
		return errors.Wrap(err, "entry.repository delete error")
	}

	err = u.tagRepository.DeleteEntryTags(id)

	if err != nil {
		return errors.Wrap(err, "entry.Usecase.tagRepository delete error")
	}

	return nil
}

func (u *usecase) GetUserEntries(userID uint64) ([]*models.Entry, error) {
	entries, err := u.entryRepository.GetUserEntries(userID)

	if err != nil {
		return nil, errors.Wrap(err, "Error in func entry.Usecase.GetUserPosts")
	}

	for idx := range entries {
		err = u.addAdditionalFieldsToEntry(entries[idx])

		if err != nil {
			return nil, errors.Wrap(err, "entry.Usecase.GetUserPosts error while add additional fields")
		}
	}

	return entries, nil
}

func (u *usecase) GetUserEntriesForDay(userID uint64, date time.Time) ([]*models.Entry, error) {
	entries, err := u.entryRepository.GetUserEntriesForDay(userID, date)

	if err != nil {
		return nil, errors.Wrap(err, "Error in func entry.Usecase.GetUserPosts")
	}

	for idx := range entries {
		err = u.addAdditionalFieldsToEntry(entries[idx])

		if err != nil {
			return nil, errors.Wrap(err, "entry.Usecase.GetUserPosts error while add additional fields")
		}
	}

	return entries, nil
}
