package conn

import (
	"context"
	"fmt"
	"log"

	"github.com/pierelucas/gorpc"
)

// Listen starts the gorpc fileserver. Call callback for every incoming request. The return of the callback function is returned
// to the calling client. We can send data back if we want
func Listen(ctx context.Context, port string, callback func(string, interface{}) interface{}) {
	addr := fmt.Sprintf(":%s", port)
	s := &gorpc.Server{
		// Accept clients on this TCP address.
		Addr: addr,

		// Request Handler - this functions is called on each incoming request
		Handler: callback,
	}

	if err := s.Start(); err != nil {
		log.Fatalf("Cannot start Server: %s", err)
	}

	log.Println("Server sucessfully started")
	log.Printf("Listening on localhost:%s\n", port)

	// wait for context
	<-ctx.Done()

	s.Stop()
	log.Println("Server sucessfully stopped")

	return
}
