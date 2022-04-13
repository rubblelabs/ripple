ripple
======

Go packages to interact with the Ripple protocol.

[![Go Reference](https://pkg.go.dev/badge/github.com/rubblelabs/ripple.svg)](https://pkg.go.dev/github.com/rubblelabs/ripple)
[![Build Status](https://github.com/rubblelabs/ripple/actions/workflows/go.yml/badge.svg)](https://github.com/rubblelabs/ripple/actions/workflows/go.yml)

The data, crypto, and websockets packages are very functional and quite well tested. Most websockets commands are implemented but not all.

The peers and ledger packages are the least polished packages currently, and they are very much unfinished (and the tests might be non-existent or non-functional), but better to get the code out in the open.

We've included command-line tools to show how to apply the library:

* listener: connects to rippled servers with the peering protocol and displays the traffic
* subscribe: tracks ledgers and transactions via websockets and explains each transaction's metadata
* tx: creates transactions, signs them, and submits them via websockets
* vanity: generates new ripple wallets in search of vanity addresses

The hope is one day that these packages might lay the foundations for an alternative implementation of the [Ripple daemon](https://github.com/ripple/rippled). This is, however, a long way off!

Please bear in mind that this has been an exercise that has taken a lot of time, so if you want to help and are not a developer, bounties and thanks are more than welcome. Please see the [AUTHORS](https://github.com/rubblelabs/ripple/blob/master/AUTHORS) file. If you'd like to chat about the code, have a look here:

[![Gitter chat](https://badges.gitter.im/rubblelabs/ripple.png)](https://gitter.im/rubblelabs/ripple)
