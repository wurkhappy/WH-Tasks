package main

import (
	"encoding/json"
	"github.com/wurkhappy/WH-Tasks/WH-TasksUpdates/models"
	"log"
)

func UpdateTask(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var task *models.Task
	json.Unmarshal(body, &task)

	key := "WHTU_" + task.VersionID

	c := redisPool.Get()
	if _, err := c.Do("HMSET", key, task.ID, body); err != nil {
		log.Panic(err)
	}

	if _, err := c.Do("EXPIRE", key, 60*60*24); err != nil {
		log.Panic(err)
	}
	return nil, nil, 200
}

func UpdateSubTask(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var message struct {
		VersionID string         `json:"versionID"`
		TaskID    string         `json:"taskID"`
		SubTasks  []*models.Task `json:"subTasks"`
	}
	json.Unmarshal(body, &message)

	key := "WHTU_" + message.VersionID
	for _, subTask := range message.SubTasks {
		subTask.ParentID = message.TaskID
		subTask.VersionID = message.VersionID
		c := redisPool.Get()
		t, _ := json.Marshal(subTask)
		if _, err := c.Do("HMSET", key, subTask.ID, t); err != nil {
			log.Panic(err)
		}
		if _, err := c.Do("EXPIRE", key, 60*60*24); err != nil {
			log.Panic(err)
		}
	}

	return nil, nil, 200
}
