package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/LYH2263/go-delayqueue"
	"github.com/LYH2263/go-delayqueue/console"
)

func main() {
	addr := flag.String("addr", ":8220", "listen address")
	journal := flag.String("journal", "", "journal path")
	flag.Parse()
	opts := delayqueue.Options{}
	if *journal != "" {
		opts.JournalPath = *journal
	}
	b, err := delayqueue.New(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()
	api := &console.API{Broker: b}
	srv := console.New(api)
	log.Printf("delayqueue console on %s", *addr)
	if err := http.ListenAndServe(*addr, srv.Handler); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
