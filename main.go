package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nasonawa/quay-yala/cli"
)

func main() {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
