package postgres

import (
	"time"
	"timetracker/internal/Entry/repository"
	"timetracker/models"
	"timetracker/pkg"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Entry struct {
	ID          uint64    `gorm:"column:id"`
	UserID      *uint64   `gorm:"column:user_id"`
	ProjectID   *uint64   `gorm:"column:project_id;default:null"`
	Description string    `gorm:"column:description"`
	TimeStart   time.Time `gorm:"column:time_start"`
	TimeEnd     time.Time `gorm:"column:time_end"`
}

func (Entry) TableName() string {
	return "entry"
}

func toPostgresEntry(e *models.Entry) *Entry {
	return &Entry{
		ID:          e.ID,
		UserID:      e.UserID,
		ProjectID:   e.ProjectID,
		Description: e.Description,
		TimeStart:   e.TimeStart,
		TimeEnd:     e.TimeEnd,
	}
}

func toModelEntry(e *Entry) *models.Entry {
	return &models.Entry{
		ID:          e.ID,
		UserID:      e.UserID,
		ProjectID:   e.ProjectID,
		Description: e.Description,
		TimeStart:   e.TimeStart,
		TimeEnd:     e.TimeEnd,
	}
}

func toModelEntries(entries []*Entry) []*models.Entry {
	out := make([]*models.Entry, len(entries))

	for i, b := range entries {
		out[i] = toModelEntry(b)
	}

	return out
}

type entryRepository struct {
	db *gorm.DB
}

func (er *entryRepository) CreateEntry(e *models.Entry) error {
	postgresEntry := toPostgresEntry(e)

	tx := er.db.Create(postgresEntry)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table entry)")
	}

	e.ID = postgresEntry.ID
	return nil
}

func (er *entryRepository) UpdateEntry(e *models.Entry) error {
	postgresEntry := toPostgresEntry(e)

	tx := er.db.Omit("id").Updates(postgresEntry)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table entry)")
	}

	return nil
}

func (er *entryRepository) GetEntry(id uint64) (*models.Entry, error) {
	var entry Entry

	tx := er.db.Where("id = ?", id).Take(&entry)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table entry)")
	}

	return toModelEntry(&entry), nil
}

func (er *entryRepository) DeleteEntry(id uint64) error {
	tx := er.db.Delete(&Entry{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table entry)")
	}

	return nil
}

func (er *entryRepository) GetUserEntries(userID uint64) ([]*models.Entry, error) {
	entries := make([]*Entry, 0, 10)

	tx := er.db.Where(&Entry{UserID: &userID}).Find(&entries)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table entry)")
	}

	return toModelEntries(entries), nil
}

func (er *entryRepository) GetUserEntriesForDay(userID uint64, date time.Time) ([]*models.Entry, error) {
	entries := make([]*Entry, 0, 10)

	todayStart, todayEnd := pkg.GetDayInterval(date)
	tx := er.db.Where(&Entry{UserID: &userID}).Where("time_start BETWEEN ? AND ?", todayStart, todayEnd).Find(&entries)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table entry)")
	}

	return toModelEntries(entries), nil
}

func NewEntryRepository(db *gorm.DB) repository.RepositoryI {
	return &entryRepository{
		db: db,
	}
}
