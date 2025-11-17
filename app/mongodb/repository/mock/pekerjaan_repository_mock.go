package repository

import (
	"errors"
	"time"
	"alumni-app/app/mongodb/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockPekerjaanRepository struct {
	Data map[string]model.PekerjaanAlumni
}

func NewMockPekerjaanRepository() *MockPekerjaanRepository {
	return &MockPekerjaanRepository{
		Data: make(map[string]model.PekerjaanAlumni),
	}
}

func (m *MockPekerjaanRepository) GetAll(search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	var list []model.PekerjaanAlumni
	for _, p := range m.Data {
		if p.DeletedAt == nil { 
			list = append(list, p)
		}
	}
	return list, nil
}

func (m *MockPekerjaanRepository) Count(search string) (int64, error) {
	return int64(len(m.Data)), nil
}

func (m *MockPekerjaanRepository) GetByID(id primitive.ObjectID) (model.PekerjaanAlumni, error) {
	p, ok := m.Data[id.Hex()]
	if !ok || p.DeletedAt != nil {
		return model.PekerjaanAlumni{}, errors.New("not found")
	}
	return p, nil
}

func (m *MockPekerjaanRepository) GetByAlumniID(alumniID primitive.ObjectID) ([]model.PekerjaanAlumni, error) {
	var list []model.PekerjaanAlumni
	for _, p := range m.Data {
		if p.AlumniID == alumniID && p.DeletedAt == nil {
			list = append(list, p)
		}
	}
	return list, nil
}

func (m *MockPekerjaanRepository) Create(p *model.PekerjaanAlumni) error {
	m.Data[p.ID.Hex()] = *p
	return nil
}

func (m *MockPekerjaanRepository) Update(id primitive.ObjectID, p *model.PekerjaanAlumni) error {
	_, ok := m.Data[id.Hex()]
	if !ok {
		return errors.New("not found")
	}
	m.Data[id.Hex()] = *p
	return nil
}

func (m *MockPekerjaanRepository) SoftDelete(id primitive.ObjectID) error {
	p, ok := m.Data[id.Hex()]
	if !ok {
		return errors.New("not found")
	}
	now := time.Now()
	p.DeletedAt = &now
	m.Data[id.Hex()] = p
	return nil
}

func (m *MockPekerjaanRepository) GetTrashed() ([]model.PekerjaanAlumni, error) {
	var list []model.PekerjaanAlumni
	for _, p := range m.Data {
		if p.DeletedAt != nil {
			list = append(list, p)
		}
	}
	return list, nil
}

func (m *MockPekerjaanRepository) GetTrashedByAlumniIDs(ids []primitive.ObjectID) ([]model.PekerjaanAlumni, error) {
	var result []model.PekerjaanAlumni
	for _, p := range m.Data {
		for _, id := range ids {
			if p.AlumniID == id && p.DeletedAt != nil {
				result = append(result, p)
			}
		}
	}
	return result, nil
}

func (m *MockPekerjaanRepository) Restore(id primitive.ObjectID) error {
	p, ok := m.Data[id.Hex()]
	if !ok {
		return errors.New("not found")
	}
	p.DeletedAt = nil
	m.Data[id.Hex()] = p
	return nil
}

func (m *MockPekerjaanRepository) HardDelete(id primitive.ObjectID) error {
	_, ok := m.Data[id.Hex()]
	if !ok {
		return errors.New("not found")
	}
	delete(m.Data, id.Hex())
	return nil
}
