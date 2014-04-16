package DB

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	// "log"
)

var SaveTask *sql.Stmt
var UpdateTask *sql.Stmt
var UpsertTask *sql.Stmt
var FindTasksByVersionID *sql.Stmt
var FindTaskByID *sql.Stmt

func CreateStatements() {
	var err error
	SaveTask, err = DB.Prepare("INSERT INTO task(id, data) VALUES($1, $2)")
	if err != nil {
		panic(err)
	}

	UpdateTask, err = DB.Prepare("UPDATE task SET data = $2 WHERE id = $1")
	if err != nil {
		panic(err)
	}

	UpsertTask, err = DB.Prepare("SELECT upsert_task($1, $2)")
	if err != nil {
		panic(err)
	}

	FindTasksByVersionID, err = DB.Prepare("SELECT data FROM task WHERE data->>'versionID' = $1")
	if err != nil {
		panic(err)
	}

	FindTaskByID, err = DB.Prepare("SELECT data FROM task WHERE id = $1")
	if err != nil {
		panic(err)
	}
}
