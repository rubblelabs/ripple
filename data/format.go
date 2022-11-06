package data

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type NodeType uint8
type NodeFormat uint8
type HashPrefix uint32
type LedgerNamespace uint16

const (
	// Hash Prefixes
	HP_TRANSACTION_ID   HashPrefix = 0x54584E00 // 'TXN' transaction
	HP_TRANSACTION_NODE HashPrefix = 0x534E4400 // 'SND' transaction plus metadata (probably should have been TND!)
	HP_LEAF_NODE        HashPrefix = 0x4D4C4E00 // 'MLN' account state
	HP_INNER_NODE       HashPrefix = 0x4D494E00 // 'MIN' inner node in tree
	HP_LEDGER_MASTER    HashPrefix = 0x4C575200 // 'LWR' ledger master data for signing (probably should have been LGR!)
	HP_TRANSACTION_SIGN HashPrefix = 0x53545800 // 'STX' inner transaction to sign
	HP_VALIDATION       HashPrefix = 0x56414C00 // 'VAL' validation for signing
	HP_PROPOSAL         HashPrefix = 0x50525000 // 'PRP' proposal for signing

	// Node Types
	NT_UNKNOWN          NodeType = 0
	NT_LEDGER           NodeType = 1
	NT_TRANSACTION      NodeType = 2
	NT_ACCOUNT_NODE     NodeType = 3
	NT_TRANSACTION_NODE NodeType = 4

	// Node Formats
	NF_PREFIX NodeFormat = 1
	NF_HASH   NodeFormat = 2
	NF_WIRE   NodeFormat = 3

	// Ledger index NameSpaces
	NS_ACCOUNT         LedgerNamespace = 'a'
	NS_DIRECTORY_NODE  LedgerNamespace = 'd'
	NS_RIPPLE_STATE    LedgerNamespace = 'r'
	NS_OFFER           LedgerNamespace = 'o' // Entry for an offer
	NS_OWNER_DIRECTORY LedgerNamespace = 'O' // Directory of things owned by an account
	NS_BOOK_DIRECTORY  LedgerNamespace = 'B' // Directory of order books
	NS_SKIP_LIST       LedgerNamespace = 's'
	NS_AMENDMENT       LedgerNamespace = 'f'
	NS_FEE             LedgerNamespace = 'e'
	NS_SUSPAY          LedgerNamespace = 'u'
	NS_TICKET          LedgerNamespace = 'T'
	NS_SIGNER_LIST     LedgerNamespace = 'S'
	NS_XRPU_CHANNEL    LedgerNamespace = 'x'
	NS_CHECK           LedgerNamespace = 'C'
	NS_DEPOSIT_PREAUTH LedgerNamespace = 'p'
	NS_NEGATIVE_UNL    LedgerNamespace = 'N'
)

var nodeTypes = [...]string{
	NT_UNKNOWN:          "Unknown",
	NT_LEDGER:           "Ledger",
	NT_TRANSACTION:      "Transaction",
	NT_ACCOUNT_NODE:     "Account Node",
	NT_TRANSACTION_NODE: "Transaction Node",
}

type NodeHeader struct {
	LedgerSequence uint32
	_              uint32 //padding for repeated LedgerIndex
	NodeType       NodeType
}

type enc struct {
	typ, field uint8
}

const (
	ST_UINT16    uint8 = 1
	ST_UINT32    uint8 = 2
	ST_UINT64    uint8 = 3
	ST_HASH128   uint8 = 4
	ST_HASH256   uint8 = 5
	ST_AMOUNT    uint8 = 6
	ST_VL        uint8 = 7
	ST_ACCOUNT   uint8 = 8
	ST_OBJECT    uint8 = 14
	ST_ARRAY     uint8 = 15
	ST_UINT8     uint8 = 16
	ST_HASH160   uint8 = 17
	ST_PATHSET   uint8 = 18
	ST_VECTOR256 uint8 = 19
	ST_HASH96    uint8 = 20
	ST_HASH192   uint8 = 21
	ST_HASH384   uint8 = 22
	ST_HASH512   uint8 = 23
)

// See rippled's SField.cpp for the strings and corresponding encoding values.
var encodings = map[enc]string{
	// 16-bit unsigned integers (common)
	{ST_UINT16, 1}: "LedgerEntryType",
	{ST_UINT16, 2}: "TransactionType",
	{ST_UINT16, 3}: "SignerWeight",
	{ST_UINT16, 4}: "TransferFee",
	// 16-bit unsigned integers (uncommon)
	{ST_UINT16, 16}: "Version",
	// 32-bit unsigned integers (common)
	{ST_UINT32, 2}:  "Flags",
	{ST_UINT32, 3}:  "SourceTag",
	{ST_UINT32, 4}:  "Sequence",
	{ST_UINT32, 5}:  "PreviousTxnLgrSeq",
	{ST_UINT32, 6}:  "LedgerSequence",
	{ST_UINT32, 7}:  "CloseTime",
	{ST_UINT32, 8}:  "ParentCloseTime",
	{ST_UINT32, 9}:  "SigningTime",
	{ST_UINT32, 10}: "Expiration",
	{ST_UINT32, 11}: "TransferRate",
	{ST_UINT32, 12}: "WalletSize",
	{ST_UINT32, 13}: "OwnerCount",
	{ST_UINT32, 14}: "DestinationTag",
	// 32-bit unsigned integers (uncommon)
	{ST_UINT32, 16}: "HighQualityIn",
	{ST_UINT32, 17}: "HighQualityOut",
	{ST_UINT32, 18}: "LowQualityIn",
	{ST_UINT32, 19}: "LowQualityOut",
	{ST_UINT32, 20}: "QualityIn",
	{ST_UINT32, 21}: "QualityOut",
	{ST_UINT32, 22}: "StampEscrow",
	{ST_UINT32, 23}: "BondAmount",
	{ST_UINT32, 24}: "LoadFee",
	{ST_UINT32, 25}: "OfferSequence",
	{ST_UINT32, 26}: "FirstLedgerSequence",
	{ST_UINT32, 27}: "LastLedgerSequence",
	{ST_UINT32, 28}: "TransactionIndex",
	{ST_UINT32, 29}: "OperationLimit",
	{ST_UINT32, 30}: "ReferenceFeeUnits",
	{ST_UINT32, 31}: "ReserveBase",
	{ST_UINT32, 32}: "ReserveIncrement",
	{ST_UINT32, 33}: "SetFlag",
	{ST_UINT32, 34}: "ClearFlag",
	{ST_UINT32, 35}: "SignerQuorum",
	{ST_UINT32, 36}: "CancelAfter",
	{ST_UINT32, 37}: "FinishAfter",
	{ST_UINT32, 38}: "SignerListID",
	{ST_UINT32, 39}: "SettleDelay",
	{ST_UINT32, 40}: "TicketCount",
	{ST_UINT32, 41}: "TicketSequence",
	{ST_UINT32, 42}: "NFTokenTaxon",
	{ST_UINT32, 43}: "MintedNFTokens",
	{ST_UINT32, 44}: "BurnedNFTokens",
	// 64-bit unsigned integers (common)
	{ST_UINT64, 1}:  "IndexNext",
	{ST_UINT64, 2}:  "IndexPrevious",
	{ST_UINT64, 3}:  "BookNode",
	{ST_UINT64, 4}:  "OwnerNode",
	{ST_UINT64, 5}:  "BaseFee",
	{ST_UINT64, 6}:  "ExchangeRate",
	{ST_UINT64, 7}:  "LowNode",
	{ST_UINT64, 8}:  "HighNode",
	{ST_UINT64, 9}:  "DestinationNode",
	{ST_UINT64, 10}: "Cookie",
	{ST_UINT64, 11}: "ServerVersion",
	{ST_UINT64, 12}: "NFTokenOfferNode",
	// 128-bit (common)
	{ST_HASH128, 1}: "EmailHash",
	// 256-bit (common)
	{ST_HASH256, 1}:  "LedgerHash",
	{ST_HASH256, 2}:  "ParentHash",
	{ST_HASH256, 3}:  "TransactionHash",
	{ST_HASH256, 4}:  "AccountHash",
	{ST_HASH256, 5}:  "PreviousTxnID",
	{ST_HASH256, 6}:  "LedgerIndex",
	{ST_HASH256, 7}:  "WalletLocator",
	{ST_HASH256, 8}:  "RootIndex",
	{ST_HASH256, 9}:  "AccountTxnID",
	{ST_HASH256, 10}: "NFTokenID",
	// 256-bit (uncommon)
	{ST_HASH256, 16}: "BookDirectory",
	{ST_HASH256, 17}: "InvoiceID",
	{ST_HASH256, 18}: "Nickname",
	{ST_HASH256, 19}: "Amendment",
	{ST_HASH256, 20}: "TicketID",
	{ST_HASH256, 21}: "Digest",
	{ST_HASH256, 22}: "Channel",
	{ST_HASH256, 24}: "CheckID",
	{ST_HASH256, 25}: "ValidatedHash",
	{ST_HASH256, 26}: "PreviousPageMin",
	{ST_HASH256, 27}: "NextPageMin",
	{ST_HASH256, 28}: "NFTokenBuyOffer",
	{ST_HASH256, 29}: "NFTokenSellOffer",
	// currency amount (common)
	{ST_AMOUNT, 1}:  "Amount",
	{ST_AMOUNT, 2}:  "Balance",
	{ST_AMOUNT, 3}:  "LimitAmount",
	{ST_AMOUNT, 4}:  "TakerPays",
	{ST_AMOUNT, 5}:  "TakerGets",
	{ST_AMOUNT, 6}:  "LowLimit",
	{ST_AMOUNT, 7}:  "HighLimit",
	{ST_AMOUNT, 8}:  "Fee",
	{ST_AMOUNT, 9}:  "SendMax",
	{ST_AMOUNT, 10}: "DeliverMin",
	// currency amount (uncommon)
	{ST_AMOUNT, 16}: "MinimumOffer",
	{ST_AMOUNT, 17}: "RippleEscrow",
	{ST_AMOUNT, 18}: "DeliveredAmount",
	{ST_AMOUNT, 19}: "NFTokenBrokerFee",
	// variable length (common)
	{ST_VL, 1}:  "PublicKey",
	{ST_VL, 2}:  "MessageKey",
	{ST_VL, 3}:  "SigningPubKey",
	{ST_VL, 4}:  "TxnSignature",
	{ST_VL, 5}:  "URI",
	{ST_VL, 6}:  "Signature",
	{ST_VL, 7}:  "Domain",
	{ST_VL, 8}:  "FundCode",
	{ST_VL, 9}:  "RemoveCode",
	{ST_VL, 10}: "ExpireCode",
	{ST_VL, 11}: "CreateCode",
	{ST_VL, 12}: "MemoType",
	{ST_VL, 13}: "MemoData",
	{ST_VL, 14}: "MemoFormat",
	// variable length (uncommon)
	{ST_VL, 16}: "Fulfillment",
	{ST_VL, 17}: "Condition",
	{ST_VL, 18}: "MasterSignature",
	{ST_VL, 19}: "UNLModifyValidator",
	{ST_VL, 20}: "ValidatorToDisable",
	{ST_VL, 21}: "ValidatorToReEnable",
	// account
	{ST_ACCOUNT, 1}: "Account",
	{ST_ACCOUNT, 2}: "Owner",
	{ST_ACCOUNT, 3}: "Destination",
	{ST_ACCOUNT, 4}: "Issuer",
	{ST_ACCOUNT, 5}: "Authorize",
	{ST_ACCOUNT, 6}: "Unauthorize",
	{ST_ACCOUNT, 7}: "Target",
	{ST_ACCOUNT, 8}: "RegularKey",
	{ST_ACCOUNT, 9}: "NFTokenMinter",
	// inner object
	{ST_OBJECT, 1}:  "EndOfObject",
	{ST_OBJECT, 2}:  "TransactionMetaData",
	{ST_OBJECT, 3}:  "CreatedNode",
	{ST_OBJECT, 4}:  "DeletedNode",
	{ST_OBJECT, 5}:  "ModifiedNode",
	{ST_OBJECT, 6}:  "PreviousFields",
	{ST_OBJECT, 7}:  "FinalFields",
	{ST_OBJECT, 8}:  "NewFields",
	{ST_OBJECT, 9}:  "TemplateEntry",
	{ST_OBJECT, 10}: "Memo",
	{ST_OBJECT, 11}: "SignerEntry",
	{ST_OBJECT, 12}: "NFToken",
	// inner object (uncommon)
	{ST_OBJECT, 16}: "Signer",
	{ST_OBJECT, 18}: "Majority",
	{ST_OBJECT, 19}: "DisabledValidator",
	// array of objects
	{ST_ARRAY, 1}:  "EndOfArray",
	{ST_ARRAY, 2}:  "SigningAccounts",
	{ST_ARRAY, 3}:  "Signers",
	{ST_ARRAY, 4}:  "SignerEntries",
	{ST_ARRAY, 5}:  "Template",
	{ST_ARRAY, 6}:  "Necessary",
	{ST_ARRAY, 7}:  "Sufficient",
	{ST_ARRAY, 8}:  "AffectedNodes",
	{ST_ARRAY, 9}:  "Memos",
	{ST_ARRAY, 10}: "NFTokens",
	// array of objects (uncommon)
	{ST_ARRAY, 16}: "Majorities",
	{ST_ARRAY, 17}: "DisabledValidators",
	// 8-bit unsigned integers (common)
	{ST_UINT8, 1}: "CloseResolution",
	{ST_UINT8, 2}: "Method",
	{ST_UINT8, 3}: "TransactionResult",
	// 8-bit unsigned integers (uncommon)
	{ST_UINT8, 16}: "TickSize",
	{ST_UINT8, 17}: "UNLModifyDisabling",
	// 160-bit (common)
	{ST_HASH160, 1}: "TakerPaysCurrency",
	{ST_HASH160, 2}: "TakerPaysIssuer",
	{ST_HASH160, 3}: "TakerGetsCurrency",
	{ST_HASH160, 4}: "TakerGetsIssuer",
	// path set
	{ST_PATHSET, 1}: "Paths",
	// vector of 256-bit
	{ST_VECTOR256, 1}: "Indexes",
	{ST_VECTOR256, 2}: "Hashes",
	{ST_VECTOR256, 3}: "Amendments",
	{ST_VECTOR256, 4}: "NFTokenOffers",
}

var reverseEncodings map[string]enc
var signingFields map[enc]struct{}

func init() {
	reverseEncodings = make(map[string]enc)
	signingFields = make(map[enc]struct{})
	for e, name := range encodings {
		reverseEncodings[name] = e
		if strings.Contains(name, "Signature") {
			signingFields[e] = struct{}{}
		}
	}
}

func (h HashPrefix) String() string {
	return string(h.Bytes())
}

func (h HashPrefix) Bytes() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(h))
	return b
}

func (n NodeType) String() string {
	return nodeTypes[n]
}

func (e enc) Priority() uint32 {
	return uint32(e.typ)<<16 | uint32(e.field)
}

func (e enc) SigningField() bool {
	_, ok := signingFields[e]
	return ok
}

func readEncoding(r Reader) (*enc, error) {
	var e enc
	if b, err := r.ReadByte(); err != nil {
		return nil, err
	} else {
		e.typ = b >> 4
		e.field = b & 0xF
	}
	var err error
	if e.typ == 0 {
		if e.typ, err = r.ReadByte(); err != nil {
			return nil, err
		}
	}
	if e.field == 0 {
		if e.field, err = r.ReadByte(); err != nil {
			return nil, err
		}
	}
	return &e, nil
}

func writeEncoding(w io.Writer, e enc) error {
	var err error
	switch {
	case e.typ < 16 && e.field < 16:
		_, err = w.Write([]uint8{e.typ<<4 | e.field})
	case e.typ < 16:
		_, err = w.Write([]uint8{e.typ << 4, e.field})
	case e.field < 16:
		_, err = w.Write([]uint8{e.field, e.typ})
	default:
		_, err = w.Write([]uint8{0, e.typ, e.field})
	}
	return err
}

func write(w io.Writer, v interface{}) error {
	return binary.Write(w, binary.BigEndian, v)
}

func writeValues(w io.Writer, values []interface{}) error {
	for _, v := range values {
		if err := binary.Write(w, binary.BigEndian, v); err != nil {
			return err
		}
	}
	return nil
}

func read(r Reader, dest interface{}) error {
	return binary.Read(r, binary.BigEndian, dest)
}

func writeVariableLength(w io.Writer, b []byte) error {
	n := len(b)
	var err error
	switch {
	case n < 0 || n > 918744:
		return fmt.Errorf("Unsupported Variable Length encoding: %d", n)
	case n <= 192:
		_, err = w.Write([]uint8{uint8(n)})
	case n <= 12480:
		n -= 193
		_, err = w.Write([]uint8{193 + uint8(n>>8), uint8(n)})
	case n <= 918744:
		n -= 12481
		_, err = w.Write([]uint8{241 + uint8(n>>16), uint8(n >> 8), uint8(n)})
	}
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func readVariableLength(r Reader) (int, error) {
	var first, second, third byte
	var err error
	if first, err = r.ReadByte(); err != nil {
		return 0, err
	}
	switch {
	case first <= 192:
		return int(first), nil
	case first <= 240:
		if second, err = r.ReadByte(); err != nil {
			return 0, nil
		}
		return 193 + int(first-193)*256 + int(second), nil
	case first <= 254:
		if second, err = r.ReadByte(); err != nil {
			return 0, nil
		}
		if third, err = r.ReadByte(); err != nil {
			return 0, nil
		}
		return 12481 + int(first-241)*65536 + int(second)*256 + int(third), nil
	}
	return 0, fmt.Errorf("Unsupported Variable Length encoding")
}

func unmarshalSlice(s []byte, r Reader, prefix string) error {
	n, err := r.Read(s)
	if n != len(s) {
		return fmt.Errorf("%s: short read: %d expected: %d", prefix, n, len(s))
	}
	if err != nil {
		return fmt.Errorf("%s: %s", prefix, err.Error())
	}
	return nil
}
