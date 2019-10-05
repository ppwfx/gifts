package main

import (
	"fmt"
	"github.com/21stio/go-playground/gifts/pkg/boot"
	"log"
	"os"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	d, err := boot.GetDependencies(requireEnv("EMPLOYEES_JSON"), requireEnv("GIFTS_JSON"), requireEnv("ADDR"))
	if err != nil {
	    return
	}

	return d.Server.Serve()
}

func requireEnv(k string) (s string) {
	s = os.Getenv(k)
	if s == "" {
		panic(fmt.Sprintf("ENV %v required", k))
	}

	return
}


