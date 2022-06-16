package ttmtron

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/mr-tron/base58"
)

// HexToAddress converts a hex representation of a Tron address
// into a Base58 string with a 4 byte checksum.
func HexToAddress(hexAddr string) (b58 string, err error) {
	bytes, err := hex.DecodeString(hexAddr)
	if err != nil {
		return "", err
	}
	var checksum [32]byte
	checksum = sha256.Sum256(bytes)
	checksum = sha256.Sum256(checksum[:])
	bytes = append(bytes, checksum[:4]...)
	b58 = base58.EncodeAlphabet(bytes, base58.BTCAlphabet)

	return
}

func Base58ToHex(in string) string {
	if len(in) != 34 { //nolint:gomnd // length of base58 address
		return in // already hex
	}

	res, _ := base58.Decode(in)

	return hex.EncodeToString(res[:len(res)-4])
}

func HexToInt256(hex string) (*big.Int, error) {
	var bignum, ok = new(big.Int).SetString(hex, 16)

	if ok {
		return bignum, nil
	}

	return nil, errors.New("failed to convert")
}

func EncodeAddressToParameter(in string) (string, error) {
	// base58 address need to convert
	if len(in) == 34 { //nolint:gomnd // base58 length
		in = Base58ToHex(in)
	}

	if len(in) == 42 { //nolint:gomnd // hex address starting with 41
		in = in[2:]
	}

	return fmt.Sprintf("%064s", in), nil
}

func DecodeConstantToSymbol(in string) (string, error) {
	if len(in)%64 != 0 {
		return "", errors.New("input should divide by 64")
	}

	length, err := HexToInt256(in[64:128])

	if err != nil {
		return "", err
	}

	symbolLength := length.Uint64()

	bs, err := hex.DecodeString(in[128 : 192-64+symbolLength*2])
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func AddParameterAddress(in string) string {
	addr := Base58ToHex(in)
	addr = strings.TrimPrefix(addr, "41")

	return fmt.Sprintf("%064s", addr)
}

func AddParameterAmount(in uint64) string {
	return fmt.Sprintf("%064x", in)
}
