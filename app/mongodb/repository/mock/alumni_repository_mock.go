package repository

import (
	"errors"
	"time"
	"alumni-app/app/mongodb/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAlumniRepository struct {
	Data map[string]model.Alumni
}

func NewMockAlumniRepository() *MockAlumniRepository {
	return &MockAlumniRepository{
		Data: make(map[string]model.Alumni),
	}
}

func (m *MockAlumniRepository) GetAll(search, sortBy, order string, page, limit int) ([]model.Alumni, error) {
	var list []model.Alumni
	for _, a := range m.Data {
		list = append(list, a)
	}
	return list, nil
}

func (m *MockAlumniRepository) Count(search string) (int64, error) {
	return int64(len(m.Data)), nil
}

func (m *MockAlumniRepository) GetByID(id primitive.ObjectID) (model.Alumni, error) {
	a, ok := m.Data[id.Hex()]
	if !ok {
		return model.Alumni{}, errors.New("not found")
	}
	return a, nil
}

func (m *MockAlumniRepository) GetByUserID(userID primitive.ObjectID) (model.Alumni, error) {
	for _, a := range m.Data {
		if a.UserID == userID {
			return a, nil
		}
	}
	return model.Alumni{}, errors.New("not found")
}

func (m *MockAlumniRepository) GetAllByUserID(userID primitive.ObjectID) ([]model.Alumni, error) {
	var list []model.Alumni
	for _, a := range m.Data {
		if a.UserID == userID {
			list = append(list, a)
		}
	}
	return list, nil
}

func (m *MockAlumniRepository) Create(a *model.Alumni) error {
	m.Data[a.ID.Hex()] = *a
	return nil
}

func (m *MockAlumniRepository) Update(id primitive.ObjectID, a *model.Alumni) error {
	_, ok := m.Data[id.Hex()]
	if !ok {
		return errors.New("not found")
	}
	m.Data[id.Hex()] = *a
	return nil
}

func (m *MockAlumniRepository) SoftDelete(id primitive.ObjectID) error {
	a, ok := m.Data[id.Hex()]
	if !ok {
		return errors.New("not found")
	}

	now := time.Now()
	a.DeletedAt = &now
	m.Data[id.Hex()] = a
	return nil
}
