package testutils

import "timetracker/models"

type TagMother struct{}

func (tm *TagMother) CreateDefaultTag() models.Tag {
	return models.Tag{
		ID:     1,
		UserID: 1,
		Name:   "Default models.Tag",
		About:  "This is a default models.Tag",
		Color:  "blue",
	}
}

func (tm *TagMother) CreateCustomTag(id, userID uint64, name, about, color string) models.Tag {
	return models.Tag{
		ID:     id,
		UserID: userID,
		Name:   name,
		About:  about,
		Color:  color,
	}
}
