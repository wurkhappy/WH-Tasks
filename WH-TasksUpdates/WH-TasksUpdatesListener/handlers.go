package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
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

func UpdateSubTasks(params map[string]interface{}, body []byte) ([]byte, error, int) {
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

type PaymentItem struct {
	TaskID    string `json:"taskID"`
	SubTaskID string `json:"subtaskID"`
}

type PaymentItems []*PaymentItem

func (p PaymentItems) GetByTaskID(taskID string) *PaymentItem {
	for _, paymentItem := range p {
		if paymentItem.TaskID == taskID {
			return paymentItem
		}
	}
	return nil
}

func CheckPayment(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var message struct {
		VersionID    string       `json:"versionID"`
		PaymentItems PaymentItems `json:"paymentItems"`
	}
	json.Unmarshal(body, &message)

	key := "WHTU_" + message.VersionID
	c := redisPool.Get()
	vals, err := redis.Values(c.Do("HGETALL", key))
	if err != nil {
		return nil, fmt.Errorf("%s", "There was an error finding that token"), 401
	}

	for i, val := range vals {
		if (i+1)%2 == 0 {
			var t *models.Task
			json.Unmarshal(val.([]byte), &t)
			var idToFetch string = t.ParentID
			if t.ParentID == "" {
				idToFetch = t.ID
			}
			paymentItem := message.PaymentItems.GetByTaskID(idToFetch)
			if paymentItem == nil {
				continue
			}
			if paymentItem.SubTaskID == t.ID || paymentItem.TaskID == t.ID {
				if _, err := c.Do("HDEL", key, t.ID); err != nil {
					log.Panic(err)
				}
			}
			if paymentItem.SubTaskID == "" && paymentItem.TaskID == t.ParentID {
				if _, err := c.Do("HDEL", key, t.ID); err != nil {
					log.Panic(err)
				}
			}
		}

	}

	return nil, nil, 200
}
