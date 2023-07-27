package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch = int64(82800000) // Go release date epoch (Nov. 10, 2009)
	/* "epoch" refers to a specific point in time from which time is
	measured or calculated. It serves as a reference point for expressing
	dates and timestamps. Without an epoch there would be no consistent
	starting point for the timestamp portion of the unique ID.
	*/
	workerBits   = 10
	sequenceBits = 12
	maxWorkerID  = 1024 // 2^10 = 1024
	maxSequence  = 4096 // 2^12 = 4096
	/* Since a bit can be either 0 or 1, you take 2 to the power of the
	number of bits to get the maximum value that can be represented by to
	find the maximum value that can be represented by the number of bits.
	*/
)

// IDGenerator represents a Snowflake ID generator.
type IDGenerator struct {
	mu       sync.Mutex
	workerID int64
	lastTime int64
	sequence int64
}

// NewIDGenerator creates a new instance of IDGenerator with the given worker ID.
func NewIDGenerator(workerID int64) (*IDGenerator, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}
	return &IDGenerator{
		workerID: workerID,
	}, nil
}

// GenerateID generates a new unique ID.
func (g *IDGenerator) GenerateID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixNano() / 1e6 // convert to milliseconds, the number of nanoseconds elapsed since January 1, 1970 UTC

	if g.lastTime == now {
		g.sequence = (g.sequence + 1) & maxSequence
		/*
			GitHub Copilot: The bitwise AND operation is used to mask off any bits in the "sequence" variable that are set to 1 beyond the maximum value represented by the "maxSequence" variable.
			For example, if "maxSequence" is set to 15 (binary 1111), then the bitwise AND operation will ensure that the value of "sequence" is always between 0 and 15. If "sequence" is currently 14 (binary 1110) and the operation adds 1 to it, the result would be 15 (binary 1111). However, the bitwise AND operation with "maxSequence" would then mask off the leftmost bit, resulting in a final value of 15 & 15 = 15.
			In this way, the bitwise AND operation effectively "wraps around" the value of "sequence" once it reaches the maximum value represented by "maxSequence". This helps ensure that the generated IDs or sequences do not exceed a certain range and can be safely stored or used in other parts of the program.
		*/
		if g.sequence == 0 {
			// Sequence rollover within the same millisecond, wait until next millisecond
			for now <= g.lastTime {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		g.sequence = 0 // sequence is only used when ids are generated in the same millisecond.
	}

	g.lastTime = now

	// ID is constructed by shifting the timestamp, worker ID, and sequence bits into place.
	id := (now-epoch)<<workerBits | (g.workerID << sequenceBits) | g.sequence

	return id
}

// func main() {
// 	// Create a new ID generator with worker ID 1.
// 	generator, err := NewIDGenerator(1)
// 	if err != nil {
// 		fmt.Println("Error creating ID generator:", err)
// 		return
// 	}

// 	// Generate and print 10 unique IDs.
// 	for i := 0; i < 10; i++ {
// 		id := generator.GenerateID()
// 		fmt.Println("Generated ID:", id)
// 	}
// 	// create instance of IDGenerator and call GenerateID method on the same line
// }
