package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/errors"
	"github.com/RichardKnop/machinery/v1/signatures"
	"github.com/fatih/color"
)

func ctrlc() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		color.Set(color.FgGreen)
		fmt.Println("\nExecution stopped by", sig)
		color.Unset()
		os.Exit(0)
	}()
}

func main() {
	var cnf = config.Config{
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
	}

	server, err := machinery.NewServer(&cnf)
	if err != nil {
		errors.Fail(err, "Can not create server!")
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		ctrlc()
		cmd := strings.Trim(scanner.Text(), "\n")
		execTask := signatures.TaskSignature{
			Name: "OSCommand",
			Args: []signatures.TaskArg{
				signatures.TaskArg{
					Type:  "string",
					Value: cmd,
				},
			},
		}
		asyncResult, err := server.SendTask(&execTask)
		if err != nil {
			color.Set(color.FgRed)
			errors.Fail(err, "Could not send task")
		}

		result, err := asyncResult.Get()
		if err != nil {
			color.Set(color.FgRed)
			errors.Fail(err, "Getting task state failed with error")
		}
		color.Set(color.FgGreen)
		fmt.Print(result.Interface())
		color.Unset()
	}

}
