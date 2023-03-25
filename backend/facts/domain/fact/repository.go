package fact

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type NotFoundError struct {
	FactUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("fact '%s' not found", e.FactUUID)
}

type Repository interface {
	AddFact(ctx context.Context, f *Fact) error
	GetFact(ctx context.Context, factUUID uuid.UUID) (*Fact, error)
	DeleteFact(ctx context.Context, factUUID uuid.UUID) error
	UpdateFact(
		ctx context.Context,
		factUUID uuid.UUID,
		updateFn func(ctx context.Context, tr *Fact) (*Fact, error),
	) error
}
