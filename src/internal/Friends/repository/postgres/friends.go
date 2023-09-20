package postgres

import (
	"fmt"
	"timetracker/internal/Friends/repository"
	"timetracker/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type FriendRelation struct {
	SubscriberID *uint64 `gorm:"column:subscriber_id"`
	UserID       *uint64 `gorm:"column:user_id"`
}

func (FriendRelation) TableName() string {
	return "friend_relation"
}

type friendRepository struct {
	db *gorm.DB
}

func toPostgresFriendRelation(t *models.FriendRelation) *FriendRelation {
	return &FriendRelation{
		SubscriberID: t.SubscriberID,
		UserID:       t.UserID,
	}
}

func toModelFriendRelation(t *FriendRelation) *models.FriendRelation {
	return &models.FriendRelation{
		SubscriberID: t.SubscriberID,
		UserID:       t.UserID,
	}
}

func toModelFriendRelations(friends []*FriendRelation) []*models.FriendRelation {
	out := make([]*models.FriendRelation, len(friends))

	for i, b := range friends {
		out[i] = toModelFriendRelation(b)
	}

	return out
}

func (fr friendRepository) CreateFriendRelation(t *models.FriendRelation) error {
	postgresFriend := toPostgresFriendRelation(t)

	tx := fr.db.Create(postgresFriend)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table friend_relation)")
	}

	return nil
}

func (fr friendRepository) CheckFriends(t *models.FriendRelation) (bool, error) {
	postgresFriend := toPostgresFriendRelation(t)
	fmt.Println("postgresFriend: ", *postgresFriend.SubscriberID, "   ", *postgresFriend.UserID)
	tx := fr.db.Where(postgresFriend).Take(&FriendRelation{})

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return false, nil
	} else if tx.Error != nil {
		return false, errors.Wrap(tx.Error, "database error (table friend_relation)")
	}

	return true, nil
}

func (fr friendRepository) DeleteFriendRelation(friendRel *models.FriendRelation) error {
	relation := toPostgresFriendRelation(friendRel)
	tx := fr.db.Where(relation).Delete(&FriendRelation{})

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table friend_relation)")
	}

	return nil
}

func (fr friendRepository) GetUserSubs(userID uint64) ([]uint64, error) {
	userIDs := make([]uint64, 0, 10)
	tx := fr.db.Table(FriendRelation{}.TableName()+" f1").
		Select("f1.subscriber_id").
		Joins("left join friend_relation f2 on f2.user_id = f1.subscriber_id and f2.subscriber_id = f1.user_id").
		Where("f1.user_id = ? and f2.user_id is null", userID).Find(&FriendRelation{}).Pluck("subscriber_id", &userIDs)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table friend_relation)")
	}

	return userIDs, nil
}

func (fr friendRepository) GetUserFriends(userID uint64) ([]uint64, error) {
	userIDs := make([]uint64, 0, 10)
	tx := fr.db.Table(FriendRelation{}.TableName()+" f1").
		Select("f1.subscriber_id").
		Joins("join friend_relation f2 on f2.user_id = f1.subscriber_id and f2.subscriber_id = f1.user_id").
		Where("f1.user_id = ?", userID).Find(&FriendRelation{}).Pluck("subscriber_id", &userIDs)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table friend_relation)")
	}

	return userIDs, nil
}

func NewFriendRepository(db *gorm.DB) repository.RepositoryI {
	return &friendRepository{
		db: db,
	}
}
