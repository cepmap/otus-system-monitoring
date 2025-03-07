package main

import (
	"fmt"

	"github.com/cepmap/otus-system-monitoring/internal/stats/disksLoad"
)

func main() {
	res, err := disksLoad.GetStats()
	if err != nil {
		return
	}

	fmt.Println(res)

	//ticker := time.NewTicker(1 * time.Second)
	//defer ticker.Stop() // Ensure the ticker is stopped when done
	//
	//// Run a loop that executes a function every second
	//for range ticker.C {
	//	// Call your function here
	//	res, err := diskStat.GetStats()
	//	if err != nil {
	//		return
	//	}
	//	fmt.Println(res)
	//}
}
