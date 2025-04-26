package writer

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

type Wrapper struct {
	isFinalMerge bool
	file         *os.File
	writer       *bufio.Writer
}

func New(isFinalMerge bool) *Wrapper {
	outputWriter := &Wrapper{isFinalMerge: isFinalMerge}
	if isFinalMerge {
		outputWriter.writer = bufio.NewWriter(os.Stdout)
	} else {
		file, err := os.CreateTemp(".", "chunk_*.txt")
		if err != nil {
			log.Fatal(err)
		}
		outputWriter.file = file
		outputWriter.writer = bufio.NewWriter(file)
	}
	return outputWriter
}

func (o *Wrapper) Write(value int) {
	if _, err := o.writer.WriteString(strconv.Itoa(value) + "\n"); err != nil {
		log.Fatal(err)
	}
}

func (o *Wrapper) CloseAndGetFileName() string {
	if err := o.writer.Flush(); err != nil {
		log.Fatal("Error flushing stdout:", err)
	}
	if o.isFinalMerge {
		return ""
	} else {
		if err := o.file.Close(); err != nil {
			log.Fatalf("failed to close file %s: %v", o.file.Name(), err)
		}
		return o.file.Name()
	}
}
