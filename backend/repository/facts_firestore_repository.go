package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FactsFirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewFactsFirestoreRepository(firestoreClient *firestore.Client) FactsRepository {
	return FactsFirestoreRepository{
		firestoreClient: firestoreClient,
	}
}

func (r FactsFirestoreRepository) factsCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("facts")
}

func (r FactsFirestoreRepository) CreateFact(ctx context.Context, f *FactModel) error {
	collection := r.factsCollection()

	return r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Create(collection.Doc(uuid.NewString()), f)
	})
}

func (r FactsFirestoreRepository) ReadFact(ctx context.Context, factUUID string) (*FactModel, error) {
	firestoreFact, err := r.factsCollection().Doc(factUUID).Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get fact with uuid %s", factUUID)
	}

	f, err := r.unmarshalFact(firestoreFact)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r FactsFirestoreRepository) ReadRandomFact(ctx context.Context) (*FactModel, error) {
	// TODO: Get all UUIDs from facts collection, select random one
	firestoreFact, err := r.factsCollection().Doc(uuid.NewString()).Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get random fact")
	}

	f, err := r.unmarshalFact(firestoreFact)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r FactsFirestoreRepository) UpdateFact(
	ctx context.Context,
	factUUID string,
	updateFn func(ctx context.Context, f *FactModel) (*FactModel, error),
) error {
	factsCollection := r.factsCollection()

	return r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		documentRef := factsCollection.Doc(factUUID)

		firestoreFact, err := tx.Get(documentRef)
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}

		f, err := r.unmarshalFact(firestoreFact)
		if err != nil {
			return err
		}

		updatedFact, err := updateFn(ctx, f)
		if err != nil {
			return err
		}

		return tx.Set(documentRef, updatedFact)
	})
}

func (r FactsFirestoreRepository) DeleteFact(ctx context.Context, factUUID string) error {
	_, err := r.factsCollection().Doc(factUUID).Delete(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to delete fact with uuid %s", factUUID)
	}

	return nil
}

func (r FactsFirestoreRepository) unmarshalFact(doc *firestore.DocumentSnapshot) (*FactModel, error) {
	FactModel := FactModel{}
	err := doc.DataTo(&FactModel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load document")
	}

	return &FactModel, nil
}
