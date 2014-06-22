#TODO

##Data
* Write good tests for metadata interpretation
* Use Freeform type for _some_ memos and Previous/New/Final fields

##Peers
* Implement all handlers
* Clean out ripple.proto

##Ledger
* Allow subscribing to incoming Proposals/Validations/Transactions for use in listener

##Terminal
* Add pathset output
* Add proposal and validation output

##Websockets
* Add missing commands
* Make connection resilient via reconnection strategy (r.ripple.com?)
* Allow connection to multiple endpoints?

##Tools

###tx
* Implement OfferCreate, OfferCancel, AccountSet and TrustSet commands
* Use websockets to optionally acquire correct sequence number for account derived from seed 
* Add memo support
