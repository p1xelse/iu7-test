package testutils

import "timetracker/models"

type FriendRelationBuilder struct {
	fRel models.FriendRelation
}

func NewFriendRelationBuilder() *FriendRelationBuilder {
	return &FriendRelationBuilder{}
}

func (b *FriendRelationBuilder) WithSubID(subID uint64) *FriendRelationBuilder {
	b.fRel.SubscriberID = &subID
	return b
}

func (b *FriendRelationBuilder) WithUserID(userID uint64) *FriendRelationBuilder {
	b.fRel.UserID = &userID
	return b
}

func (b *FriendRelationBuilder) Build() models.FriendRelation {
	return b.fRel
}
