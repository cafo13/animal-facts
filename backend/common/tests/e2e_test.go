package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateFact(t *testing.T) {
	t.Parallel()

	hour := RelativeDate(12, 12)

	userID := "TestCreateFact-user"
	trainerJWT := FakeTrainerJWT(t, uuid.New().String())
	attendeeJWT := FakeAttendeeJWT(t, userID)
	factsHTTPClient := NewFactsHTTPClient(t, attendeeJWT)

	// Cancel the fact if exists and make the hour available
	facts := factsHTTPClient.GetFacts(t)
	for _, fact := range facts.Facts {
		if fact.Time.Equal(hour) {
			factsTrainerHTTPClient := NewFactsHTTPClient(t, trainerJWT)
			factsTrainerHTTPClient.CancelFact(t, fact.Uuid, 200)
			break
		}
	}

	factUUID := factsHTTPClient.CreateFact(t, "some note", hour)

	factsResponse := factsHTTPClient.GetFacts(t)
	require.Len(t, factsResponse.Facts, 1)
	require.Equal(t, factUUID, factsResponse.Facts[0].Uuid, "Attendee should see the fact")
}
