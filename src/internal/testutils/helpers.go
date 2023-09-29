package testutils

import "timetracker/models"

func BuildUsersByIDs(userBuilder *UserBuilder, ids []uint64) []*models.User {
	users := make([]*models.User, 0)
	for _, id := range ids {
		user := userBuilder.WithID(id).Build()
		users = append(users, &user)
	}

	return users
}

func MakePointerSlice[T any](src []T) []*T {
	resSlice := make([]*T, len(src))
	for idx := range src {
		val := new(T)
		*val = src[idx]
		resSlice[idx] = val
	}

	return resSlice
}

