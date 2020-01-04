package main

import (
	"flag"
	"fmt"
	"os"
	gp "test_grpc/internal/api/server"
	pdb "test_grpc/internal/db"
)

func main() {
	//fmt.Printf("Server init \n")
	// Load init parameters
	var configDB pdb.Config
	dbHostName := flag.String("dbhost", "", "host name")
	dbUser := flag.String("dbuser", "", "user db name")
	dbPass := flag.String("dbpass", "", "user db pass")
	dbBase := flag.String("dbbase", "", "database name")
	dbPort := flag.String("dbport", "", "port for database connect")

	grpcPort := flag.String("grpcport", "", "port for grpc connect")

	flag.Parse()

	for _, a := range os.Args {
		fmt.Printf("%s ", a)
	}
	fmt.Println("")

	// Check existing obligatory http and db parameters
	if *grpcPort == "" || *dbHostName == "" || *dbUser == "" || *dbPass == "" ||
		*dbBase == "" || *dbPort == "" {
		fmt.Println("Init server error: set not all obligatory parameters.")
		flag.PrintDefaults()
		fmt.Println("")
		os.Exit(1)
	}

	// Set db parameters
	configDB.Port = *dbPort
	configDB.Host = *dbHostName
	configDB.User = *dbUser
	configDB.Pass = *dbPass
	configDB.Db = *dbBase

	// Connect to DB
	//fmt.Printf("Base start\n")
	parDB, err := pdb.ConnectToDB(configDB)
	if err != nil {
		fmt.Printf("Base not connect\n")
		os.Exit(1)
	}
	defer parDB.Base.Close()

	//Start gRPC server
	g := gp.New(parDB.Base, *grpcPort)
	fmt.Printf("Server start\n")

	err = g.Start()
	if err != nil {
		fmt.Printf("Server gRPC error: %s\n", err.Error())
	}
}
