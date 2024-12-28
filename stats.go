package main

import (
	"fmt"
	"sort"
	"time"
)

func getKeyList(dates map[int]int) []int {
	var keys []int
	for key := range dates {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	return keys
}

func determineDayOfWeek(currentDayOfWeek int, daysAgo int) int {
	dayOfTheWeek := (currentDayOfWeek - daysAgo) % 7
	if dayOfTheWeek < 0 {
		dayOfTheWeek = (dayOfTheWeek + 7) % 7
	}
	return int(dayOfTheWeek)
}

func handleCutoffs(cols map[int]column, cutoff, start time.Weekday, week int) {
	if len(cols[week]) > 0 {
		for start < cutoff {
			cols[week][start+1] = -1
			start = start + 1
		}
	} else {
		cols[week] = make(column, 7)
		for start < cutoff {
			cols[week][start+1] = -1
			start = start + 1
		}
	}

}

// TODO: will check tomorrow whether the offsets are working or not they seem well for now
func calculateOffset() int {
	var offset int
	switch time.Now().Weekday() {
	case time.Sunday:
		offset = 6
	case time.Monday:
		offset = 5
	case time.Tuesday:
		offset = 4
	case time.Wednesday:
		offset = 3
	case time.Thursday:
		offset = 2
	case time.Friday:
		offset = 1
	case time.Saturday:
		offset = 0
	}
	return offset
}

// TODO: Figure out a way to get the first and last day of consideration and mark them a negative(invalid value) so that once they are found when rendering it can be switched out
type column []int

func genCols(keys []int, dates map[int]int) map[int]column {
	cols := make(map[int]column)
	today := time.Now().Weekday()
	offset := calculateOffset()
	for _, k := range keys {
		week := int((k + offset) / 7)

		if len(cols[week]) == 0 {
			cols[week] = make(column, 7)
		}
		day := determineDayOfWeek(int(today), k)

		cols[week][day] += dates[k]
	}
	handleCutoffs(cols, time.Saturday, today, 0)
	handleCutoffs(cols, time.Weekday(determineDayOfWeek(int(today), 183))-1, time.Sunday-1, 26)
	return cols
}

func determineAndPrintColour(commit int) {
	if commit >= 10 {
		fmt.Print("\033[48;5;34m \033[0m")
	} else if commit >= 5 {
		fmt.Print("\033[48;5;28m \033[0m")
	} else if commit >= 1 {
		fmt.Print("\033[48;5;22m \033[0m")
	} else if commit == 0 {
		fmt.Print("\033[48;5;17m \033[0m")
	} else {
		fmt.Print(" ")
	}
}

func printStats(cols map[int]column) {
	for i := 0; i < 7; i++ {
		switch i {
		case 1:
			fmt.Print("Mon ")
		case 3:
			fmt.Print("Wed ")
		case 5:
			fmt.Print("Fri ")
		default:
			fmt.Print("    ")
		}
		for j := 26; j >= 0; j-- {
			if len(cols[j]) == 0 {
				determineAndPrintColour(0)
			} else {
				determineAndPrintColour(cols[j][i])
			}
		}
		fmt.Println()
	}
}

// NOTE: to get the week mod the key with 26(weeks in 6 months) or maybe 27 (to adjust for the days of the week) to get the first day mod the

// TODO: Add a mapping to read store all the commit dates to generate image of contributions
func GetStats(email string, repo string) {
	dates := make(map[int]int)
	dates = genDatesMap(email, repo, dates)
	keys := getKeyList(dates)
	cols := genCols(keys, dates)
	printStats(cols)
}
