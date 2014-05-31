package main

import (
	"flag"
	"github.com/donovanhide/ripple/crypto"
	"github.com/donovanhide/ripple/ledger"
	"github.com/donovanhide/ripple/peers"
	"github.com/donovanhide/ripple/storage"
	"github.com/golang/glog"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/influxdb"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"time"
)

var trusted = flag.String("trusted", "r.ripple.com:51235", "trusted hosts separated by commas")
var maxPeers = flag.Int("maxpeers", 1, "maximum number of peers to connect to")
var name = flag.String("name", "RippleListener", "name to connect to the peer network as")
var port = flag.String("port", "51235", "port to use to connect to the peer network")

func checkErr(err error) {
	if err != nil {
		glog.Fatalln(err)
	}
}

func servePeers(m *peers.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := make(chan []byte)
		m.Status <- c
		w.Write(<-c)
	}
}

func main() {
	flag.Parse()
	go influxdb.Influxdb(metrics.DefaultRegistry, time.Second*5, &influxdb.Config{
		Host:     "192.168.0.109:8086",
		Database: "ripple",
		Username: "root",
		Password: "root",
	})
	runtime.GOMAXPROCS(runtime.NumCPU())
	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, os.Kill)
	key, err := crypto.GenerateRootDeterministicKey(nil)
	checkErr(err)
	db := storage.NewEmptyMemoryDB()
	mgr, err := ledger.NewManager(db)
	checkErr(err)
	go mgr.Start()
	config := &peers.Config{
		Key:      key,
		Name:     *name,
		Port:     *port,
		Sync:     mgr,
		MaxPeers: *maxPeers,
		Trusted:  *trusted,
	}
	peerManager, err := peers.NewManager(config)
	checkErr(err)
	http.Handle("/peers", servePeers(peerManager))
	go http.ListenAndServe(":8000", nil)
	<-kill
	peerManager.Quit <- true
}
