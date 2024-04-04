package postrgresql

import (
	"database/sql"

	"georgie5.net/QUOTEBOX/pkg/models"
)

type QuoteModel struct {
	DB *sql.DB
}

func (m *QuoteModel) Insert(author, category, body string) (int, error) {
	var id int

	s := `
	INSERT INTO quotations(author_name, category, quote)
	VALUES ($1, $2, $3)
	RETURNING quotation_id 
	`
	err := m.DB.QueryRow(s, author, category, body).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *QuoteModel) Read() ([]*models.Quote, error) {

	//SQL statement
	s := `
		SELECT author_name, category, quote
		FROM quotations
		LIMIT 10
	`
	rows, err := m.DB.Query(s)
	if err != nil {
		return nil, err
	}

	// cleanup before we leave Read()
	defer rows.Close()

	quotes := []*models.Quote{}

	for rows.Next() {
		q := &models.Quote{}
		err = rows.Scan(&q.Author_name, &q.Category, &q.Body)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, q)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return quotes, nil

}
