package integrationutils

import (
	"timetracker/models"
	"timetracker/models/dto"
)

func EntryModel2Dto (entry *models.Entry) dto.ReqCreateUpdateEntry {
	dtoObj := dto.ReqCreateUpdateEntry{
		ID          : entry.ID,
		ProjectID   : entry.ProjectID,
		Description : entry.Description,
		TimeStart   : entry.TimeStart,
		TimeEnd     : entry.TimeEnd,
		TagList: []uint64{},
	}

	for _, tag := range entry.TagList {
		dtoObj.TagList = append(dtoObj.TagList, tag.ID)
	}

	return dtoObj
}