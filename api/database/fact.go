package database

import "gorm.io/gorm"

type Fact struct {
	Database
	gorm.Model
	Text     string `json:"text"`
	Category string `json:"category"`
	Source   string `json:"source"`
	Approved bool   `json:"approved"`
}

func (f *Fact) Create() error {
	err := f.db.Create(&f).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *Fact) Read() error {
	err := f.db.First(&f, f.ID).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *Fact) Update() error {
	err := f.db.Updates(&f).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *Fact) Delete() error {
	err := f.db.Delete(&f).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *Fact) Count() (int64, error) {
	var count int64
	err := f.db.Model(&Fact{}).Count(&count).Error
	if err != nil {
		return count, err
	}

	return count, nil
}
