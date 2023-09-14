package errors

type NotFoundError struct {
  Msg string
}

func (e NotFoundError) Error() string {
  return e.Msg
}
