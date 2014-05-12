package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/wurkhappy/WH-Tasks/models"
	"net/http"
)

func CreateTasksByVersionID(params map[string]interface{}, body []byte) ([]byte, error, int) {
	versionID := params["id"].(string)
	var tasks []*models.Task

	err := json.Unmarshal(body, &tasks)
	if err != nil {
		return nil, fmt.Errorf("%s", "Wrong value types"), http.StatusBadRequest
	}

	for _, task := range tasks {
		task.VersionID = versionID

		//TODO: this should really be a transaction
		//because if one save goes bad and others have already been saved then it could
		//lead to weird zombie tasks
		//Actually I should build a dynamic query with this like I do with model.TasksForIDs
		err = task.Upsert()
		if err != nil {
			return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
		}
	}

	a, _ := json.Marshal(tasks)

	events := Events{&Event{"tasks.created", a}}
	go events.Publish()

	return a, nil, http.StatusOK

}

func GetTasksByVersionID(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)
	tasks, err := models.FindTasksByVersionID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding tasks"), http.StatusBadRequest
	}

	p, _ := json.Marshal(tasks)
	return p, nil, http.StatusOK
}

func GetTasks(params map[string]interface{}, body []byte) ([]byte, error, int) {
	var err error
	var tasks models.Tasks
	if _, ok := params["versionID"]; ok {
		ids := params["versionID"].([]string)
		tasks, err = models.TasksForIDs(ids)
		if err != nil {
			return nil, fmt.Errorf("Error getting tasks %s", err), http.StatusBadRequest
		}
	}

	p, _ := json.Marshal(tasks)
	return p, nil, http.StatusOK
}

func UpdateTask(params map[string]interface{}, body []byte) ([]byte, error, int) {
	id := params["id"].(string)

	var updatedTask *models.Task
	json.Unmarshal(body, &updatedTask)

	task, err := models.FindTaskByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", "Error finding task"), http.StatusBadRequest
	}

	var updatedSubTasks models.Tasks
	for _, newSubTask := range updatedTask.SubTasks {
		oldSubTask := task.SubTasks.GetByID(newSubTask.ID)
		if (oldSubTask.LastAction != nil && newSubTask.LastAction == nil) ||
			(oldSubTask.LastAction == nil && newSubTask.LastAction != nil) ||
			(oldSubTask.LastAction != nil && newSubTask.LastAction != nil && oldSubTask.LastAction.Name != newSubTask.LastAction.Name) {
			updatedSubTasks = append(updatedSubTasks, newSubTask)
		}
	}

	task.SubTasks = updatedTask.SubTasks

	var subTasksComplete bool = true
	for _, subTask := range task.SubTasks {
		if subTask.LastAction == nil || (subTask.LastAction != nil && subTask.LastAction.Name != models.ActionCompleted) {
			subTasksComplete = false
		}
	}
	if subTasksComplete {
		task.LastAction = models.CompletedActionForUser(params["userID"].(string))
	}

	err = task.Update()
	if err != nil {
		return nil, fmt.Errorf("%s %s", "Error saving: ", err.Error()), http.StatusBadRequest
	}

	m := map[string]interface{}{
		"versionID": task.VersionID,
		"taskID":    task.ID,
		"subTasks":  updatedSubTasks,
	}
	jsonSubTasks, _ := json.Marshal(m)
	events := Events{&Event{"task.subTasks.updated", jsonSubTasks}}
	go events.Publish()

	jsonString, _ := json.Marshal(task)
	return jsonString, nil, http.StatusOK
}
