package entity

import (
	"bytes"
	"time"

	uuid "github.com/google/uuid"
)

// In case of changing uuid library, it's enough to change it in one place
type Id = uuid.UUID

type File struct {
	Name      string
	Content   *bytes.Buffer
	CreatedAt time.Time
	UpdatedAt time.Time
}
