package Helpers

import (
	"NullOps/CLI_Handlers"
	"NullOps/Interface"
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"runtime"
	"strconv"
	"sync"
)

type Method func(string)

func ThreadingWithSingleProducer(method Method, maxWorkers int, lines []string) {
	var wg sync.WaitGroup
	workChan := make(chan string, maxWorkers*5)
	seen := make(map[string]struct{})

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range workChan {
				method(line)
			}
		}()
	}

	for _, line := range lines {
		if _, exists := seen[line]; !exists {
			seen[line] = struct{}{}
			workChan <- line
		}
	}

	close(workChan)
	wg.Wait()

	Interface.Write("Press enter to go Main Menu")
	_, err := fmt.Scanln()
	CLI_Handlers.LogError(err)
}

func ThreadingWithBatches(method Method, maxWorkers int, lines []string) {
	batchSize := CalculateBatchSize(len(lines), runtime.NumCPU())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)

	batchCh := make(chan []string, maxWorkers)

	Interface.Option("Lines per Thread:", strconv.Itoa(batchSize))
	Interface.Option("Threading Type:", ThreadingType)
	Interface.Option("Threads", strconv.Itoa(maxWorkers))

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			for batch := range batchCh {
				for _, line := range batch {
					select {
					case sem <- struct{}{}:
						go func(line string) {
							defer func() {
								<-sem
							}()
							method(line)
						}(line)
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	for i := 0; i < len(lines); i += batchSize {
		end := i + batchSize
		if end > len(lines) {
			end = len(lines)
		}
		batchCh <- lines[i:end]
	}

	close(batchCh)

	wg.Wait()

	Interface.Write("Press enter to go Main Menu")
	_, err := fmt.Scanln()
	CLI_Handlers.LogError(err)
}

func ThreadingWithSemaphore(method Method, maxWorkers int, lines []string) {
	batchSize := CalculateBatchSize(len(lines), runtime.NumCPU())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	Interface.Option("Lines per Thread", strconv.Itoa(batchSize))
	Interface.Option("Threading Type", ThreadingType)
	Interface.Option("Threads", strconv.Itoa(maxWorkers))

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(int64(maxWorkers))

	batchCh := make(chan []string, maxWorkers)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			for batch := range batchCh {
				for _, line := range batch {
					if err := sem.Acquire(ctx, 1); err != nil {
						return
					}
					go func(line string) {
						defer sem.Release(1)
						method(line)
					}(line)
				}
			}
		}()
	}

	for i := 0; i < len(lines); i += batchSize {
		end := i + batchSize
		if end > len(lines) {
			end = len(lines)
		}
		batchCh <- lines[i:end]
	}

	close(batchCh)

	wg.Wait()

	Interface.Write("Press enter to go Main Menu")
	_, err := fmt.Scanln()
	CLI_Handlers.LogError(err)
}

func ThreadingWithParallelProducers(method Method, maxWorkers int, lines []string) {
	var wg sync.WaitGroup
	workChan := make(chan string, maxWorkers*5)
	seen := make(map[string]struct{})

	for _, line := range lines {
		seen[line] = struct{}{}
	}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := start; j < len(lines); j += maxWorkers {
				line := lines[j]
				if _, exists := seen[line]; exists {
					delete(seen, line)
					workChan <- line
				}
			}
		}(i)
	}

	var workerWg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for line := range workChan {
				method(line)
			}
		}()
	}

	wg.Wait()
	close(workChan)
	workerWg.Wait()

	Interface.Write("Press enter to go Main Menu")
	_, err := fmt.Scanln()
	CLI_Handlers.LogError(err)
}

func Threading(method Method, maxWorkers int, lines []string) {
	if ThreadingType == "Sentry" {
		ThreadingWithSemaphore(method, maxWorkers, lines)
	} else if ThreadingType == "Guardian" {
		ThreadingWithBatches(method, maxWorkers, lines)
	} else if ThreadingType == "Lympia" {
		ThreadingWithSingleProducer(method, maxWorkers, lines)
	} else if ThreadingType == "Vortex" {
		ThreadingWithParallelProducers(method, maxWorkers, lines)
	} else {
		Interface.ErrorMSG("NullOps", "Invalid Threading Type: "+ThreadingType)
	}
}
