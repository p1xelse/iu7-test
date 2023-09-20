package usecase

import (
	"github.com/pkg/errors"
	goalRep "timetracker/internal/Goal/repository"
	"timetracker/models"
)

type UsecaseI interface {
	CreateGoal(e *models.Goal) error
	UpdateGoal(e *models.Goal) error
	GetGoal(id uint64) (*models.Goal, error)
	DeleteGoal(id uint64, userID uint64) error
	GetUserGoals(userID uint64) ([]*models.Goal, error)
}

type usecase struct {
	goalRepository goalRep.RepositoryI
}

func New(gRep goalRep.RepositoryI) UsecaseI {
	return &usecase{
		goalRepository: gRep,
	}
}

func (u *usecase) CreateGoal(e *models.Goal) error {
	err := u.goalRepository.CreateGoal(e)

	if err != nil {
		return errors.Wrap(err, "Error in func goal.Usecase.CreateGoal")
	}

	return nil
}

func (u *usecase) UpdateGoal(goal *models.Goal) error {
	_, err := u.goalRepository.GetGoal(goal.ID)

	if err != nil {
		return errors.Wrap(err, "Error in func goal.Usecase.Update.GetGoal")
	}

	err = u.goalRepository.UpdateGoal(goal)

	if err != nil {
		return errors.Wrap(err, "Error in func goal.Usecase.CreateGoal")
	}

	return nil
}

func (u *usecase) GetGoal(id uint64) (*models.Goal, error) {
	resGoal, err := u.goalRepository.GetGoal(id)

	if err != nil {
		return nil, errors.Wrap(err, "goal.usecase.GetGoal error while get goal info")
	}

	return resGoal, nil
}

func (u *usecase) DeleteGoal(id uint64, userID uint64) error {
	existedGoal, err := u.goalRepository.GetGoal(id)
	if err != nil {
		return err
	}

	if existedGoal == nil {
		return errors.New("Goal not found") //TODO models error
	}

	if *existedGoal.UserID != userID {
		return errors.New("Permission denied")
	}

	err = u.goalRepository.DeleteGoal(id)

	if err != nil {
		return errors.Wrap(err, "Error in func goal.Usecase.DeleteGoal repository")
	}

	return nil
}

func (u *usecase) GetUserGoals(userID uint64) ([]*models.Goal, error) {
	Goals, err := u.goalRepository.GetUserGoals(userID)

	if err != nil {
		return nil, errors.Wrap(err, "Error in func goal.Usecase.GetUserPosts")
	}

	return Goals, nil
}
