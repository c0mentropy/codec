# codec

**一个跨平台的编码、解码和哈希CLI工具**  
零依赖，使用Go编写——支持base64、base32、hex、base85、base58、URL编码以及多种哈希算法。

---

## 目录

- [特性](#特性)  
- [支持的算法](#支持的算法)  
- [安装](#安装)  
- [使用方法](#使用方法)  
    - [基本用法](#基本用法)  
    - [选项](#选项)  
    - [示例](#示例)  
- [高级功能](#高级功能)  
    - [重复编码/解码](#重复编码解码)  
    - [输出到文件](#输出到文件)  
    - [详细模式](#详细模式)  
- [版本和作者信息](#版本和作者信息)  
- [贡献](#贡献)  
- [许可证](#许可证)  

---

## 特性

- 用于编码、解码和哈希的跨平台CLI工具  
- 支持多种编码方案（base64、base32、base85、base58、hex、URL编码）  
- 支持常见和高级哈希算法（MD5、SHA系列、SHA3、Blake2、CRC32）  
- 支持从标准输入（stdin）或文件读取输入  
- 支持输出到标准输出（stdout）或文件  
- 支持多次重复编码/解码  
- 详细模式可显示操作的详细信息  
- 零外部运行时依赖（仅使用Go标准库 + golang.org/x/crypto）  

---

## 支持的算法

### 编码 / 解码

- `base64`      — 标准Base64  
- `base64url`   — URL安全的Base64  
- `base32`      — Base32  
- `hex`         — 十六进制  
- `base85`      — ASCII85 / Base85  
- `base58`      — 比特币Base58  
- `url`         — URL百分比编码  

### 哈希

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

## 安装

确保你已安装[Go](https://golang.org/dl/)（建议版本1.16+）。

克隆仓库并构建：
git clone https://github.com/c0mentropy/codec.git
cd codec
go build -o codec main.go
或从[发布页](https://github.com/c0mentropy/codec/releases)下载预编译二进制文件（即将推出）。

---

## 使用方法

codec <操作> <算法> [数据] [选项]

- `<操作>`: `encode`（编码）、`decode`（解码）或 `hash`（哈希）  
- `<算法>`: 参见上面的支持的算法  
- `[数据]`: 输入字符串或文件路径。如果省略，将从标准输入（stdin）读取输入。  
- `[选项]`: 见下文  

### 选项

| 标志                  | 描述                                  |
| --------------------- | -------------------------------------------- |
| `-o, --output <文件>` | 将输出写入指定文件           |
| `-v, --verbose`       | 启用详细输出                        |
| `-r, --repeat <n>`    | 重复编码/解码n次（默认1次） |
| `-h, --help`          | 显示帮助信息                        |
| `-V, --version`       | 显示版本和作者信息          |

---

## 示例

- 将字符串编码为base64：
codec encode base64 "Hello, World!"
- 解码base64字符串：
codec decode base64 "SGVsbG8sIFdvcmxkIQ=="
- 使用sha256哈希一个文件：
codec hash sha256 ./file.txt
- 编码二进制文件并保存输出：
codec encode base64 ./image.png -o encoded.txt
- 解码并保存到文件：
codec decode base64 encoded.txt -o image_decoded.png
- 多次编码一个字符串：
codec encode base64 "data" --repeat 3
- 详细模式查看详细信息：
codec encode base64 "hello" -v

---

## 高级功能

### 重复编码/解码

使用`--repeat <n>`可以多次重复编码或解码步骤。这在处理嵌套编码时非常有用。

示例：对字符串进行三次Base64编码
codec encode base64 "secret" --repeat 3
### 输出到文件

使用`-o`或`--output`将结果保存到文件而不是标准输出。

示例：
codec encode base64 ./file.bin -o encoded.txt
### 详细模式

使用`-v`或`--verbose`启用详细操作日志。
codec decode base64 encoded.txt -v
这会显示操作、算法、输入/输出长度、重复次数和其他信息。

---

## 贡献

欢迎贡献和提交bug报告！  
请在[GitHub](https://github.com/c0mentropy/codec)上创建issue或提交拉取请求。

---

## 许可证

本项目基于MIT许可证授权——详见[LICENSE](LICENSE)文件了解详情。
    