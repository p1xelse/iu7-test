package dto

import (
	"timetracker/models"
)

type ReqCreateUpdateProject struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name" validate:"required"`
	About     string `json:"about"`
	Color     string `json:"color"`
	IsPrivate bool   `json:"is_private"`
}

func (req *ReqCreateUpdateProject) ToModelProject() *models.Project {
	return &models.Project{
		ID:        req.ID,
		Name:      req.Name,
		About:     req.About,
		Color:     req.Color,
		IsPrivate: req.IsPrivate,
	}
}

type RespProject struct {
	ID              uint64  `json:"id"`
	UserID          *uint64 `json:"user_id"`
	Name            string  `json:"name"`
	About           string  `json:"about"`
	Color           string  `json:"color"`
	IsPrivate       bool    `json:"is_private"`
	TotalCountHours float64 `json:"total_count_hours"`
}

func GetResponseFromModelProject(project *models.Project) *RespProject {
	return &RespProject{
		ID:              project.ID,
		UserID:          project.UserID,
		Name:            project.Name,
		About:           project.About,
		Color:           project.Color,
		IsPrivate:       project.IsPrivate,
		TotalCountHours: project.TotalCountHours,
	}
}

func GetResponseFromModelProjects(entries []*models.Project) []*RespProject {
	result := make([]*RespProject, 0, 10)
	for _, project := range entries {
		result = append(result, GetResponseFromModelProject(project))
	}

	return result
}
