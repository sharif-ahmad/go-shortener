package commands

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	cmderrs "go-shortener/commands/errors"
	"go-shortener/models"
	"go-shortener/repositories"
	"math/rand"
)

const (
	MAX_TRY = 5
)

type CreateShortURLCommand struct {
	Url     string
	UrlRepo repositories.URLRepository
}

func NewCreateShortURLCommand(url string, urlRepo repositories.URLRepository) *CreateShortURLCommand {
	return &CreateShortURLCommand{url, urlRepo}
}

func (c *CreateShortURLCommand) ValidateParams() error {
	return validation.Validate(c.Url,
		validation.Required,       // not empty
		validation.Length(5, 100), // length between 5 and 100
		is.URL,                    // is a valid URL
	)
}

func (c *CreateShortURLCommand) Call() (string, error) {
	if err := c.ValidateParams(); err != nil {
		return "", cmderrs.InvalidParamsError{Msg: err.Error()}
	}

	shortCode := randStr(10)
	for i := 0; i < MAX_TRY; i++ {
		res, err := c.UrlRepo.FindBy(map[string]any{"short_code": shortCode})
		if err != nil {
			return "", err
		}

		if res == nil {
			break
		}
		shortCode = randStr(10) // regenerate
	}

	if _, err := c.UrlRepo.Create(&models.URL{Original: c.Url, ShortCode: shortCode}); err != nil {
		return "", err
	}

	return shortCode, nil
}

// Taken from https://www.golinuxcloud.com/golang-generate-random-string/
// define the given charset, char only
var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var number = []byte("0123456789")
var alphaNumeric = append(charset, number...)

// n is the length of random string we want to generate
func randStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		// randomly select 1 character from given charset
		b[i] = alphaNumeric[rand.Intn(len(alphaNumeric))]
	}
	return string(b)
}
