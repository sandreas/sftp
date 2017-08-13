// An example SFTP server implementation using the golang SSH package.
// Serves the whole filesystem visible to the user, and has a hard-coded username and password,
// so not for real use!
package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/sandreas/sftp/examples/sftp-server-simple/sftpd"
)

// Based on example server code from golang.org/x/crypto/ssh and server_standalone
func main() {

	homeDir, err := createHomeDirectoryIfNotExists()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	basePath := "./examples"
	files := []string{}
	err = filepath.Walk(basePath, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	// fixed file list
	//files := []string{
	//	"examples",
	//	"LICENSE",
	//}

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	pathMapper := sftpd.NewPathMapper(files, basePath)
	sftpd.NewSimpleSftpServer(homeDir, "0.0.0.0", 2022, "test", "test", pathMapper)
}

func createHomeDirectoryIfNotExists() (string, error) {
	u, _ := user.Current()
	homeDir := u.HomeDir + "/.pkg-sftp"
	if _, err := os.Stat(homeDir); err != nil {
		if err := os.Mkdir(homeDir, os.FileMode(0755)); err != nil {
			return homeDir, err
		}
	}
	return homeDir, nil
}
