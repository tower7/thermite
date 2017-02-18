package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/errors"
	"github.com/tower7/thermite/tasks"
)

func main() {
	var NWorkers = flag.Int("n", 1, "The number of workers to start")
	var cnf = config.Config{
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
	}

	flag.Parse()

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		errors.Fail(err, "Could not create server")
	}
	server.RegisterTask("OSCommand", tasks.OSCommand)
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	euid := os.Geteuid()
	prefix := "OSCommand" + "-" + hostname + "-" + strconv.Itoa(euid)
	suffix := runtime.GOOS + "_" + runtime.GOARCH
	for i := 0; i < *NWorkers; i++ {
		fmt.Println(prefix + "_" + suffix + "_" + string(i))
		worker := server.NewWorker(prefix + "_" + suffix + "_" + string(i))
		go worker.Launch()
	}
	runtime.Goexit()
}
