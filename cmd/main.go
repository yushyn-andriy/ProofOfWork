package main

import (
	"context"
	"flag"
	"hashcash/utilities"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

var (
	difficulty = flag.Int("d", 7, "difficulty [1-9]")
	profiling  = flag.Bool("p", false, "write profiling")
	auth       = flag.String("a", "auth", "prefix data")

	strLen       = flag.Int("l", 8, "string length")
	maxGorotines = flag.Int("g", 8, "max goroutines num")
)

func configure() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

func logRuntimeInfo() {
	log.Println("difficulty:", *difficulty)
	log.Println("string length:", *strLen)
	log.Println("max goroutines:", *maxGorotines)
	log.Println("profiling status:", *profiling)
	log.Println()

}

func main() {
	configure()
	flag.Parse()

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

	if *difficulty < 0 || *difficulty > 9 {
		log.Fatal("difficulty must be in range 1-9")
	}

	logRuntimeInfo()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Make sure it's called to release resources even if no errors

	out := make(chan string, (*maxGorotines)*5)
	var wg sync.WaitGroup
	for gid := 1; gid <= *maxGorotines; gid++ {
		src := rand.NewSource(time.Now().UnixNano())
		go func(out chan<- string, auth string, difficulty, strLen int, gid int, src rand.Source, ctx context.Context) {
			defer wg.Done()
			utilities.Hashcash(out, auth, difficulty, strLen, gid, src, ctx)
		}(out, *auth, *difficulty, *strLen, gid, src, ctx)
		wg.Add(1)
		time.Sleep(2 * time.Nanosecond)
	}

	start := time.Now()

	s := <-out
	log.Println("Result:", s)

	log.Println("Waiting for all goroutines to stop...")

	cancel()
	wg.Wait()
	close(out)

	log.Println("Elapsed time:", time.Since(start))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
