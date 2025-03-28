package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	first := true
	maxVal := 0
	defer func() {
		if first {
			os.Stdout.WriteString("no values found\n")
		} else {
			os.Stdout.WriteString(fmt.Sprintf("max value is %d\n", maxVal))
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
			os.Stderr.WriteString(fmt.Sprintf("%v\n", err))
		}

		if first || val > maxVal {
			maxVal = val
			first = false
		}
	}
}
