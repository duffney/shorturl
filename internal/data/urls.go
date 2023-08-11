package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Url struct {
	Id         int64     `json:"id"`
	Long_url   string    `json:"long_url"`
	Short_url  string    `json:"short_url"`
	Created_at time.Time `json:"-"`
}

// Define a MovieModel struct type which wraps a sql.DB connection. TODO: Why?
type UrlModel struct {
	DB *sql.DB
}

func (u UrlModel) GetById(id int64) (*Url, error) {
	query := `
		SELECT id, long_url, short_url, created_at
		FROM urls
		WHERE id = $1`

	var url Url

	err := u.DB.QueryRow(query, id).Scan(
		&url.Id,
		&url.Long_url,
		&url.Short_url,
		&url.Created_at,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no record found")
		} else {
			return nil, err
		}
	}

	return &url, nil
}

func (u UrlModel) GetByLongUrl(longurl string) (*Url, error) {
	query := `
		SELECT id, long_url, short_url, created_at
		FROM urls
		WHERE long_url = $1`

	var url Url

	err := u.DB.QueryRow(query, longurl).Scan(&url.Id, &url.Long_url, &url.Short_url, &url.Created_at)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no record found")
		} else {
			return nil, err
		}
	}

	return &url, nil
}

/*
func (u UrlModel) GetByShortUrl(shortUrl string) (*Url, error) {
	query := `
        SELECT id, long_url, short_url, created_at
        FROM urls
        WHERE short_url = $1`

	var url Url

	err := u.DB.QueryRow(query, shortUrl).Scan(&url.Id, &url.Long_url, &url.Short_url, &url.Created_at)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("no record found")
		} else {
			return nil, err
		}
	}

	return &url, nil
}
*/

/*
func (u UrlModel) Get(id interface{}) (*Url, error) {
	var query string
	var args []interface{}

	switch id := id.(type) {
	case int64:
		query = `
			SELECT id, long_url, short_url, created_at
			FROM urls
			WHERE id = $1`
		args = []interface{}{id}
	case string: // #TODO: Add logic for short_url & hash
		query = `
			SELECT id, long_url, short_url, created_at
			FROM urls
			WHERE long_url = $1`
		args = []interface{}{id}
	default:
		return nil, errors.New("invalid id type")
	}

	var url Url
	err := u.DB.QueryRow(query, args...).Scan(
		&url.Id,
		&url.Long_url,
		&url.Short_url,
		&url.Created_at,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("no record found")
		default:
			return nil, err
		}
	}

	return &url, nil
}
*/

func (u UrlModel) Insert(url *Url) error {
	query := `
		INSERT INTO urls (id, long_url, short_url)	
		VALUES ($1, $2, $3)
		RETURNING created_at`

	args := []any{url.Id, url.Long_url, url.Short_url}

	return u.DB.QueryRow(query, args...).Scan(&url.Created_at)
}

func (u UrlModel) GetAll() ([]*Url, error) {
	query := `
		SELECT id, long_url, short_url, created_at
		FROM urls
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var urls []*Url

	for rows.Next() {
		var url Url

		err = rows.Scan(
			&url.Id,
			&url.Long_url,
			&url.Short_url,
			&url.Created_at,
		)

		if err != nil {
			return nil, err
		}

		urls = append(urls, &url)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (u UrlModel) LongUrlExists(url string) bool {
	query := `
		SELECT EXISTS (
			SELECT id
			FROM urls
			WHERE long_url = $1
		)`

	var exists bool

	u.DB.QueryRow(query, url).Scan(&exists)

	return exists
}

/*
Support for mocking the database in unit tests
type MockUrlModel struct {}

func (m *MockUrlModel) Insert(url *Url) error {
	return nil
}

func (m *MockUrlModel) GetAll() ([]*Url, error) {
	return nil, nil
}

*/
