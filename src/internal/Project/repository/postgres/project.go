package postgres

import (
	"timetracker/internal/Project/repository"
	"timetracker/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Project struct {
	ID              uint64  `gorm:"column:id"`
	UserID          *uint64 `gorm:"column:user_id"`
	Name            string  `gorm:"column:name"`
	About           string  `gorm:"column:about"`
	Color           string  `gorm:"column:color"`
	IsPrivate       bool    `gorm:"column:is_private"`
	TotalCountHours float64 `gorm:"column:total_count_hours"`
}

func (Project) TableName() string {
	return "project"
}

func toPostgresProject(p *models.Project) *Project {
	return &Project{
		ID:              p.ID,
		UserID:          p.UserID,
		Name:            p.Name,
		About:           p.About,
		Color:           p.Color,
		IsPrivate:       p.IsPrivate,
		TotalCountHours: p.TotalCountHours,
	}
}

func toModelProject(p *Project) *models.Project {
	return &models.Project{
		ID:              p.ID,
		UserID:          p.UserID,
		Name:            p.Name,
		About:           p.About,
		Color:           p.Color,
		IsPrivate:       p.IsPrivate,
		TotalCountHours: p.TotalCountHours,
	}
}

func toModelProjects(projects []*Project) []*models.Project {
	out := make([]*models.Project, len(projects))

	for i, b := range projects {
		out[i] = toModelProject(b)
	}

	return out
}

type projectRepository struct {
	db *gorm.DB
}

func (pr projectRepository) CreateProject(e *models.Project) error {
	postgresProject := toPostgresProject(e)
	tx := pr.db.Create(postgresProject)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table project)")
	}

	e.ID = postgresProject.ID
	return nil
}

func (pr projectRepository) UpdateProject(e *models.Project) error {
	postgresProject := toPostgresProject(e)

	tx := pr.db.Omit("id").Updates(postgresProject)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table project)")
	}

	return nil
}

func (pr projectRepository) GetProject(id uint64) (*models.Project, error) {
	var project Project

	tx := pr.db.Where("id = ?", id).Take(&project)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table project)")
	}

	return toModelProject(&project), nil
}

func (pr projectRepository) DeleteProject(id uint64) error {
	tx := pr.db.Delete(&Project{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table project)")
	}

	return nil
}

func (pr projectRepository) GetUserProjects(userID uint64) ([]*models.Project, error) {
	projects := make([]*Project, 0, 10)

	tx := pr.db.Where(&Project{UserID: &userID}).Find(&projects)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table project)")
	}

	return toModelProjects(projects), nil
}

func NewProjectRepository(db *gorm.DB) repository.RepositoryI {
	return &projectRepository{
		db: db,
	}
}
