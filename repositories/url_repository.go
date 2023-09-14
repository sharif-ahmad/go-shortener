package repositories

import (
  "go-shortener/models"
)

type URLRepository interface {
  FindAll(filters map[string]any) ([]*models.URL, error)
  Find(id int) (*models.URL, error)
  FindBy(params map[string]any) (*models.URL, error)
  Create(url *models.URL) (*models.URL, error)
  Update(id int, fields map[string]any) (*models.URL, error)
  UpdateBy(params map[string]any, fields map[string]any) (*models.URL, error)
}
