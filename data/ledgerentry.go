package data

type LedgerEntrySlice []LedgerEntry

type leBase struct {
	LedgerEntryType   LedgerEntryType
	LedgerIndex       *Hash256 `json:"index,omitempty"`
	PreviousTxnID     *Hash256 `json:",omitempty"`
	PreviousTxnLgrSeq *uint32  `json:",omitempty"`
	Hash              Hash256  `json:"-"`
	Id                Hash256  `json:"-"`
}

type AccountRoot struct {
	leBase
	Flags          *LedgerEntryFlag `json:",omitempty"`
	Account        *Account         `json:",omitempty"`
	Sequence       *uint32          `json:",omitempty"`
	Balance        *Value           `json:",omitempty"`
	OwnerCount     *uint32          `json:",omitempty"`
	AccountTxnID   *Hash256         `json:",omitempty"`
	RegularKey     *RegularKey      `json:",omitempty"`
	EmailHash      *Hash128         `json:",omitempty"`
	WalletLocator  *Hash256         `json:",omitempty"`
	WalletSize     *uint32          `json:",omitempty"`
	MessageKey     *VariableLength  `json:",omitempty"`
	TransferRate   *uint32          `json:",omitempty"`
	Domain         *VariableLength  `json:",omitempty"`
	TickSize       *uint8           `json:",omitempty"`
	TicketCount    *uint32          `json:",omitempty"`
	NFTokenMinter  *Account         `json:",omitempty"`
	MintedNFTokens *uint32          `json:",omitempty"`
	BurnedNFTokens *uint32          `json:",omitempty"`
	AMMID          *Hash256         `json:",omitempty"`
}

type RippleState struct {
	leBase
	Flags          *LedgerEntryFlag `json:",omitempty"`
	LowLimit       *Amount          `json:",omitempty"`
	HighLimit      *Amount          `json:",omitempty"`
	Balance        *Amount          `json:",omitempty"`
	LowNode        *NodeIndex       `json:",omitempty"`
	HighNode       *NodeIndex       `json:",omitempty"`
	LowQualityIn   *uint32          `json:",omitempty"`
	LowQualityOut  *uint32          `json:",omitempty"`
	HighQualityIn  *uint32          `json:",omitempty"`
	HighQualityOut *uint32          `json:",omitempty"`
}

type Offer struct {
	leBase
	Flags         *LedgerEntryFlag `json:",omitempty"`
	Account       *Account         `json:",omitempty"`
	Sequence      *uint32          `json:",omitempty"`
	TakerPays     *Amount          `json:",omitempty"`
	TakerGets     *Amount          `json:",omitempty"`
	BookDirectory *Hash256         `json:",omitempty"`
	BookNode      *NodeIndex       `json:",omitempty"`
	OwnerNode     *NodeIndex       `json:",omitempty"`
	Expiration    *uint32          `json:",omitempty"`
}

type Directory struct {
	leBase
	Flags             *LedgerEntryFlag `json:",omitempty"`
	RootIndex         *Hash256         `json:",omitempty"`
	Indexes           *Vector256       `json:",omitempty"`
	Owner             *Account         `json:",omitempty"`
	TakerPaysCurrency *Hash160         `json:",omitempty"`
	TakerPaysIssuer   *Hash160         `json:",omitempty"`
	TakerGetsCurrency *Hash160         `json:",omitempty"`
	TakerGetsIssuer   *Hash160         `json:",omitempty"`
	ExchangeRate      *ExchangeRate    `json:",omitempty"`
	IndexNext         *NodeIndex       `json:",omitempty"`
	IndexPrevious     *NodeIndex       `json:",omitempty"`
	NFTokenID         *Hash256         `json:",omitempty"`
}

type LedgerHashes struct {
	leBase
	Flags               *LedgerEntryFlag `json:",omitempty"`
	FirstLedgerSequence *uint32          `json:",omitempty"`
	LastLedgerSequence  *uint32          `json:",omitempty"`
	Hashes              *Vector256       `json:",omitempty"`
}

type Majority struct {
	Amendment *Hash256 `json:",omitempty"`
	CloseTime *uint32  `json:",omitempty"`
}

type Amendments struct {
	leBase
	Flags      *LedgerEntryFlag `json:",omitempty"`
	Amendments *Vector256       `json:",omitempty"`
	Majorities []Majority       `json:",omitempty"`
}

type FeeSettings struct {
	leBase
	Flags             *LedgerEntryFlag `json:",omitempty"`
	BaseFee           *Uint64Hex       `json:",omitempty"`
	ReferenceFeeUnits *uint32          `json:",omitempty"`
	ReserveBase       *uint32          `json:",omitempty"`
	ReserveIncrement  *uint32          `json:",omitempty"`
}

type DisabledValidator struct {
	FirstLedgerSequence *uint32    `json:",omitempty"`
	PublicKey           *PublicKey `json:",omitempty"`
}

type NegativeUNL struct {
	leBase
	Flags               *LedgerEntryFlag    `json:",omitempty"`
	DisabledValidators  []DisabledValidator `json:",omitempty"`
	ValidatorToDisable  *VariableLength     `json:",omitempty"`
	ValidatorToReEnable *VariableLength     `json:",omitempty"`
}

type Escrow struct {
	leBase
	Flags           *LedgerEntryFlag `json:",omitempty"`
	Account         Account          `json:",omitempty"`
	Destination     Account          `json:",omitempty"`
	Amount          Amount           `json:",omitempty"`
	Condition       *VariableLength  `json:",omitempty"`
	CancelAfter     *uint32          `json:",omitempty"`
	FinishAfter     *uint32          `json:",omitempty"`
	SourceTag       *uint32          `json:",omitempty"`
	DestinationTag  *uint32          `json:",omitempty"`
	OwnerNode       *NodeIndex       `json:",omitempty"`
	DestinationNode *NodeIndex       `json:",omitempty"`
}

type SignerEntryItem struct {
	Account       *Account `json:",omitempty"`
	SignerWeight  *uint16  `json:",omitempty"`
	WalletLocator *Hash256 `json:",omitempty"`
}

type SignerEntry struct {
	SignerEntry SignerEntryItem
}

type SignerList struct {
	leBase
	Flags         *LedgerEntryFlag `json:",omitempty"`
	OwnerNode     *NodeIndex       `json:",omitempty"`
	SignerQuorum  *uint32          `json:",omitempty"`
	SignerEntries []SignerEntry    `json:",omitempty"`
	SignerListID  *uint32          `json:",omitempty"`
}

type Ticket struct {
	leBase
	Flags          *LedgerEntryFlag `json:",omitempty"`
	Account        *Account         `json:",omitempty"`
	TicketSequence *uint32          `json:",omitempty"`
	OwnerNode      *NodeIndex       `json:",omitempty"`
}

type PayChannel struct {
	leBase
	Flags           *LedgerEntryFlag `json:",omitempty"`
	Account         *Account         `json:",omitempty"`
	Destination     *Account         `json:",omitempty"`
	Amount          *Amount          `json:",omitempty"`
	Balance         *Amount          `json:",omitempty"`
	PublicKey       *PublicKey       `json:",omitempty"`
	SettleDelay     *uint32          `json:",omitempty"`
	OwnerNode       *NodeIndex       `json:",omitempty"`
	DestinationNode *NodeIndex       `json:",omitempty"`
	Expiration      *uint32          `json:",omitempty"`
	CancelAfter     *uint32          `json:",omitempty"`
	DestinationTag  *uint32          `json:",omitempty"`
	SourceTag       *uint32          `json:",omitempty"`
}

type Check struct {
	leBase
	Flags           *LedgerEntryFlag `json:",omitempty"`
	Account         *Account         `json:",omitempty"`
	Destination     *Account         `json:",omitempty"`
	OwnerNode       *NodeIndex       `json:",omitempty"`
	SendMax         *Amount          `json:",omitempty"`
	Sequence        *uint32          `json:",omitempty"`
	DestinationNode *NodeIndex       `json:",omitempty"`
	DestinationTag  *uint32          `json:",omitempty"`
	Expiration      *uint32          `json:",omitempty"`
	SourceTag       *uint32          `json:",omitempty"`
	InvoiceID       *Hash256         `json:",omitempty"`
}

type DepositPreAuth struct {
	leBase
	Flags     *LedgerEntryFlag `json:",omitempty"`
	Account   *Account         `json:",omitempty"`
	Authorize *Account         `json:",omitempty"`
	OwnerNode *NodeIndex       `json:",omitempty"`
}

type NFTokenPage struct {
	leBase
	Flags           *LedgerEntryFlag `json:",omitempty"`
	PreviousPageMin *Hash256         `json:",omitempty"`
	NextPageMin     *Hash256         `json:",omitempty"`
	NFTokens        []NFToken        `json:",omitempty"`
}

type NFTokenOffer struct {
	leBase
	Flags            *LedgerEntryFlag `json:",omitempty"`
	Owner            *Account         `json:",omitempty"`
	NFTokenID        *Hash256         `json:",omitempty"`
	Amount           *Amount          `json:",omitempty"`
	OwnerNode        *NodeIndex       `json:",omitempty"`
	NFTokenOfferNode *NodeIndex       `json:",omitempty"`
	Destination      *Account         `json:",omitempty"`
	Expiration       *uint32          `json:",omitempty"`
}
type VoteEntry struct {
	Account    Account `json:",omitempty"`
	TradingFee uint32  `json:",omitempty"`
	VoteWeight uint32  `json:",omitempty"`
}

type AuctionSlot struct {
	Account       Account `json:",omitempty"`
	DiscountedFee uint32  `json:",omitempty"`
	Expiration    *uint32 `json:",omitempty"`
	Price         Amount  `json:",omitempty"`
}

type AMM struct {
	leBase
	Flags          *LedgerEntryFlag `json:",omitempty"`
	Asset          Asset            `json:",omitempty"`
	Asset2         Asset            `json:",omitempty"`
	Account        *Account         `json:",omitempty"`
	LPTokenBalance *Amount          `json:",omitempty"`
	TradingFee     uint32           `json:",omitempty"`
	VoteSlots      []VoteEntry      `json:",omitempty"`
	AuctionSlot    *AuctionSlot     `json:",omitempty"`
	LedgerIndex    *Hash256         `json:",omitempty"`
}

func (a *AccountRoot) Affects(account Account) bool {
	return a.Account != nil && a.Account.Equals(account)
}
func (r *RippleState) Affects(account Account) bool {
	return r.LowLimit.Issuer.Equals(account) || r.HighLimit.Issuer.Equals(account)
}
func (o *Offer) Affects(account Account) bool        { return o.Account.Equals(account) }
func (d *Directory) Affects(account Account) bool    { return false }
func (l *LedgerHashes) Affects(account Account) bool { return false }
func (a *Amendments) Affects(account Account) bool   { return false }
func (f *FeeSettings) Affects(account Account) bool  { return false }
func (n *NegativeUNL) Affects(account Account) bool  { return false }
func (s *Escrow) Affects(account Account) bool {
	return s.Account.Equals(account) || s.Destination.Equals(account)
}
func (s *SignerList) Affects(account Account) bool {
	for _, entry := range s.SignerEntries {
		if entry.SignerEntry.Account != nil && entry.SignerEntry.Account.Equals(account) {
			return true
		}
	}
	return false
}
func (t *Ticket) Affects(account Account) bool { return t.Account != nil && t.Account.Equals(account) }
func (p *PayChannel) Affects(account Account) bool {
	return (p.Account != nil && p.Account.Equals(account)) || (p.Destination != nil && p.Destination.Equals(account))
}
func (p *Check) Affects(account Account) bool {
	return (p.Account != nil && p.Account.Equals(account)) || (p.Destination != nil && p.Destination.Equals(account))
}

func (d *DepositPreAuth) Affects(account Account) bool {
	return (d.Account != nil && d.Account.Equals(account)) || (d.Authorize != nil && d.Authorize.Equals(account))
}

func (tp *NFTokenPage) Affects(account Account) bool {
	return false
}

func (to *NFTokenOffer) Affects(account Account) bool {
	return (to.Owner != nil && to.Owner.Equals(account)) || (to.Destination != nil && to.Destination.Equals(account))
}

func (p *AMM) Affects(account Account) bool { return p.Account != nil }

func (le *leBase) GetType() string                     { return ledgerEntryNames[le.LedgerEntryType] }
func (le *leBase) GetLedgerEntryType() LedgerEntryType { return le.LedgerEntryType }
func (le *leBase) Prefix() HashPrefix                  { return HP_LEAF_NODE }
func (le *leBase) NodeType() NodeType                  { return NT_ACCOUNT_NODE }
func (le *leBase) Ledger() uint32                      { return 0 }
func (le *leBase) GetHash() *Hash256                   { return &le.Hash }
func (le *leBase) NodeId() *Hash256                    { return &le.Id }
func (le *leBase) GetLedgerIndex() *Hash256            { return le.LedgerIndex }
func (le *leBase) GetPreviousTxnId() *Hash256          { return le.PreviousTxnID }

func (o *Offer) Ratio() *Value {
	return o.TakerPays.Ratio(*o.TakerGets)
}

func (a *AMM) AMMID() *Hash256 {
	return a.LedgerIndex
}
