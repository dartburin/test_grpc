package books

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	lg "github.com/sirupsen/logrus"
)

// DataBase info
type ParamDB struct {
	Base *gorm.DB
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
type BookRecord struct {
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

	db, err := gorm.Open("postgres", connstr)
	if err != nil {
		err = errors.New("Bad connection to user db")
		return &par, err
	}
	par.Base = db

	db.AutoMigrate(&BookRecord{})

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

// Start transaction in database
func (par *ParamDB) StartTX(event string) (*gorm.DB, error) {
	if par.Base == nil {
		par.Log.Errorf("DB not opened (%s).", event)
		return nil, errors.New("DB not opened")
	}

	tx := par.Base.Begin()
	par.Log.Printf("Start transaction (%s).", event)
	return tx, nil
}

// Stop transaction in database
func (par *ParamDB) StopTX(tx *gorm.DB, result *bool, event string) {
	if *result {
		par.Log.Printf("Commit transaction (%s).", event)
		_ = tx.Commit()
	} else {
		par.Log.Printf("Rollback transaction (%s).", event)
		_ = tx.Rollback()
	}
}

// Insert book info into database
func (b *BookRecord) InsertBook(par *ParamDB) (int64, error) {
	tx, err := par.StartTX("Insert Book")
	if err != nil {
		return 0, err
	}
	txEnd := false
	defer par.StopTX(tx, &txEnd, "Insert Book")

	err = tx.Create(b).Error
	if err != nil {
		return (*b).Id, errors.New("Error insert book")
	}

	txEnd = true
	return (*b).Id, nil
}

// Delete book info from database
func (b *BookRecord) DeleteBook(par *ParamDB) error {
	tx, err := par.StartTX("Delete Book")
	if err != nil {
		return err
	}
	txEnd := false
	defer par.StopTX(tx, &txEnd, "Delete Book")

	if b.Id <= 0 {
		return errors.New("Id not set")
	}

	err = tx.Delete(b).Error
	if err != nil {
		return errors.New("Error delete book")
	}

	txEnd = true
	return nil
}

// Update book info into database
func (b *BookRecord) UpdateBook(par *ParamDB) error {
	tx, err := par.StartTX("Update Book")
	if err != nil {
		return err
	}
	txEnd := false
	defer par.StopTX(tx, &txEnd, "Update Book")

	if b.Id <= 0 {
		return errors.New("Id not set")
	}

	err = tx.Save(b).Error
	if err != nil {
		return errors.New("Error update book")
	}

	txEnd = true
	return nil
}

// Select books info from database
func (b *BookRecord) SelectBooks(par *ParamDB) (*[]BookRecord, error) {
	bb := make([]BookRecord, 0)
	tx, err := par.StartTX("Select Books")
	if err != nil {
		return &bb, err
	}
	txEnd := false
	defer par.StopTX(tx, &txEnd, "Select Books")

	tx.Order("id").Find(&bb)
	txEnd = true
	return &bb, nil
}

// Select book info from database
func (b *BookRecord) SelectBook(par *ParamDB) (*BookRecord, error) {
	bb := &BookRecord{}
	tx, err := par.StartTX("Select Book")
	if err != nil {
		return nil, err
	}
	txEnd := false
	defer par.StopTX(tx, &txEnd, "Select Book")

	if b.Id <= 0 {
		return nil, errors.New("Id not set")
	}

	tx.Where("id = ?", b.Id).Find(&bb)
	txEnd = true
	return bb, nil
}
