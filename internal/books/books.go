package books

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	lg "github.com/sirupsen/logrus"
)

// DataBase info
type ParamDB struct {
	Base *sql.DB
	Log  lg.FieldLogger
}

// Data for configuration of DB connecttion
type Config struct {
	User string
	Pass string
	Db   string
	Host string
	Port string
}

// Book info
type Book struct {
	Id     int64
	Author string
	Title  string
}

// Create connection to db
func ConnectToDB(conf Config, log lg.FieldLogger) (*ParamDB, error) {
	var err error
	var par ParamDB
	par.Base = nil
	par.Log = log

	par.Log.Println("Start connection to database.")

	if conf.User == "" || conf.Pass == "" || conf.Db == "" ||
		conf.Host == "" || conf.Port == "" {
		err = errors.New("Bad connection parameters")
		return &par, err
	}

	// Connect to user database
	connstr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		conf.User, conf.Pass, conf.Host, conf.Port, conf.Db)
	db, err := openDB(connstr)
	if err != nil {
		err = errors.New("Bad connection to user db")
		return &par, err
	}
	par.Base = db

	par.Log.Println("Set connection to database.")
	return &par, nil
}

// Close user database
func (par *ParamDB) Close() (err error) {
	if par.Base == nil {
		return nil
	}
	par.Log.Printf("Close database: %v.", par)

	if err = par.Base.Close(); err != nil {
		return err
	}

	par.Base = nil
	return nil
}

// Check connect to database
func openDB(confstr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", confstr)

	if err != nil {
		err = errors.New("Bad connection to db")
		return nil, err
	}

	// Ping the database
	if err = db.Ping(); err != nil {
		err = errors.New("Couldn't ping to database")
		db.Close()
		return nil, err
	}
	return db, nil
}

// Start transaction in database
func (par *ParamDB) StartTX(event string) (*sql.Tx, error) {
	if par.Base == nil {
		par.Log.Errorf("DB not opened (%s).", event)
		return nil, errors.New("DB not opened")
	}

	tx, err := par.Base.Begin()
	if err != nil {
		par.Log.Errorf("Error start transaction (%s) (%s).", event, err.Error())
		return nil, err
	}

	par.Log.Printf("Start transaction (%s).", event)
	return tx, nil
}

// Stop transaction in database
func (par *ParamDB) StopTX(tx *sql.Tx, result *bool, event string) {
	if *result {
		par.Log.Printf("Commit transaction (%s).", event)
		if err := tx.Commit(); err != nil {
			par.Log.Errorf("Error commit transaction (%s) (%s).", event, err.Error())
		}
	} else {
		par.Log.Printf("Rollback transaction (%s).", event)
		if err := tx.Rollback(); err != nil {
			par.Log.Errorf("Error rollback transaction (%s) (%s).", event, err.Error())
		}
	}
}

// Insert book info into database
func (b *Book) InsertBook(par *ParamDB) (int64, error) {
	tx, err := par.StartTX("Insert Book")
	if err != nil {
		return 0, err
	}
	txEnd := false
	defer par.StopTX(tx, &txEnd, "Insert Book")

	query := fmt.Sprintf("INSERT INTO BookInfo (title, author) VALUES('%s', '%s') RETURNING id;", b.Title, b.Author)
	row := tx.QueryRow(query)

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, errors.New("Error insert book: - bad id")
	}

	txEnd = true
	return id, nil
}

// Delete book info from database
func (b *Book) DeleteBook(par *ParamDB) error {
	tx, err := par.StartTX("Delete Book")
	if err != nil {
		return err
	}
	txEnd := false
	defer par.StopTX(tx, &txEnd, "Delete Book")

	if b.Id <= 0 {
		return errors.New("Id not set")
	}

	query := fmt.Sprintf("DELETE FROM BookInfo where id = %v;", b.Id)
	_, err = tx.Exec(query)

	if err != nil {
		return errors.New("Error delete book")
	}

	txEnd = true
	return nil
}

// Update book info into database
func (b *Book) UpdateBook(par *ParamDB) error {
	tx, err := par.StartTX("Update Book")
	if err != nil {
		return err
	}
	txEnd := false
	defer par.StopTX(tx, &txEnd, "Update Book")

	if b.Id <= 0 {
		return errors.New("Id not set")
	}

	var query string
	switch {
	case b.Author != "" && b.Title != "":
		query = fmt.Sprintf("UPDATE BookInfo SET author = '%s', title = '%s' WHERE id = %v;", b.Author, b.Title, b.Id)
	case b.Author != "":
		query = fmt.Sprintf("UPDATE BookInfo SET author = '%s' WHERE id = %v;", b.Author, b.Id)
	case b.Title != "":
		query = fmt.Sprintf("UPDATE BookInfo SET title = '%s' WHERE id = %v;", b.Title, b.Id)
	default:
		return nil
	}
	fmt.Println(query)
	_, err = tx.Exec(query)

	if err != nil {
		return errors.New("Error update book")
	}

	txEnd = true
	return nil
}

// Select book info from database
func (b *Book) SelectBook(par *ParamDB) (*[]Book, error) {
	bb := make([]Book, 0)
	tx, err := par.StartTX("Select Book")
	if err != nil {
		return &bb, err
	}
	txEnd := false
	defer par.StopTX(tx, &txEnd, "Select Book")

	var query string
	if b.Id > 0 {
		query = fmt.Sprintf("SELECT id, title, author FROM BookInfo WHERE id = %v ORDER BY id;", b.Id)
	} else {
		query = fmt.Sprintf("SELECT id, title, author FROM BookInfo ORDER BY id;")
	}

	rows, err := tx.Query(query)
	if err != nil {
		return &bb, errors.New("Select error")
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Id, &book.Title, &book.Author); err != nil {
			return &bb, errors.New("Select scan error")
		}
		bb = append(bb, book)
	}

	if err != nil {
		return &bb, errors.New("Error select book")
	}

	txEnd = true
	return &bb, nil
}
