package tests

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	utils "sqlview/backend/utils"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

var logger = utils.Logger{
	ShowTime:  true,
	TypeOfLog: 2,
}

type TestDatabase struct {
	DB *sql.DB
}

func TestDB(t *testing.T) {
	fpath, err := os.Getwd()
	fmt.Print(fpath)
	if err != nil {
		logger.Log("error on test function TestDB", err)
	}

	fpath = filepath.Join(fpath, "debug.db")
	db, err := PullDatabase(fpath)
	if err != nil {
		logger.Log("error on test function TestDB - PullDatabase", err)
	}

	assert := assert.New(t)
	var match_string = "NOM"
	test_string := db.SVQuery("SELECT CODSEN FROM NCMTAB WHERE CODNCM = 1 AND CODNCM = 2")
	assert.Equal(match_string, test_string, "the two words are not the same")
}

func PullDatabase(filepath string) (*TestDatabase, error) {
	dbfile := filepath

	logger.Log("trying to connect to database: "+dbfile, nil)
	db, err := sql.Open("sqlite", dbfile)
	if err != nil {
		logger.Log("failed to connect to database: "+dbfile, err)
	}

	logger.Log("succesfully created local database", nil)
	return &TestDatabase{
		DB: db,
	}, nil
}

func (db *TestDatabase) SVQuery(query string) string {
	rows, err := db.DB.Query(query)
	if err != nil {
		logger.Log("", err)
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var v string
		var b string

		err := rows.Scan(&v, &b)
		if err != nil {
			logger.Log("", err)
			panic(err)
		}

		logger.Log("v is equal to: "+v+" b is equal to: "+b, nil)
	}
	return ""
}
