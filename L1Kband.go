package main

import (
	"fmt"
	"sync"
)

// Monitor which protects a line of symbols
type Monitor struct {
	line  string
	mutex *sync.Mutex
	cond  *sync.Cond
	done  bool
}

func (monitor *Monitor) insert(symbol byte) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()
	for (symbol == 'B' || symbol == 'C') && len(monitor.line) <= 3 {
		monitor.cond.Wait()
	}
	if monitor.done {
		return
	}
	monitor.line += string([]byte{symbol})
	if len(monitor.line) > 3 {
		monitor.cond.Broadcast()
	}
}

func inserter(monitor *Monitor, waiter *sync.WaitGroup, symbol byte) {
	defer waiter.Done()
	for i := 0; i < 15; i++ {
		monitor.insert(symbol)
		if monitor.done {
			return
		}
	}
	monitor.mutex.Lock()
	monitor.done = true
	monitor.mutex.Unlock()
}

func initializeMonitor() Monitor {
	mutex := sync.Mutex{}
	cond := sync.NewCond(&mutex)
	return Monitor{"*", &mutex, cond, false}
}

func main() {
	monitor := initializeMonitor()
	waiter := sync.WaitGroup{}
	waiter.Add(3)
	go inserter(&monitor, &waiter, 'A')
	go inserter(&monitor, &waiter, 'B')
	go inserter(&monitor, &waiter, 'C')
	for !monitor.done {
		fmt.Println(monitor.line)
	}
	waiter.Wait()
	fmt.Println(monitor.line)
}
