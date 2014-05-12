package main

import (
	"github.com/ant0ine/go-urlrouter"
)

var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "task.completed",
			Dest:    UpdateTask,
		},
		urlrouter.Route{
			PathExp: "task.subTasks.updated",
			Dest:    UpdateSubTask,
		},
	},
}
