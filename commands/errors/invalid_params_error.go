package errors

type InvalidParamsError struct {
  Msg string
}

func (e InvalidParamsError) Error() string {
  return e.Msg
}
