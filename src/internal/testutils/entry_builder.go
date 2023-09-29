package testutils

import (
	"time"
	"timetracker/models"
)

type EntryBuilder struct {
	entry models.Entry
}

func NewEntryBuilder() *EntryBuilder {
	return &EntryBuilder{}
}

func (b *EntryBuilder) WithID(id uint64) *EntryBuilder {
	b.entry.ID = id
	return b
}

func (b *EntryBuilder) WithUserID(userID uint64) *EntryBuilder {
	b.entry.UserID = &userID
	return b
}

func (b *EntryBuilder) WithProjectID(projectID uint64) *EntryBuilder {
	b.entry.ProjectID = &projectID
	return b
}

func (b *EntryBuilder) WithDescription(description string) *EntryBuilder {
	b.entry.Description = description
	return b
}

func (b *EntryBuilder) WithTagList(tagList []models.Tag) *EntryBuilder {
	b.entry.TagList = tagList
	return b
}

func (b *EntryBuilder) WithTimeStart(timeStart time.Time) *EntryBuilder {
	b.entry.TimeStart = timeStart
	return b
}

func (b *EntryBuilder) WithTimeEnd(timeEnd time.Time) *EntryBuilder {
	b.entry.TimeEnd = timeEnd
	return b
}

func (b *EntryBuilder) WithDuration(duration string) *EntryBuilder {
	b.entry.Duration = duration
	return b
}

func (b *EntryBuilder) Build() models.Entry {
	return b.entry
}
