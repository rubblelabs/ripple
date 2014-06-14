#TODO

##Data
* Write good tests for metadata interpretation
* Improve string formatting for Transaction types and create a short summary string

##Websockets
* Add missing commands
* Make connection resilient via reconnection strategy (r.ripple.com?)
* Allow connection to multiple endpoints?

##Tools

###tx
* Implement OfferCreate, OfferCancel, AccountSet and TrustSet commands
* Use websockets to optionally acquire correct sequence number for account derived from seed 
* Add memo support
