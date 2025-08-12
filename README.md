# codec

**A Cross-Platform Encoding, Decoding, and Hashing CLI Tool**  
Zero dependencies, written in Go — support base64, base32, hex, base85, base58, URL encoding, and multiple hash algorithms.  

---

## Table of Contents

- [Features](#features)  
- [Supported Algorithms](#supported-algorithms)  
- [Installation](#installation)  
- [Usage](#usage)  
    - [Basic Usage](#basic-usage)  
    - [Options](#options)  
    - [Examples](#examples)  
- [Advanced](#advanced)  
    - [Repeat Encoding/Decoding](#repeat-encodingdecoding)  
    - [Output to File](#output-to-file)  
    - [Verbose Mode](#verbose-mode)  
- [Version and Author Info](#version-and-author-info)  
- [Contributing](#contributing)  
- [License](#license)  

---

## Features

- Cross-platform CLI tool for encoding, decoding, and hashing  
- Supports multiple encoding schemes (base64, base32, base85, base58, hex, URL encoding)  
- Supports common and advanced hash algorithms (MD5, SHA family, SHA3, Blake2, CRC32)  
- Supports reading input from stdin or files  
- Supports output to stdout or files  
- Supports repeated encoding/decoding multiple times  
- Verbose mode to display detailed operation info  
- Zero external runtime dependencies (only Go standard lib + golang.org/x/crypto)  

---

## Supported Algorithms

### Encoding / Decoding

- `base64`      — Standard Base64  
- `base64url`   — URL-safe Base64  
- `base32`      — Base32  
- `hex`         — Hexadecimal  
- `base85`      — ASCII85 / Base85  
- `base58`      — Bitcoin Base58  
- `url`         — URL Percent Encoding  

### Hashing

- `md5`  
- `sha1`  
- `sha256`  
- `sha512`  
- `sha3-224`  
- `sha3-256`  
- `sha3-384`  
- `sha3-512`  
- `crc32-ieee`  
- `crc32-castagnoli`  
- `crc32-koopman`  
- `blake2b-256`  
- `blake2b-512`  
- `blake2s-256`  

---

## Installation

Make sure you have [Go](https://golang.org/dl/) installed (version 1.16+ recommended).

Clone the repo and build:

```bash
git clone https://github.com/yourgithub/codec.git
cd codec
go build -o codec main.go
```

Or download prebuilt binaries from [Releases](https://github.com/yourgithub/codec/releases) (coming soon).

---

## Usage

codec <action> <algorithm> [data] [options]

- `<action>`: `encode`, `decode`, or `hash`  
- `<algorithm>`: see Supported Algorithms above  
- `[data]`: input string or file path. If omitted, input is read from standard input (stdin).  
- `[options]`: see below  

### Options

| Flag                  | Description                                  |
| --------------------- | -------------------------------------------- |
| `-o, --output <file>` | Write output to the specified file           |
| `-v, --verbose`       | Enable verbose output                        |
| `-r, --repeat <n>`    | Repeat encoding/decoding n times (default 1) |
| `-h, --help`          | Show help message                            |
| `-V, --version`       | Show version and author information          |

---

## Examples

- Encode a string to base64:

```bash
codec encode base64 "Hello, World!"
```

- Decode base64 string:

```bash
codec decode base64 "SGVsbG8sIFdvcmxkIQ=="
```

- Hash a file with sha256:

```bash
codec hash sha256 ./file.txt
```

- Encode a binary file and save output:

```bash
codec encode base64 ./image.png -o encoded.txt
```

- Decode and save to file:

```bash
codec decode base64 encoded.txt -o image_decoded.png
```

- Encode a string multiple times:

```bash
codec encode base64 "data" --repeat 3
```

- Verbose mode to see detailed info:

```bash
codec encode base64 "hello" -v
```

---

## Advanced

### Repeat Encoding/Decoding

Use `--repeat <n>` to repeat the encoding or decoding step multiple times. This is useful if you have nested encodings.

Example: triple Base64 encode a string

```bash
codec encode base64 "secret" --repeat 3
```

### Output to File

Use `-o` or `--output` to save the result to a file instead of stdout.

Example:

```bash
codec encode base64 ./file.bin -o encoded.txt
```

### Verbose Mode

Enable detailed operation logging with `-v` or `--verbose`.

```bash
codec decode base64 encoded.txt -v
```

This shows the action, algorithm, input/output lengths, repeat counts, and other info.

---

## Version and Author Info

Show version, author, and contact info with:

```bash
codec -V
```

Output example:

```bash
codec version v0.0.1
Author: Ckyan Comentroy
Email: comentropy@foxmail.com
GitHub: https://github.com/c0mentropy/codec

```

---

## Contributing

Contributions and bug reports are welcome!  
Feel free to open issues or submit pull requests on [GitHub](https://github.com/yourgithub/codec).

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.