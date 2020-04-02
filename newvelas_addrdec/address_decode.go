package newvelas_addrdec

import (
	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/mr-tron/base58"
	"strings"
)

var (
	alphabet = addressEncoder.BTCAlphabet
)

//AddressDecoderV2
type AddressDecoderV2 struct {
	*openwallet.AddressDecoderV2Base
}

//NewAddressDecoder 地址解析器
func NewAddressDecoder() *AddressDecoderV2 {
	decoder := AddressDecoderV2{}
	return &decoder
}

//AddressDecode 地址解析
func (dec *AddressDecoderV2) AddressDecode(addr string, opts ...interface{}) ([]byte, error) {

	strippedAddress := strings.Replace(addr, "V", "", 1)

	return base58.Decode(strippedAddress)
}

//AddressEncode 地址编码
func (dec *AddressDecoderV2) AddressEncode(hash []byte, opts ...interface{}) (string, error) {

	if len(hash) > 32 {
		//公钥hash处理
		publicKey := owcrypt.PointDecompress(hash, owcrypt.ECC_CURVE_SECP256K1)
		hash = owcrypt.Hash(publicKey[1:len(publicKey)], 0, owcrypt.HASH_ALG_KECCAK256)
		hash = hash[12:]
	} else if len(hash) == 32 {
		hash = hash[12:]
	}

	encodedAddress := base58.Encode(hash)

	return "V" + string(encodedAddress), nil
}

// AddressVerify 地址校验
func (dec *AddressDecoderV2) AddressVerify(address string, opts ...interface{}) bool {
	strippedAddress := strings.Replace(address, "V", "", 1)

	hash, err := base58.Decode(strippedAddress)
	if err != nil {
		return false
	}
	if len(hash) != 20 {
		return false
	}
	return true
}
