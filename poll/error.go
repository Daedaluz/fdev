package poll

type Error struct {
	msg string
	err error
}

func (e Error) Error() string {
	return e.msg
}

func (e Error) String() string {
	return e.msg
}

func (e Error) Unwrap() error {
	return e.err
}

func wrap(e error) error {
	return Error{
		msg: e.Error(),
		err: e,
	}
}

var (
	ErrInvalidFd = Error{"invalid fd", nil}
	ErrTimeout   = Error{"timeout", nil}
)
