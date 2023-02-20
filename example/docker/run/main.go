package main

import "github.com/go-zoox/docker/container"

func main() {
	if err := container.Run(&container.RunConfig{
		Name:  "gozoox_docker_dd5003339e011",
		Image: "whatwewant/node-builder:v16-1.5.3",
		Commands: []string{
			`
			set -e
			
			errasd asdasd

			git clone --progress --verbose https://github.com/zcorky/storage --depth 1

			cd storage
			
			npm i --verbose --registry https://registry.npmmirror.com
			
			./node_modules/.bin/tant-intl build --verbose && ./node_modules/.bin/nx run-many --target=build --projects=portal,ta,va --verbose`,
		},
		// Volumes: map[string]string{
		// 	"/tests/go-zoox/docker": "/tests",
		// },
		WorkingDir: struct {
			Host      string
			Container string
		}{
			Host:      "/tmp/go-zoox/docker",
			Container: "/app",
		},
		Env: map[string]string{
			"FOO": "BAR",
		},
	}); err != nil {
		panic("run err:" + err.Error())
	}

	// // Dockerfile
	// if err := container.Run(&container.RunConfig{
	// 	Name:  "gozoox_docker_dd5003339e01",
	// 	Image: "whatwewant/node-builder:v16-1.5.3",
	// 	Dockerfile: `
	// FROM node:16-alpine

	// CMD echo "hello world"
	// 	`,
	// 	// Volumes: map[string]string{
	// 	// 	"/tests/go-zoox/docker": "/tests",
	// 	// },
	// 	WorkingDir: struct {
	// 		Host      string
	// 		Container string
	// 	}{
	// 		Host:      "/tmp/go-zoox/docker",
	// 		Container: "/app",
	// 	},
	// 	Env: map[string]string{
	// 		"FOO": "BAR",
	// 	},
	// 	// AutoRemove: true,
	// }); err != nil {
	// 	panic(err)
	// }
}
