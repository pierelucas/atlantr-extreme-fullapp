package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/pierelucas/atlantr-extreme-license-server/data"
	"github.com/pierelucas/atlantr-extreme-license-server/license"

	"github.com/pierelucas/atlantr-extreme-license-server/conn"

	"github.com/pierelucas/atlantr-extreme-license-server/utils"

	"github.com/tidwall/buntdb"
)

func runDB(ctx context.Context) {
	// open database
	log.Printf("parsing database %s\n", conf.DBName.String())
	db, err := buntdb.Open(conf.DBName.String())
	utils.CheckErrorFatal(err)

	defer db.Close()

	// start listener
	mutex := &sync.Mutex{}
	go conn.Listen(ctx, conf.Port.String(), func(clientAddr string, request interface{}) interface{} {
		var err error

		req := utils.Base64Decode(request.(string))

		pair, err := license.NewPair()
		utils.CheckError(err)

		err = pair.Unmarshal(req)
		utils.CheckError(err)

		// log connection
		log.Printf("IP: [%s] | Connection from ID: [%s] with License: [%s] \n", clientAddr, pair.ID.String(), pair.LICENSE.String())

		// First of all, we check if the appID is matching with the appID in our config file
		if conf.AppID.String() != pair.APPID.String() {
			log.Printf("IP: [%s] | The client ID: [%s] uses a different program version\n", clientAddr, pair.ID.String())
			return 1
		}

		// read database and check if the machineID is connected to the license key or if the license key is not connected to any machineID
		// when the machineID is not connected to any license key then connect the machineID to the license key
		// when the machineID is connected to the same license key as in database then the client connected from the same device
		var val string
		if err := db.View(func(tx *buntdb.Tx) error {
			var err error

			val, err = tx.Get(pair.LICENSE.String())
			if err != nil {
				return err
			}

			return nil
		}); err == nil {
			if val == "" {
				// write to database if the machineID value is empty and return success to the client
				mutex.Lock()
				err = db.Update(func(tx *buntdb.Tx) error {
					_, _, err := tx.Set(pair.LICENSE.String(), pair.ID.String(), nil)
					return err
				})
				if err != nil {
					log.Print(err)
				}
				mutex.Unlock()

				// log database write
				log.Printf("IP: [%s] | ID: [%s] is sucessfully registered with License: [%s]\n", clientAddr, pair.ID.String(), pair.LICENSE.String())

				return 0
			}

			// check if the clients machineID is the same machineID stored in database under the client's license key
			// when the machineID is the same, return success to the client
			if val == pair.ID.String() {
				// log database found
				log.Printf("IP: [%s] | ID: [%s] is sucessfully found in database\n", clientAddr, pair.ID.String())
				return 0
			}

			// if the client machineID is not the same a in database means that the client shared his key or trying to use a different machine
			// in this case the client have to buy a machineID reset. Return failure to the client.
			// log database write
			log.Printf("IP: [%s] | ID: [%s] is not the same as the ID in database under License: [%s]\n", clientAddr, pair.ID.String(), pair.LICENSE.String())
			return 1
		}

		// return failure if the license is not found in database. The client properly tries to use a random key.
		log.Printf("IP: [%s] | License: [%s] not found in database\n", clientAddr, pair.LICENSE.String())
		return 1
	})

	<-ctx.Done() // wait for context close
}

func init() {
	// open config.json
	conf = data.NewConf()
	conf.Open(configname)
}

func main() {
	// set-up logging
	log.SetFlags(log.LstdFlags)
	log.SetPrefix("The-Old-jew license Server\t")
	if *flagLOGOUTPUT != "" {
		f, err := os.OpenFile(*flagLOGOUTPUT, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		utils.CheckErrorFatal(err)
		defer f.Close() // defer = LIFO

		log.SetOutput(f) // set logfile
	}

	// Show that the programm is startign properly
	log.Println("The-Old-Jew-License-Server is starting...")

	// context
	ctx, cancel := context.WithCancel(context.Background())

	// Catch ctrl+c interrupt
	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("ABORT: you pressed ctrl+c")
		cancel()
	}()

	// run database server
	go runDB(ctx)

	// block routine on context
	<-ctx.Done()
	log.Println("shutting down server")
	time.Sleep(time.Second) // wait for server shutdown

	return // EXIT_SUCCESS
}
