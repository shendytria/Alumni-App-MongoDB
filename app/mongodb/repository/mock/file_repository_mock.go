package repository

import (
	"errors"
	"alumni-app/app/mongodb/model"
)

type MockFileRepository struct {
	files []model.File
	err   error
}

func NewMockFileRepository() *MockFileRepository {
	return &MockFileRepository{
		files: []model.File{},
		err:   nil,
	}
}

func (m *MockFileRepository) Create(file *model.File) error {
	if m.err != nil {
		return m.err
	}
	file.ID = file.ID
	m.files = append(m.files, *file)
	return nil
}

func (m *MockFileRepository) FindAll() ([]model.File, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.files, nil
}

func (m *MockFileRepository) FindByID(id string) (*model.File, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, f := range m.files {
		if f.ID.Hex() == id {
			return &f, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockFileRepository) Delete(id string) error {
	if m.err != nil {
		return m.err
	}
	for i, f := range m.files {
		if f.ID.Hex() == id {
			m.files = append(m.files[:i], m.files[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
