package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	utils "sqlview/backend/utils"

	"database/sql"

	_ "modernc.org/sqlite"
)

type SVServer struct {
	WorkingDir *string
}

type LocalDatabase struct {
	DB *sql.DB
}

type ResponseObject struct {
	Fields []string `json:"fields"`
}

var localServer = SVServer{}
var logger = utils.Logger{
	ShowTime:  true,
	TypeOfLog: 0,
}
var batchLogger = utils.Logger{
	ShowTime:  false,
	TypeOfLog: 1,
}

// Defaults
func Start() {
	defaultDir, err := os.Getwd()
	fmt.Print("defaultDir: " + defaultDir)
	if err != nil {
		logger.Log("error getting defaultDir", err)
		panic(err)
	}

	workDir := filepath.Join(defaultDir, "backend")
	m, err := path.Match("'*'", workDir)
	if err != nil {
		logger.Log("error checking if workDir matches the shell pattern", err)
	}
	if !m {
		logger.Log("workDir does not match the shell pattern", nil)
	}

	localServer.WorkingDir = &workDir

	if localServer.WorkingDir == nil {
		logger.Log("localServer Working Dir is empty", nil)
		panic("exit 0")
	}

	logger.Log("Server working directory: "+*localServer.WorkingDir, nil)
}

// Initializes a reference to the SQLite database, useful for communicating with the frontend
func PullDatabase(filepath string) (*LocalDatabase, error) {
	dbfile := filepath

	logger.Log("trying to connect to database: "+dbfile, nil)
	db, err := sql.Open("sqlite", dbfile)
	if err != nil {
		logger.Log("failed to connect to database: "+dbfile, err)
	}

	logger.Log("succesfully created local database", nil)
	return &LocalDatabase{
		DB: db,
	}, nil
}

// Query engine, receives a plain string (Eg. "Select * from foo where bar = x")
// and returns whatever it finds
func (db *LocalDatabase) SVQuery(query string) string {
	rows, err := db.DB.Query(query)
	if err != nil {
		logger.Log(query, nil)
		logger.Log("query error", err)
		retErr, err := json.Marshal("Incorrect Query")
		if err != nil {
			return string(retErr)
		}
		return string(retErr)
	}

	cols, err := rows.Columns()
	if err != nil {
		logger.Log("error scanning number of columns in row", err)
	}

	values := make([]interface{}, len(cols))
	valuePointers := make([]interface{}, len(cols))
	valuesArray := []string{}

	for rows.Next() {
		for i := range cols {
			valuePointers[i] = &values[i] // Creates two identical arrays,
			// with one of them storing pointers and the other storing values

		}
		defer rows.Close()
		err = rows.Scan(valuePointers...)
		if err != nil {
			logger.Log("", err)
		}

		for i := range cols { // Magic (empty interfaces)
			val := values[i] //https://stackoverflow.com/questions/17845619/how-to-call-the-scan-variadic-function-using-reflection
			b, ok := val.([]byte)
			var v interface{}
			if ok { // no idea how this works
				v = string(b)
			} else {
				v = val
			}
			valuesArray = append(valuesArray, v.(string)) // valuesArray will later be converted to JSON and returned
			batchLogger.Log("scanned row "+v.(string), nil)
		}
	}

	if err != nil {
		return err.Error()
	}

	separatedValuesArray := SeparateObjects(valuesArray, len(cols))
	v, err := json.Marshal(separatedValuesArray) // Convert to JSON so it can be interpreted directly by the frontend
	if err != nil {
		logger.Log("", err)
		return ""
	}
	return string(v)
}

// Since the values obtained from the query/reflection come in a single array containing everything, we have to split them
// this function takes the amount of columns and strings inside of the valuesArray, and separates them.
func SeparateObjects(arr_to_split []string, cols int) []ResponseObject {
	valArr := []ResponseObject{}
	obj := ResponseObject{}
	for i := range arr_to_split {
		obj.Fields = append(obj.Fields, arr_to_split[i])
		if (i+1)%cols == 0 {
			valArr = append(valArr, obj)
			obj.Fields = []string{}
		}
	}
	return valArr
}
