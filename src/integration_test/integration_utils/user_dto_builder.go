package integrationutils

import (
	"encoding/json"
	"timetracker/models/dto"
)

type ReqUpdateUserBuilder struct {
	user dto.ReqUpdateUser
}

// Name     string `json:"name"`
// 	Email    string `json:"email"`
// 	About    string `json:"about"`
// 	Role     string `json:"role"`
// 	Password string `json:"password"`

func NewReqUpdateUserBuilder() *ReqUpdateUserBuilder {
	return &ReqUpdateUserBuilder{}
}

func (b *ReqUpdateUserBuilder) WithName(name string) *ReqUpdateUserBuilder {
	b.user.Name = name
	return b
}

func (b *ReqUpdateUserBuilder) WithAbout(about string) *ReqUpdateUserBuilder {
	b.user.About = about
	return b
}

func (b *ReqUpdateUserBuilder) WithRole(role string) *ReqUpdateUserBuilder {
	b.user.Role = role
	return b
}

func (b *ReqUpdateUserBuilder) WithPassword(password string) *ReqUpdateUserBuilder {
	b.user.Password = password
	return b
}

func (b *ReqUpdateUserBuilder) Build() dto.ReqUpdateUser {
	return b.user
}

func (b *ReqUpdateUserBuilder) Json() ([]byte, error) {
	return json.Marshal(b.user)
}
