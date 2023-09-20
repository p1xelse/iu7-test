package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
	"timetracker/internal/Goal/repository"
	"timetracker/models"
)

type Goal struct {
	ID          uint64    `gorm:"column:id"`
	UserID      *uint64   `gorm:"column:user_id"`
	Name        string    `gorm:"column:name"`
	ProjectID   *uint64   `gorm:"column:project_id"`
	Description string    `gorm:"column:description"`
	TimeStart   time.Time `gorm:"column:time_start"`
	TimeEnd     time.Time `gorm:"column:time_end"`
	HoursCount  float64   `gorm:"column:hours_count"`
}

func (Goal) TableName() string {
	return "goal"
}

func toPostgresGoal(g *models.Goal) *Goal {
	return &Goal{
		ID:          g.ID,
		UserID:      g.UserID,
		Name:        g.Name,
		ProjectID:   g.ProjectID,
		Description: g.Description,
		TimeStart:   g.TimeStart,
		TimeEnd:     g.TimeEnd,
		HoursCount:  g.HoursCount,
	}
}

func toModelGoal(g *Goal) *models.Goal {
	return &models.Goal{
		ID:          g.ID,
		UserID:      g.UserID,
		Name:        g.Name,
		ProjectID:   g.ProjectID,
		Description: g.Description,
		TimeStart:   g.TimeStart,
		TimeEnd:     g.TimeEnd,
		HoursCount:  g.HoursCount,
	}
}

func toModelGoals(goals []*Goal) []*models.Goal {
	out := make([]*models.Goal, len(goals))

	for i, b := range goals {
		out[i] = toModelGoal(b)
	}

	return out
}

type goalRepository struct {
	db *gorm.DB
}

func (gr goalRepository) CreateGoal(g *models.Goal) error {
	postgresGoal := toPostgresGoal(g)

	tx := gr.db.Create(postgresGoal)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table goal)")
	}

	g.ID = postgresGoal.ID
	return nil
}

func (gr goalRepository) UpdateGoal(g *models.Goal) error {
	postgresGoal := toPostgresGoal(g)

	tx := gr.db.Omit("id").Updates(postgresGoal)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table goal)")
	}

	return nil
}

func (gr goalRepository) GetGoal(id uint64) (*models.Goal, error) {
	var goal Goal

	tx := gr.db.Where("id = ?", id).Take(&goal)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table goal)")
	}

	return toModelGoal(&goal), nil
}

func (gr goalRepository) DeleteGoal(id uint64) error {
	tx := gr.db.Delete(&Goal{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table goal)")
	}

	return nil
}

func (gr goalRepository) GetUserGoals(userID uint64) ([]*models.Goal, error) {
	goals := make([]*Goal, 0, 10)

	tx := gr.db.Where(&Goal{UserID: &userID}).Find(&goals)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table goal)")
	}

	return toModelGoals(goals), nil
}

func NewGoalRepository(db *gorm.DB) repository.RepositoryI {
	return &goalRepository{
		db: db,
	}
}
