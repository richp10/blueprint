// Package note provides access to the note table in the MySQL database.
package note

import (
	"database/sql"
	"fmt"

	"github.com/gocraft/dbr"
)

var (
	// table is the table name.
	table = "note"
)

// Item defines the model.
type Item struct {
	ID        uint32       `db:"id"`
	Name      string       `db:"name"`
	UserID    uint32       `db:"user_id"`
	CreatedAt dbr.NullTime `db:"created_at"`
	UpdatedAt dbr.NullTime `db:"updated_at"`
	DeletedAt dbr.NullTime `db:"deleted_at"`
}

// ByID gets an item by ID.
func ByID(db *dbr.Session, ID string, userID string) (Item, bool, error) {
	result := Item{}
	err := db.
		Select("*").
		From(table).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", ID, userID).
		Limit(1).
		LoadStruct(&result)
	return result, err == dbr.ErrNotFound, err
}

// ByUserID gets all items for a user.
func ByUserID(db *dbr.Session, userID string) ([]Item, bool, error) {
	var result []Item
	_, err := db.
		Select("*").
		From(table).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		LoadStructs(&result)
	return result, err == dbr.ErrNotFound, err
}

// ByUserIDPaginate gets items for a user based on page and max variables.
func ByUserIDPaginate(db *dbr.Session, userID string, max int, page int) ([]Item, bool, error) {
	var result []Item
	_, err := db.
		Select("*").
		From(table).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Limit(uint64(max)).
		Offset(uint64(page)).
		LoadStructs(&result)
	return result, err == dbr.ErrNotFound, err
}

// ByUserIDCount counts the number of items for a user.
func ByUserIDCount(db *dbr.Session, userID string) (int, error) {
	var result int
	err := db.
		Select("count(*)").
		From(table).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		LoadValue(&result)
	return result, err
}

// Create adds an item.
func Create(db *dbr.Session, name string, userID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(name, user_id)
		VALUES
		(?,?)
		`, table),
		name, userID)
	return result, err
}

// Update makes changes to an existing item.
func Update(db *dbr.Session, name string, ID string, userID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET name = ?
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		name, ID, userID)
	return result, err
}

// DeleteHard removes an item.
func DeleteHard(db *dbr.Session, ID string, userID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		DELETE FROM %v
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		`, table),
		ID, userID)
	return result, err
}

// DeleteSoft marks an item as removed.
func DeleteSoft(db *dbr.Session, ID string, userID string) (sql.Result, error) {
	result, err := db.Exec(fmt.Sprintf(`
		UPDATE %v
		SET deleted_at = NOW()
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ID, userID)
	return result, err
}
