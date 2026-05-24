package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ChatMessage struct {
	bun.BaseModel `bun:"table:messages"`

	ID        uuid.UUID  `bun:"id,pk,nullzero" json:"id"`
	Content   string     `bun:"content,nullzero" json:"content"`
	CreatedAt time.Time  `bun:"created_at,nullzero" json:"created_at"`
	UpdatedAt *time.Time `bun:"updated_at,nullzero" json:"updated_at,omitempty"`
}
