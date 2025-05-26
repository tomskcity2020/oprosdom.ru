package dz6

import (
	"errors"
)

func CreateModel(model string) (ToSliceInterface, error) {

	switch model {
	case "member":
		return &Member{}, nil
	case "kvartira":
		return &Kvartira{}, nil
	default:
		return nil, errors.New("no model found")
	}

}
