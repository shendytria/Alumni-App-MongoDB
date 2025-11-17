package repository

import (
	"alumni-app/app/mongodb/model"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ==================== MOCK REPOSITORY ====================

type MockUserRepository struct {
	users map[string]model.User // pakai string key = hex ID
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]model.User),
	}
}

func (m *MockUserRepository) Create(user *model.User) error {
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}

	m.users[user.ID.Hex()] = *user
	return nil
}

func (m *MockUserRepository) GetByUsername(username string) (model.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return model.User{}, errors.New("user not found")
}

func (m *MockUserRepository) GetByID(id interface{}) (model.User, error) {
	strID := id.(primitive.ObjectID).Hex()
	user, exists := m.users[strID]
	if !exists {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}
