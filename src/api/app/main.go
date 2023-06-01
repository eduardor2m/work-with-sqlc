package main

import "github.com/eduardor2m/work-with-sqlc/src/api"

func main() {
	api := api.NewAPI(&api.Options{})

	api.Serve()
}
