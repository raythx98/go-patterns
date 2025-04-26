package main

import (
	"bufio"
	"container/heap"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/raythx98/go-patterns/sort-big-file/minheap"
	"github.com/raythx98/go-patterns/sort-big-file/writer"
)

type ChunkFile struct {
	Name    string
	File    *os.File
	Scanner *bufio.Scanner
}

func main() {
	// Adjust chunk size and max files open based on available memory
	const (
		chunkSize    = 2
		maxFilesOpen = 2
	)

	chunkFiles := make([]*ChunkFile, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Read large files into manageable chunk
		chunk := make([]int, 0, chunkSize)
		for i := 0; i < chunkSize; i++ {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()

			val, err := strconv.Atoi(line)
			if err != nil {
				log.Fatal(err)
			}

			chunk = append(chunk, val)
		}

		if len(chunk) == 0 {
			break
		}

		// sort chunk
		slices.Sort(chunk)

		// write chunk to temp file
		chunkWriter := writer.New(false)
		for _, val := range chunk {
			chunkWriter.Write(val)
		}
		chunkFiles = append(chunkFiles, &ChunkFile{Name: chunkWriter.CloseAndGetFileName()})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Repeatedly merge sorted chunks
	for len(chunkFiles) > 0 {
		// Split chunk files into batches
		batches := make([][]*ChunkFile, 0)
		for i := 0; i < len(chunkFiles); i += maxFilesOpen {
			end := i + maxFilesOpen
			if end > len(chunkFiles) {
				end = len(chunkFiles)
			}
			batches = append(batches, chunkFiles[i:end])
		}

		mergedFile := make([]*ChunkFile, 0)
		// Merge batch of sorted chunks
		for _, batch := range batches { // execute batch
			// open file & scanners
			for _, chunkFile := range batch {
				file, err := os.Open(chunkFile.Name)
				if err != nil {
					log.Fatal(err)
				}
				chunkFile.File = file
				chunkFile.Scanner = bufio.NewScanner(file)
			}

			// create output writer
			// if last pass, we can write to stdout directly, otherwise write to file
			isFinalMerge := len(chunkFiles) <= maxFilesOpen
			outputWriter := writer.New(isFinalMerge)

			// create & initialize min heap
			minHeap := &minheap.MinHeap{}
			heap.Init(minHeap)

			// Push minimum item from each chunk into min heap
			for i, chunkFile := range batch {
				if chunkFile.Scanner.Scan() {
					val, err := strconv.Atoi(chunkFile.Scanner.Text())
					if err != nil {
						log.Fatal(err)
					}

					heap.Push(minHeap, minheap.Item{Value: val, FileIndex: i})
				}
			}

			// Repeatedly pop from min heap, write to output, and push next smallest value from the same file
			for minHeap.Len() > 0 {
				minItem := heap.Pop(minHeap).(minheap.Item)
				outputWriter.Write(minItem.Value)
				if batch[minItem.FileIndex].Scanner.Scan() {
					val, err := strconv.Atoi(batch[minItem.FileIndex].Scanner.Text())
					if err != nil {
						log.Fatal(err)
					}
					heap.Push(minHeap, minheap.Item{Value: val, FileIndex: minItem.FileIndex})
				}
			}

			// Close all files in the batch
			for _, chunkFile := range batch {
				if err := chunkFile.File.Close(); err != nil {
					log.Fatalf("failed to close file %s: %v", chunkFile.File.Name(), err)
				}
			}

			// Close output writer, add newly merged chunk for next round of merge
			if newChunkFileName := outputWriter.CloseAndGetFileName(); newChunkFileName != "" {
				mergedFile = append(mergedFile, &ChunkFile{Name: newChunkFileName})
			}
		}

		// Remove old chunk files
		for _, chunkFile := range chunkFiles {
			if err := os.Remove(chunkFile.Name); err != nil {
				log.Fatalf("failed to remove file %s: %v", chunkFile.Name, err)
			}
		}

		// Update chunk files for next round of merge
		chunkFiles = mergedFile
	}
}
