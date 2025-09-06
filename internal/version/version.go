package version

import "fmt"

const (
	Version = "v0.0.4"
	Author  = "Ckyan Comentroy"
	Email   = "comentropy@foxmail.com"
	GitHub  = "https://github.com/c0mentropy/codec"
)

func PrintHelp() {
	fmt.Println(`codec - Cross-platform encoding/decoding and hashing tool

Usage:
  codec <action> <algorithm> [data] [options]

Where:

  <action>      The operation to perform:
                  encode   Encode data using the specified algorithm
                  decode   Decode data using the specified algorithm
                  hash     Calculate the hash of data or file
                  compare  Compare the hashes of several files to see if they are the same

  <algorithm>   The algorithm to use (see the supported list below)

  [data]        The input data or file path.
                If omitted, input is read from stdin.

  [options]     Additional flags:
                  -o, --output [file]   Write output to specified file (default Output to: filename.algo)
                  -v, --verbose         Enable verbose output
                  -r, --repeat <n>      Repeat encoding/decoding n times (default 1)
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
    codec encode base64 "data" --repeat 3`)
}

func PrintVersion() {
	fmt.Printf("codec version %s\n", Version)
	fmt.Printf("Author: %s\n", Author)
	fmt.Printf("Email: %s\n", Email)
	fmt.Printf("GitHub: %s\n", GitHub)
}

func PrintList() {
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
  If you need an algorithm, you can submit a PR or issue to me or send me an email.`)
}
