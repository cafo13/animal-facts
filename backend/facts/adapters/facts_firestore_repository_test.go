package adapters_test

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/cafo13/animal-facts/backend/facts/adapters"
	"github.com/cafo13/animal-facts/backend/facts/app/query"
	"github.com/cafo13/animal-facts/backend/facts/domain/fact"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// todo - make tests parallel after fix of emulator: https://github.com/firebase/firebase-tools/issues/2452

func TestFactsFirestoreRepository_AddFact(t *testing.T) {
	t.Parallel()
	repo := newFirebaseRepository(t)

	testCases := []struct {
		Name                string
		FactConstructor func(t *testing.T) *fact.Fact
	}{
		{
			Name:                "standard_fact",
			FactConstructor: newExampleFact,
		},
		{
			Name:                "cancelled_fact",
			FactConstructor: newCanceledFact,
		},
		{
			Name:                "fact_with_note",
			FactConstructor: newFactWithNote,
		},
		{
			Name:                "fact_with_proposed_reschedule",
			FactConstructor: newFactWithProposedReschedule,
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			expectedFact := c.FactConstructor(t)

			err := repo.AddFact(ctx, expectedFact)
			require.NoError(t, err)

			assertPersistedFactEquals(t, repo, expectedFact)
		})
	}
}

func TestFactsFirestoreRepository_UpdateFact(t *testing.T) {
	t.Parallel()
	repo := newFirebaseRepository(t)
	ctx := context.Background()

	expectedFact := newExampleFact(t)

	err := repo.AddFact(ctx, expectedFact)
	require.NoError(t, err)

	var updatedFact *fact.Fact

	err = repo.UpdateFact(
		ctx,
		expectedFact.UUID(),
		fact.MustNewUser(expectedFact.UserUUID(), fact.Attendee),
		func(ctx context.Context, tr *fact.Fact) (*fact.Fact, error) {
			assertFactsEquals(t, expectedFact, tr)

			err := tr.UpdateNotes("note")
			require.NoError(t, err)

			updatedFact = tr

			return tr, nil
		},
	)
	require.NoError(t, err)

	assertPersistedFactEquals(t, repo, updatedFact)
}

func TestFactsFirestoreRepository_GetFact_not_exists(t *testing.T) {
	t.Parallel()
	repo := newFirebaseRepository(t)

	factUUID := uuid.New().String()

	tr, err := repo.GetFact(
		context.Background(),
		factUUID,
		fact.MustNewUser(uuid.New().String(), fact.Attendee),
	)
	assert.Nil(t, tr)
	assert.EqualError(t, err, fact.NotFoundError{factUUID}.Error())
}

func TestFactsFirestoreRepository_get_and_update_another_users_fact(t *testing.T) {
	t.Parallel()
	repo := newFirebaseRepository(t)

	ctx := context.Background()
	tr := newExampleFact(t)

	err := repo.AddFact(ctx, tr)
	require.NoError(t, err)

	assertPersistedFactEquals(t, repo, tr)

	requestingUser := fact.MustNewUser(uuid.New().String(), fact.Attendee)

	_, err = repo.GetFact(
		context.Background(),
		tr.UUID(),
		requestingUser,
	)
	assert.EqualError(
		t,
		err,
		fact.ForbiddenToSeeFactError{
			RequestingUserUUID: requestingUser.UUID(),
			FactOwnerUUID:  tr.UserUUID(),
		}.Error(),
	)

	err = repo.UpdateFact(
		ctx,
		tr.UUID(),
		requestingUser,
		func(ctx context.Context, tr *fact.Fact) (*fact.Fact, error) {
			return nil, nil
		},
	)
	assert.EqualError(
		t,
		err,
		fact.ForbiddenToSeeFactError{
			RequestingUserUUID: requestingUser.UUID(),
			FactOwnerUUID:  tr.UserUUID(),
		}.Error(),
	)
}

func TestFactsFirestoreRepository_AllFacts(t *testing.T) {
	t.Parallel()
	repo := newFirebaseRepository(t)

	// AllFacts returns all documents, because of that we need to do exception and do DB cleanup
	// In general, I recommend to do it before test. In that way you are sure that cleanup is done.
	// Thanks to that tests are more stable.
	// More about why it is important you can find in https://threedots.tech/post/database-integration-testing/
	err := repo.RemoveAllFacts(context.Background())
	require.NoError(t, err)

	ctx := context.Background()

	exampleFact := newExampleFact(t)
	canceledFact := newCanceledFact(t)
	factWithNote := newFactWithNote(t)
	factWithProposedReschedule := newFactWithProposedReschedule(t)

	factsToAdd := []*fact.Fact{
		exampleFact,
		canceledFact,
		factWithNote,
		factWithProposedReschedule,
	}

	for _, tr := range factsToAdd {
		err = repo.AddFact(ctx, tr)
		require.NoError(t, err)
	}

	facts, err := repo.AllFacts(context.Background())
	require.NoError(t, err)

	proposedNewTime := factWithProposedReschedule.ProposedNewTime()
	proposer := factWithProposedReschedule.MovedProposedBy().String()

	expectedFacts := []query.Fact{
		{
			UUID:           exampleFact.UUID(),
			UserUUID:       exampleFact.UserUUID(),
			User:           "User",
			Time:           exampleFact.Time(),
			Notes:          "",
			CanBeCancelled: true,
		},
		{
			UUID:           factWithNote.UUID(),
			UserUUID:       factWithNote.UserUUID(),
			User:           "User",
			Time:           factWithNote.Time(),
			Notes:          factWithNote.Notes(),
			CanBeCancelled: true,
		},
		{
			UUID:           factWithProposedReschedule.UUID(),
			UserUUID:       factWithProposedReschedule.UserUUID(),
			User:           "User",
			Time:           factWithProposedReschedule.Time(),
			Notes:          "",
			ProposedTime:   &proposedNewTime,
			MoveProposedBy: &proposer,
			CanBeCancelled: true,
		},
	}

	var filteredFacts []query.Fact
	for _, tr := range facts {
		for _, ex := range expectedFacts {
			if tr.UUID == ex.UUID {
				filteredFacts = append(filteredFacts, tr)
			}
		}
	}

	assertQueryFactsEquals(t, expectedFacts, filteredFacts)
}

func TestFactsFirestoreRepository_FindFactsForUser(t *testing.T) {
	t.Parallel()
	repo := newFirebaseRepository(t)

	ctx := context.Background()

	userUUID := uuid.New().String()

	tr1, err := fact.NewFact(
		uuid.New().String(),
		userUUID,
		"User",
		time.Now(),
	)
	require.NoError(t, err)

	err = repo.AddFact(ctx, tr1)
	require.NoError(t, err)

	tr2, err := fact.NewFact(
		uuid.New().String(),
		userUUID,
		"User",
		time.Now(),
	)
	require.NoError(t, err)

	err = repo.AddFact(ctx, tr2)
	require.NoError(t, err)

	// this fact should be not in the list
	canceledFact, err := fact.NewFact(
		uuid.New().String(),
		userUUID,
		"User",
		time.Now(),
	)
	require.NoError(t, err)

	err = canceledFact.Cancel()
	require.NoError(t, err)

	err = repo.AddFact(ctx, canceledFact)
	require.NoError(t, err)

	facts, err := repo.FindFactsForUser(context.Background(), userUUID)
	require.NoError(t, err)

	assertQueryFactsEquals(t, facts, []query.Fact{
		{
			UUID:     tr1.UUID(),
			UserUUID: userUUID,
			User:     "User",
			Time:     tr1.Time(),
		},
		{
			UUID:     tr2.UUID(),
			UserUUID: userUUID,
			User:     "User",
			Time:     tr2.Time(),
		},
	})
}

func newRandomFactTime() time.Time {
	min := time.Now().AddDate(0, 0, 5).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func newExampleFact(t *testing.T) *fact.Fact {
	t.Helper()
	tr, err := fact.NewFact(
		uuid.New().String(),
		uuid.New().String(),
		"User",
		newRandomFactTime(),
	)
	require.NoError(t, err)

	return tr
}

func newCanceledFact(t *testing.T) *fact.Fact {
	t.Helper()
	tr, err := fact.NewFact(
		uuid.New().String(),
		uuid.New().String(),
		"User",
		newRandomFactTime(),
	)
	require.NoError(t, err)

	err = tr.Cancel()
	require.NoError(t, err)

	return tr
}

func newFactWithNote(t *testing.T) *fact.Fact {
	t.Helper()
	tr := newExampleFact(t)
	err := tr.UpdateNotes("foo")
	require.NoError(t, err)

	return tr
}

func newFactWithProposedReschedule(t *testing.T) *fact.Fact {
	t.Helper()
	tr := newExampleFact(t)
	tr.ProposeReschedule(time.Now().AddDate(0, 0, 14), fact.Trainer)

	return tr
}

func assertPersistedFactEquals(t *testing.T, repo adapters.FactsFirestoreRepository, tr *fact.Fact) {
	t.Helper()
	persistedFact, err := repo.GetFact(
		context.Background(),
		tr.UUID(),
		fact.MustNewUser(tr.UserUUID(), fact.Attendee),
	)
	require.NoError(t, err)

	assertFactsEquals(t, tr, persistedFact)
}

// Firestore is not storing time with same precision, so we need to round it a bit
var cmpRoundTimeOpt = cmp.Comparer(func(x, y time.Time) bool {
	return x.Truncate(time.Millisecond).Equal(y.Truncate(time.Millisecond))
})

func assertFactsEquals(t *testing.T, tr1, tr2 *fact.Fact) {
	t.Helper()
	cmpOpts := []cmp.Option{
		cmpRoundTimeOpt,
		cmp.AllowUnexported(
			fact.UserType{},
			time.Time{},
			fact.Fact{},
		),
	}

	assert.True(
		t,
		cmp.Equal(tr1, tr2, cmpOpts...),
		cmp.Diff(tr1, tr2, cmpOpts...),
	)
}

func assertQueryFactsEquals(t *testing.T, expectedFacts, facts []query.Fact) bool {
	t.Helper()
	cmpOpts := []cmp.Option{
		cmpRoundTimeOpt,
		cmpopts.SortSlices(func(x, y query.Fact) bool {
			return x.Time.After(y.Time)
		}),
	}
	return assert.True(t,
		cmp.Equal(expectedFacts, facts, cmpOpts...),
		cmp.Diff(expectedFacts, facts, cmpOpts...),
	)
}

func newFirebaseRepository(t *testing.T) adapters.FactsFirestoreRepository {
	t.Helper()
	firestoreClient, err := firestore.NewClient(context.Background(), os.Getenv("GCP_PROJECT"))
	require.NoError(t, err)

	return adapters.NewFactsFirestoreRepository(firestoreClient)
}
