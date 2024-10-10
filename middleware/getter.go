package getter

import (
	"fmt"
	server "sqlview/backend/src"
)

//type Converter struct {

//}

var LocalDatabase = server.LocalDatabase{}

func GetDatabase(filepath string) {
	db, err := server.PullDatabase(filepath)
	if err != nil {
		panic(err)
	}
	LocalDatabase = *db
}

func ProcessQuery(query string) (string, error) {
	if LocalDatabase.DB == nil {
		return "", nil
	}
	ret := LocalDatabase.SVQuery(query)
	fmt.Println(string(ret))
	return string(ret), nil
}
