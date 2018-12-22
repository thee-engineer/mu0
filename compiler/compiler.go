package compiler

import (
	"log"
	"os"
	"time"

	"github.com/thee-engineer/mu0-vm/mu0"
)

// Compile ...
func Compile(inFile, outFile string) []mu0.Word {
	log.Printf("Compile: %s -> %s\n", inFile, outFile)
	t := time.Now()

	// Lexical analysis of source code
	tree := lex(inFile)

	// Compiled binary for MU0
	var binary []mu0.Word
	var instruction mu0.Word

	// Iterate lexical tree
	for _, tkn := range tree {
		// Skip EQU
		if tkn.t == tokenTypeEQU {
			continue
		}

		// If define, create word, else parse token
		if tkn.t == tokenTypeDEF {
			// If link to label, else word define
			if addr, ok := labels[tkn.arg]; ok {
				instruction = mu0.Word(addr)
			} else {
				instruction = parseArg(tkn, tree)
			}
		} else {
			// Extract instruction op code
			instruction = tokenTypeToOPC[tkn.t]

			// Extract instruction argument
			instruction |= parseArg(tkn, tree)
		}

		binary = append(binary, instruction)
	}

	// Create output file
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// Create byte array
	byteStream := make([]byte, len(binary)*2)

	// Move binary words to byte array
	for idx, w := range binary {
		// ! DEBUG
		// fmt.Printf("%X %s\n", idx, decompileInstruction(w))

		byteStream[idx*2] = byte(w >> 8)
		byteStream[idx*2+1] = byte(w & 0x00FF)
	}

	// Write binary data to file
	_, err = f.Write(byteStream)
	if err != nil {
		log.Fatalln(err)
	}

	// Log details
	log.Printf("Compile: finished OK, in %d ns\n",
		time.Now().Nanosecond()-t.Nanosecond())
	log.Printf("Compile: wrote %d bytes to %s\n", len(byteStream), outFile)

	// Return binary as word array
	return binary
}
