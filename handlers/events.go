package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-Tasks/models"
	"log"
	"net/http"
)

type Event struct {
	Name string
	Body []byte
}

type Events []*Event

func (e Events) Publish() {
	ch := getChannel()
	defer ch.Close()
	for _, event := range e {
		event.PublishOnChannel(ch)
	}
}

func (e *Event) PublishOnChannel(ch *amqp.Channel) {
	if ch == nil {
		ch = getChannel()
		defer ch.Close()
	}

	ch.ExchangeDeclare(
		config.MainExchange, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	)

	ch.Publish(
		config.MainExchange, // exchange
		e.Name,              // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        e.Body,
		})
}

func getChannel() *amqp.Channel {
	ch, err := Connection.Channel()
	if err != nil {
		dialRMQ()
		ch, err = Connection.Channel()
		if err != nil {
			log.Print(err.Error())
		}
	}

	return ch
}

type PaymentItem struct {
	TaskID    string `json:"taskID"`
	SubTaskID string `json:"subtaskID"`
}

func PaymentAccepted(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var data struct {
		VersionID    string         `json:"versionID"`
		UserID       string         `json:"userID"`
		PaymentItems []*PaymentItem `json:"paymentItems"`
	}
	json.Unmarshal(body, &data)
	tasks, err := models.FindTasksByVersionID(data.VersionID)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding tasks"), http.StatusBadRequest
	}
	updateTasks := []*models.Task{}
	for _, paymentItem := range data.PaymentItems {
		task := tasks.GetByID(paymentItem.TaskID)
		if paymentItem.SubTaskID == "" {
			task.LastAction = models.PaidActionForUser(data.UserID)
			for _, subTask := range task.SubTasks {
				subTask.LastAction = models.PaidActionForUser(data.UserID)
			}
		} else {
			subTask := task.SubTasks.GetByID(paymentItem.SubTaskID)
			subTask.LastAction = models.PaidActionForUser(data.UserID)
		}
		updateTasks = append(updateTasks, task)
	}
	for _, task := range updateTasks {
		task.Upsert()
	}

	return nil, nil, 200
}
