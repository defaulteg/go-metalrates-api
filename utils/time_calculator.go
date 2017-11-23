package utils

import (
	"time"
	"fmt"
)

func ExecutionTime(functionName string, start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("%s method took %s\n", functionName, elapsed)
}

