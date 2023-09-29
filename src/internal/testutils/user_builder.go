package testutils

import "timetracker/models"

type UserBuilder struct {
	user models.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (b *UserBuilder) WithID(id uint64) *UserBuilder {
	b.user.ID = id
	return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.user.Name = name
	return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.user.Email = email
	return b
}

func (b *UserBuilder) WithAbout(about string) *UserBuilder {
	b.user.About = about
	return b
}

func (b *UserBuilder) WithRole(role string) *UserBuilder {
	b.user.Role = role
	return b
}

func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	b.user.Password = password
	return b
}

func (b *UserBuilder) Build() models.User {
	return b.user
}
