package integrationutils

import (
	"encoding/json"
	"timetracker/models/dto"
)

type ReqCreateUpdateProjectBuilder struct {
	project dto.ReqCreateUpdateProject
}

func NewReqCreateUpdateProjectBuilder() *ReqCreateUpdateProjectBuilder {
	return &ReqCreateUpdateProjectBuilder{}
}

func (b *ReqCreateUpdateProjectBuilder) WithID(id uint64) *ReqCreateUpdateProjectBuilder {
	b.project.ID = id
	return b
}

func (b *ReqCreateUpdateProjectBuilder) WithName(name string) *ReqCreateUpdateProjectBuilder {
	b.project.Name = name
	return b
}

func (b *ReqCreateUpdateProjectBuilder) WithAbout(about string) *ReqCreateUpdateProjectBuilder {
	b.project.About = about
	return b
}

func (b *ReqCreateUpdateProjectBuilder) WithColor(color string) *ReqCreateUpdateProjectBuilder {
	b.project.Color = color
	return b
}

func (b *ReqCreateUpdateProjectBuilder) WithIsPrivate(isPrivate bool) *ReqCreateUpdateProjectBuilder {
	b.project.IsPrivate = isPrivate
	return b
}

func (b *ReqCreateUpdateProjectBuilder) Build() dto.ReqCreateUpdateProject {
	return b.project
}

func (b *ReqCreateUpdateProjectBuilder) Json() ([]byte, error) {
	return json.Marshal(b.project)
}
