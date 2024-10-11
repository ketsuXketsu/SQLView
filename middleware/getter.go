package getter

import (
	"fmt"
	server "sqlview/backend/src"
)

var LocalDatabase = server.LocalDatabase{}

// Receives a string path and asks the server for a LocalDatabase using said path
func GetDatabase(filepath string) {
	db, err := server.PullDatabase(filepath)
	if err != nil {
		panic(err)
	}
	LocalDatabase = *db
}

// Middleman between the front and backend
func ProcessQuery(query string) (string, error) {
	if LocalDatabase.DB == nil {
		return "", nil
	}
	ret := LocalDatabase.SVQuery(query)
	fmt.Println(string(ret))
	return string(ret), nil
}
