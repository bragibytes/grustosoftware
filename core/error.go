package core

type err struct {
	msg  string
	code int
}

func NewError(msg string, code int) err {
	x := err{
		msg,
		code,
	}
	return x
}

func (e err) Error() string {
	return e.msg
}

func (e err) IsNil() bool {
	if e.msg == "" {
		return true
	}
	return false
}
