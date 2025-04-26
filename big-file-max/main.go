package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	first := true
	maxVal := math.MinInt
	defer func() {
		output := fmt.Sprintf("max value is %d\n", maxVal)
		if first {
			output = "no values found\n"
		}
		_, err := os.Stdout.WriteString(output)
		if err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		next := scanner.Text()
		if len(next) == 0 {
			break
		}

		val, err := strconv.Atoi(next)
		if err != nil {
			log.Fatal(err)
		}

		if first || val > maxVal {
			maxVal = val
			first = false
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
