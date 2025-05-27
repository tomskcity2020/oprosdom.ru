package dz6

import "errors"

var members []Member
var kvartirs []Kvartira

type ToSliceInterface interface {
	ToSlice()
}

func (obj *Member) ToSlice()   {}
func (obj *Kvartira) ToSlice() {}

func ToSlice(m ToSliceInterface) (interface{}, error) {

	switch getType := m.(type) {
	case *Member:
		members = append(members, *getType)
		return &members, nil
	case *Kvartira:
		kvartirs = append(kvartirs, *getType)
		return &kvartirs, nil
	default:
		return nil, errors.New("incorrect type")
	}

}
