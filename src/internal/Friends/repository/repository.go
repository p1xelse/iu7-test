package repository

import "timetracker/models"

type RepositoryI interface {
	CreateFriendRelation(t *models.FriendRelation) error
	DeleteFriendRelation(friendRel *models.FriendRelation) error
	CheckFriends(t *models.FriendRelation) (bool, error)
	GetUserSubs(userID uint64) ([]uint64, error)
	GetUserFriends(userID uint64) ([]uint64, error)
}
