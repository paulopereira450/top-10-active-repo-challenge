# Challenge Rank and List the most active repositories

## Considerations or choices taken during the development of the script

1. **Definition of the Activity Score**: To calculate the score of the repositories I used this three metrics:


**Commit Amount**: The sum of all the commits done to the repository.

**Activity Decay Rate**: Applied a decay rate to each activity, making older commits less influential on the score.

**Unique Contributor Count**: Sums all the contributors for each repository.

**Activity Weight Calculation**: Calculates a weighted score for each repository based on the changes applied on the commits. In order to not give as much value to activities with large changes, since the repositories could have different sizes, a logarithmic scaling was applied to reduce the value of this metric. The procedure is applied to the files changed, additions and deletions.


2. **Usage of a heap**: Since it was a large amount of repositories to validate and sort, and to avoid unnecessary procedures I used an heap to get and order the repositories with most score. This avoids the process of sorting all the repositories and will only manage the required amount of repositories to show. If the list exceeds the required size, it will remove the repo with the lowest score on the heap.

# Improvements
- Since the recency factor is using the current time I needed to reduce a lot the decay rate. Since the script was using the sample provided, I could get the last commit date and use it as the date to calculate the recency factor. Instead I opted to use the current date for the script handle different samples.

## Usage

### Prerequisites

- **Go 1.16+** installed

### Running the Script

To run the script, use the following command:

```bash
go run main.go
```

### Script result with provided file (commits.csv)
```bash
Top 10 Most Active Repositories:

        Repo            Score
1.      repo250         277394.26
2.      repo518         111425.49
3.      repo126         99416.24
4.      repo740         85066.82
5.      repo795         52641.48
6.      repo127         45470.33
7.      repo476         45072.71
8.      repo982         43048.05
9.      repo703         33716.89
10.     repo831         24396.54
```

### Note

> **Important**: The score values will change since the calculation of the scores is done using the current time to calculate the recency factor. Although the order will maintain if the file was not changed.