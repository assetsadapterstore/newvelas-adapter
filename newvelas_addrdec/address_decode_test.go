package newvelas_addrdec

import (
	"encoding/hex"
	"testing"
)

func TestAddressDecoder_AddressEncode(t *testing.T) {

	addrdec := NewAddressDecoder()
	p2pk, _ := hex.DecodeString("0000000000000000000000000000000000000000")
	p2pkAddr, _ := addrdec.AddressEncode(p2pk)
	t.Logf("p2pkAddr: %s", p2pkAddr)
}

func TestAddressDecoder_AddressDecode(t *testing.T) {

	addrdec := NewAddressDecoder()
	p2pkAddr := "V8aq1aCzKRZsqJb9Apsbdtm1wPkgtQXVdZ"
	p2pkHash, _ := addrdec.AddressDecode(p2pkAddr)
	//checksum: d9a537a9
	t.Logf("p2pkHash: %s", hex.EncodeToString(p2pkHash))
}

func TestAddressDecoderV2_AddressVerify(t *testing.T) {
	addrdec := NewAddressDecoder()
	p2pkAddr := "V8aq1aCzKRZsqJb9Apsbdtm1wPkgtQXVdZ"
	flag := addrdec.AddressVerify(p2pkAddr)
	t.Logf("flag: %v", flag)
}