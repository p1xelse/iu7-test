package integrationutils

import (
	"encoding/json"
	"time"
	"timetracker/models/dto"
)

type ReqCreateUpdateEntryBuilder struct {
	entry dto.ReqCreateUpdateEntry
}

func NewReqCreateUpdateEntryBuilder() *ReqCreateUpdateEntryBuilder {
	return &ReqCreateUpdateEntryBuilder{}
}

func (b *ReqCreateUpdateEntryBuilder) WithID(id uint64) *ReqCreateUpdateEntryBuilder {
	b.entry.ID = id
	return b
}

func (b *ReqCreateUpdateEntryBuilder) WithProjectID(projectID uint64) *ReqCreateUpdateEntryBuilder {
	b.entry.ProjectID = &projectID
	return b
}

func (b *ReqCreateUpdateEntryBuilder) WithDescription(description string) *ReqCreateUpdateEntryBuilder {
	b.entry.Description = description
	return b
}

func (b *ReqCreateUpdateEntryBuilder) WithTagList(tagList []uint64) *ReqCreateUpdateEntryBuilder {
	b.entry.TagList = tagList
	return b
}

func (b *ReqCreateUpdateEntryBuilder) WithTimeStart(timeStart time.Time) *ReqCreateUpdateEntryBuilder {
	b.entry.TimeStart = timeStart
	return b
}

func (b *ReqCreateUpdateEntryBuilder) WithTimeEnd(timeEnd time.Time) *ReqCreateUpdateEntryBuilder {
	b.entry.TimeEnd = timeEnd
	return b
}

func (b *ReqCreateUpdateEntryBuilder) Build() dto.ReqCreateUpdateEntry {
	return b.entry
}

func (b *ReqCreateUpdateEntryBuilder) Json() ([]byte, error) {
	return json.Marshal(b.entry)
}
