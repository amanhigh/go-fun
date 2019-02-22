package go_metrics

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	metrics "github.com/rcrowley/go-metrics"
)

func MetricFun() {
	rand.Seed(time.Now().UnixNano())

	go metrics.Log(metrics.DefaultRegistry, 2*time.Second, log.New(os.Stdout, "metrics: ", log.Lmicroseconds))

	name := "com.aman.work.time"
	//t := metrics.NewTimer()
	//_ = metrics.Register(name, t)
	t := metrics.GetOrRegister(name, metrics.NewTimer()).(metrics.Timer)

	for i := 0; i < 100; i++ {
		//fmt.Println("Doing Work & Sleeping")
		t.Time(func() {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		})
	}

	fmt.Println(metrics.Get(name))
	time.Sleep(4 * time.Second)
}
