package main

import (
	"log"

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

// Make an enqueuer with a particular namespace
var client = work.NewClient("academy-booking", redisPool)

func main() {

	if jobs, total, err := client.DeadJobs(1); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Total=%d\n", total)
		for _, job := range jobs {
			log.Printf("Job=%s\n", job.Args["userSerial"])
		}
	}
}
