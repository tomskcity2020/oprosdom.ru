package dz6

type ToSliceInterface interface {
	ToSlice() (interface{}, error)
}

func (obj *Member) ToSlice() (interface{}, error) {

	slice := []Member{*obj}

	// слайс создается внутри функции поэтому передаем не указатель
	return slice, nil

}

func (obj *Kvartira) ToSlice() (interface{}, error) {

	slice := []Kvartira{*obj}

	return slice, nil

}
