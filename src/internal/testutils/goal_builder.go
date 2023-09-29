package testutils

import (
	"time"
	"timetracker/models"
)

type GoalBuilder struct {
	goal models.Goal
}

func NewGoalBuilder() *GoalBuilder {
	return &GoalBuilder{}
}

func (b *GoalBuilder) WithID(id uint64) *GoalBuilder {
	b.goal.ID = id
	return b
}

func (b *GoalBuilder) WithUserID(userID uint64) *GoalBuilder {
	b.goal.UserID = &userID
	return b
}

func (b *GoalBuilder) WithName(name string) *GoalBuilder {
	b.goal.Name = name
	return b
}

func (b *GoalBuilder) WithProjectID(projectID uint64) *GoalBuilder {
	b.goal.ProjectID = &projectID
	return b
}

func (b *GoalBuilder) WithHoursCount(hoursCount float64) *GoalBuilder {
	b.goal.HoursCount = hoursCount
	return b
}

func (b *GoalBuilder) WithDescription(description string) *GoalBuilder {
	b.goal.Description = description
	return b
}

func (b *GoalBuilder) WithTimeStart(timeStart time.Time) *GoalBuilder {
	b.goal.TimeStart = timeStart
	return b
}

func (b *GoalBuilder) WithTimeEnd(timeEnd time.Time) *GoalBuilder {
	b.goal.TimeEnd = timeEnd
	return b
}

func (b *GoalBuilder) Build() models.Goal {
	return b.goal
}
