package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var namespaceWorker = "ea-booking-api"
var enqueuer = work.NewEnqueuer(namespaceWorker, redisPool)

var redisHost = ":6380"

// var redisHost = "ea-booking-api-redis.persistence-id-prod.svc.cluster.local:6379"

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", redisHost)
	},
}

func main() {
	meetingSerialArg := flag.String("m", "LSMT-TEST", "meeting serial")
	classSerialArg := flag.String("c", "CLASS-TEST", "class serial")
	bookingPeriodEndTimeArg := flag.String("t", "2020-01-01T00:00:00+07:00", "when the job run") // default immediately run
	flag.Parse()

	meetingSerial := *meetingSerialArg
	classSerial := *classSerialArg
	bookingPeriodEndTime, err := time.Parse(time.RFC3339, *bookingPeriodEndTimeArg)
	if err != nil {
		fmt.Println("error", err)
	}

	secondsFromNow := countSecondsFromNow(&bookingPeriodEndTime)

	_, err = enqueuer.EnqueueUniqueIn("auto-book", secondsFromNow, work.Q{
		"data": map[string]interface{}{
			"meetingSerial": meetingSerial,
			"classSerial":   classSerial,
		},
		"object_id_": fmt.Sprintf("%s:unique-handler:%s", namespaceWorker, meetingSerial),
	})
	if err != nil {
		fmt.Println("error", err)
	}
}

func countSecondsFromNow(later *time.Time) int64 {
	secondsNow := time.Now().Unix()
	scondsLater := later.Unix()

	return scondsLater - secondsNow
}
