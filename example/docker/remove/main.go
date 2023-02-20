package main

import (
	"github.com/go-zoox/docker/container"
)

func main() {
	err := container.Remove("gozoox_docker_dd5003339e01")
	if err != nil {
		panic(err)
	}
}
