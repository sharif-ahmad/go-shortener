package errors

type InternalError struct {
  Msg string
}

func (e InternalError) Error() string {
  return e.Msg
}
