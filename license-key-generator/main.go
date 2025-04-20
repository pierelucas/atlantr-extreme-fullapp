package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/google/uuid"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	pbar "github.com/schollz/progressbar/v3"
	"github.com/tidwall/buntdb"
)

func main() {
	// read num of keys from commandline
	read := bufio.NewReader(os.Stdin)

	fmt.Printf("Please enter num of keys\n\n-> ")
	numStr, err := read.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// delete the delimeter and convert string to int
	num, err := strconv.Atoi(strings.Replace(numStr, "\n", "", -1))

	// create progressbar
	bar := pbar.NewOptions(int(num),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][reset]Creating license keys..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	// open client file
	f1, err := os.OpenFile("clientkeys.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f1.Close()

	// open database
	db, err := buntdb.Open("license.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	// create client and trimmed keys and write to file
	for i := 0; i < num; i++ {
		KEY := uuid.New().String()
		f1.WriteString(KEY + "\r\n")

		// remove punctuation and control characters from the key
		trimmedKEY := strings.Map(func(r rune) rune {
			if unicode.IsPunct(r) || unicode.IsControl(r) {
				return -1
			}
			return r
		}, KEY)

		// write trimmed keys to database with the empty value ""
		mutex := sync.Mutex{}
		if err := db.View(func(tx *buntdb.Tx) error {
			var err error

			_, err = tx.Get(trimmedKEY)
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			// write the license to database.
			mutex.Lock()
			err = db.Update(func(tx *buntdb.Tx) error {
				_, _, err := tx.Set(trimmedKEY, "", nil)

				bar.Add(1) // adding one more key to progressbar

				return err
			})
			if err != nil {
				log.Print(err)
			}
			mutex.Unlock()
		}
	}
}
