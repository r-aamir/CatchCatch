package execfunc

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"
)

// RecoverWrapper is http middleware to recover from panics
func RecoverWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := WithRecover(func() error {
			h.ServeHTTP(w, r)
			return nil
		})
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
	})
}

// WithRecover wraps functions and turns its panics into errors
func WithRecover(fn func() error) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%v", r)
			log.Printf("[panic withRecover] %v", err)
			debug.PrintStack()
		}
	}()
	return fn()
}

// OnExit traps application to run function before exit
func OnExit(fn func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		fn()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()
}
