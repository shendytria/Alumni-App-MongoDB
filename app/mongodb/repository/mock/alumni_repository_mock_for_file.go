package repository

import (
	"alumni-app/app/mongodb/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAlumniRepoForFile struct {
	data []model.Alumni
}

func NewMockAlumniRepoForFile() *MockAlumniRepoForFile {
	return &MockAlumniRepoForFile{
		data: []model.Alumni{},
	}
}

func (m *MockAlumniRepoForFile) GetAllByUserID(userID primitive.ObjectID) ([]model.Alumni, error) {
	return m.data, nil
}
