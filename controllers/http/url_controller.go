package http

import (
  "go-shortener/commands"
  "go-shortener/repositories"
  "net/http"
)

type URLController struct {
  BaseController
  BaseUrl string
  urlRepo repositories.URLRepository
}

func NewURLController(
  baseUrl string, resp http.ResponseWriter, req *http.Request, urlRepo repositories.URLRepository,
) *URLController {
  return &URLController{BaseController{resp, req}, baseUrl, urlRepo}
}

func (c *URLController) Put() {
  url := c.GetQueryParams("url")
  shortenUrlCmd := commands.NewShortenURL(url, c.urlRepo)

  hash, err := shortenUrlCmd.Call()
  if err != nil {
    c.HandleError(err)
    return
  }

  c.resp.WriteHeader(http.StatusOK)
  c.resp.Write([]byte(c.BaseUrl + "/u/" + hash))
}
