package stats

import (
	"database/sql"
	"errors"

	sqlite3 "github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate		= errors.New("record already exists")
	ErrNotExists		= errors.New("row does not exist")
	ErrUpdateFailed		= errors.New("update has failed")
	ErrDeleteFailed		= errors.New("delete has failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS stats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		active INTEGER NOT NULL,
		total INTEGER NOT NULL
	);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(stats Stats) (*Stats, error) {
	res, err := r.db.Exec("INSERT INTO stats (active, total) values (?, ?)", stats.Active, stats.Total)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	stats.ID = id
	return &stats, nil
}

func (r *SQLiteRepository) All() ([]Stats, error) {
	rows, err := r.db.Query("SELECT * FROM stats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Stats
	for rows.Next() {
		var stats Stats
		if err := rows.Scan(&stats.ID, &stats.Active, &stats.Total); err != nil {
			return nil, err
		}
		all = append(all, stats)
	}
	return all, err
}

func (r *SQLiteRepository) GetById(id int64) (*Stats, error) {
	row := r.db.QueryRow("SELECT * FROM stats where id = ?", id)
	var stats Stats
	if err := row.Scan(&stats.ID, &stats.Active, &stats.Total); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &stats, nil
}

func (r *SQLiteRepository) GetFirst() (*Stats, error) {
	row := r.db.QueryRow("SELECT * FROM stats limit 1")
	var stats Stats
	if err := row.Scan(&stats.ID, &stats.Active, &stats.Total); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &stats, nil
}