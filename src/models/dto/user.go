package dto

import "timetracker/models"

type ReqUpdateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	About    string `json:"about"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

func (req *ReqUpdateUser) ToModelUser() *models.User {
	return &models.User{
		Name:     req.Name,
		Email:    req.Email,
		About:    req.About,
		Role:     req.Role,
		Password: req.Password,
	}
}

func GetResponseFromModelUsers(users []*models.User) []*RespUser {
	result := make([]*RespUser, 0, 10)
	for _, user := range users {
		result = append(result, GetResponseFromModelUser(user))
	}

	return result
}
