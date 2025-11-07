package model

import "time"

type PekerjaanRequestBase struct {
	NamaPerusahaan     string `json:"nama_perusahaan"`
	PosisiJabatan      string `json:"posisi_jabatan"`
	BidangIndustri     string `json:"bidang_industri"`
	LokasiKerja        string `json:"lokasi_kerja"`
	GajiRange          int64  `json:"gaji_range"`
	TanggalMulaiKerja  string `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja string `json:"tanggal_selesai_kerja"`
	StatusPekerjaan    string `json:"status_pekerjaan"`
	DeskripsiPekerjaan string `json:"deskripsi_pekerjaan"`
}

type CreatePekerjaanRequest struct {
	PekerjaanRequestBase
	AlumniID int `json:"alumni_id"`
}

type UpdatePekerjaanRequest struct {
	PekerjaanRequestBase
}

type PekerjaanAlumni struct {
	ID                  int        `json:"id"`
	AlumniID            int        `json:"alumni_id"`
	NamaPerusahaan      string     `json:"nama_perusahaan"`
	PosisiJabatan       string     `json:"posisi_jabatan"`
	BidangIndustri      string     `json:"bidang_industri"`
	LokasiKerja         string     `json:"lokasi_kerja"`
	GajiRange           int64      `json:"gaji_range"`
	TanggalMulaiKerja   time.Time  `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string     `json:"status_pekerjaan"`
	DeskripsiPekerjaan  string     `json:"deskripsi_pekerjaan"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

type AlumniPekerjaanResponse struct {
	ID            int    `json:"id"`
	Nama          string `json:"nama"`
	Jurusan       string `json:"jurusan"`
	TahunLulus    int    `json:"tahun_lulus"`
	BidangIndustri string `json:"bidang_industri"`
	NamaPerusahaan string `json:"nama_perusahaan"`
	PosisiJabatan  string `json:"posisi_jabatan"`
	GajiRange      int64  `json:"gaji_range"`
}
