package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PekerjaanAlumni struct {
	ID                  primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	AlumniID            primitive.ObjectID  `bson:"alumni_id" json:"alumni_id"`
	NamaPerusahaan      string              `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string              `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string              `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string              `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           int64               `bson:"gaji_range" json:"gaji_range"`
	TanggalMulaiKerja   time.Time           `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time          `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string              `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  string              `bson:"deskripsi_pekerjaan" json:"deskripsi_pekerjaan"`
	CreatedAt           time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time           `bson:"updated_at" json:"updated_at"`
	DeletedAt           *time.Time          `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type PekerjaanRequestBase struct {
	NamaPerusahaan      string `json:"nama_perusahaan"`
	PosisiJabatan       string `json:"posisi_jabatan"`
	BidangIndustri      string `json:"bidang_industri"`
	LokasiKerja         string `json:"lokasi_kerja"`
	GajiRange           int64  `json:"gaji_range"`
	TanggalMulaiKerja   string `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja string `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string `json:"status_pekerjaan"`
	DeskripsiPekerjaan  string `json:"deskripsi_pekerjaan"`
}

type CreatePekerjaanRequest struct {
	PekerjaanRequestBase
	AlumniID primitive.ObjectID `json:"alumni_id"`
}

type UpdatePekerjaanRequest struct {
	PekerjaanRequestBase
}