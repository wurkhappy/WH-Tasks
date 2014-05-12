package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/wurkhappy/WH-Config"
	"github.com/wurkhappy/WH-Tasks/WH-TasksUpdates/models"
	"time"
)

const INTERVAL_PERIOD time.Duration = 24 * time.Hour

const HOUR_TO_TICK int = 10 //this will be UTC time since servers are on UTC
const MINUTE_TO_TICK int = 0
const SECOND_TO_TICK int = 00

var nextTick time.Time
var production = flag.Bool("production", false, "Production settings")

func main() {
	flag.Parse()
	if *production {
		config.Prod()
	} else {
		config.Test()
	}
	dialRMQ()
	setupRedis()
	ticker := updateTicker()
	for {
		<-ticker.C
		fmt.Println(time.Now(), "- just ticked")
		ticker = updateTicker()
	}
}

func updateTicker() *time.Ticker {
	if nextTick.IsZero() {
		nextTick = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), HOUR_TO_TICK, MINUTE_TO_TICK, SECOND_TO_TICK, 0, time.Local)
	} else {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}

	fmt.Println(nextTick, "- next tick")
	ProcessTasks()
	diff := nextTick.Sub(time.Now())
	if diff.Seconds() < 0 {
		return updateTicker()
	}
	return time.NewTicker(diff)
}

func ProcessTasks() error {
	c := redisPool.Get()
	v, err := redis.Values(c.Do("KEYS", "WHTU_*"))
	if err != nil {
		return fmt.Errorf("%s", "There was an error finding that token")
	}

	//redis returns an array of keys so let's iterate through that
	for _, value := range v {
		tasks := []*models.Task{}
		key := string(value.([]byte))

		//for the given key let's get the fields and values
		vals, err := redis.Values(c.Do("HGETALL", key))
		if err != nil {
			return fmt.Errorf("%s", "There was an error finding that token")
		}

		for i := len(vals); i >= 0; i-- {
			//odd numbers are fields, evens are values in redis
			if (i+1)%2 == 0 {
				var t *models.Task
				json.Unmarshal(vals[i].([]byte), &t)
				if t.LastAction != nil && t.LastAction.Name == "completed" {
					tasks = append(tasks, t)
				}
			}
		}

		if len(tasks) > 0 {
			ts, _ := json.Marshal(tasks)
			event := &Event{"tasks.updated.notify", ts}
			event.Publish()
		}

		//When we're done delete the key so we start fresh for the next day
		if _, err := c.Do("DEL", key); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
