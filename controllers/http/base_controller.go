package http

import (
  cmderrs "go-shortener/commands/errors"
  "log"
  "net/http"
)

type BaseController struct {
  resp http.ResponseWriter
  req  *http.Request
}

func (c *BaseController) GetQueryParams(key string) string {
  return c.req.URL.Query().Get(key)
}

func (c *BaseController) ReadBodyParams(any) error {

  return nil
}

func (c *BaseController) HandleError(err error) {
  log.Println(err)

  switch err.(type) {
  case cmderrs.InvalidParamsError:
    c.resp.WriteHeader(http.StatusBadRequest)
    c.resp.Write([]byte(err.Error()))
  default:
    c.resp.WriteHeader(http.StatusInternalServerError)
    c.resp.Write([]byte("internal server error"))
  }
}
