package utilities

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/bits"
	"math/rand"

	shac "github.com/TvdW/gotor/sha1"
)

const EVERY_COUNTER = 100000000

// Hashcash tries to found sha1 hash with difficulty
// leading zeros at the start of hash hexadecimal representation
func Hashcash(
	out chan<- string,
	auth string, difficulty int,
	strLen int,
	gid int,
	src rand.Source,
	ctx context.Context,
) {
	log.Println("Goroutine started:", gid)

	leadingZeros := difficulty * 4
	var (
		number uint64
		i      int
		data   []byte
	)
	row := make([]byte, strLen)

	hasher := shac.New()
	hasher.Write([]byte(auth))

	for {
		i += 1
		if i%EVERY_COUNTER == 0 {
			log.Printf("GR %d : i : %d", gid, i)
		}

		select {
		case <-ctx.Done():
			log.Println("Quit goroutine:", gid)
			return
		default:
			suffix := RandStringBytesMaskImprSrcUnsafe(strLen, src)
			row = []byte(suffix)
			hasher1 := hasher.Clone()
			hasher1.Write(row)
			data = hasher1.Sum(nil)
			number = binary.BigEndian.Uint64(data[:8])

			if bits.LeadingZeros64(number) >= leadingZeros {
				log.Println("Found hash GR:", gid)
				log.Printf("%s|%x|GR:%d\n\r", auth+suffix, data, gid)

				out <- fmt.Sprintf("%s", suffix)
				log.Println("Written to the channel GR:", gid)
				return
			}
		}
	}

}
