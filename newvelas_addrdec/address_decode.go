package newvelas_addrdec

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/mr-tron/base58"
	"regexp"
	"strings"
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

	strippedAddress := strings.TrimPrefix(addr, "V")
	raw, err := base58.Decode(strippedAddress)
	if err != nil {
		return nil, err
	}

	//无checksum的编码
	if len(raw) == 20 {
		return raw, nil
	}

	decodedAddress := hex.EncodeToString(raw)

	regex := regexp.MustCompile(`([0-9abcdef]+)([0-9abcdef]{8})$`)
	if !regex.MatchString(decodedAddress) {
		return nil, fmt.Errorf("invalid decoded address")
	}

	matches := regex.FindStringSubmatch(decodedAddress)
	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid address")
	}

	for len(matches[1]) > 40 {
		if strings.HasPrefix(matches[1], "0") {
			matches[1] = strings.TrimPrefix(matches[1], "0")
		} else {
			return nil, fmt.Errorf("invalid match")
		}
	}

	checksum := sha(sha(matches[1]))[:8]

	if matches[2] != checksum {
		return nil, fmt.Errorf("invalid checksum")
	}

	if len(matches[1]) != 40 {
		return nil, fmt.Errorf("failed to get eth address")
	}

	//log.Debugf("checksum: %s", checksum)
	return hex.DecodeString(matches[1])
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

	checksum := sha(sha(hex.EncodeToString(hash)))[:8]
	//checksum := owcrypt.Hash(hash, 0, owcrypt.HASH_ALG_DOUBLE_SHA256)[:4]
	//log.Debugf("checksum: %s", checksum)
	checksumBytes, _ := hex.DecodeString(checksum)
	raw := append(hash, checksumBytes...)
	encodedAddress := base58.Encode(raw)
	if len(encodedAddress) < 33 {
		encodedAddress = fmt.Sprintf("%s%s", strings.Repeat("1", 33-len(encodedAddress)), encodedAddress)
	}
	return "V" + string(encodedAddress), nil
}

// AddressVerify 地址校验
func (dec *AddressDecoderV2) AddressVerify(address string, opts ...interface{}) bool {

	strippedAddress := strings.TrimPrefix(address, "V")
	raw, err := base58.Decode(strippedAddress)
	if err != nil {
		return false
	}

	//无checksum的编码
	if len(raw) <= 20 {
		return false
	}

	_, err = dec.AddressDecode(address, opts)
	if err != nil {
		return false
	}

	return true
}

func sha(raw string) string {
	hasher := sha256.New()
	hasher.Write([]byte(raw))
	return hex.EncodeToString(hasher.Sum(nil))
}
