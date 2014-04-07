package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Tasks/models"
	"net/http"
)

func UpdateAction(params map[string]interface{}, body []byte, userID string) ([]byte, error, int) {
	id := params["id"].(string)
	var task *models.Task
	task, err := models.FindTaskByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding task"), http.StatusBadRequest
	}

	var action *models.Action
	json.Unmarshal(body, &action)

	task.LastAction = action

	task.Update()

	a, _ := json.Marshal(action)
	return a, nil, http.StatusOK
}
