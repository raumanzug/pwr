/*
Pwr program weaves java .properties files into another java .properties file.

Usage:

	pwr [ --from <source .properties file> ] [ --to <target .properties file> ]

the flag mean:

	--from source java .properties file from which key/value pairs are taken in order to update
		key/value settings in target java .properties file
	--to target java .properties file whose key/value pairs are updated
*/
package main

import (
	"io"
	"log"
	"os"
	"regexp"

	"golang.org/x/text/encoding/charmap"
)

const (
	// key/value pair representation in java .properties file
	propertiesRegex = "(?m:^\\s*(\\S+)\\s*[:=]\\s*((?:.*\\\\$\\n)*.*$))"
)

var (
	decoder                  = charmap.ISO8859_1.NewDecoder()
	pCompiledPropertiesRegex *regexp.Regexp // compiled from propertiesRegex
)

func convert_iso_8859_1_to_utf_8_string(iso_8859_1 []byte) (utf_8_string string, err error) {
	var utf_8 []byte

	utf_8, err = decoder.Bytes(iso_8859_1)
	if err != nil {
		return
	}
	utf_8_string = string(utf_8)

	return
}

// reads a file and stores its content into result
func getProps(filename string) (result []byte, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	result, err = io.ReadAll(file)
	if err != nil {
		return
	}

	return
}

// weaves parameter setting determined by a dict into the target .properties file
func outProps(dict map[string]*[]byte, pCliParams *cliParams) (err error) {
	TargetProps, err := getProps(pCliParams.targetFilename)
	if err != nil {
		return
	}

	bakFile := pCliParams.targetFilename + ".bak"
	err = os.Rename(pCliParams.targetFilename, bakFile)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			err = os.Rename(bakFile, pCliParams.targetFilename)
		}
	}()

	targetFile, err := os.Create(pCliParams.targetFilename)
	if err != nil {
		return
	}
	defer targetFile.Close()

	curTargetPropsIndex := 0
	for _, submatches := range pCompiledPropertiesRegex.FindAllSubmatchIndex(TargetProps, -1) {
		k_iso_8859_1 := TargetProps[submatches[2]:submatches[3]]
		var key string
		key, err = convert_iso_8859_1_to_utf_8_string(k_iso_8859_1)
		if err != nil {
			return
		}
		pV := dict[key]
		if pV == nil {
			continue
		}
		_, err = targetFile.Write(TargetProps[curTargetPropsIndex:submatches[4]])
		if err != nil {
			return
		}
		_, err = targetFile.Write(*pV)
		if err != nil {
			return
		}
		curTargetPropsIndex = submatches[5]
	}
	_, err = targetFile.Write(TargetProps[curTargetPropsIndex:])

	return
}

/*
How does it works?

 1. read in source .properties file into a []byte (using procedure getProps).
 2. parses this []byte and find the key/value pairs, create a dict[string]*[]byte
    which assigns keys to values.
 3. backup target .properties file
 4. read in backup target .properties file into a []byte (using procedure getProps).
 5. create target .properties file from backup .properties file and the beforementioned
    modified dict.
*/
func main() {
	var err error
	pCompiledPropertiesRegex, err = regexp.Compile(propertiesRegex)
	if err != nil {
		log.Fatal(err)
	}

	cliParams := cliParams{}
	cli(&cliParams)

	// 1.
	sourceProps, err := getProps(cliParams.sourceFilename)
	if err != nil {
		log.Fatal(err)
	}

	// 2.
	dict := make(map[string]*[]byte)
	for _, submatches := range pCompiledPropertiesRegex.FindAllSubmatchIndex(sourceProps, -1) {
		k_iso_8859_1 := sourceProps[submatches[2]:submatches[3]]
		var key string
		key, err = convert_iso_8859_1_to_utf_8_string(k_iso_8859_1)
		if err != nil {
			log.Fatal(err)
		}
		value := (sourceProps[submatches[4]:submatches[5]])
		dict[key] = &value
	}

	// 3., 4., 5.
	err = outProps(dict, &cliParams)
	if err != nil {
		log.Fatal(err)
	}

}
