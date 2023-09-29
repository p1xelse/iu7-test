package testutils

import "timetracker/models"

type TagBuilder struct {
	tag models.Tag
}

func NewTagBuilder() *TagBuilder {
	return &TagBuilder{}
}

func (b *TagBuilder) WithID(id uint64) *TagBuilder {
	b.tag.ID = id
	return b
}

func (b *TagBuilder) WithUserID(userID uint64) *TagBuilder {
	b.tag.UserID = userID
	return b
}

func (b *TagBuilder) WithName(name string) *TagBuilder {
	b.tag.Name = name
	return b
}

func (b *TagBuilder) WithAbout(about string) *TagBuilder {
	b.tag.About = about
	return b
}

func (b *TagBuilder) WithColor(color string) *TagBuilder {
	b.tag.Color = color
	return b
}

func (b *TagBuilder) Build() models.Tag {
	return b.tag
}
