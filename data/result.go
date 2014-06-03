package data

type TransactionResult uint8

const (
	tesSUCCESS               TransactionResult = 0
	tecCLAIM                 TransactionResult = 100
	tecPATH_PARTIAL          TransactionResult = 101
	tecUNFUNDED_ADD          TransactionResult = 102
	tecUNFUNDED_OFFER        TransactionResult = 103
	tecUNFUNDED_PAYMENT      TransactionResult = 104
	tecFAILED_PROCESSING     TransactionResult = 105
	tecDIR_FULL              TransactionResult = 121
	tecINSUF_RESERVE_LINE    TransactionResult = 122
	tecINSUF_RESERVE_OFFER   TransactionResult = 123
	tecNO_DST                TransactionResult = 124
	tecNO_DST_INSUF_XRP      TransactionResult = 125
	tecNO_LINE_INSUF_RESERVE TransactionResult = 126
	tecNO_LINE_REDUNDANT     TransactionResult = 127
	tecPATH_DRY              TransactionResult = 128
	tecUNFUNDED              TransactionResult = 129
	tecMASTER_DISABLED       TransactionResult = 130
	tecNO_REGULAR_KEY        TransactionResult = 131
	tecOWNERS                TransactionResult = 132
)

var resultNames = map[TransactionResult]string{
	tesSUCCESS:               "tesSUCCESS",
	tecCLAIM:                 "tecCLAIM",
	tecPATH_PARTIAL:          "tecPATH_PARTIAL",
	tecUNFUNDED_ADD:          "tecUNFUNDED_ADD",
	tecUNFUNDED_OFFER:        "tecUNFUNDED_OFFER",
	tecUNFUNDED_PAYMENT:      "tecUNFUNDED_PAYMENT",
	tecFAILED_PROCESSING:     "tecFAILED_PROCESSING",
	tecDIR_FULL:              "tecDIR_FULL",
	tecINSUF_RESERVE_LINE:    "tecINSUF_RESERVE_LINE",
	tecINSUF_RESERVE_OFFER:   "tecINSUF_RESERVE_OFFER",
	tecNO_DST:                "tecNO_DST",
	tecNO_DST_INSUF_XRP:      "tecNO_DST_INSUF_XRP",
	tecNO_LINE_INSUF_RESERVE: "tecNO_LINE_INSUF_RESERVE",
	tecNO_LINE_REDUNDANT:     "tecNO_LINE_REDUNDANT",
	tecPATH_DRY:              "tecPATH_DRY",
	tecUNFUNDED:              "tecUNFUNDED",
	tecMASTER_DISABLED:       "tecMASTER_DISABLED",
	tecNO_REGULAR_KEY:        "tecNO_REGULAR_KEY",
	tecOWNERS:                "tecOWNERS",
}

var reverseResults map[string]TransactionResult

func init() {
	reverseResults = make(map[string]TransactionResult)
	for result, name := range resultNames {
		reverseResults[name] = result
	}
}
func (r TransactionResult) String() string {
	return resultNames[r]
}
