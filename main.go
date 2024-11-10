package main

import (
	"container/heap"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	commitFilePath       = "commits.csv"
	decayRate            = 0.0001
	expectedColumnsCount = 6
	topListNumber        = 10
)

type RepoStats struct {
	Name  string
	Score float64
}

type RepoMetric struct {
	Commits            int
	UniqueContributors map[string]bool
	ActivityWeight     float64
}

type RepoStatsHeap []RepoStats

func (h RepoStatsHeap) Len() int {
	return len(h)
}

func (h RepoStatsHeap) Less(i, j int) bool {
	return h[i].Score < h[j].Score
}

func (h RepoStatsHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *RepoStatsHeap) Push(x interface{}) {
	*h = append(*h, x.(RepoStats))
}

func (h *RepoStatsHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	file, err := os.Open(commitFilePath)
	if err != nil {
		log.Fatalf("Error while reading the commits file: %s", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read header of commits file: %s", err)
	}

	if len(header) != expectedColumnsCount {
		log.Fatalf("Header with unexpected amount of columns: %s", strings.Join(header, ", "))
	}

	currentTime := time.Now()
	repoMetrics := make(map[string]*RepoMetric)

	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("Error reading row: %s", err)
		}

		if len(row) != expectedColumnsCount {
			fmt.Printf("Row with unexpected amount of columns: %s. Skipping row.\n", strings.Join(row, ", "))
			continue
		}

		repoName := row[2]
		if repoName == "" {
			fmt.Printf("Invalid repository on row: %s. Skipping row.\n", strings.Join(row, ", "))
			continue
		}

		if _, found := repoMetrics[repoName]; !found {
			repoMetrics[repoName] = &RepoMetric{
				Commits:            0,
				UniqueContributors: make(map[string]bool),
				ActivityWeight:     0,
			}
		}

		i, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			fmt.Printf("Invalid timestamp received: %s. Skipping row.\n", strings.Join(row, ", "))
			continue
		}
		timeDiff := currentTime.Sub(time.Unix(i, 0)).Hours()
		recencyFactor := math.Exp(-decayRate * timeDiff)

		repoMetrics[repoName].Commits++
		repoMetrics[repoName].UniqueContributors[row[1]] = true
		repoMetrics[repoName].ActivityWeight += getCommitWeight(row, recencyFactor)
	}

	h := &RepoStatsHeap{}
	heap.Init(h)

	for repoName, metrics := range repoMetrics {
		heap.Push(h, RepoStats{
			Name:  repoName,
			Score: (float64(metrics.Commits) + float64(len(metrics.UniqueContributors))) * metrics.ActivityWeight,
		})

		if h.Len() > topListNumber {
			heap.Pop(h)
		}
	}

	var topRepos []RepoStats
	for h.Len() > 0 {
		topRepos = append(topRepos, heap.Pop(h).(RepoStats))
	}

	fmt.Printf("Top 10 Most Active Repositories:\n\n")
	fmt.Println("\tRepo\t\tScore")
	for i := 0; i < len(topRepos); i++ {
		repoStats := topRepos[len(topRepos)-i-1]
		fmt.Printf("%d.\t%s\t\t%.2f\n", i+1, repoStats.Name, repoStats.Score)
	}
}

func getCommitWeight(row []string, recencyFactor float64) float64 {
	return (1 + getWeightOfActivityChange(row[3]) + getWeightOfActivityChange(row[4]) + getWeightOfActivityChange(row[5])) * recencyFactor
}

func getWeightOfActivityChange(s string) float64 {
	value, _ := strconv.Atoi(s)

	return math.Log(1 + float64(value))
}
