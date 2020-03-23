module github.com/rubblelabs/ripple

go 1.13

require (
	github.com/agl/ed25519 v0.0.0-00010101000000-000000000000
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/fatih/color v1.9.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/gorilla/websocket v1.4.2
	github.com/juju/loggo v0.0.0-20190526231331-6e530bcce5d8 // indirect
	github.com/juju/testing v0.0.0-20191001232224-ce9dec17d28b
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/willf/bitset v1.1.10
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
)

replace github.com/agl/ed25519 => github.com/inn4science/ed25519 v1.0.0
