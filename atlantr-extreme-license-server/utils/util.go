package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Print(err)
	}
	return string(data)
}

// CheckError and print to log if err != nil
func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// CheckErrorFatal print to log and calls os.Exit(1) is error is != nil
func CheckErrorFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// CheckErrorPrint and print if err != nil
func CheckErrorPrint(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// CheckErrorPrintFatal and calls os.Exit(1) is error is != nil
func CheckErrorPrintFatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
