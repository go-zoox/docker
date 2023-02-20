package main

import (
	"os"
	"path"

	"github.com/go-zoox/docker/image"
)

func main() {
	pwd, _ := os.Getwd()

	err := image.Build(&image.BuildConfig{
		Name:    "testgozoox:latest",
		Context: path.Join(pwd, "example/docker/build"),
		// Dockerfile: path.Join(pwd, "example/docker/build", "Dockerfile"),
	})
	if err != nil {
		panic(err)
	}
}
