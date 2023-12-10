package handler_test

import (
	"github.com/cafo13/animal-facts/pkg/repository"
	"github.com/cafo13/animal-facts/public-api/handler"
	"reflect"
	"testing"
	"time"
)

var exampleFact = repository.Fact{
	ID:        1,
	Fact:      "The Blue Whale is the largest animal that has ever lived.",
	Source:    "https://factanimal.com/blue-whale/",
	Approved:  true,
	CreatedAt: time.Now(),
	CreatedBy: "some.user",
	UpdatedAt: time.Now(),
	UpdatedBy: "some.user",
}

func TestFactsHandler_Get(t *testing.T) {
	type fields struct {
		factsRepository repository.FactsRepository
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *handler.Fact
		wantErr bool
	}{
		{
			name: "get fact works",
			fields: fields{
				factsRepository: repository.NewMockFactsRepository(
					map[int]*repository.Fact{
						1: &exampleFact,
					},
					false,
				),
			},
			args: args{
				id: 1,
			},
			want: &handler.Fact{
				Fact:   "The Blue Whale is the largest animal that has ever lived.",
				Source: "https://factanimal.com/blue-whale/",
			},
			wantErr: false,
		},
		{
			name: "get fact errors on fact not found",
			fields: fields{
				factsRepository: repository.NewMockFactsRepository(
					map[int]*repository.Fact{},
					false,
				),
			},
			args: args{
				id: 1,
			},
			wantErr: true,
		},
		{
			name: "get fact errors on repository get fact failure",
			fields: fields{
				factsRepository: repository.NewMockFactsRepository(
					map[int]*repository.Fact{
						1: &exampleFact,
					},
					true,
				),
			},
			args: args{
				id: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := handler.NewFactsHandler(tt.fields.factsRepository)
			got, err := f.Get(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactsHandler_GetRandomApproved(t *testing.T) {
	type fields struct {
		factsRepository repository.FactsRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    *handler.Fact
		wantErr bool
	}{
		{
			name: "get random fact success",
			fields: fields{
				factsRepository: repository.NewMockFactsRepository(
					map[int]*repository.Fact{
						1: &exampleFact,
					},
					false,
				),
			},
			want: &handler.Fact{
				Fact:   "The Blue Whale is the largest animal that has ever lived.",
				Source: "https://factanimal.com/blue-whale/",
			},
			wantErr: false,
		},
		{
			name: "get random fact errors due to no facts",
			fields: fields{
				factsRepository: repository.NewMockFactsRepository(
					map[int]*repository.Fact{},
					false,
				),
			},
			wantErr: true,
		},
		{
			name: "get random fact errors due to repository error",
			fields: fields{
				factsRepository: repository.NewMockFactsRepository(
					map[int]*repository.Fact{
						1: &exampleFact,
					},
					true,
				),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := handler.NewFactsHandler(tt.fields.factsRepository)
			got, err := f.GetRandomApproved()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRandomApproved() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRandomApproved() got = %v, want %v", got, tt.want)
			}
		})
	}
}
