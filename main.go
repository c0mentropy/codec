package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/ascii85"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"math/big"
	"net/url"
	"os"
	"strings"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/sha3"
)

var verbose bool
var outputFile string
var repeatCount int = 1

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	// Version check
	for _, arg := range os.Args[1:] {
		if arg == "-V" || arg == "--version" {
			printVersion()
			return
		}
	}

	if os.Args[1] == "--list" {
		printList()
		return
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		printHelp()
		return
	}

	args := parseFlags(os.Args[1:])

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Missing action and algorithm arguments")
		printHelp()
		os.Exit(1)
	}

	action := strings.ToLower(args[0])
	algo := strings.ToLower(args[1])
	var data string
	if len(args) >= 3 {
		data = args[2]
	} else {
		in, err := io.ReadAll(bufio.NewReader(os.Stdin))
		checkErr(err)
		data = string(in)
	}

	var result string

	switch action {
	case "encode":
		result = doEncode(algo, data, repeatCount)
	case "decode":
		result = doDecode(algo, data, repeatCount)
	case "hash":
		result = doHash(algo, data)
	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s\n", action)
		os.Exit(1)
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(result), 0644)
		checkErr(err)
		if verbose {
			fmt.Printf("[INFO] Result saved to file: %s\n", outputFile)
		}
	} else {
		if verbose {
			fmt.Println("[RESULT]")
		}
		fmt.Print(result)
		if !strings.HasSuffix(result, "\n") {
			fmt.Println()
		}
	}
}

func parseFlags(args []string) []string {
	parsed := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-v", "--verbose":
			verbose = true
		case "-o", "--output":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "Error: -o/--output requires a file path")
				os.Exit(1)
			}
			outputFile = args[i+1]
			i++
		case "-r", "--repeat":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "Error: --repeat requires a number")
				os.Exit(1)
			}
			fmt.Sscanf(args[i+1], "%d", &repeatCount)
			if repeatCount < 1 {
				repeatCount = 1
			}
			i++
		default:
			parsed = append(parsed, args[i])
		}
	}
	return parsed
}

func doEncode(algo, input string, repeat int) string {
	data := readData(input)
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
			buf := make([]byte, len(data)*5)
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
		fmt.Printf("[INFO] Operation : encode\n[INFO] Algorithm : %s\n[INFO] Repeat    : %d\n[INFO] InputLen  : %d\n[INFO] OutputLen : %d\n",
			algo, repeat, len(readData(input)), len(out))
	}
	return out
}

func doDecode(algo, input string, repeat int) string {
	data := readData(input)
	var out []byte
	var err error
	for i := 0; i < repeat; i++ {
		switch algo {
		case "base64":
			data += strings.Repeat("=", (4 - len(data) % 4) % 4)
			out, err = base64.StdEncoding.DecodeString(data)
		case "base64url":
			data += strings.Repeat("=", (4 - len(data) % 4) % 4)
			out, err = base64.URLEncoding.DecodeString(data)
		case "base32":
			data += strings.Repeat("=", (8 - len(data) % 8) % 8)
			out, err = base32.StdEncoding.DecodeString(data)
		case "hex":
			if len(data) % 2 != 0 {
				data = "0" + data
			}
			out, err = hex.DecodeString(data)
		case "base85":
			dst := make([]byte, len(data))
			n, _, err2 := ascii85.Decode(dst, []byte(data), true)
			err = err2
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
		checkErr(err)
		data = string(out)
	}
	if verbose {
		fmt.Printf("[INFO] Operation : decode\n[INFO] Algorithm : %s\n[INFO] Repeat    : %d\n[INFO] InputLen  : %d\n[INFO] OutputLen : %d\n",
			algo, repeat, len(readData(input)), len(out))
	}
	return string(out)
}

func doHash(algo, input string) string {
	var h hash.Hash
	switch algo {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	case "sha3-224":
		h = sha3.New224()
	case "sha3-256":
		h = sha3.New256()
	case "sha3-384":
		h = sha3.New384()
	case "sha3-512":
		h = sha3.New512()
	case "crc32-ieee":
		h = crc32.NewIEEE()
	case "crc32-castagnoli":
		h = crc32.New(crc32.MakeTable(crc32.Castagnoli))
	case "crc32-koopman":
		h = crc32.New(crc32.MakeTable(crc32.Koopman))
	case "blake2b-256":
		var err error
		h, err = blake2b.New256(nil)
		checkErr(err)
	case "blake2b-512":
		var err error
		h, err = blake2b.New512(nil)
		checkErr(err)
	case "blake2s-256":
		var err error
		h, err = blake2s.New256(nil)
		checkErr(err)
	default:
		fmt.Fprintf(os.Stderr, "Unknown hash algorithm: %s\n", algo)
		os.Exit(1)
	}

	if fileExists(input) {
		f, err := os.Open(input)
		checkErr(err)
		defer f.Close()
		_, err = io.Copy(h, f)
		checkErr(err)
	} else {
		io.WriteString(h, input)
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))
	if verbose {
		fmt.Printf("[INFO] Operation : hash\n[INFO] Algorithm : %s\n[INFO] InputLen  : %d\n[INFO] OutputLen : %d\n",
			algo, len(input), len(sum))
	}
	return sum
}

func readData(path string) string {
	var dataStr string
	if fileExists(path) {
		data, err := os.ReadFile(path)
		checkErr(err)
		dataStr = string(data)
	} else {
		dataStr = path
	}

	// 新增：去除首尾空白字符（包括换行符）
	dataStr = strings.TrimSpace(dataStr)

	return dataStr
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

const version = "v0.0.3"
const author = "Ckyan Comentroy"
const email = "comentropy@foxmail.com"
const github = "https://github.com/c0mentropy/codec"

func printHelp() {
	fmt.Println(`codec - Cross-platform encoding/decoding and hashing tool

Usage:
  codec <action> <algorithm> [data] [options]

Where:

  <action>      The operation to perform:
                  encode   Encode data using the specified algorithm
                  decode   Decode data using the specified algorithm
                  hash     Calculate the hash of data or file

  <algorithm>   The algorithm to use (see the supported list below)

  [data]        The input data or file path.
                If omitted, input is read from stdin.

  [options]     Additional flags:
                  -o, --output <file>   Write output to specified file
                  -v, --verbose         Enable verbose output
                  -r, --repeat <n>          Repeat encoding/decoding n times (default 1)
                  -h, --help            Show this help message and exit
                  -V, --version         Show version information and exit

Supported Algorithms:

  View via --list

Examples:

  Encode a string to base64:
    codec encode base64 "Hello, World!"

  Decode base64 string:
    codec decode base64 "SGVsbG8sIFdvcmxkIQ=="

  Hash a file with sha256:
    codec hash sha256 ./file.txt

  Encode a file and save output:
    codec encode base64 ./image.png -o encoded.txt

  Decode and save to file:
    codec decode base64 encoded.txt -o image_decoded.png

  Encode with multiple iterations:
    codec encode base64 "data" --repeat 3
`)
}

func printVersion() {
	fmt.Printf("codec version %s\n", version)
	fmt.Printf("Author: %s\n", author)
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("GitHub: %s\n", github)
}


func printList() {
	fmt.Println(`Supported Actions and Algorithms:

Encoding / Decoding algorithms:
  base64         Standard Base64 encoding/decoding
  base64url      URL-safe Base64 encoding/decoding
  base32         Base32 encoding/decoding
  hex            Hexadecimal encoding/decoding
  base85         ASCII85/Base85 encoding/decoding
  base58         Bitcoin Base58 encoding/decoding
  url            URL percent encoding/decoding

Hash algorithms:
  md5
  sha1
  sha256
  sha512
  sha3-224
  sha3-256
  sha3-384
  sha3-512
  crc32-ieee
  crc32-castagnoli
  crc32-koopman
  blake2b-256
  blake2b-512
  blake2s-256
 
Additional:
  If you need an algorithm, you can submit a PR or issue to me or send me an email.
`)
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

	// Handle leading zeros
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
			return nil, fmt.Errorf("Invalid Base58 character: %c", r)
		}
		x.Mul(x, base)
		x.Add(x, big.NewInt(int64(index)))
	}

	decoded := x.Bytes()

	// Handle leading 1s converted to 0x00
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
