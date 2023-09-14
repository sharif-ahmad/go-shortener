package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"go-shortener/models"

	repoerr "go-shortener/repositories/errors"
)

const (
	TABLE_NAME = "urls"
)

type PostgresURLRepository struct {
	db        *sql.DB
	tableName string
}

func NewPostgresURLRepository(connStr string) (*PostgresURLRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db = sqldblogger.OpenDriver(connStr, db.Driver(), zerologadapter.New(zerolog.New(os.Stdout)))

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresURLRepository{db, TABLE_NAME}, nil
}

func (r *PostgresURLRepository) FindAll(filters map[string]any) ([]*models.URL, error) {
	where, values := r.whereClause(0, filters)
	queryStr := `SELECT * FROM ` + r.tableName + " " + where

	rows, err := r.db.Query(queryStr, values...)
	if err != nil {
		return nil, repoerr.InternalError{Msg: err.Error()}
	}
	defer rows.Close()

	var records []*models.URL
	for rows.Next() {
		var url models.URL
		if err := rows.Scan(&url.ID, &url.Original, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt); err != nil {
			return nil, repoerr.InternalError{Msg: fmt.Sprintf("FindAll: %v", err)}
		}

		records = append(records, &url)
	}

	return nil, nil
}

func (r *PostgresURLRepository) Find(id int) (*models.URL, error) {
	queryStr := `SELECT * FROM ` + r.tableName + ` WHERE id = ?`

	var url models.URL
	row := r.db.QueryRow(queryStr, id)

	if err := row.Scan(&url.ID, &url.Original, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt); err != nil {
		return nil, repoerr.InternalError{Msg: fmt.Sprintf("Find %d: %v", id, err)}
	}

	return &url, nil
}

func (r *PostgresURLRepository) FindBy(params map[string]any) (*models.URL, error) {
	where, values := r.whereClause(0, params)
	queryStr := `SELECT * FROM ` + r.tableName + " " + where

	var url models.URL
	row := r.db.QueryRow(queryStr, values...)

	if err := row.Scan(&url.ID, &url.Original, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, repoerr.InternalError{Msg: fmt.Sprintf("FindBy: %v", err)}
	}

	return &url, nil
}

func (r *PostgresURLRepository) Create(url *models.URL) (*models.URL, error) {
	now := time.Now()
	url.CreatedAt, url.UpdatedAt = now, now

	queryStr := `INSERT INTO ` + r.tableName +
		`(original, short_code, created_at, updated_at) VALUES($1, $2, $3, $4) RETURNING id`
	row := r.db.QueryRow(queryStr, url.Original, url.ShortCode, url.CreatedAt, url.UpdatedAt)
	err := row.Scan(&url.ID)
	if err != nil {
		return nil, repoerr.InternalError{Msg: err.Error()}
	}

	return url, nil
}

func (r *PostgresURLRepository) Save(url *models.URL) error {
	_, err := r.Update(url.ID, map[string]any{
		"original":   url.Original,
		"short_code": url.ShortCode,
		"created_at": url.CreatedAt,
		"updated_at": url.UpdatedAt,
	})

	return err
}

func (r *PostgresURLRepository) Update(id int, fields map[string]any) (*models.URL, error) {
	return r.UpdateBy(fields, map[string]any{"id": id})
}

func (r *PostgresURLRepository) UpdateBy(fields map[string]any, filters map[string]any) (*models.URL, error) {

	setClause, setValues := r.setClause(0, fields)
	whereClause, whereValues := r.whereClause(len(setValues), filters)

	queryStr := `UPDATE ` + r.tableName + " " + setClause + " " + whereClause + `RETURNING *`

	row := r.db.QueryRow(queryStr, append(setValues, whereValues...))

	var url models.URL
	err := row.Scan(&url.ID, &url.Original, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		return nil, repoerr.InternalError{Msg: err.Error()}
	}

	return &url, nil
}

func (r *PostgresURLRepository) clause(clauseName string, counter int, params map[string]any) (string, []any) {
	if len(params) == 0 {
		return "", []any{}
	}

	clause := strings.Builder{}
	var values []any

	clause.WriteString(clauseName)
	for k, v := range params {
		counter++
		clause.WriteString(" ")
		clause.WriteString(k)
		clause.WriteString(fmt.Sprintf(" = $%d", counter))

		values = append(values, v)
	}

	return clause.String(), values
}

func (r *PostgresURLRepository) setClause(counter int, params map[string]any) (string, []any) {
	return r.clause("SET", counter, params)
}

func (r *PostgresURLRepository) whereClause(counter int, params map[string]any) (string, []any) {
	return r.clause("WHERE", counter, params)
}
