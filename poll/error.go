package poll

type Error string

func (e Error) Error() string {
	return string(e)
}

func (e Error) String() string {
	return string(e)
}

const (
	ErrInvalidFd = Error("invalid fd")
	ErrTimeout   = Error("timeout")
)
