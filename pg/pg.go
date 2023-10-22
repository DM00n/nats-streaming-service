package pg

import (
	"database/sql"
)

type DB struct {
	db *sql.DB
}

func NewDB(psqlconn string) (*DB, error) {
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (database DB) SelectAll() (*sql.Rows, error) {
	rows, err := database.db.Query(`SELECT * FROM "orders"`)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (database DB) InsertValue(orderId string, order []byte) error {
	sqlStatement := `INSERT INTO orders (id, data) VALUES ($1, $2)`
	_, err := database.db.Exec(sqlStatement, orderId, order)
	return err
}
func (database DB) UpdateValue(orderId string, order []byte) error {
	sqlStatement := `UPDATE orders SET data=$2 WHERE id=$1`
	_, err := database.db.Exec(sqlStatement, orderId, order)
	return err
}
