package main

import (
	"fmt"

	"github.com/go-zoox/docker/container"
)

func main() {
	log, err := container.Logs("go_zoox_job_caa4fefe-7b6c-4656-b9ec-63d4f78f582c")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(log))
}
