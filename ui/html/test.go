package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func printer(m ...string) {
	fmt.Println(m)
}


func main() {
	files := []string{}
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".tmpl") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
	printer(files...)
}
