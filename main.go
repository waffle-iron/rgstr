package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/castillobg/rgstr/registries"
	// A blank import so that the consul AdapterFactory registers itself.
	_ "github.com/castillobg/rgstr/registries/consul"
	"github.com/castillobg/rgstr/runtimes"
	// A blank import so that the rkt AdapterFactory registers itself.
	_ "github.com/castillobg/rgstr/runtimes/rkt"
)

func main() {
	runtimeAddr := flag.String("a", "localhost:15441", "The `address` where rkt's API service is listening.")
	registryAddr := flag.String("ra", "localhost:8500", "The `registry address`.")
	registryName := flag.String("registry", "consul", "The `registry`. rgstr currently supports Consul.")
	runtimeName := flag.String("runtime", "rkt", "The `runtime`. rgstr currently supports rkt.")
	flag.Parse()

	registryFactory, ok := registries.LookUp(*registryName)
	if !ok {
		fmt.Printf("No registry with name \"%s\" found.\n", *registryName)
		os.Exit(1)
	}
	registry, err := registryFactory.New(*registryAddr)
	if err != nil {
		fmt.Printf("Error initializing registry client for \"%s\": %s\n", *registryName, err.Error())
		os.Exit(1)
	}

	runtimeFactory, ok := runtimes.LookUp(*runtimeName)
	if !ok {
		fmt.Printf("No runtime with name \"%s\" found.\n", *runtimeName)
		os.Exit(1)
	}
	runtime, err := runtimeFactory.New(*runtimeAddr, registry)
	if err != nil {
		fmt.Printf("Error initializing runtime client for \"%s\": %s", runtime, err.Error())
		os.Exit(1)
	}

	errs := make(chan error)
	go runtime.Listen(errs)
	fmt.Printf("rgstr is listening for changes in %s...\n", *runtimeName)
	log.Fatal(<-errs)
}
