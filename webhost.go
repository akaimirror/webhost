// Command webhost runs simple web server serving static files from mapped
// host/path mapped to fs directories.
package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/artyom/autoflags"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	args := struct {
		Addr string `flag:"addr,address to listen at"`
		Conf string `flag:"map,file with host/backend mapping"`

		RTo time.Duration `flag:"rto,maximum duration before timing out read of the request"`
		WTo time.Duration `flag:"wto,maximum duration before timing out write of the response"`
	}{
		Addr: "localhost:8080",
		Conf: "mapping.yml",
		RTo:  10 * time.Second,
		WTo:  5 * time.Minute,
	}
	autoflags.Define(&args)
	flag.Parse()
	h, err := newHandler(args.Conf)
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:         args.Addr,
		Handler:      h,
		ReadTimeout:  args.RTo,
		WriteTimeout: args.WTo,
	}
	log.Fatal(srv.ListenAndServe())
}

func newHandler(file string) (http.Handler, error) {
	m, err := readMapping(file)
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	for k, v := range m {
		switch idx := strings.IndexRune(k, '/'); idx {
		case -1:
			mux.Handle(k+"/", http.FileServer(http.Dir(v)))
		default:
			mux.Handle(k, http.StripPrefix(k[idx:], http.FileServer(http.Dir(v))))
		}
	}
	return mux, nil
}

func readMapping(file string) (map[string]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lr := io.LimitReader(f, 1<<20)
	b := new(bytes.Buffer)
	if _, err := io.Copy(b, lr); err != nil {
		return nil, err
	}
	m := make(map[string]string)
	if err := yaml.Unmarshal(b.Bytes(), &m); err != nil {
		return nil, err
	}
	return m, nil
}
