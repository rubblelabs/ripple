// Simple tool to find the seeds for ripple account ids which match a regular expression
package main

import (
	"flag"
	"github.com/donovanhide/ripple/crypto"
	"log"
	"os"
	"os/signal"
	"regexp"
	"runtime"
)

var name = flag.String("name", "ripple", "desired name to appear in ripple account id")
var insensitive = flag.Bool("insenstive", true, "ignore case sensitivity")

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func search(target *regexp.Regexp) {
	for {
		key, err := crypto.GenerateRootDeterministicKey(nil)
		checkErr(err)
		account, err := key.GenerateAccountId(0)
		checkErr(err)
		if target.MatchString(account.String()) {
			log.Println(key.Seed.String(), account.String())
		}
	}
}

func main() {
	flag.Parse()
	match := *name
	if *insensitive {
		match = "(?i)" + match
	}
	target, err := regexp.Compile(match)
	checkErr(err)
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, os.Kill)
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Printf("Searching for \"%s\" with %d processors", *name, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go search(target)
	}
	<-kill
}
