package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Tasks/models"
	"net/http"
	"time"
)

func UpdateAction(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	userID := params["userID"].(string)
	var task *models.Task
	task, err := models.FindTaskByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding task"), http.StatusBadRequest
	}

	var action *models.Action
	json.Unmarshal(body, &action)

	task.LastAction = action
	if task.LastAction != nil {
		task.LastAction.UserID = userID
		task.LastAction.Date = time.Now().UTC()
	}

	task.Update()
	if task.LastAction != nil {
		t, _ := json.Marshal(task)
		events := Events{&Event{"task." + task.LastAction.Name, t}}
		go events.Publish()
	}

	a, _ := json.Marshal(action)
	return a, nil, http.StatusOK
}
