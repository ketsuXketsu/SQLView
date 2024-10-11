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
	Headers []string `json:"headers"`
	Fields  []string `json:"fields"`
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
		return err.Error()
	}

	values := make([]interface{}, len(cols))
	valuePointers := make([]interface{}, len(cols))
	valuesArray := make([][]interface{}, 0)

	// Reading value from .db file
	for rows.Next() {
		for i := range cols {
			valuePointers[i] = &values[i]
		}

		err = rows.Scan(valuePointers...)
		if err != nil {
			logger.Log("error scanning rows", err)
			return err.Error()
		}
		// Something about empty interfaces
		rowData := make([]interface{}, len(cols))
		for i, val := range values {
			b, ok := val.([]byte)
			if ok {
				rowData[i] = string(b)
			} else {
				rowData[i] = val
			}

			batchLogger.Log(fmt.Sprintf("scanned value: %v", rowData[i]), nil)
		}

		valuesArray = append(valuesArray, rowData)
	}

	if err = rows.Err(); err != nil {
		logger.Log("error after scanning rows", err)
		return err.Error()
	}

	// Convert the 2D slice to a slice of ResponseObject
	var responseArray []ResponseObject
	for _, row := range valuesArray {
		obj := ResponseObject{Fields: make([]string, len(row))}
		obj.Headers = cols
		for i, val := range row {
			obj.Fields[i] = fmt.Sprintf("%v", val)
		}
		responseArray = append(responseArray, obj)
	}

	v, err := json.Marshal(responseArray)
	if err != nil {
		logger.Log("error marshaling JSON", err)
		return err.Error()
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
