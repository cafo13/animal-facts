package adapters

import (
	"context"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/cafo13/animal-facts/backend/facts/app/query"
	"github.com/cafo13/animal-facts/backend/facts/domain/fact"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FactModel struct {
	UUID     string `firestore:"Uuid"`
	Text string `firestore:"Text"`
	Source string `firestore:"Source"`
}

type FactsFirestoreRepository struct {
	firestoreClient *firestore.Client
}

func NewFactsFirestoreRepository(
	firestoreClient *firestore.Client,
) FactsFirestoreRepository {
	return FactsFirestoreRepository{
		firestoreClient: firestoreClient,
	}
}

func (r FactsFirestoreRepository) factsCollection() *firestore.CollectionRef {
	return r.firestoreClient.Collection("facts")
}

func (r FactsFirestoreRepository) AddFact(ctx context.Context, tr *fact.Fact) error {
	collection := r.factsCollection()

	factModel := r.marshalFact(tr)

	return r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Create(collection.Doc(factModel.UUID), factModel)
	})
}

func (r FactsFirestoreRepository) GetFact(
	ctx context.Context,
	factUUID string,
	user fact.User,
) (*fact.Fact, error) {
	firestoreFact, err := r.factsCollection().Doc(factUUID).Get(ctx)

	if status.Code(err) == codes.NotFound {
		return nil, fact.NotFoundError{factUUID}
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to get actual docs")
	}

	tr, err := r.unmarshalFact(firestoreFact)
	if err != nil {
		return nil, err
	}

	if err := fact.CanUserSeeFact(user, *tr); err != nil {
		return nil, err
	}

	return tr, nil
}

func (r FactsFirestoreRepository) UpdateFact(
	ctx context.Context,
	factUUID string,
	user fact.User,
	updateFn func(ctx context.Context, tr *fact.Fact) (*fact.Fact, error),
) error {
	factsCollection := r.factsCollection()

	return r.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		documentRef := factsCollection.Doc(factUUID)

		firestoreFact, err := tx.Get(documentRef)
		if err != nil {
			return errors.Wrap(err, "unable to get actual docs")
		}

		tr, err := r.unmarshalFact(firestoreFact)
		if err != nil {
			return err
		}

		if err := fact.CanUserSeeFact(user, *tr); err != nil {
			return err
		}

		updatedFact, err := updateFn(ctx, tr)
		if err != nil {
			return err
		}

		return tx.Set(documentRef, r.marshalFact(updatedFact))
	})
}

func (r FactsFirestoreRepository) marshalFact(tr *fact.Fact) FactModel {
	factModel := FactModel{
		UUID:     tr.UUID(),
		UserUUID: tr.UserUUID(),
		User:     tr.UserName(),
		Time:     tr.Time(),
		Notes:    tr.Notes(),
		Canceled: tr.IsCanceled(),
	}

	if tr.IsRescheduleProposed() {
		proposedBy := tr.MovedProposedBy().String()
		proposedTime := tr.ProposedNewTime()

		factModel.MoveProposedBy = &proposedBy
		factModel.ProposedTime = &proposedTime
	}

	return factModel
}

func (r FactsFirestoreRepository) unmarshalFact(doc *firestore.DocumentSnapshot) (*fact.Fact, error) {
	factModel := FactModel{}
	err := doc.DataTo(&factModel)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load document")
	}

	var moveProposedBy fact.UserType
	if factModel.MoveProposedBy != nil {
		moveProposedBy, err = fact.NewUserTypeFromString(*factModel.MoveProposedBy)
		if err != nil {
			return nil, err
		}
	}

	var proposedTime time.Time
	if factModel.ProposedTime != nil {
		proposedTime = *factModel.ProposedTime
	}

	return fact.UnmarshalFactFromDatabase(
		factModel.UUID,
		factModel.UserUUID,
		factModel.User,
		factModel.Time,
		factModel.Notes,
		factModel.Canceled,
		proposedTime,
		moveProposedBy,
	)
}

func (r FactsFirestoreRepository) AllFacts(ctx context.Context) ([]query.Fact, error) {
	query := r.
		factsCollection().
		Query.
		Where("Time", ">=", time.Now().Add(-time.Hour*24)).
		Where("Canceled", "==", false)

	iter := query.Documents(ctx)

	return r.factModelsToQuery(iter)
}

func (r FactsFirestoreRepository) FindFactsForUser(ctx context.Context, userUUID string) ([]query.Fact, error) {
	query := r.factsCollection().Query.
		Where("Time", ">=", time.Now().Add(-time.Hour*24)).
		Where("UserUuid", "==", userUUID).
		Where("Canceled", "==", false)

	iter := query.Documents(ctx)

	return r.factModelsToQuery(iter)
}

// warning: RemoveAllFacts was designed for tests for doing data cleanups
func (r FactsFirestoreRepository) RemoveAllFacts(ctx context.Context) error {
	for {
		iter := r.factsCollection().Limit(100).Documents(ctx)
		numDeleted := 0

		batch := r.firestoreClient.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return errors.Wrap(err, "unable to get document")
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to remove docs")
		}
	}
}

func (r FactsFirestoreRepository) factModelsToQuery(iter *firestore.DocumentIterator) ([]query.Fact, error) {
	var facts []query.Fact

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		tr, err := r.unmarshalFact(doc)
		if err != nil {
			return nil, err
		}

		queryFact := query.Fact{
			UUID:           tr.UUID(),
			UserUUID:       tr.UserUUID(),
			User:           tr.UserName(),
			Time:           tr.Time(),
			Notes:          tr.Notes(),
			CanBeCancelled: tr.CanBeCanceledForFree(),
		}

		if tr.IsRescheduleProposed() {
			proposedTime := tr.ProposedNewTime()
			queryFact.ProposedTime = &proposedTime

			proposedBy := tr.MovedProposedBy().String()
			queryFact.MoveProposedBy = &proposedBy
		}

		facts = append(facts, queryFact)
	}

	sort.Slice(facts, func(i, j int) bool { return facts[i].Time.Before(facts[j].Time) })

	return facts, nil
}
