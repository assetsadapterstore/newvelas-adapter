package newvelas_addrdec

import (
	"encoding/hex"
	"testing"
)

func TestAddressDecoder_AddressEncode(t *testing.T) {

	addrdec := NewAddressDecoder()
	p2pk, _ := hex.DecodeString("eb057d96e2532257e47dbd8d3090c8be5030db77")
	p2pkAddr, _ := addrdec.AddressEncode(p2pk)
	t.Logf("p2pkAddr: %s", p2pkAddr)
}

func TestAddressDecoder_AddressDecode(t *testing.T) {

	addrdec := NewAddressDecoder()
	p2pkAddr := "V4GuZFgXAjqLAXgqTCcbJmyzNbp3g"
	p2pkHash, _ := addrdec.AddressDecode(p2pkAddr)
	t.Logf("p2pkHash: %s", hex.EncodeToString(p2pkHash))
}

func TestAddressDecoderV2_AddressVerify(t *testing.T) {
	addrdec := NewAddressDecoder()
	p2pkAddr := "V4GuZFgXAjqLAXgqTCcbJmyzNbp3g"
	flag := addrdec.AddressVerify(p2pkAddr)
	t.Logf("flag: %v", flag)
}