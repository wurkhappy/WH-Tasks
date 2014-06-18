package models

import (
	"database/sql"
	"encoding/json"
	_ "github.com/bmizerany/pq"
	"github.com/nu7hatch/gouuid"
	"github.com/wurkhappy/WH-Tasks/DB"
	"log"
	"strconv"
	"time"
)

type Task struct {
	ID           string    `json:"id"`
	VersionID    string    `json:"versionID"`
	IsPaid       bool      `json:"isPaid"`
	Hours        float64   `json:"hours"`
	SubTasks     Tasks     `json:"subTasks"`
	Title        string    `json:"title"`
	Index        int64     `json:"index"`
	DateExpected time.Time `json:"dateExpected,omitempty"`
	LastAction   *Action   `json:"lastAction"`
}

//for unmarshaling purposes
type task struct {
	ID           string    `json:"id"`
	VersionID    string    `json:"versionID"`
	IsPaid       bool      `json:"isPaid"`
	Hours        float64   `json:"hours"`
	SubTasks     Tasks     `json:"subTasks"`
	Title        string    `json:"title"`
	Index        int64     `json:"index"`
	DateExpected time.Time `json:"dateExpected"`
	LastAction   *Action   `json:"lastAction"`
}

type Tasks []*Task

func NewTask() *Task {
	id, _ := uuid.NewV4()
	return &Task{
		ID: id.String(),
	}
}

func (t *Task) Save() (err error) {
	jsonByte, _ := json.Marshal(t)
	_, err = DB.SaveTask.Exec(t.ID, string(jsonByte))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (t *Task) Upsert() (err error) {
	jsonByte, _ := json.Marshal(t)
	r, err := DB.UpsertTask.Query(t.ID, string(jsonByte))
	r.Close()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (t *Task) Update() (err error) {
	jsonByte, _ := json.Marshal(t)
	_, err = DB.UpdateTask.Exec(t.ID, string(jsonByte))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func FindTasksByVersionID(id string) (t Tasks, err error) {
	r, err := DB.FindTasksByVersionID.Query(id)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return dbRowsToTasks(r)
}

func FindTaskByID(id string) (t *Task, err error) {
	var s string
	err = DB.FindTaskByID.QueryRow(id).Scan(&s)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(s), &t)
	return t, nil
}

func dbRowsToTasks(r *sql.Rows) (tasks Tasks, err error) {
	for r.Next() {
		var s string
		err = r.Scan(&s)
		if err != nil {
			return nil, err
		}
		var t *Task
		json.Unmarshal([]byte(s), &t)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (t *Task) UnmarshalJSON(bytes []byte) (err error) {
	var tk *task
	err = json.Unmarshal(bytes, &tk)
	if err != nil {
		return err
	}

	if tk.ID == "" {
		id, _ := uuid.NewV4()
		tk.ID = id.String()
	}
	t.ID = tk.ID
	t.VersionID = tk.VersionID
	t.IsPaid = tk.IsPaid
	t.Hours = tk.Hours
	t.SubTasks = tk.SubTasks
	t.Title = tk.Title
	t.Index = tk.Index
	t.DateExpected = tk.DateExpected
	t.LastAction = tk.LastAction
	return nil
}

func (t Tasks) GetByID(id string) *Task {
	for _, task := range t {
		if task.ID == id {
			return task
		}
	}
	return nil
}

func TasksForIDs(ids []string) (t Tasks, err error) {
	var querylist string
	interfaceIDs := make([]interface{}, len(ids))
	for i, id := range ids {
		querylist += "data->>'versionID' = $" + strconv.Itoa(i+1)
		if i+1 != len(ids) {
			querylist += " OR "
		}
		interfaceIDs[i] = interface{}(id)
	}
	rows, err := DB.DB.Query("SELECT data FROM task WHERE "+querylist, interfaceIDs...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	return dbRowsToTasks(rows)
}
