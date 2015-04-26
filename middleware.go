package main

import (
	"log"
	"net/http"
	"time"
)

func commonHandler(next http.HandlerFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		t0 := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %q %v", r.Method, r.URL.String(), time.Now().Sub(t0))
	}

	return http.HandlerFunc(fn)
}
