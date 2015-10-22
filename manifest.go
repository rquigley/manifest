package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"log"
	"os"
	"path/filepath"
)

var fileMap map[string]string

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		checksum, _ := checksumFilepath(path)
		checksum = checksum[:8]
		fileMap[path] = checksum
	}
	return nil

}

func checksumFilepath(filePath string) (string, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return string(result), err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return string(result), err
	}

	return hex.EncodeToString(hash.Sum(result)), nil
}

func main() {
	fileMap = make(map[string]string)
	enc := json.NewEncoder(os.Stdout)

	app := cli.NewApp()
	app.Name = "manifest"
	app.Usage = "build JSON containing filepath -> checksum for every file at specified path "
	app.Action = func(c *cli.Context) {
		root := c.Args()[0]
		err := filepath.Walk(root, visit)
		if err != nil {
			fmt.Printf("filepath.Walk() returned %v\n", err)
		}
		// for k, v := range fileMap {
		//         fmt.Println("k:", k, "v:", v)
		// }
		if err := enc.Encode(&fileMap); err != nil {
			log.Println(err)
		}
	}

	app.Run(os.Args)
}
