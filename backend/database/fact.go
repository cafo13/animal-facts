package database

import (
	"math/rand"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Fact struct {
	Database
	gorm.Model
	Text     string `json:"Text,omitempty"`
	Category string `json:"Category,omitempty"`
	Source   string `json:"Source,omitempty"`
	Approved bool   `json:"Approved,omitempty"`
}

type FactId struct {
	ID uint `gorm:"primarykey"`
}

func (f *Fact) Create() error {
	err := f.Database.db.Create(&f).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *Fact) Read() error {
	err := f.Database.db.First(&f, f.ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *Fact) Update() error {
	err := f.Database.db.Updates(&f).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *Fact) Delete() error {
	err := f.Database.db.Delete(&f).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *Fact) GetRandomFactId() (uint, error) {
	factIds := []FactId{}
	err := f.Database.db.Model(&Fact{}).Where("deleted_at IS NULL").Find(&factIds).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to get all fact Ids")
	}

	factCount := len(factIds)
	if factCount < 1 {
		return 0, errors.New("no facts found, unable to select random fact")
	}

	factId := factIds[rand.Intn(factCount)]

	return factId.ID, nil
}
