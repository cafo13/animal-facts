package fact

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	FactUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("fact '%s' not found", e.FactUUID)
}

type Repository interface {
	AddFact(ctx context.Context, f *Fact) error
	GetFact(ctx context.Context, factUUID string) (*Fact, error)
	DeleteFact(ctx context.Context, factUUID string) error
	UpdateFact(
		ctx context.Context,
		factUUID string,
		updateFn func(ctx context.Context, tr *Fact) (*Fact, error),
	) error
}
