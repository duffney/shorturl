package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Url struct {
	Id         int64     `json:"id"`
	Long_url   string    `json:"long_url"`
	Short_url  string    `json:"short_url"`
	Visits     int       `json:"visits"`
	Created_at time.Time `json:"-"`
}

// Define a MovieModel struct type which wraps a sql.DB connection. TODO: Why?
type UrlModel struct {
	DB *sql.DB
}

func (u UrlModel) GetById(id int64) (*Url, error) {
	query := `
		SELECT id, long_url, short_url, visits, created_at
		FROM urls
		WHERE id = $1`

	var url Url

	err := u.DB.QueryRow(query, id).Scan(
		&url.Id,
		&url.Long_url,
		&url.Short_url,
		&url.Visits,
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
		SELECT id, long_url, short_url, visits, created_at
		FROM urls
		WHERE long_url = $1`

	var url Url

	err := u.DB.QueryRow(query, longurl).Scan(
		&url.Id,
		&url.Long_url,
		&url.Short_url,
		&url.Visits,
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
	case string: //
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
		RETURNING visits, created_at`

	args := []any{url.Id, url.Long_url, url.Short_url}

	return u.DB.QueryRow(query, args...).Scan(&url.Visits, &url.Created_at)
}

func (u UrlModel) GetAll(long_url, short_url, sort, direction string, pager Pager) ([]*Url, Metadata, error) {
	// query := `
	// 	SELECT id, long_url, short_url, visits, created_at
	// 	FROM urls
	// 	ORDER BY created_at DESC`

	// Exact match
	// query := `
	// 	SELECT id, long_url, short_url, visits, created_at
	// 	FROM urls
	// 	WHERE (LOWER(long_url) = LOWER($1) OR $1 = '')
	// 	AND (LOWER(short_url) = LOWER($2) OR $2 = '')
	// 	ORDER BY created_at DESC`

	// partical match text search
	// query := `
	// SELECT id, long_url, short_url, visits, created_at
	// FROM urls
	// WHERE long_url ILIKE '%' || $1 || '%' OR $1 = ''
	// AND short_url ILIKE '%' || $2 || '%' OR $2 = ''
	// ORDER BY created_at DESC`

	// sort
	// query := fmt.Sprintf(`
	// SELECT id, long_url, short_url, visits, created_at
	// FROM urls
	// WHERE long_url ILIKE '%%' || $1 || '%%' OR $1 = ''
	// AND short_url ILIKE '%%' || $2 || '%%' OR $2 = ''
	// ORDER BY %s %s, id ASC
	// LIMIT %d OFFSET %d`, sort, direction, limit, offset)

	// pagination metadata
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, long_url, short_url, visits, created_at
	FROM urls
	WHERE long_url ILIKE '%%' || $1 || '%%' OR $1 = ''
	AND short_url ILIKE '%%' || $2 || '%%' OR $2 = ''
	ORDER BY %s %s, id ASC
	LIMIT %d OFFSET %d`, sort, direction, pager.limit(), pager.offset())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := u.DB.QueryContext(ctx, query, long_url, short_url)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	urls := []*Url{}

	for rows.Next() {
		var url Url

		err = rows.Scan(
			&totalRecords,
			&url.Id,
			&url.Long_url,
			&url.Short_url,
			&url.Visits,
			&url.Created_at,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		urls = append(urls, &url)

	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, pager.Page, pager.PageSize)

	return urls, metadata, nil
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

func (u UrlModel) IncrementVisits(int int64) {
	query := `
		UPDATE urls
    	SET visits = visits + 1
   		WHERE id = $1;`

	u.DB.Exec(query, int)
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
