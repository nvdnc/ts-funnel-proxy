package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"tailscale.com/ipn/store/mem"
	"tailscale.com/tsnet"
)

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	var (
		hostname = flag.String("ts-hostname", getEnvOrDefault("TS_HOSTNAME", "my-funneled-service"), "hostname to use on the tailnet")

		// Use any HTTP server as the upstream
		// A tool to debug HTTP requests might be httplab, which starts an interactive web server. To run it:
		//   go run github.com/gchaincl/httplab/cmd/httplab@latest --port 4000
		upstream = flag.String("upstream", getEnvOrDefault("UPSTREAM_ENDPOINT", "http://localhost:4000"), "upstream URL to proxy to")
	)

	flag.Parse()

	srv := new(tsnet.Server)
	srv.Ephemeral = true
	srv.Store = new(mem.Store)
	srv.Hostname = *hostname
	ln, err := srv.ListenFunnel("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}

	for _, domain := range srv.CertDomains() {
		log.Printf("Access at https://%s\n", domain)
	}

	proxy, err := NewProxy(*upstream)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.Serve(ln, proxy))
}

type SimpleProxy struct {
	Proxy *httputil.ReverseProxy
}

func NewProxy(rawUrl string) (*SimpleProxy, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	s := &SimpleProxy{httputil.NewSingleHostReverseProxy(u)}

	originalDirector := s.Proxy.Director
	s.Proxy.Director = func(r *http.Request) {
		// This is where we could modify requests
		originalDirector(r)
	}

	s.Proxy.ModifyResponse = func(r *http.Response) error {
		// This is where we could modify responses
		return nil
	}

	return s, nil
}

func (s *SimpleProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(w, r)
}
