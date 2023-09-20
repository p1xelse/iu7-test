package pkg

import "gopkg.in/go-playground/validator.v9"

type Response struct {
	Body interface{} `json:"body"`
}

func IsRequestValid(entity interface{}) (bool, error) {
	validate := validator.New()
	err := validate.Struct(entity)
	if err != nil {
		return false, err
	}
	return true, nil
}
