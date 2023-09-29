package testutils

import "timetracker/models"

type ProjectBuilder struct {
	project models.Project
}

func NewProjectBuilder() *ProjectBuilder {
	return &ProjectBuilder{}
}

func (b *ProjectBuilder) WithID(id uint64) *ProjectBuilder {
	b.project.ID = id
	return b
}

func (b *ProjectBuilder) WithUserID(userID uint64) *ProjectBuilder {
	b.project.UserID = &userID
	return b
}

func (b *ProjectBuilder) WithName(name string) *ProjectBuilder {
	b.project.Name = name
	return b
}

func (b *ProjectBuilder) WithAbout(about string) *ProjectBuilder {
	b.project.About = about
	return b
}

func (b *ProjectBuilder) WithColor(color string) *ProjectBuilder {
	b.project.Color = color
	return b
}

func (b *ProjectBuilder) WithIsPrivate(isPrivate bool) *ProjectBuilder {
	b.project.IsPrivate = isPrivate
	return b
}

func (b *ProjectBuilder) WithTotalCountHours(totalCountHours float64) *ProjectBuilder {
	b.project.TotalCountHours = totalCountHours
	return b
}

func (b *ProjectBuilder) Build() models.Project {
	return b.project
}
