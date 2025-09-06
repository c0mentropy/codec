package encode

import (
	"encoding/ascii85"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strings"

	"github.com/c0mentropy/codec/internal/util"
)

func DoEncode(algo, input string, repeat int, verbose bool) string {
	data := util.ReadData(input)
	var out string
	for i := 0; i < repeat; i++ {
		switch algo {
		case "base64":
			out = base64.StdEncoding.EncodeToString([]byte(data))
		case "base64url":
			out = base64.URLEncoding.EncodeToString([]byte(data))
		case "base32":
			out = base32.StdEncoding.EncodeToString([]byte(data))
		case "hex":
			out = hex.EncodeToString([]byte(data))
		case "base85":
			buf := make([]byte, ascii85.MaxEncodedLen(len(data)))
			n := ascii85.Encode(buf, []byte(data))
			out = string(buf[:n])
		case "base58":
			out = base58Encode([]byte(data))
		case "url":
			out = url.QueryEscape(data)
		default:
			fmt.Fprintf(os.Stderr, "Unknown encoding algorithm: %s\n", algo)
			os.Exit(1)
		}
		data = out
	}
	if verbose {
		if util.FileExists(input) {
			fmt.Printf("[INFO] Operation : encode\n[INFO] Algorithm : %s\n[INFO] InputType : file (%s)\n[INFO] Repeat    : %d\n[INFO] OutputLen : %d\n",
				algo, util.BaseName(input), repeat, len(out))
		} else {
			fmt.Printf("[INFO] Operation : encode\n[INFO] Algorithm : %s\n[INFO] InputType : string\n[INFO] Repeat    : %d\n[INFO] InputLen  : %d\n[INFO] OutputLen : %d\n",
				algo, repeat, len(input), len(out))
		}
	}
	return out
}

func DoDecode(algo, input string, repeat int, verbose bool) string {
	data := util.ReadData(input)
	var out []byte
	var err error
	for i := 0; i < repeat; i++ {
		switch algo {
		case "base64":
			data += strings.Repeat("=", (4-len(data)%4)%4)
			out, err = base64.StdEncoding.DecodeString(data)
		case "base64url":
			data += strings.Repeat("=", (4-len(data)%4)%4)
			out, err = base64.URLEncoding.DecodeString(data)
		case "base32":
			data += strings.Repeat("=", (8-len(data)%8)%8)
			out, err = base32.StdEncoding.DecodeString(data)
		case "hex":
			if len(data)%2 != 0 {
				data = "0" + data
			}
			out, err = hex.DecodeString(data)
		case "base85":
			dst := make([]byte, len(data)*4)
			n, _, err2 := ascii85.Decode(dst, []byte(data), true)
			util.CheckErr(err2)
			out = dst[:n]
		case "base58":
			out, err = base58Decode(data)
		case "url":
			var s string
			s, err = url.QueryUnescape(data)
			out = []byte(s)
		default:
			fmt.Fprintf(os.Stderr, "Unknown decoding algorithm: %s\n", algo)
			os.Exit(1)
		}
		util.CheckErr(err)
		data = string(out)
	}
	if verbose {
		if util.FileExists(input) {
			fmt.Printf("[INFO] Operation : decode\n[INFO] Algorithm : %s\n[INFO] InputType : file (%s)\n[INFO] Repeat    : %d\n[INFO] OutputLen : %d\n",
				algo, util.BaseName(input), repeat, len(out))
		} else {
			fmt.Printf("[INFO] Operation : decode\n[INFO] Algorithm : %s\n[INFO] InputType : string\n[INFO] Repeat    : %d\n[INFO] InputLen  : %d\n[INFO] OutputLen : %d\n",
				algo, repeat, len(input), len(out))
		}
	}
	return string(out)
}

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func base58Encode(input []byte) string {
	x := new(big.Int).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)
	var result []byte
	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = append([]byte{base58Alphabet[mod.Int64()]}, result...)
	}
	for _, b := range input {
		if b == 0x00 {
			result = append([]byte{base58Alphabet[0]}, result...)
		} else {
			break
		}
	}
	return string(result)
}

func base58Decode(input string) ([]byte, error) {
	x := big.NewInt(0)
	base := big.NewInt(58)
	for _, r := range input {
		index := strings.IndexRune(base58Alphabet, r)
		if index < 0 {
			return nil, fmt.Errorf("invalid Base58 character: %c", r)
		}
		x.Mul(x, base)
		x.Add(x, big.NewInt(int64(index)))
	}
	decoded := x.Bytes()
	leadingZeros := 0
	for _, r := range input {
		if r == rune(base58Alphabet[0]) {
			leadingZeros++
		} else {
			break
		}
	}
	decoded = append(make([]byte, leadingZeros), decoded...)
	return decoded, nil
}
