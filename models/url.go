package models

import "time"

type URL struct {
  ID        int
  Original  string
  ShortCode string
  CreatedAt time.Time
  UpdatedAt time.Time
}
