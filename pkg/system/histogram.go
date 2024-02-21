package main

import "sort"

var histogramLength int = 50

var appHistogram map[string]*histogram = make(map[string]*histogram)

var preTime map[string]int64 = make(map[string]int64)

var unit float64 = 1000 * 60

type histogram struct {
	sum            int
	array          []int
	nonZeroIndexes []int
}

func insertSorted(s []int, e int) []int {
	i := sort.SearchInts(s, e)
	if i < len(s) && s[i] == e {
		return s
	}
	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = e
	return s
}

func updateHistogram(appID string, index int) {
	if appHistogram[appID] == nil {
		appHistogram[appID] = &histogram{
			array:          make([]int, histogramLength),
			nonZeroIndexes: make([]int, 0),
		}
	}
	if index >= histogramLength {
		index = histogramLength - 1
	}
	if appHistogram[appID].array[index] == 0 {
		appHistogram[appID].nonZeroIndexes = insertSorted(appHistogram[appID].nonZeroIndexes, index)
	}
	appHistogram[appID].sum += 1
	appHistogram[appID].array[index] += 1
}

func getWindow(app *Application) (int, int) {
	if IsFixed > 0 || appHistogram[app.AppID] == nil {
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
	if prewarmWindow == 0 && keepAliveWindow == 0 {
		keepAliveWindow += 1
	} else if prewarmWindow == keepAliveWindow {
		prewarmWindow -= 1
		keepAliveWindow += 1
	}
	if prewarmWindow < 2*Minute {
		prewarmWindow = 0
	}
	if keepAliveWindow < 2*Minute {
		keepAliveWindow = 2 * Minute
	}
	prewarmWindow = int(float64(prewarmWindow) * 0.9 * unit)
	keepAliveWindow = int(float64(keepAliveWindow) * 1.1 * unit)
	return prewarmWindow, keepAliveWindow
}

func getPercentage(appID string, time int64) float64 {
	if appHistogram[appID] == nil {
		return 0
	}
	time1 := float64(time) / unit
	sum := appHistogram[appID].sum
	if sum == 0 {
		return 0
	}

	percentage := 0.0
	nonZeroIndexes := appHistogram[appID].nonZeroIndexes
	array := appHistogram[appID].array

	for _, index := range nonZeroIndexes {
		if index <= int(time1) {
			percentage += float64(array[index])
		} else {
			break
		}
	}
	return percentage / float64(sum)
}
