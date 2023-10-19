package integrationutils

import (
	"encoding/json"
	"timetracker/models/dto"
)

type ReqCreateUpdateTagBuilder struct {
	tag dto.ReqCreateUpdateTag
}

func NewReqCreateUpdateTagBuilder() *ReqCreateUpdateTagBuilder {
	return &ReqCreateUpdateTagBuilder{}
}

func (b *ReqCreateUpdateTagBuilder) WithID(id uint64) *ReqCreateUpdateTagBuilder {
	b.tag.ID = id
	return b
}

func (b *ReqCreateUpdateTagBuilder) WithName(name string) *ReqCreateUpdateTagBuilder {
	b.tag.Name = name
	return b
}

func (b *ReqCreateUpdateTagBuilder) WithAbout(about string) *ReqCreateUpdateTagBuilder {
	b.tag.About = about
	return b
}

func (b *ReqCreateUpdateTagBuilder) WithColor(color string) *ReqCreateUpdateTagBuilder {
	b.tag.Color = color
	return b
}

func (b *ReqCreateUpdateTagBuilder) Build() dto.ReqCreateUpdateTag {
	return b.tag
}

func (b *ReqCreateUpdateTagBuilder) Json() ([]byte, error) {
	return json.Marshal(b.tag)
}
