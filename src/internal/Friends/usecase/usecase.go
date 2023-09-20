package usecase

import (
	"timetracker/models"

	friendRep "timetracker/internal/Friends/repository"
	userRep "timetracker/internal/User/repository"

	"github.com/pkg/errors"
)

type UsecaseI interface {
	CreateFriendRelation(friendRel *models.FriendRelation) error
	DeleteFriendRelation(friendRel *models.FriendRelation) error
	CheckIsFriends(userID1 uint64, userID2 uint64) (bool, error)
	GetUserSubs(id uint64) ([]*models.User, error)
	GetUserFriends(id uint64) ([]*models.User, error)
}

type usecase struct {
	friendsRepository friendRep.RepositoryI
	userRepository    userRep.RepositoryI
}

func New(fRep friendRep.RepositoryI, uRep userRep.RepositoryI) UsecaseI {
	return &usecase{
		friendsRepository: fRep,
		userRepository:    uRep,
	}
}

func (uc *usecase) CreateFriendRelation(friendRel *models.FriendRelation) error {
	if friendRel.SubscriberID == friendRel.UserID {
		return models.ErrBadRequest
	}

	friendExists, err := uc.friendsRepository.CheckFriends(friendRel)
	if err != nil {
		return errors.Wrap(err, "friends repository error")
	}
	
	if friendExists {
		return models.ErrConflictFriend
	}

	err = uc.friendsRepository.CreateFriendRelation(friendRel)
	if err != nil {
		return errors.Wrap(err, "friends repository error")
	}

	return err
}

func (uc *usecase) DeleteFriendRelation(friends *models.FriendRelation) error {
	if friends.SubscriberID == friends.UserID {
		return models.ErrBadRequest
	}

	friendExists, err := uc.friendsRepository.CheckFriends(friends)
	if err != nil {
		return errors.Wrap(err, "friends repository error")
	}
	if !friendExists {
		return models.ErrNotFound
	}

	err = uc.friendsRepository.DeleteFriendRelation(friends)
	if err != nil {
		return errors.Wrap(err, "friends repository error")
	}

	return nil
}

func (uc *usecase) GetUserFriends(id uint64) ([]*models.User, error) {
	friendsIDs, err := uc.friendsRepository.GetUserFriends(id)

	if err != nil {
		return nil, errors.Wrap(err, "friends repository error")
	}

	if len(friendsIDs) == 0 {
		return nil, nil
	}

	friends, err := uc.userRepository.GetUsersByIDs(friendsIDs)

	if err != nil {
		return nil, errors.Wrap(err, "friends repository error")
	}

	return friends, nil
}

func (uc *usecase) GetUserSubs(id uint64) ([]*models.User, error) {
	subIDs, err := uc.friendsRepository.GetUserSubs(id)

	if err != nil {
		return nil, errors.Wrap(err, "friends repository error")
	}

	if len(subIDs) == 0 {
		return nil, nil
	}

	subs, err := uc.userRepository.GetUsersByIDs(subIDs)

	if err != nil {
		return nil, errors.Wrap(err, "friends repository error")
	}

	return subs, nil
}

func (uc *usecase) CheckIsFriends(userID1 uint64, userID2 uint64) (bool, error) {
	if userID1 == userID2 {
		return false, models.ErrBadRequest
	}

	friendRel1 := models.FriendRelation{SubscriberID: &userID1, UserID: &userID2}
	friendExists1, err := uc.friendsRepository.CheckFriends(&friendRel1)
	if err != nil {
		return false, errors.Wrap(err, "friends repository error")
	}

	friendRel2 := models.FriendRelation{SubscriberID: &userID2, UserID: &userID1}
	friendExists2, err := uc.friendsRepository.CheckFriends(&friendRel2)
	if err != nil {
		return false, errors.Wrap(err, "friends repository error")
	}

	isBilateralFriendship := friendExists1 && friendExists2

	return isBilateralFriendship, err
}
