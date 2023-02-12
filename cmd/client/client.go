package main

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"hashcash/utilities"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	strLen       = flag.Int("l", 8, "string length")
	maxGorotines = flag.Int("g", 8, "max goroutines num")
	profiling    = flag.Bool("p", false, "write profiling")

	addr = flag.String("addr", "localhost:8080", "host:port")

	cert = flag.String("cert", ".\\cert\\cert1.pem", "cert location")
	key  = flag.String("key", ".\\cert\\key.pem", "key location")
)

func configure() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

func logRuntimeInfo() {
	log.Println("addr:", *addr)
	log.Println("cert:", *cert)
	log.Println("key:", *key)
	log.Println("string length:", *strLen)
	log.Println("max goroutines:", *maxGorotines)
	log.Println("profiling status", *profiling)

	log.Println()

}

func main() {
	flag.Parse()
	configure()
	logRuntimeInfo()

	if *profiling {
		cpu, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer cpu.Close()

		f, err := os.Create("mem.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()

		pprof.StartCPUProfile(cpu)
		defer pprof.StopCPUProfile()
	}

	cert, err := tls.LoadX509KeyPair(*cert, *key)
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", *addr, &config)

	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println(v.Subject)
	}
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	var auth string
	for {

		reply := make([]byte, 256)
		n, err := conn.Read(reply)
		if err != nil {
			log.Printf("client: read: %s", err)
			break
		}
		log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)

		args := strings.Split(strings.TrimSuffix(string(reply[:n]), "\n"), " ")
		fmt.Println(args)
		switch args[0] {
		case "HELO":
			n, err = conn.Write([]byte("EHLO\n"))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)
		case "ERROR":
			log.Printf("client: errros: wrote %s", strings.Join(args[1:], " "))
			return
		case "END":
			n, err = conn.Write([]byte("OK\n"))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)
			return
		case "POW":
			if len(args) != 3 {
				log.Printf("client: command switch: POW not enouth args got %d", len(args))
				return
			}

			auth = args[1]
			difficulty, err := strconv.Atoi(args[2])
			if err != nil {
				log.Printf("client: command switch: parse difficulty %q", err)
				return
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel() // Make sure it's called to release resources even if no errors

			out := make(chan string, *maxGorotines*2)
			var wg sync.WaitGroup
			for gid := 1; gid <= *maxGorotines; gid++ {
				src := rand.NewSource(time.Now().UnixNano())
				go func(
					out chan<- string,
					auth string,
					difficulty int,
					strLen int,
					gid int,
					src rand.Source,
					ctx context.Context,
				) {
					defer wg.Done()
					log.Printf("GR: %d, auth: %s, diff: %d, strLen: %d", gid, auth, difficulty, strLen)
					utilities.Hashcash(out, auth, difficulty, strLen, gid, src, ctx)
				}(out, auth, difficulty, *strLen, gid, src, ctx)
				wg.Add(1)
				time.Sleep(2 * time.Nanosecond)
			}

			start := time.Now()

			suffix := <-out
			cancel()
			wg.Wait()
			close(out)

			log.Println("Elapsed time:", time.Since(start))
			log.Printf("Auth: %q, Suffix: %q", auth, suffix)

			n, err = conn.Write([]byte(fmt.Sprintf("%s\n", suffix)))
			if err != nil {
				log.Printf("client: write: suffix %s", err)
				return
			}
			log.Printf("client: conn: wrote suffix %d bytes", n)
		case "NAME":
			resp := fmt.Sprintf("%x %s\n",
				sha1.Sum([]byte(auth+args[1])),
				"NAME SURNAME")

			log.Printf("client: conn: NAME resp %s", resp)

			n, err = conn.Write([]byte(resp))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)
		case "MAILNUM":
			resp := fmt.Sprintf("%x %s\n",
				sha1.Sum([]byte(auth+args[1])),
				"1")

			log.Printf("client: conn: MAILNUM resp %s", resp)

			n, err = conn.Write([]byte(resp))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)
		case "MAIL1":
			resp := fmt.Sprintf("%x %s\n",
				sha1.Sum([]byte(auth+args[1])),
				"mail@example.com")

			log.Printf("client: conn: MAIL1 resp %s", resp)

			n, err = conn.Write([]byte(resp))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)
		case "MAIL2":
			resp := fmt.Sprintf("%x %s\n",
				sha1.Sum([]byte(auth+args[1])),
				"mail@example.com")

			log.Printf("client: conn: MAIL2 resp %s", resp)

			n, err = conn.Write([]byte(resp))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)

		case "SKYPE":
			resp := fmt.Sprintf("%x %s\n",
				sha1.Sum([]byte(auth+args[1])),
				"N/A")

			log.Printf("client: conn: SKYPE resp %s", resp)

			n, err = conn.Write([]byte(resp))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)
		case "BIRTHDATE":
			resp := fmt.Sprintf("%x %s\n",
				sha1.Sum([]byte(auth+args[1])),
				"23.03.1763")

			log.Printf("client: conn: BIRTHDATE resp %s", resp)

			n, err = conn.Write([]byte(resp))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)
		case "COUNTRY":
			resp := fmt.Sprintf("%x %s\n",
				sha1.Sum([]byte(auth+args[1])),
				"COUNTRY")

			log.Printf("client: conn: COUNTRY resp %s", resp)

			n, err = conn.Write([]byte(resp))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)

		case "ADDRNUM":
			resp := fmt.Sprintf("%x %s\n",
				sha1.Sum([]byte(auth+args[1])),
				"1")

			log.Printf("client: conn: ADDRNUM resp %s", resp)

			n, err = conn.Write([]byte(resp))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)

		case "ADDRLINE1", "ADDRLINE2":
			resp := fmt.Sprintf("%x %s\n",
				sha1.Sum([]byte(auth+args[1])),
				"ADDR")

			log.Printf("client: conn: ADDRLINE1,2 resp %s", resp)

			n, err = conn.Write([]byte(resp))
			if err != nil {
				log.Printf("client: write: %s", err)
				return
			}
			log.Printf("client: conn: wrote %d bytes", n)

		default:
			log.Printf("client: command switch: unknown command %q %v", args[0], args[0] == "POW")
			return
		}

	}

}

func init() {
	rand.Seed(time.Now().UnixNano())
}
