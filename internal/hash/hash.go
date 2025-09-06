package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/sha3"

	"github.com/c0mentropy/codec/internal/util"
)

func DoHash(algo, input string, verbose bool) string {
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
		util.CheckErr(err)
	case "blake2b-512":
		var err error
		h, err = blake2b.New512(nil)
		util.CheckErr(err)
	case "blake2s-256":
		var err error
		h, err = blake2s.New256(nil)
		util.CheckErr(err)
	default:
		fmt.Fprintf(os.Stderr, "Unknown hash algorithm: %s\n", algo)
		os.Exit(1)
	}

	if util.FileExists(input) {
		f, err := os.Open(input)
		util.CheckErr(err)
		defer f.Close()
		_, err = io.Copy(h, f)
		util.CheckErr(err)
		sum := fmt.Sprintf("%x", h.Sum(nil))
		if verbose {
			fmt.Printf("\n[INFO] Operation : hash\n[INFO] Algorithm : %s\n[INFO] InputType : file (%s)\n[INFO] OutputLen : %d\n[INFO] Result : ",
				algo, util.BaseName(input), len(sum))
		}
		return sum
	}

	io.WriteString(h, input)
	sum := fmt.Sprintf("%x", h.Sum(nil))
	if verbose {
		fmt.Printf("\n[INFO] Operation : hash\n[INFO] Algorithm : %s\n[INFO] InputType : string\n[INFO] InputLen  : %d\n[INFO] OutputLen : %d\n[INFO] Result : ",
			algo, len(input), len(sum))
	}
	return sum
}

func DoCompare(algo string, inputs []string, verbose bool, output string) bool {
	if len(inputs) < 2 {
		fmt.Fprintln(os.Stderr, "compare requires at least 2 files")
		os.Exit(1)
	}

	files := inputs
	isSame := true
	results := []string{}

	hashes := make([]string, len(files))
	for i, f := range files {
		hashes[i] = DoHash(algo, f, false)
	}

	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if hashes[i] == hashes[j] {
				results = append(results, fmt.Sprintf("[ OK ] %s == %s", files[i], files[j]))
			} else {
				results = append(results, fmt.Sprintf("[FAIL] %s != %s", files[i], files[j]))
				isSame = false
			}
		}
	}

	outputText := ""
	for _, line := range results {
		outputText += line + "\n"
	}

	if output != "" {
		f, err := os.OpenFile(output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot open output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		_, err = f.WriteString(outputText)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot write to output file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Print(outputText)
	}

	return isSame
}
