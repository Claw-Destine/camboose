package datatypes

type WrongTypeError struct {
	What string
}

func (wte *WrongTypeError) Error() string {
	return wte.What
}
