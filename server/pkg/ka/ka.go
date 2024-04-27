package ka

type BoomError struct {
	Msg string
}

var _ error = BoomError{}

func Boom(msg string) BoomError {
	return BoomError{Msg: msg}
}

func (k BoomError) Error() string {
	return k.Msg
}
