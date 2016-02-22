package main

import (
	"flag"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/titan-x/titan"
)

const (
	addr     = "127.0.0.1:3000"
	testAddr = "127.0.0.1:3001"
)

var (
	start = flag.Bool("start", false, "Start a Titan server at address: "+addr)
	ext   = flag.Bool("ext", false, "Run external client test case. Titan server will run at address: "+testAddr)
)

func main() {
	flag.Parse()
	switch {
	case *start:
		startServer(addr)
	case *ext:
		startExtTest(testAddr)
	default:
		flag.PrintDefaults()
	}
}

func startServer(addr string) {
	s, err := titan.NewServer(addr)
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}
	defer s.Close()

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("error closing server: %v", err)
	}
}

func startExtTest(addr string) {
	log.Printf("-ext flag is provided, starting external client test case.")
	titan.InitConf("test")

	now := time.Now().Unix()
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims["userid"] = "1"
	t.Claims["created"] = now
	ts, err := t.SignedString([]byte(titan.Conf.App.JWTPass()))
	if err != nil {
		log.Fatalf("failed to sign JWT token: %v", err)
	}
	log.Printf("Sample valid user JWT token for testing: %v", ts)

	startServer(addr)
}