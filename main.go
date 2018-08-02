package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/k0kubun/pp"
)

var (
	version   string
	buildDate string
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	confFile := flag.String("conf", "config.yml", "specify config file")
	var showVersion, showConfig bool
	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showConfig, "x", false, "show config")
	flag.Parse()

	if showVersion {
		log.Printf("mirage %v (%v)\n", version, buildDate)
		return
	}

	log.Println("Launch succeeded!")

	cfg := NewConfig(*confFile)

	if showConfig {
		log.Println("mirage config:")
		pp.Print(cfg)
		log.Println("") // add linebreak
	}

	Setup(cfg)
	Run()
}
