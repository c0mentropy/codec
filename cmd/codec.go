package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/c0mentropy/codec/internal/encode"
	"github.com/c0mentropy/codec/internal/hash"
	"github.com/c0mentropy/codec/internal/util"
	"github.com/c0mentropy/codec/internal/version"
)

type Config struct {
	Verbose     bool
	IsOutput    bool
	OutputFile  string
	RepeatCount int
}

func Main() {
	cfg := &Config{
		Verbose:     false,
		IsOutput:    false,
		OutputFile:  "",
		RepeatCount: 1,
	}
	run(os.Args[1:], cfg)
}

func run(args []string, cfg *Config) {
	if len(args) < 1 {
		version.PrintHelp()
		return
	}

	for _, arg := range args {
		switch arg {
		case "-V", "--version":
			version.PrintVersion()
			return
		case "--list":
			version.PrintList()
			return
		case "-h", "--help":
			version.PrintHelp()
			return
		}
	}

	args = parseFlags(args, cfg)

	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Missing action and algorithm arguments")
		version.PrintHelp()
		os.Exit(1)
	}

	action := strings.ToLower(args[0])
	algo := strings.ToLower(args[1])
	inputs := args[2:]

	// stdin
	if len(inputs) == 0 {
		in, err := io.ReadAll(bufio.NewReader(os.Stdin))
		util.CheckErr(err)
		inputs = []string{string(in)}
	}

	// * (all)
	expandedInputs := []string{}
	for _, in := range inputs {
		matches, err := filepath.Glob(in)
		if err != nil || len(matches) == 0 {
			expandedInputs = append(expandedInputs, in)
		} else {
			expandedInputs = append(expandedInputs, matches...)
		}
	}
	inputs = expandedInputs

	for _, input := range inputs {
		var result string
		switch action {
		case "encode":
			result = encode.DoEncode(algo, input, cfg.RepeatCount, cfg.Verbose)
		case "decode":
			result = encode.DoDecode(algo, input, cfg.RepeatCount, cfg.Verbose)
		case "hash":
			result = hash.DoHash(algo, input, cfg.Verbose)
		case "compare":
			hash.DoCompare(algo, inputs, cfg.Verbose, cfg.OutputFile)
			return
		default:
			fmt.Fprintf(os.Stderr, "Unknown action: %s\n", action)
			os.Exit(1)
		}

		if cfg.OutputFile != "" {
			f, err := os.OpenFile(cfg.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			util.CheckErr(err)
			defer f.Close()

			_, err = f.Write([]byte(result))
			util.CheckErr(err)
		} else if cfg.IsOutput {
			base := util.BaseName(input)
			outName := fmt.Sprintf("%s.%s", base, algo)
			err := os.WriteFile(outName, []byte(result), 0644)
			util.CheckErr(err)
		} else {
			fmt.Println(result)
		}

	}
}

func parseFlags(args []string, cfg *Config) []string {
	parsed := []string{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-v", "--verbose":
			cfg.Verbose = true
		case "-o", "--output":
			cfg.IsOutput = true
			if i+1 >= len(args) {
				fmt.Println("If no output file is specified, the default output will be to: filename.algo")
			} else {
				cfg.OutputFile = args[i+1]
			}
			i++
		case "-r", "--repeat":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "Error: --repeat requires a number")
				os.Exit(1)
			}
			fmt.Sscanf(args[i+1], "%d", &cfg.RepeatCount)
			if cfg.RepeatCount < 1 {
				cfg.RepeatCount = 1
			}
			i++
		default:
			parsed = append(parsed, args[i])
		}
	}
	return parsed
}
