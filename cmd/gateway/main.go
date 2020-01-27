package main

import (
	"flag"
	"fmt"
	"os"
	pr "test_grpc/internal/api/gateway"
	lg "test_grpc/internal/logger"
)

func main() {
	//fmt.Printf("Proxy init \n")
	// Load init parameters
	httpPort := flag.String("httpport", "8080", "port for http connect")
	grpcHostName := flag.String("grpchost", "", "grpc host name")
	grpcPort := flag.String("grpcport", "8086", "port for grpc connect")

	logLevel := flag.String("loglvl", "", "logging message level")
	logFile := flag.String("logfile", "", "logging message to file")
	flag.Parse()

	//Init log system
	log := lg.LogInit(*logLevel, *logFile)
	log.Println("Server log system init.")
	lg.PrintOsArgs(log)

	// Check existing obligatory http and db parameters
	if *grpcHostName == "" {
		flag.PrintDefaults()
		fmt.Println("")
		log.Fatalln("Init server error: set not all obligatory parameters.")
		os.Exit(1)
	}

	//Start gRPC + HTTP server
	g := pr.New(*grpcHostName, *grpcPort, *httpPort, log)
	err := g.Start()
	if err != nil {
		log.Fatalf("HTTP gateway error: %s.", err.Error())
	}
}
