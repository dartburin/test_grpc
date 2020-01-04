package main

import (
	"flag"
	"fmt"
	"os"
	pr "test_grpc/internal/api/proxy"
)

func main() {
	//fmt.Printf("Proxy init \n")
	// Load init parameters
	httpPort := flag.String("httpport", "", "port for http connect")
	grpcHostName := flag.String("grpchost", "", "grpc host name")
	grpcPort := flag.String("grpcport", "", "port for grpc connect")

	flag.Parse()

	for _, a := range os.Args {
		fmt.Printf("%s ", a)
	}
	fmt.Println("")

	// Check existing obligatory http and db parameters
	if *grpcHostName == "" || *httpPort == "" || *grpcPort == "" {
		fmt.Println("Init proxy error: set not all obligatory parameters.")
		flag.PrintDefaults()
		fmt.Println("")
		os.Exit(1)
	}

	//Start gRPC + HTTP server
	g := pr.New(*grpcHostName, *grpcPort, *httpPort)
	fmt.Printf("HTTP proxy start\n")

	err := g.Start()

	if err != nil {
		fmt.Printf("HTTP proxy error: %s\n", err.Error())
	}
}
