package main

import "sort"

var histogramLength int = 100

type histogram struct {
	sum            int
	array          []int
	nonZeroIndexes []int
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func insertSorted(s []int, e int) []int {
	i := sort.SearchInts(s, e)
	if i < len(s) && s[i] == e {
		// e already in s at s[i], don't insert again
		return s
	}
	// Insert e at s[i], move others to the right
	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = e
	return s
}

func updateHistogram(appID string, index int, value int) {
	if appHistogram[appID] == nil {
		appHistogram[appID] = &histogram{
			array:          make([]int, histogramLength),
			nonZeroIndexes: make([]int, 0),
		}
	}

	appHistogram[appID].array[index] += value
	if appHistogram[appID].array[index] > 0 && !contains(appHistogram[appID].nonZeroIndexes, index) {
		// 将新的非零索引插入到有序列表中
		appHistogram[appID].nonZeroIndexes = insertSorted(appHistogram[appID].nonZeroIndexes, index)
	}
}

func getWindow(app *Application) (int, int) {
	if IsFixed > 0 {
		return defaultPreWarmTime, defaultKeepAliveTime
	}
	sum := appHistogram[app.AppID].sum
	if sum <= SumLimit {
		return defaultPreWarmTime, defaultKeepAliveTime
	}
	prewarmWindow, keepAliveWindow := 0, 0
	sum1, sum2 := 0, 0
	nonZeroIndexes := appHistogram[app.AppID].nonZeroIndexes
	array := appHistogram[app.AppID].array

	for _, index := range nonZeroIndexes {
		sum1 += array[index]
		if float64(sum1) >= leftBound*float64(sum) {
			prewarmWindow = index
			if float64(sum1) >= leftBound2*float64(sum) {
				prewarmWindow = 0
			}
			break
		}
	}
	for i := len(nonZeroIndexes) - 1; i >= 0; i-- {
		index := nonZeroIndexes[i]
		sum2 += array[index]
		if float64(sum-sum2) <= (1.0-rightBound)*float64(sum) {
			keepAliveWindow = index // 注意：这里应该是 index 而不是 i
		} else {
			break
		}
	}
	return prewarmWindow * 60 * 1000, keepAliveWindow * 60 * 1000
}
