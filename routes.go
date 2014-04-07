package main

import (
	"github.com/ant0ine/go-urlrouter"
	"github.com/wurkhappy/WH-Tasks/handlers"
)

//order matters so most general should go towards the bottom
var router urlrouter.Router = urlrouter.Router{
	Routes: []urlrouter.Route{
		urlrouter.Route{
			PathExp: "/agreements/v/:id/tasks",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"POST": handlers.CreateTasksByVersionID,
				"GET":  handlers.GetTasksByVersionID,
			},
		},
		urlrouter.Route{
			PathExp: "/tasks/:id",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"PUT": handlers.UpdateTask,
			},
		},
		urlrouter.Route{
			PathExp: "/tasks/:id/action",
			Dest: map[string]func(map[string]interface{}, []byte, string) ([]byte, error, int){
				"POST": handlers.UpdateAction,
			},
		},
	},
}
