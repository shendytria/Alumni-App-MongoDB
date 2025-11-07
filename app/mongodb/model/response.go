package model

// MetaInfo digunakan untuk informasi pagination dan query
type MetaInfo struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Total  int    `json:"total"`
	Pages  int    `json:"pages"`
	SortBy string `json:"sortBy"`
	Order  string `json:"order"`
	Search string `json:"search"`
}

// Response untuk data Alumni
type AlumniResponse struct {
	Data []Alumni `json:"data"`
	Meta MetaInfo `json:"meta"`
}

// Response untuk data Pekerjaan
type PekerjaanResponse struct {
	Data []PekerjaanAlumni `json:"data"`
	Meta MetaInfo          `json:"meta"`
}

// Response generik (bisa dipakai kalau butuh custom)
type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

// ErrorResponse digunakan untuk menampilkan pesan error standar
type ErrorResponse struct {
	Error string `json:"error" example:"Pesan kesalahan"`
}

// FileUploadResponse digunakan untuk response upload file (foto/sertifikat)
type FileUploadResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"File uploaded successfully"`
	Path    string `json:"path" example:"/uploads/foto/FOTO_12345_uuid.png"`
}
