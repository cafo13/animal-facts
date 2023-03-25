package fact_test

import (
	"strings"
	"testing"
	"time"

	"github.com/cafo13/animal-facts/backend/facts/domain/fact"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFact(t *testing.T) {
	t.Parallel()
	factUUID := uuid.New().String()
	userUUID := uuid.New().String()
	userName := "user name"
	factTime := time.Now().Round(time.Hour)

	tr, err := fact.NewFact(factUUID, userUUID, userName, factTime)
	require.NoError(t, err)

	assert.Equal(t, factUUID, tr.UUID())
	assert.Equal(t, userUUID, tr.UserUUID())
	assert.Equal(t, factTime, tr.Time())
	assert.Equal(t, userName, tr.UserName())
}

func TestNewFact_invalid(t *testing.T) {
	t.Parallel()
	factUUID := uuid.New().String()
	userUUID := uuid.New().String()
	factTime := time.Now().Round(time.Hour)
	userName := "user name"

	_, err := fact.NewFact("", userUUID, userName, factTime)
	assert.Error(t, err)

	_, err = fact.NewFact(factUUID, "", userName, factTime)
	assert.Error(t, err)

	_, err = fact.NewFact(factUUID, userUUID, userName, time.Time{})
	assert.Error(t, err)

	_, err = fact.NewFact(factUUID, userUUID, "", time.Time{})
	assert.Error(t, err)
}

func TestFact_UpdateNotes(t *testing.T) {
	t.Parallel()
	tr := newExampleFact(t)
	// it's always a good idea to ensure about pre-conditions in the test ;-)
	require.Equal(t, "", tr.Notes())

	err := tr.UpdateNotes("foo")
	require.NoError(t, err)
	assert.Equal(t, "foo", tr.Notes())
}

func TestFact_UpdateNotes_too_long(t *testing.T) {
	t.Parallel()
	tr := newExampleFact(t)

	err := tr.UpdateNotes(strings.Repeat("x", 1001))
	assert.EqualError(t, err, fact.ErrNoteTooLong.Error())
}

func TestFact_MoreThanDayUntilFact(t *testing.T) {
	t.Parallel()
	factNow := newExampleFactWithTime(t, time.Now())
	assert.False(t, factNow.CanBeCanceledForFree())

	factInTwoDays := newExampleFactWithTime(t, time.Now().AddDate(0, 0, 2))
	assert.True(t, factInTwoDays.CanBeCanceledForFree())
}

func newExampleFact(t *testing.T) *fact.Fact {
	tr, err := fact.NewFact(
		uuid.New().String(),
		uuid.New().String(),
		"user name",
		time.Now().AddDate(0, 0, 5).Round(time.Hour),
	)
	require.NoError(t, err)

	return tr
}

func newExampleFactWithTime(t *testing.T, factTime time.Time) *fact.Fact {
	tr, err := fact.NewFact(
		uuid.New().String(),
		uuid.New().String(),
		"user name",
		factTime,
	)
	require.NoError(t, err)

	return tr
}

func newCanceledFact(t *testing.T) *fact.Fact {
	tr := newExampleFact(t)
	require.NoError(t, tr.Cancel())

	return tr
}
