package model

import "time"

type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    PasswordHash  string    `json:"-"` 
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
    Username string `json:"username" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    PasswordHash string `json:"password" validate:"required,min=6"`
    Role     string `json:"role"`
}

type LoginRequest struct {
    Username string `json:"username"`
    PasswordHash string `json:"password"`
}

type LoginResponse struct {
    User  User   `json:"user"`
    Token string `json:"token"`
}
