package integrationutils

import (
	"encoding/json"
	"time"
	"timetracker/models/dto"
)

type ReqCreateUpdateGoalBuilder struct {
	goal dto.ReqCreateUpdateGoal
}

func NewReqCreateUpdateGoalBuilder() *ReqCreateUpdateGoalBuilder {
	return &ReqCreateUpdateGoalBuilder{}
}

func (b *ReqCreateUpdateGoalBuilder) WithID(id uint64) *ReqCreateUpdateGoalBuilder {
	b.goal.ID = id
	return b
}

func (b *ReqCreateUpdateGoalBuilder) WithName(name string) *ReqCreateUpdateGoalBuilder {
	b.goal.Name = name
	return b
}

func (b *ReqCreateUpdateGoalBuilder) WithProjectID(projectID uint64) *ReqCreateUpdateGoalBuilder {
	b.goal.ProjectID = &projectID
	return b
}

func (b *ReqCreateUpdateGoalBuilder) WithHoursCount(hoursCount float64) *ReqCreateUpdateGoalBuilder {
	b.goal.HoursCount = hoursCount
	return b
}

func (b *ReqCreateUpdateGoalBuilder) WithDescription(description string) *ReqCreateUpdateGoalBuilder {
	b.goal.Description = description
	return b
}

func (b *ReqCreateUpdateGoalBuilder) WithTimeStart(timeStart time.Time) *ReqCreateUpdateGoalBuilder {
	b.goal.TimeStart = timeStart
	return b
}

func (b *ReqCreateUpdateGoalBuilder) WithTimeEnd(timeEnd time.Time) *ReqCreateUpdateGoalBuilder {
	b.goal.TimeEnd = timeEnd
	return b
}

func (b *ReqCreateUpdateGoalBuilder) Build() dto.ReqCreateUpdateGoal {
	return b.goal
}

func (b *ReqCreateUpdateGoalBuilder) Json() ([]byte, error) {
	return json.Marshal(b.goal)
}
