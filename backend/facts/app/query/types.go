package query

import "github.com/google/uuid"

type Fact struct {
	UUID   uuid.UUID
	Text   string
	Source string
}
