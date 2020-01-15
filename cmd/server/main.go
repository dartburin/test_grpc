package main

import (
	"flag"
	"fmt"
	"os"
	gp "test_grpc/internal/api/server"
	pdb "test_grpc/internal/books"
	lg "test_grpc/internal/logger"
)

func main() {
	// Load init parameters
	dbHostName := flag.String("dbhost", "", "host name")
	dbUser := flag.String("dbuser", "", "user db name")
	dbPass := flag.String("dbpass", "", "user db pass")
	dbBase := flag.String("dbbase", "", "database name")
	dbPort := flag.String("dbport", "", "port for database connect")

	logLevel := flag.String("loglvl", "", "logging message level")
	logFile := flag.String("logfile", "", "logging message to file")
	grpcPort := flag.String("grpcport", "", "port for grpc connect")
	flag.Parse()

	//Init log system
	log := lg.LogInit(*logLevel, *logFile)
	log.Println("Server log system init.")
	lg.PrintOsArgs(log)

	// Check existing obligatory http and db parameters
	if *grpcPort == "" || *dbHostName == "" || *dbUser == "" || *dbPass == "" ||
		*dbBase == "" || *dbPort == "" {
		flag.PrintDefaults()
		fmt.Println("")
		log.Fatalln("Init server error: set not all obligatory parameters.")
		os.Exit(1)
	}

	// Set db parameters
	var configDB pdb.Config
	configDB.Port = *dbPort
	configDB.Host = *dbHostName
	configDB.User = *dbUser
	configDB.Pass = *dbPass
	configDB.Db = *dbBase

	// Connect to DB
	parDB, err := pdb.ConnectToDB(configDB, log)
	if err != nil {
		log.Fatalf("Don`t connect for database (%s).", err.Error())
		os.Exit(1)
	}
	defer parDB.Base.Close()

	//Start gRPC server
	g := gp.New(parDB.Base, *grpcPort, log)
	err = g.Start()
	if err != nil {
		log.Fatalf("Server gRPC error: %s.", err.Error())
	}
}
