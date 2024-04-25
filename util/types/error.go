package types

type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return e.Msg
}

func (e Error) RuntimeError() {
	//TODO implement me
	panic("implement me")
}

var ErrInvalidInput = Error{
	Code: 1,
	Msg:  "invalid input",
}
