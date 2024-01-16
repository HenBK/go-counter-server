package counter

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const (
	// The time window in seconds
	timeWindow = 60
	// The file where the counter state is persisted
	persistentFile = "data.txt"
)

type InMemoryCounter struct {
	mutex        sync.Mutex
	requestTimes map[int64]int
}

// Creates a new InMemoryCounter
func NewInMemoryCounter() *InMemoryCounter {
	return &InMemoryCounter{
		requestTimes: make(map[int64]int),
	}
}

// Increments the counter for a given request and returns the current counter value
func (c *InMemoryCounter) Increment(requestTime time.Time) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	currentTime := requestTime.Unix()
	c.requestTimes[currentTime]++

	count := 0

	for time, requestCount := range c.requestTimes {

		if currentTime - time <= timeWindow {
			count += requestCount
		} else {
			delete(c.requestTimes, time)
		}
	}

	return count
}

// Writes the current state to a file
func (c *InMemoryCounter) Persist() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	file, err := os.Create(persistentFile)

	if err != nil {
		return err
	}

	defer file.Close()

	for time, requestCount := range c.requestTimes {
		if _, err := file.WriteString(fmt.Sprintf("%d %d\n", time, requestCount)); err != nil {
			return err
		}
	}

	return nil
}

// Load reads the counter state from a file
func (c *InMemoryCounter) Load() error {
	file, err := os.Open(persistentFile)

	if err != nil {
		if os.IsNotExist(err) {
			// The file does not exist, so we don't have to load anything
			// it will be created when Persist is called
			return nil
		}
		return err
	}

	defer file.Close()

	var time int64
	var requestCount int

	for {
		_, err := fmt.Fscanf(file, "%d %d\n", &time, &requestCount)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		c.requestTimes[time] = requestCount
	}

	return nil
}
