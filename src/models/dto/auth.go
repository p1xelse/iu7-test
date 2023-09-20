package dto

import (
	"timetracker/models"
)

type ReqUserSignIn struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ReqUserSignUp struct {
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required"`
	About      string `json:"about"`
	Role       string `json:"role"`
	Password   string `json:"password" validate:"required"`
	AdminToken string `json:"admin_token"`
}

func (req *ReqUserSignIn) ToModelUser() *models.User {
	return &models.User{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (req *ReqUserSignUp) ToModelUser() *models.User {
	if req.Role == "" {
		req.Role = models.DefaultUser.String()
	}
	
	return &models.User{
		Name:     req.Name,
		Email:    req.Email,
		About:    req.About,
		Role:     req.Role,
		Password: req.Password,
	}
}

type RespUser struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	About string `json:"about"`
	Role  string `json:"role"`
}

func GetResponseFromModelUser(user *models.User) *RespUser {
	return &RespUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		About: user.About,
		Role:  user.Role,
	}
}

//
//func GetResponseFromModelEntries(entries []*models.Entry) []*RespEntry {
//	result := make([]*RespEntry, 0, 10)
//	for _, entry := range entries {
//		result = append(result, GetResponseFromModelEntry(entry))
//	}
//
//	return result
//}
