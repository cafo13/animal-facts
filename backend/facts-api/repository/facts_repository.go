package repository

import (
	"context"

	"github.com/google/uuid"
)

type FactModel struct {
	UUID uuid.UUID `firestore:"Uuid" json:"uuid"`

	Text     string `firestore:"Text" json:"text"`
	Source   string `firestore:"Source" json:"source"`
	Approved bool   `firestore:"Approved" json:"approved"`
}

type FactsRepository interface {
	CreateFact(ctx context.Context, f *FactModel) error
	ReadFact(ctx context.Context, factUUID string) (*FactModel, error)
	ReadRandomFact(ctx context.Context) (*FactModel, error)
	UpdateFact(
		ctx context.Context,
		factUUID string,
		updateFn func(ctx context.Context, f *FactModel) (*FactModel, error),
	) error
	DeleteFact(ctx context.Context, factUUID string) error
}
