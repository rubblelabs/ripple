package crypto

import (
	"encoding/hex"
	"fmt"
	. "launchpad.net/gocheck"
)

type KeySuite struct{}

var _ = Suite(&KeySuite{})

func checkHash(h Hash, err error) string {
	if err != nil {
		panic(err)
	}
	return h.ToJSON()
}

func checkSignature(sender, receiver Key, hash []byte) bool {
	sig, err := sender.Sign(hash)
	if err != nil {
		panic(err)
	}
	ok, err := Verify(receiver.PublicCompressed(), sig, hash)
	if err != nil {
		panic(err)
	}
	return ok
}

func checkHex(b []byte, err error) string {
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

func hexToBytes(s string) []byte {
	h, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return h
}

// Examples from https://ripple.com/wiki/Account_Family
func (s *KeySuite) TestWikiVectors(c *C) {
	zero, err := NewRippleHash("0")
	c.Check(err, IsNil)
	c.Check(zero.ToJSON(), Equals, ACCOUNT_ZERO)
	c.Check(checkHex(Sha512Half(zero.PayloadTrimmed())), Equals, "B8244D028981D693AF7B456AF8EFA4CAD63D282E19FF14942C246E50D9351D22")

	seed := hexToBytes("71ED064155FFADFA38782C5E0158CB26")
	key, err := GenerateRootDeterministicKey(seed)
	c.Check(err, IsNil)
	c.Check(checkHex(key.PrivateBytes(), nil), Equals, "7CFBA64F771E93E817E15039215430B53F7401C34931D111EAB3510B22DBB0D8")
	c.Check(key.Seed.ToJSON(), Equals, "shHM53KPZ87Gwdqarm1bAmPeXg8Tn")
	c.Check(checkHex(key.Seed.Value().Bytes(), nil), Equals, "71ED064155FFADFA38782C5E0158CB26")
	c.Check(checkHash(key.PublicGenerator()), Equals, "fht5yrLWh3P8DrJgQuVNDPQVXGTMyPpgRHFKGQzFQ66o3ssesk3o")
	c.Check(checkHash(key.GenerateAccountId(0)), Equals, "rhcfR9Cg98qCxHpCcPBmMonbDBXo84wyTn")
	c.Check(checkHash(key.PublicNodeKey()), Equals, "n9MXXueo837zYH36DvMc13BwHcqtfAWNJY5czWVbp7uYTj7x17TH")
	c.Check(checkHash(key.PrivateNodeKey()), Equals, "pa91wmE8V8K63SAMGMpdFpik8wGAcbUdSmHABccV9jFfqhTijH1")

	account, err := key.GenerateAccountKey(0)
	c.Check(err, IsNil)
	c.Check(checkHash(account.PublicAccountKey()), Equals, "aBRoQibi2jpDofohooFuzZi9nEzKw9Zdfc4ExVNmuXHaJpSPh8uJ")
	c.Check(checkHash(account.PrivateAccountKey()), Equals, "pwMPbuE25rnajigDPBEh9Pwv8bMV2ebN9gVPTWTh4c3DtB14iGL")
}

// Examples from https://github.com/ripple/rippled/blob/develop/src/ripple_data/protocol/RippleAddress.cpp
func (s *KeySuite) TestRippledVectors(c *C) {
	testMessage := []byte("Hello, nurse!")
	seed, err := GenerateFamilySeed("masterpassphrase")
	c.Check(err, IsNil)
	c.Check(seed.ToJSON(), Equals, "snoPBrXtMeMyMHUVTgbuqAfg1SUTb")
	key, err := GenerateRootDeterministicKey(seed.PayloadTrimmed())
	c.Check(err, IsNil)
	c.Check(checkHash(key.PublicNodeKey()), Equals, "n94a1u4jAz288pZLtw6yFWVbi89YamiC6JBXPVUj5zmExe5fTVg9")
	c.Check(checkHash(key.PrivateNodeKey()), Equals, "pnen77YEeUd4fFKG7iycBWcwKpTaeFRkW2WFostaATy1DSupwXe")
	hash, err := Sha512Half(testMessage)
	c.Check(err, IsNil)
	c.Check(checkSignature(key, key, hash), Equals, true)
	c.Check(checkHash(key.PublicGenerator()), Equals, "fhuJKrhSDzV2SkjLn9qbwm5AaRmrxDPfFsHDCP6yfDZWcxDFz4mt")
	first, err := key.GenerateAccountKey(0)
	c.Check(err, IsNil)
	second, err := key.GenerateAccountKey(1)
	c.Check(err, IsNil)
	c.Check(checkHash(key.GenerateAccountId(0)), Equals, "rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh")
	c.Check(checkHash(first.PublicAccountKey()), Equals, "aBQG8RQAzjs1eTKFEAQXr2gS4utcDiEC9wmi7pfUPTi27VCahwgw")
	c.Check(checkHash(first.PrivateAccountKey()), Equals, "p9JfM6HHi64m6mvB6v5k7G2b1cXzGmYiCNJf6GHPKvFTWdeRVjh")
	c.Check(checkHash(key.GenerateAccountId(1)), Equals, "r4bYF7SLUMD7QgSLLpgJx38WJSY12ViRjP")
	c.Check(checkHash(second.PublicAccountKey()), Equals, "aBPXpTfuLy1Bhk3HnGTTAqnovpKWQ23NpFMNkAF6F1Atg5vDyPrw")
	c.Check(checkHash(second.PrivateAccountKey()), Equals, "p9JEm822LMrzJii1k7TvdphfENTp6G5jr253Xa5rkzUWVr8ogQt")
	c.Check(checkSignature(first, first, hash), Equals, true)
	c.Check(checkSignature(first, second, hash), Equals, false)
	c.Check(checkSignature(second, second, hash), Equals, true)
	// Skipped message encryption - doesn't appear to be used in rippled's codebase...
}
