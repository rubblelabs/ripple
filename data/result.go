package data

type TransactionResult int16

const (
	// 0: S Success (success)
	// Causes:
	// - Success.
	// Implications:
	// - Applied
	// - Forwarded
	tesSUCCESS TransactionResult = 0
)
const (
	// 100 .. 159 C Claim fee only (ripple transaction with no good paths, pay to non-existent account, no path)
	// Causes:
	// - Success, but does not achieve optimal result.
	// - Invalid transaction or no effect, but claim fee to use the sequence number.
	// Implications:
	// - Applied
	// - Forwarded
	// Only allowed as a return code of appliedTransaction when !tapRetry. Otherwise, treated as terRETRY.
	//
	// DO NOT CHANGE THESE NUMBERS: They appear in ledger meta data.
	tecCLAIM TransactionResult = iota + 100
	tecPATH_PARTIAL
	tecUNFUNDED_ADD
	tecUNFUNDED_OFFER
	tecUNFUNDED_PAYMENT
	tecFAILED_PROCESSING
)
const (
	tecDIR_FULL TransactionResult = iota + 121
	tecINSUF_RESERVE_LINE
	tecINSUF_RESERVE_OFFER
	tecNO_DST
	tecNO_DST_INSUF_XRP
	tecNO_LINE_INSUF_RESERVE
	tecNO_LINE_REDUNDANT
	tecPATH_DRY
	tecUNFUNDED
	tecNO_ALTERNATIVE_KEY
	tecNO_REGULAR_KEY
	tecOWNERS
	tecNO_ISSUER
	tecNO_AUTH
	tecNO_LINE
	tecINSUFF_FEE
	tecFROZEN
	tecNO_TARGET
	tecNO_PERMISSION
	tecNO_ENTRY
	tecINSUFFICIENT_RESERVE
	tecNEED_MASTER_KEY
	tecDST_TAG_NEEDED
	tecINTERNAL
	tecOVERSIZE
	tecCRYPTOCONDITION_ERROR
	tecINVARIANT_FAILED
	tecEXPIRED
	tecDUPLICATE
	tecKILLED
	tecHAS_OBLIGATIONS
	tecTOO_SOON
	tecHOOK_ERROR
	tecMAX_SEQUENCE_REACHED
	tecNO_SUITABLE_NFTOKEN_PAGE
	tecNFTOKEN_BUY_SELL_MISMATCH
	tecNFTOKEN_OFFER_TYPE_MISMATCH
	tecCANT_ACCEPT_OWN_NFTOKEN_OFFER
	tecINSUFFICIENT_FUNDS
	tecOBJECT_NOT_FOUND
	tecINSUFFICIENT_PAYMENT
	tecUNFUNDED_AMM
	tecAMM_BALANCE
	tecAMM_FAILED
	tecAMM_INVALID_TOKENS
	tecAMM_EMPTY
	tecAMM_NOT_EMPTY
	tecAMM_ACCOUNT
	tecINCOMPLETE
	tecXCHAIN_BAD_TRANSFER_ISSUE
	tecXCHAIN_NO_CLAIM_ID
	tecXCHAIN_BAD_CLAIM_ID
	tecXCHAIN_CLAIM_NO_QUORUM
	tecXCHAIN_PROOF_UNKNOWN_KEY
	tecXCHAIN_CREATE_ACCOUNT_NONXRP_ISSUE
	tecXCHAIN_WRONG_CHAIN
	tecXCHAIN_REWARD_MISMATCH
	tecXCHAIN_NO_SIGNERS_LIST
	tecXCHAIN_SENDING_ACCOUNT_MISMATCH
	tecXCHAIN_INSUFF_CREATE_AMOUNT
	tecXCHAIN_ACCOUNT_CREATE_PAST
	tecXCHAIN_ACCOUNT_CREATE_TOO_MANY
	tecXCHAIN_PAYMENT_FAILED
	tecXCHAIN_SELF_COMMIT
	tecXCHAIN_BAD_PUBLIC_KEY_ACCOUNT_PAIR
	tecXCHAIN_CREATE_ACCOUNT_DISABLED
	tecEMPTY_DID
	tecINVALID_UPDATE_TIME
	tecTOKEN_PAIR_NOT_FOUND
	tecARRAY_EMPTY
	tecARRAY_TOO_LARGE
)

const (
	// Transaction Errors
	// -399 .. -300: L Local error (transaction fee inadequate, exceeds local limit)
	// Only valid during non-consensus processing.
	// Implications:
	// - Not forwarded
	// - No fee check
	telLOCAL_ERROR TransactionResult = iota - 399
	telBAD_DOMAIN
	telBAD_PATH_COUNT
	telBAD_PUBLIC_KEY
	telFAILED_PROCESSING
	telINSUF_FEE_P
	telNO_DST_PARTIAL
	telCAN_NOT_QUEUE
	telCAN_NOT_QUEUE_BALANCE
	telCAN_NOT_QUEUE_BLOCKS
	telCAN_NOT_QUEUE_BLOCKED
	telCAN_NOT_QUEUE_FEE
	telCAN_NOT_QUEUE_FULL
)
const (
	// -299 .. -200: M Malformed (bad signature)
	// Causes:
	// - Transaction corrupt.
	// Implications:
	// - Not applied
	// - Not forwarded
	// - Reject
	// - Can not succeed in any imagined ledger.
	temMALFORMED TransactionResult = iota - 299
	temBAD_AMOUNT
	temBAD_CURRENCY
	temBAD_EXPIRATION
	temBAD_FEE
	temBAD_ISSUER
	temBAD_LIMIT
	temBAD_OFFER
	temBAD_PATH
	temBAD_PATH_LOOP
	temBAD_SEND_XRP_LIMIT
	temBAD_SEND_XRP_MAX
	temBAD_SEND_XRP_NO_DIRECT
	temBAD_SEND_XRP_PARTIAL
	temBAD_SEND_XRP_PATHS
	temBAD_SEQUENCE
	temBAD_SIGNATURE
	temBAD_SRC_ACCOUNT
	temBAD_TRANSFER_RATE
	temDST_IS_SRC
	temDST_NEEDED
	temINVALID
	temINVALID_FLAG
	temREDUNDANT
	temRIPPLE_EMPTY
	temDISABLED
	temBAD_SIGNER
	temBAD_QUORUM
	temBAD_WEIGHT
	temBAD_TICK_SIZE
	temINVALID_ACCOUNT_ID
	temCANNOT_PREAUTH_SELF
	temUNCERTAIN
	temUNKNOWN
	temSEQ_AND_TICKET
	temBAD_NFTOKEN_TRANSFER_FEE
	temBAD_AMM_TOKENS
	temXCHAIN_EQUAL_DOOR_ACCOUNTS
	temXCHAIN_BAD_PROOF
	temXCHAIN_BRIDGE_BAD_ISSUES
	temXCHAIN_BRIDGE_NONDOOR_OWNER
	temXCHAIN_BRIDGE_BAD_MIN_ACCOUNT_CREATE_AMOUNT
	temXCHAIN_BRIDGE_BAD_REWARD_AMOUNT
	temEMPTY_DID
	temARRAY_EMPTY
	temARRAY_TOO_LARGE
)
const (
	// -199 .. -100: F Failure (sequence number previously used)
	// Causes:
	// - Transaction cannot succeed because of ledger state.
	// - Unexpected ledger state.
	// - C++ exception.
	// Implications:
	// - Not applied
	// - Not forwarded
	// - Could succeed in an imagined ledger.
	tefFAILURE TransactionResult = iota - 199
	tefALREADY
	tefBAD_ADD_AUTH
	tefBAD_AUTH
	tefBAD_CLAIM_ID
	tefBAD_GEN_AUTH
	tefBAD_LEDGER
	tefCLAIMED
	tefCREATED
	tefDST_TAG_NEEDED
	tefEXCEPTION
	tefGEN_IN_USE
	tefINTERNAL
	tefNO_AUTH_REQUIRED // Can't set auth if auth is not required.
	tefPAST_SEQ
	tefWRONG_PRIOR
	tefMASTER_DISABLED
	tefMAX_LEDGER
	tefBAD_SIGNATURE
	tefBAD_QUORUM
	tefNOT_MULTI_SIGNING
	tefBAD_AUTH_MASTER
	tefINVARIANT_FAILED
	tefTOO_BIG
	tefNO_TICKET
	tefNFTOKEN_IS_NOT_TRANSFERABLE
)
const (
	// -99 .. -1: R Retry (sequence too high, no funds for txn fee, originating account non-existent)
	// Causes:
	// - Prior application of another, possibly non-existant, another transaction could allow this transaction to succeed.
	// Implications:
	// - Not applied
	// - Not forwarded
	// - Might succeed later
	// - Hold
	// - Makes hole in sequence which jams transactions.
	terRETRY       TransactionResult = iota - 99
	terFUNDS_SPENT                   // This is a free transaction, therefore don't burden network.
	terINSUF_FEE_B                   // Can't pay fee TransactionError     = -99 therefore don't burden network.
	terNO_ACCOUNT                    // Can't pay fee, therefore don't burden network.
	terNO_AUTH                       // Not authorized to hold IOUs.
	terNO_LINE                       // Internal flag.
	terOWNERS                        // Can't succeed with non-zero owner count.
	terPRE_SEQ                       // Can't pay fee, no point in forwarding, therefore don't burden network.
	terLAST                          // Process after all other transactions
	terNO_RIPPLE                     // Rippling not allowed
	terQUEUED                        // Transaction is being held in TxQ until fee drops
	terPRE_TICKET                    //Ticket is not yet in ledger but might be on its way
	terNO_AMM                        // AMM doesn't exist for the asset pair
)

var resultNames = map[TransactionResult]struct {
	Token string
	Human string
}{
	tesSUCCESS:                            {"tesSUCCESS", "The transaction was applied."},
	tecCLAIM:                              {"tecCLAIM", "Fee claimed. Sequence used. No action."},
	tecDIR_FULL:                           {"tecDIR_FULL", "Can not add entry to full directory."},
	tecFAILED_PROCESSING:                  {"tecFAILED_PROCESSING", "Failed to correctly process transaction."},
	tecINSUF_RESERVE_LINE:                 {"tecINSUF_RESERVE_LINE", "Insufficient reserve to add trust line."},
	tecINSUF_RESERVE_OFFER:                {"tecINSUF_RESERVE_OFFER", "Insufficient reserve to create offer."},
	tecNO_DST:                             {"tecNO_DST", "Destination does not exist. Send XRP to create it."},
	tecNO_DST_INSUF_XRP:                   {"tecNO_DST_INSUF_XRP", "Destination does not exist. Too little XRP sent to create it."},
	tecNO_LINE_INSUF_RESERVE:              {"tecNO_LINE_INSUF_RESERVE", "No such line. Too little reserve to create it."},
	tecNO_LINE_REDUNDANT:                  {"tecNO_LINE_REDUNDANT", "Can't set non-existant line to default."},
	tecPATH_DRY:                           {"tecPATH_DRY", "Path could not send partial amount."},
	tecPATH_PARTIAL:                       {"tecPATH_PARTIAL", "Path could not send full amount."},
	tecNO_ALTERNATIVE_KEY:                 {"tecNO_ALTERNATIVE_KEY", "The operation would remove the ability to sign transactions with the account."},
	tecNO_REGULAR_KEY:                     {"tecNO_REGULAR_KEY", "Regular key is not set."},
	tecUNFUNDED:                           {"tecUNFUNDED", "One of _ADD, _OFFER, or _SEND. Deprecated."},
	tecUNFUNDED_ADD:                       {"tecUNFUNDED_ADD", "Insufficient XRP balance for WalletAdd."},
	tecUNFUNDED_OFFER:                     {"tecUNFUNDED_OFFER", "Insufficient balance to fund created offer."},
	tecUNFUNDED_PAYMENT:                   {"tecUNFUNDED_PAYMENT", "Insufficient XRP balance to send."},
	tecOWNERS:                             {"tecOWNERS", "Non-zero owner count."},
	tecNO_ISSUER:                          {"tecNO_ISSUER", "Issuer account does not exist."},
	tecNO_AUTH:                            {"tecNO_AUTH", "Not authorized to hold asset."},
	tecNO_LINE:                            {"tecNO_LINE", "No such line."},
	tecINSUFF_FEE:                         {"tecINSUFF_FEE", "Insufficient balance to pay fee."},
	tecFROZEN:                             {"tecFROZEN", "Asset is frozen."},
	tecNO_TARGET:                          {"tecNO_TARGET", "Target account does not exist."},
	tecNO_PERMISSION:                      {"tecNO_PERMISSION", "No permission to perform requested operation."},
	tecNO_ENTRY:                           {"tecNO_ENTRY", "No matching entry found."},
	tecINSUFFICIENT_RESERVE:               {"tecINSUFFICIENT_RESERVE", "Insufficient reserve to complete requested operation."},
	tecNEED_MASTER_KEY:                    {"tecNEED_MASTER_KEY", "The operation requires the use of the Master Key."},
	tecDST_TAG_NEEDED:                     {"tecDST_TAG_NEEDED", "A destination tag is required."},
	tecINTERNAL:                           {"tecINTERNAL", "An internal error has occurred during processing."},
	tecCRYPTOCONDITION_ERROR:              {"tecCRYPTOCONDITION_ERROR", "Malformed, invalid, or mismatched conditional or fulfillment."},
	tecINVARIANT_FAILED:                   {"tecINVARIANT_FAILED", "One or more invariants for the transaction were not satisfied."},
	tecOVERSIZE:                           {"tecOVERSIZE", "Object exceeded serialization limits"},
	tecEXPIRED:                            {"tecEXPIRED", "Expiration time is passed."},
	tecDUPLICATE:                          {"tecDUPLICATE", "Ledger object already exists."},
	tecKILLED:                             {"tecKILLED", "FillOrKill offer killed."},
	tecHAS_OBLIGATIONS:                    {"tecHAS_OBLIGATIONS", "The account cannot be deleted since it has obligations."},
	tecTOO_SOON:                           {"tecTOO_SOON", "It is too early to attempt the requested operation. Please wait."},
	tecMAX_SEQUENCE_REACHED:               {"tecMAX_SEQUENCE_REACHED", "The maximum sequence number was reached."},
	tecNO_SUITABLE_NFTOKEN_PAGE:           {"tecNO_SUITABLE_NFTOKEN_PAGE", "A suitable NFToken page could not be located."},
	tecNFTOKEN_BUY_SELL_MISMATCH:          {"tecNFTOKEN_BUY_SELL_MISMATCH", "The 'Buy' and 'Sell' NFToken offers are mismatched."},
	tecNFTOKEN_OFFER_TYPE_MISMATCH:        {"tecNFTOKEN_OFFER_TYPE_MISMATCH", "The type of NFToken offer is incorrect."},
	tecCANT_ACCEPT_OWN_NFTOKEN_OFFER:      {"tecCANT_ACCEPT_OWN_NFTOKEN_OFFER", "An NFToken offer cannot be claimed by its owner."},
	tecINSUFFICIENT_FUNDS:                 {"tecINSUFFICIENT_FUNDS", "Not enough funds available to complete requested transaction."},
	tecOBJECT_NOT_FOUND:                   {"tecOBJECT_NOT_FOUND", "A requested object could not be located."},
	tecINSUFFICIENT_PAYMENT:               {"tecINSUFFICIENT_PAYMENT", "The payment is not sufficient."},
	tecUNFUNDED_AMM:                       {"tecUNFUNDED_AMM", "Insufficient balance to fund AMM."},
	tecAMM_BALANCE:                        {"tecAMM_BALANCE", "AMM has invalid balance."},
	tecAMM_FAILED:                         {"tecAMM_FAILED", "AMM transaction failed."},
	tecAMM_INVALID_TOKENS:                 {"tecAMM_INVALID_TOKENS", "AMM invalid LP tokens."},
	tecAMM_EMPTY:                          {"tecAMM_EMPTY", "AMM is in empty state."},
	tecAMM_NOT_EMPTY:                      {"tecAMM_NOT_EMPTY", "AMM is not in empty state."},
	tecAMM_ACCOUNT:                        {"tecAMM_ACCOUNT", "This operation is not allowed on an AMM Account."},
	tecINCOMPLETE:                         {"tecINCOMPLETE", "Some work was completed, but more submissions required to finish."},
	tecXCHAIN_BAD_TRANSFER_ISSUE:          {"tecXCHAIN_BAD_TRANSFER_ISSUE", "Bad xchain transfer issue."},
	tecXCHAIN_NO_CLAIM_ID:                 {"tecXCHAIN_NO_CLAIM_ID", "No such xchain claim id."},
	tecXCHAIN_BAD_CLAIM_ID:                {"tecXCHAIN_BAD_CLAIM_ID", "Bad xchain claim id."},
	tecXCHAIN_CLAIM_NO_QUORUM:             {"tecXCHAIN_CLAIM_NO_QUORUM", "Quorum was not reached on the xchain claim."},
	tecXCHAIN_PROOF_UNKNOWN_KEY:           {"tecXCHAIN_PROOF_UNKNOWN_KEY", "Unknown key for the xchain proof."},
	tecXCHAIN_CREATE_ACCOUNT_NONXRP_ISSUE: {"tecXCHAIN_CREATE_ACCOUNT_NONXRP_ISSUE", "Only XRP may be used for xchain create account."},
	tecXCHAIN_WRONG_CHAIN:                 {"tecXCHAIN_WRONG_CHAIN", "XChain Transaction was submitted to the wrong chain."},
	tecXCHAIN_REWARD_MISMATCH:             {"tecXCHAIN_REWARD_MISMATCH", "The reward amount must match the reward specified in the xchain bridge."},
	tecXCHAIN_NO_SIGNERS_LIST:             {"tecXCHAIN_NO_SIGNERS_LIST", "The account did not have a signers list."},
	tecXCHAIN_SENDING_ACCOUNT_MISMATCH:    {"tecXCHAIN_SENDING_ACCOUNT_MISMATCH", "The sending account did not match the expected sending account."},
	tecXCHAIN_INSUFF_CREATE_AMOUNT:        {"tecXCHAIN_INSUFF_CREATE_AMOUNT", "Insufficient amount to create an account."},
	tecXCHAIN_ACCOUNT_CREATE_PAST:         {"tecXCHAIN_ACCOUNT_CREATE_PAST", "The account create count has already passed."},
	tecXCHAIN_ACCOUNT_CREATE_TOO_MANY:     {"tecXCHAIN_ACCOUNT_CREATE_TOO_MANY", "There are too many pending account create transactions to submit a new one."},
	tecXCHAIN_PAYMENT_FAILED:              {"tecXCHAIN_PAYMENT_FAILED", "Failed to transfer funds in a xchain transaction."},
	tecXCHAIN_SELF_COMMIT:                 {"tecXCHAIN_SELF_COMMIT", "Account cannot commit funds to itself."},
	tecXCHAIN_BAD_PUBLIC_KEY_ACCOUNT_PAIR: {"tecXCHAIN_BAD_PUBLIC_KEY_ACCOUNT_PAIR", "Bad public key account pair in an xchain transaction."},
	tecXCHAIN_CREATE_ACCOUNT_DISABLED:     {"tecXCHAIN_CREATE_ACCOUNT_DISABLED", "This bridge does not support account creation."},
	tecEMPTY_DID:                          {"tecEMPTY_DID", "The DID object did not have a URI or DIDDocument field."},
	tecINVALID_UPDATE_TIME:                {"tecINVALID_UPDATE_TIME", "The Oracle object has invalid LastUpdateTime field."},
	tecTOKEN_PAIR_NOT_FOUND:               {"tecTOKEN_PAIR_NOT_FOUND", "Token pair is not found in Oracle object."},
	tecARRAY_EMPTY:                        {"tecARRAY_EMPTY", "Array is empty."},
	tecARRAY_TOO_LARGE:                    {"tecARRAY_TOO_LARGE", "Array is too large."},

	tefFAILURE:                     {"tefFAILURE", "Failed to apply."},
	tefALREADY:                     {"tefALREADY", "The exact transaction was already in this ledger."},
	tefBAD_ADD_AUTH:                {"tefBAD_ADD_AUTH", "Not authorized to add account."},
	tefBAD_AUTH:                    {"tefBAD_AUTH", "Transaction's public key is not authorized."},
	tefBAD_CLAIM_ID:                {"tefBAD_CLAIM_ID", "Malformed: Bad claim id."},
	tefBAD_GEN_AUTH:                {"tefBAD_GEN_AUTH", "Not authorized to claim generator."},
	tefBAD_LEDGER:                  {"tefBAD_LEDGER", "Ledger in unexpected state."},
	tefCLAIMED:                     {"tefCLAIMED", "Can not claim a previously claimed account."},
	tefCREATED:                     {"tefCREATED", "Can't add an already created account."},
	tefDST_TAG_NEEDED:              {"tefDST_TAG_NEEDED", "Destination tag required."},
	tefEXCEPTION:                   {"tefEXCEPTION", "Unexpected program state."},
	tefGEN_IN_USE:                  {"tefGEN_IN_USE", "Generator already in use."},
	tefINTERNAL:                    {"tefINTERNAL", "Internal error."},
	tefNO_AUTH_REQUIRED:            {"tefNO_AUTH_REQUIRED", "Auth is not required."},
	tefPAST_SEQ:                    {"tefPAST_SEQ", "This sequence number has already past."},
	tefWRONG_PRIOR:                 {"tefWRONG_PRIOR", "This previous transaction does not match."},
	tefMASTER_DISABLED:             {"tefMASTER_DISABLED", "Master key is disabled."},
	tefMAX_LEDGER:                  {"tefMAX_LEDGER", "Ledger sequence too high."},
	tefBAD_SIGNATURE:               {"tefBAD_SIGNATURE", "The transaction was multi-signed, but contained a signature for an address not part of a SignerList associated with the sending account."},
	tefBAD_AUTH_MASTER:             {"tefBAD_AUTH_MASTER", "Auth for unclaimed account needs correct master key."},
	tefINVARIANT_FAILED:            {"tefINVARIANT_FAILED", "Fee claim violated invariants for the transaction."},
	tefTOO_BIG:                     {"tefTOO_BIG", "Transaction affects too many items."},
	tefNO_TICKET:                   {"tefNO_TICKET", "Ticket is not in ledger."},
	tefNFTOKEN_IS_NOT_TRANSFERABLE: {"tefNFTOKEN_IS_NOT_TRANSFERABLE", "The specified NFToken is not transferable."},

	telLOCAL_ERROR:           {"telLOCAL_ERROR", "Local failure."},
	telBAD_DOMAIN:            {"telBAD_DOMAIN", "Domain too long."},
	telBAD_PATH_COUNT:        {"telBAD_PATH_COUNT", "Malformed: Too many paths."},
	telBAD_PUBLIC_KEY:        {"telBAD_PUBLIC_KEY", "Public key too long."},
	telFAILED_PROCESSING:     {"telFAILED_PROCESSING", "Failed to correctly process transaction."},
	telINSUF_FEE_P:           {"telINSUF_FEE_P", "Fee insufficient."},
	telNO_DST_PARTIAL:        {"telNO_DST_PARTIAL", "Partial payment to create account not allowed."},
	telCAN_NOT_QUEUE:         {"telCAN_NOT_QUEUE", "Can not queue at this time."},
	telCAN_NOT_QUEUE_BALANCE: {"telCAN_NOT_QUEUE_BALANCE", "Can not queue at this time: insufficient balance to pay all queued fees."},
	telCAN_NOT_QUEUE_BLOCKS:  {"telCAN_NOT_QUEUE_BLOCKS", "Can not queue at this time: would block later queued transaction(s)."},
	telCAN_NOT_QUEUE_BLOCKED: {"telCAN_NOT_QUEUE_BLOCKED", "Can not queue at this time: blocking transaction in queue."},
	telCAN_NOT_QUEUE_FEE:     {"telCAN_NOT_QUEUE_FEE", "Can not queue at this time: fee insufficient to replace queued transaction."},
	telCAN_NOT_QUEUE_FULL:    {"telCAN_NOT_QUEUE_FULL", "Can not queue at this time: queue is full."},

	temMALFORMED:                   {"temMALFORMED", "Malformed transaction."},
	temBAD_AMOUNT:                  {"temBAD_AMOUNT", "Can only send positive amounts."},
	temBAD_CURRENCY:                {"temBAD_CURRENCY", "Malformed: Bad currency."},
	temBAD_FEE:                     {"temBAD_FEE", "Invalid fee, negative or not XRP."},
	temBAD_EXPIRATION:              {"temBAD_EXPIRATION", "Malformed: Bad expiration."},
	temBAD_ISSUER:                  {"temBAD_ISSUER", "Malformed: Bad issuer."},
	temBAD_LIMIT:                   {"temBAD_LIMIT", "Limits must be non-negative."},
	temBAD_OFFER:                   {"temBAD_OFFER", "Malformed: Bad offer."},
	temBAD_PATH:                    {"temBAD_PATH", "Malformed: Bad path."},
	temBAD_PATH_LOOP:               {"temBAD_PATH_LOOP", "Malformed: Loop in path."},
	temBAD_SIGNATURE:               {"temBAD_SIGNATURE", "Malformed: Bad signature."},
	temBAD_SRC_ACCOUNT:             {"temBAD_SRC_ACCOUNT", "Malformed: Bad source account."},
	temBAD_TRANSFER_RATE:           {"temBAD_TRANSFER_RATE", "Malformed: Transfer rate must be >= 1.0"},
	temBAD_SEQUENCE:                {"temBAD_SEQUENCE", "Malformed: Sequence is not in the past."},
	temBAD_SEND_XRP_LIMIT:          {"temBAD_SEND_XRP_LIMIT", "Malformed: Limit quality is not allowed for XRP to XRP."},
	temBAD_SEND_XRP_MAX:            {"temBAD_SEND_XRP_MAX", "Malformed: Send max is not allowed for XRP to XRP."},
	temBAD_SEND_XRP_NO_DIRECT:      {"temBAD_SEND_XRP_NO_DIRECT", "Malformed: No Ripple direct is not allowed for XRP to XRP."},
	temBAD_SEND_XRP_PARTIAL:        {"temBAD_SEND_XRP_PARTIAL", "Malformed: Partial payment is not allowed for XRP to XRP."},
	temBAD_SEND_XRP_PATHS:          {"temBAD_SEND_XRP_PATHS", "Malformed: Paths are not allowed for XRP to XRP."},
	temDST_IS_SRC:                  {"temDST_IS_SRC", "Destination may not be source."},
	temDST_NEEDED:                  {"temDST_NEEDED", "Destination not specified."},
	temINVALID:                     {"temINVALID", "The transaction is ill-formed."},
	temINVALID_FLAG:                {"temINVALID_FLAG", "The transaction has an invalid flag."},
	temREDUNDANT:                   {"temREDUNDANT", "Sends same currency to self."},
	temRIPPLE_EMPTY:                {"temRIPPLE_EMPTY", "PathSet with no paths."},
	temUNCERTAIN:                   {"temUNCERTAIN", "In process of determining result. Never returned."},
	temUNKNOWN:                     {"temUNKNOWN", "The transactions requires logic not implemented yet."},
	temDISABLED:                    {"temDISABLED", "The transaction requires logic that is currently disabled."},
	temBAD_TICK_SIZE:               {"temBAD_TICK_SIZE", "Malformed: Tick size out of range."},
	temINVALID_ACCOUNT_ID:          {"temINVALID_ACCOUNT_ID", "Malformed: A field contains an invalid account ID."},
	temCANNOT_PREAUTH_SELF:         {"temCANNOT_PREAUTH_SELF", "Malformed: An account may not preauthorize itself."},
	temSEQ_AND_TICKET:              {"temSEQ_AND_TICKET", "Transaction contains a TicketSequence and a non-zero Sequence"},
	temBAD_NFTOKEN_TRANSFER_FEE:    {"temBAD_NFTOKEN_TRANSFER_FEE", "Malformed: The NFToken transfer fee must be between 1 and 5000, inclusive."},
	temBAD_WEIGHT:                  {"temBAD_WEIGHT", "The SignerListSet transaction includes a SignerWeight that is invalid, for example a zero or negative value."},
	temBAD_AMM_TOKENS:              {"temBAD_AMM_TOKENS", ""},
	temXCHAIN_EQUAL_DOOR_ACCOUNTS:  {"temXCHAIN_EQUAL_DOOR_ACCOUNTS", ""},
	temXCHAIN_BAD_PROOF:            {"temXCHAIN_BAD_PROOF", ""},
	temXCHAIN_BRIDGE_BAD_ISSUES:    {"temXCHAIN_BRIDGE_BAD_ISSUES", ""},
	temXCHAIN_BRIDGE_NONDOOR_OWNER: {"temXCHAIN_BRIDGE_NONDOOR_OWNER", ""},
	temXCHAIN_BRIDGE_BAD_MIN_ACCOUNT_CREATE_AMOUNT: {"temXCHAIN_BRIDGE_BAD_MIN_ACCOUNT_CREATE_AMOUNT", ""},
	temXCHAIN_BRIDGE_BAD_REWARD_AMOUNT:             {"temXCHAIN_BRIDGE_BAD_REWARD_AMOUNT", ""},
	temEMPTY_DID:                                   {"temEMPTY_DID", ""},
	temARRAY_EMPTY:                                 {"temARRAY_EMPTY", ""},
	temARRAY_TOO_LARGE:                             {"temARRAY_TOO_LARGE", ""},

	terRETRY:       {"terRETRY", "Retry transaction."},
	terFUNDS_SPENT: {"terFUNDS_SPENT", "Can't set password, password set funds already spent."},
	terINSUF_FEE_B: {"terINSUF_FEE_B", "Account balance can't pay fee."},
	terLAST:        {"terLAST", "Process last."},
	terNO_RIPPLE:   {"terNO_RIPPLE", "Path does not permit rippling."},
	terNO_ACCOUNT:  {"terNO_ACCOUNT", "The source account does not exist."},
	terNO_AUTH:     {"terNO_AUTH", "Not authorized to hold IOUs."},
	terNO_LINE:     {"terNO_LINE", "No such line."},
	terPRE_SEQ:     {"terPRE_SEQ", "Missing/inapplicable prior transaction."},
	terOWNERS:      {"terOWNERS", "Non-zero owner count."},
	terQUEUED:      {"terQUEUED", "Held until escalated fee drops."},
	terPRE_TICKET:  {"terPRE_TICKET", "Ticket is not yet in ledger."},
	terNO_AMM:      {"terNO_AMM", "AMM doesn't exist for the asset pair."},
}

var reverseResults map[string]TransactionResult

func init() {
	reverseResults = make(map[string]TransactionResult)
	for result, name := range resultNames {
		reverseResults[name.Token] = result
	}
}

func (r TransactionResult) String() string {
	return resultNames[r].Token
}

func (r TransactionResult) Human() string {
	return resultNames[r].Human
}

func (r TransactionResult) Success() bool {
	return r == tesSUCCESS
}

func (r TransactionResult) Queued() bool {
	return r == terQUEUED
}

func (r TransactionResult) Symbol() string {
	switch r {
	case tesSUCCESS, tecCLAIM:
		return "✓"
	case tecPATH_PARTIAL, tecPATH_DRY:
		return "½"
	case tecUNFUNDED, tecUNFUNDED_ADD, tecUNFUNDED_OFFER, tecUNFUNDED_PAYMENT:
		return "$"
	default:
		return "✗"
	}
}
