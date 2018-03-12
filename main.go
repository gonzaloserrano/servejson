package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var (
		route = flag.String("route", "/", "json file path")
		file  = flag.String("file", "", "json file path")
		port  = flag.Int("port", 8080, "server port")
	)
	flag.Parse()

	logger := log.New(os.Stderr, "", log.LstdFlags)
	if *file == "" {
		logger.Fatal("should provide a json file with -file")
	}
	if r := *route; string(r[0]) != "/" {
		*route = fmt.Sprintf("/%s", *route)
	}

	err := serve(*port, *file, *route, logger)
	if err != nil {
		logger.Fatal(err)
	}
}

func serve(port int, file, route string, logger *log.Logger) error {
	http.DefaultServeMux.HandleFunc(route, jsonFileWithCORSOptionsHandlerFunc(file, logger))
	log.Printf("Serving %s at %s in port %d\n", file, route, port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func jsonFileWithCORSOptionsHandlerFunc(file string, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		for name, value := range headers {
			h.Set(name, value)
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			logger.Printf("Served OPTIONS to %s\n", r.RemoteAddr)
			return
		}
		h.Set("Content-Type", "application/json")
		logger.Printf("Served file to %s\n", r.RemoteAddr)
		http.ServeFile(w, r, file)
	}
}

var headers = map[string]string{
	"Access-Control-Allow-Methods":  allMethodsStr,
	"Access-Control-Allow-Origin":   "*",
	"Access-Control-Allow-Headers":  "content-type, authorization, x-custom-auth",
	"Access-Control-Expose-Headers": "authorization",
}

var allMethods = []string{
	http.MethodHead,
	http.MethodGet,
	http.MethodPut,
	http.MethodPost,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodOptions,
	http.MethodConnect,
	http.MethodTrace,
}

var allMethodsStr = strings.Join(allMethods, ",")
