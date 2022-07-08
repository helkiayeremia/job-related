package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// Make a redis pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "academy-booking-api-redis.persistence-id-stag.svc.cluster.local:6379")
	},
}

type Context struct {
	customerID int64
}

type TriggerAutoBookReq struct {
	MeetingSerial string
	ClassSerial   string
}

var namespaceWorker = "academy-booking"

func main() {
	// Make a new pool. Arguments:
	// Context{} is a struct that will be the context for the request.
	// 10 is the max concurrency
	// "my_app_namespace" is the Redis namespace
	// redisPool is a Redis pool
	pool := work.NewWorkerPool(Context{}, 10, namespaceWorker, redisPool)

	// Add middleware that will be executed for each job
	pool.Middleware((*Context).Log)

	// Map the name of jobs to handler functions
	pool.Job("auto-book", (*Context).AutoBook)
	pool.Job("auto-book-by-source-serial", (*Context).AutoBookBySourceSerial)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job: ", job.Name)
	return next()
}

func (c *Context) AutoBook(job *work.Job) error {
	// Extract arguments:
	data := job.Args["data"]
	if err := job.ArgError(); err != nil {
		return err
	}

	byteData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error marshal: ", err)
	}

	var triggerAutoBookReq *TriggerAutoBookReq
	if err = json.Unmarshal(byteData, &triggerAutoBookReq); err != nil {
		fmt.Println("error unmarshal: ", err)
	}

	fmt.Println("meetingSerial : ", triggerAutoBookReq.MeetingSerial)
	fmt.Println("classSerial : ", triggerAutoBookReq.ClassSerial)
	fmt.Println("Job enqueued at : ", job.EnqueuedAt)
	fmt.Println("Job id : ", job.ID)

	return nil
}

func (c *Context) AutoBookBySourceSerial(job *work.Job) error {
	fmt.Println(job.ArgString("userSerial"))
	return nil
}
