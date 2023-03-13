package facts_test

import (
	"testing"

	"github.com/cafo13/animal-facts/api/database"
	"github.com/cafo13/animal-facts/api/models"
)

func Test_Facts_DatabaseFactToModelsFactWithID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer cmockCtrltrl.Finish()

	mockCtrl.NewMockFactHandler()

	var tests = []struct {
		name               string
		dbFact             *database.Fact
		expectedModelsFact *models.FactWithID
	}{
		{
			"Full db fact to mode fact",
			&database.Fact{},
			&models.FactWithID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func (t *t.testing.T){
			expectedFact := tt.expectedModelsFact

			if !want.MatchString(msg) || err != nil {
				t.Errorf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
			}
		})
	}
}
