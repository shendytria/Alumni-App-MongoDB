package model

import "time"

type Alumni struct {
	ID         int        `json:"id"`
	NIM        string     `json:"nim"`
	Nama       string     `json:"nama"`
	Jurusan    string     `json:"jurusan"`
	Angkatan   int        `json:"angkatan"`
	TahunLulus int        `json:"tahun_lulus"`
	Email      string     `json:"email"`
	NoTelepon  string     `json:"no_telepon"`
	Alamat     string     `json:"alamat"`
	UserID     int        `json:"user_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`	
}

type CreateAlumniRequest struct {
	NIM, Nama, Jurusan, Email, NoTelepon, Alamat string
	Angkatan, TahunLulus, UserID                 int
}

type UpdateAlumniRequest struct {
	Nama, Jurusan, Email, NoTelepon, Alamat string
	Angkatan, TahunLulus                   int
}
