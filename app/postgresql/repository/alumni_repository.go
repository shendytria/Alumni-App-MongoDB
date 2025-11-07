package repository

import (
	"alumni-app/app/postgresql/model"
	"alumni-app/database/postgresql"
	"database/sql"
	"fmt"
	"time"
)

type AlumniRepository struct{}

func NewAlumniRepository() *AlumniRepository {
	return &AlumniRepository{}
}

func (r *AlumniRepository) GetAll(search, sortBy, order string, page, limit int) ([]model.Alumni, error) {
    allowedSort := map[string]bool{
        "id": true, "nama": true, "angkatan": true,
        "tahun_lulus": true, "created_at": true,
    }
    if !allowedSort[sortBy] {
        sortBy = "id" 
    }
    if order != "desc" {
        order = "asc" 
    }

    offset := (page - 1) * limit

    q := fmt.Sprintf(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, user_id, created_at, updated_at, deleted_at
        FROM alumni
        WHERE (nama ILIKE $1 OR jurusan ILIKE $1 OR email ILIKE $1) AND deleted_at IS NULL
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

    rows, err := database.DB.Query(q, "%"+search+"%", limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.Alumni
    for rows.Next() {
        var a model.Alumni
        if err := rows.Scan(
            &a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
            &a.Email, &a.NoTelepon, &a.Alamat, &a.UserID, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
        ); err != nil {
            return nil, err
        }
        list = append(list, a)
    }
    return list, nil
}

func (r *AlumniRepository) Count(search string) (int, error) {
	var total int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM alumni
		WHERE (nama ILIKE $1 OR jurusan ILIKE $1 OR email ILIKE $1)
		AND deleted_at IS NULL
	`, "%"+search+"%").Scan(&total)
	return total, err
}

func (r *AlumniRepository) GetByID(id int) (model.Alumni, error) {
	var a model.Alumni
	err := database.DB.QueryRow(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, user_id, created_at, updated_at, deleted_at
		FROM alumni
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(
		&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
		&a.Email, &a.NoTelepon, &a.Alamat, &a.UserID, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return model.Alumni{}, err
	}
	return a, err
}

func (r *AlumniRepository) Create(a *model.Alumni) error {
	now := time.Now()
	return database.DB.QueryRow(`
		INSERT INTO alumni (
			nim, nama, jurusan, angkatan, tahun_lulus, email,
			no_telepon, alamat, created_at, updated_at, user_id
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id
	`,
		a.NIM, a.Nama, a.Jurusan, a.Angkatan, a.TahunLulus,
		a.Email, a.NoTelepon, a.Alamat, now, now, a.UserID,
	).Scan(&a.ID)
}

func (r *AlumniRepository) Update(id int, a *model.Alumni) error {
	_, err := database.DB.Exec(`
		UPDATE alumni
		SET nama=$1, jurusan=$2, angkatan=$3, tahun_lulus=$4,
		    email=$5, no_telepon=$6, alamat=$7, updated_at=$8
		WHERE id=$9 AND deleted_at IS NULL
	`,
		a.Nama, a.Jurusan, a.Angkatan, a.TahunLulus,
		a.Email, a.NoTelepon, a.Alamat, time.Now(), id,
	)
	return err
}

func (r *AlumniRepository) DeleteByID(id int) error {
	res, err := database.DB.Exec(`
		UPDATE alumni SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
