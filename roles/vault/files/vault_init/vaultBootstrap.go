package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	mode = flag.String("mode", "init", "init or unseal")
	port = flag.Int("port", 8200, "Port for vault")
	file = flag.String("file", "/tmp/vaultkeys", "Port for vault")
)

func main() {
	// parse flags
	flag.Parse()

	if *mode == "init" {
		keys := vaultInit(*port)
		err := ioutil.WriteFile(*file, []byte(keys), 0666)
		if err != nil {
			os.Exit(1)
		}

	} else if *mode == "unseal" {
		raw, err := ioutil.ReadFile(*file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		lines := strings.Split(string(raw), "\n")
		var keys []string
		for i := range lines {
			keypair := strings.Split(lines[i], ": ")
			keys = append(keys, keypair[1])
		}
		unsealed := vaultUnseal(8200, keys)
		if !unsealed {
			fmt.Println("Not able to unseal vault")
			os.Exit(1)
		}
	}
}
