package adapters

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/cafo13/animal-facts/backend/facts/app/query"
	"github.com/cafo13/animal-facts/backend/facts/domain/fact"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type FactModel struct {
	UUID uuid.UUID `firestore:"Uuid"`

	Text   string `firestore:"Text"`
	Source string `firestore:"Source"`
}

type FactsFirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewFactsFirestoreRepository(firestoreClient *firestore.Client) FactsFirestoreRepository {
	return FactsFirestoreRepository{
		firestoreClient: firestoreClient,
	}
}

func (r FactsFirestoreRepository) factsCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("facts")
}

func (r FactsFirestoreRepository) AddFact(ctx context.Context, f *fact.Fact) error {
	collection := r.factsCollection()

	factModel := r.marshalFact(f)

	return r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Create(collection.Doc(factModel.UUID.String()), factModel)
	})
}

func (r FactsFirestoreRepository) GetFact(ctx context.Context, factUUID uuid.UUID) (*fact.Fact, error) {
	firestoreFact, err := r.factsCollection().Doc(factUUID.String()).Get(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get fact")
	}

	f, err := r.unmarshalFact(firestoreFact)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r FactsFirestoreRepository) UpdateFact(
	ctx context.Context,
	factUUID uuid.UUID,
	updateFn func(ctx context.Context, f *fact.Fact) (*fact.Fact, error),
) error {
	factsCollection := r.factsCollection()

	return r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		documentRef := factsCollection.Doc(factUUID.String())

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

		return tx.Set(documentRef, r.marshalFact(updatedFact))
	})
}

func (r FactsFirestoreRepository) DeleteFact(ctx context.Context, factUUID uuid.UUID) error {
	_, err := r.factsCollection().Doc(factUUID.String()).Delete(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to delete fact")
	}

	return nil
}

func (r FactsFirestoreRepository) FindFactByID(ctx context.Context, factUUID string) (query.Fact, error) {
	q := r.
		factsCollection().
		Query.
		Where("Uuid", "==", factUUID)

	doc, err := q.Documents(ctx).Next()
	if err != nil {
		return query.Fact{}, err
	}

	f, err := r.unmarshalFact(doc)
	if err != nil {
		return query.Fact{}, err
	}

	return query.Fact{
		UUID:   f.UUID(),
		Text:   f.Text(),
		Source: f.Source(),
	}, nil
}

func (r FactsFirestoreRepository) FindRandomFact(ctx context.Context) (query.Fact, error) {
	// query := r.
	// 	factsCollection().
	// 	Query.
	// 	Where("Uuid", "==", "1")

	// iter := query.Documents(ctx)

	return query.Fact{
		UUID:   uuid.New(),
		Text:   "getting random fact is WIP needs to be implemented",
		Source: "getting random fact is WIP needs to be implemented",
	}, nil
}

func (r FactsFirestoreRepository) marshalFact(f *fact.Fact) FactModel {
	factModel := FactModel{
		UUID:   f.UUID(),
		Text:   f.Text(),
		Source: f.Source(),
	}

	return factModel
}

func (r FactsFirestoreRepository) unmarshalFact(doc *firestore.DocumentSnapshot) (*fact.Fact, error) {
	factModel := FactModel{}
	err := doc.DataTo(&factModel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load document")
	}

	return fact.NewFact(
		factModel.UUID,
		factModel.Text,
		factModel.Source,
	)
}
