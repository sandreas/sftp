package main

import (
	"log"
	"os/user"
	"os"
	"github.com/sandreas/sftp"
)


func main() {
	u, _ := user.Current()
	homeDir := u.HomeDir + "/.graft"
	//logFileName := homeDir + "/graft.log"
	//os.Remove(logFileName)
	//logFile, err := os.OpenFile(logFileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	//if err != nil {
	//	fmt.Println("could not open logfile", err)
	//	os.Exit(1)
	//}
	//defer logFile.Close()
	//
	//mw := io.MultiWriter(os.Stdout, logFile)
	//log.SetOutput(mw)

	log.SetOutput(os.Stdout)

	var matchingPaths []string
	matchingPaths = append(matchingPaths, "examples")
	matchingPaths = append(matchingPaths, "LICENSE")

	sftp.NewSimpleServer(homeDir, "0.0.0.0", 2022, "graft", "graft", matchingPaths)
}