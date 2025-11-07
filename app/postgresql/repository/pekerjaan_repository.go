package repository

import (
	"alumni-app/app/postgresql/model"
	"alumni-app/database/postgresql"
	"database/sql"
	"time"
)

type PekerjaanRepository struct{}

func NewPekerjaanRepository() *PekerjaanRepository {
	return &PekerjaanRepository{}
}

func (r *PekerjaanRepository) GetAll(search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
		       deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan
		WHERE deleted_at IS NULL 
		  AND (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1)
		ORDER BY ` + sortBy + ` ` + order + `
		LIMIT $2 OFFSET $3
	`

	rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PekerjaanRepository) Count(search string) (int, error) {
	var total int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM pekerjaan
		WHERE deleted_at IS NULL
		  AND (nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1)
	`, "%"+search+"%").Scan(&total)
	return total, err
}

func (r *PekerjaanRepository) GetByID(id int) (model.PekerjaanAlumni, error) {
	var p model.PekerjaanAlumni
	err := database.DB.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
		&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
		&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return model.PekerjaanAlumni{}, err
	}
	return p, err
}

func (r *PekerjaanRepository) GetByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
	rows, err := database.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
		       deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan
		WHERE alumni_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PekerjaanRepository) Create(p *model.PekerjaanAlumni) error {
	now := time.Now()
	return database.DB.QueryRow(`
		INSERT INTO pekerjaan (
			alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
			gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
			deskripsi_pekerjaan, created_at, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$11)
		RETURNING id
	`, p.AlumniID, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja,
		p.GajiRange, p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan,
		p.DeskripsiPekerjaan, now).Scan(&p.ID)
}

func (r *PekerjaanRepository) Update(id int, p *model.PekerjaanAlumni) error {
	_, err := database.DB.Exec(`
		UPDATE pekerjaan
		SET nama_perusahaan=$1, posisi_jabatan=$2, bidang_industri=$3, lokasi_kerja=$4,
		    gaji_range=$5, tanggal_mulai_kerja=$6, tanggal_selesai_kerja=$7,
		    status_pekerjaan=$8, deskripsi_pekerjaan=$9, updated_at=$10
		WHERE id=$11 AND deleted_at IS NULL
	`, p.NamaPerusahaan, p.PosisiJabatan, p.BidangIndustri, p.LokasiKerja, p.GajiRange,
		p.TanggalMulaiKerja, p.TanggalSelesaiKerja, p.StatusPekerjaan, p.DeskripsiPekerjaan, time.Now(), id)
	return err
}

func (r *PekerjaanRepository) Delete(id, userID int, isAdmin bool) error {
	query := `
		UPDATE pekerjaan p
		SET deleted_at = NOW()
		FROM alumni a
		WHERE p.id = $1 AND p.deleted_at IS NULL
	`
	args := []interface{}{id}

	if !isAdmin {
		query += " AND p.alumni_id = a.id AND a.user_id = $2"
		args = append(args, userID)
	}

	res, err := database.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PekerjaanRepository) GetTrashed(userID int, isAdmin bool) ([]model.PekerjaanAlumni, error) {
	query := `
		SELECT p.id, p.alumni_id, p.nama_perusahaan, p.posisi_jabatan, p.bidang_industri, p.lokasi_kerja,
		       p.gaji_range, p.tanggal_mulai_kerja, p.tanggal_selesai_kerja, p.status_pekerjaan,
		       p.deskripsi_pekerjaan, p.created_at, p.updated_at
		FROM pekerjaan p
		JOIN alumni a ON a.id = p.alumni_id
		WHERE p.deleted_at IS NOT NULL
	`
	args := []interface{}{}

	if !isAdmin {
		query += " AND a.user_id = $1"
		args = append(args, userID)
	}

	query += " ORDER BY p.deleted_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *PekerjaanRepository) Restore(id, userID int, isAdmin bool) error {
	query := `
		UPDATE pekerjaan p
		SET deleted_at = NULL
		FROM alumni a
		WHERE p.id = $1 AND p.alumni_id = a.id AND p.deleted_at IS NOT NULL
	`
	args := []interface{}{id}

	if !isAdmin {
		query += " AND a.user_id = $2"
		args = append(args, userID)
	}

	res, err := database.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PekerjaanRepository) HardDelete(id, userID int, isAdmin bool) error {
	query := `
		DELETE FROM pekerjaan p
		USING alumni a
		WHERE p.id = $1 AND p.deleted_at IS NOT NULL
	`
	args := []interface{}{id}

	if !isAdmin {
		// Batasi hanya boleh menghapus pekerjaan miliknya sendiri
		query += " AND p.alumni_id = a.id AND a.user_id = $2"
		args = append(args, userID)
	}

	res, err := database.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}


func (r *PekerjaanRepository) GetByTahunLulusWithGaji(tahun int, minGaji int64) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT a.id, a.nama, a.jurusan, a.tahun_lulus,
		       p.bidang_industri, p.nama_perusahaan, p.posisi_jabatan, p.gaji_range
		FROM alumni a
		JOIN pekerjaan p ON a.id = p.alumni_id
		WHERE a.tahun_lulus = $1 AND p.gaji_range >= $2 AND p.deleted_at IS NULL
	`, tahun, minGaji)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var (
			id, tahunLulus int
			nama, jurusan, bidang, perusahaan, posisi string
			gaji int64
		)
		if err := rows.Scan(&id, &nama, &jurusan, &tahunLulus, &bidang, &perusahaan, &posisi, &gaji); err != nil {
			return nil, err
		}
		result = append(result, map[string]interface{}{
			"id": id, "nama": nama, "jurusan": jurusan, "tahun_lulus": tahunLulus,
			"bidang_industri": bidang, "nama_perusahaan": perusahaan,
			"posisi_jabatan": posisi, "gaji_range": gaji,
		})
	}
	return result, nil
}
