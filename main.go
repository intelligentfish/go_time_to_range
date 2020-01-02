package main

import (
	"fmt"
	"time"
)

// Time2Range 时间区间切割
// 将开始时间切割为整天(快照区间)和时间片段(非快照区间)
func Time2Range(start, end time.Time) (snapshotTimeRange /*快照区间*/, nonSnapshotTimeRange /*非快照区间*/ []time.Time) {
	now := time.Now()
	time000000 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if start == time000000 || start.After(time000000) {
		// 开始时间大于等于今日凌晨，均不能使用快照
		nonSnapshotTimeRange = append(nonSnapshotTimeRange, start, end)
		return
	}

	// 开始日期所在的起点
	start000000 := time.Date(start.Year(), start.Month(), start.Day(),
		0, 0, 0, 0, time.Local)
	// 开始日期所在的终点
	start235959 := time.Date(start.Year(), start.Month(), start.Day(),
		23, 59, 59, 0, time.Local)
	// 结束日期所在的起点
	end000000 := time.Date(end.Year(), end.Month(), end.Day(),
		0, 0, 0, 0, time.Local)
	// 结束日期所在的终点
	end235959 := time.Date(end.Year(), end.Month(), end.Day(),
		23, 59, 59, 0, time.Local)
	// 是否可使用快照
	canSnapshot := func() bool {
		return start.Year()+1 < end.Year() /*跨年*/ ||
			(start000000.AddDate(0, 0, 1).Year() <= end000000.Year() &&
				start000000.AddDate(0, 0, 1).YearDay() < end000000.YearDay()) /*超过一天*/
	}
	if start == start000000 {
		// 开始日期对齐
		if end == end235959 {
			// 结束日期也对齐
			snapshotTimeRange = append(snapshotTimeRange, start, end)
		} else if end.Before(end235959) {
			// 结束日期未对齐
			if canSnapshot() {
				// 有快照可用
				snapshotTimeRange = append(snapshotTimeRange, start,
					end235959.AddDate(0, 0, -1))
			}
			// 末尾区间查询原始数据
			nonSnapshotTimeRange = append(nonSnapshotTimeRange, end000000, end)
		}
	} else if start.After(start000000) {
		// 开始日期未对齐
		if end == end235959 {
			// 结束日期对齐
			if canSnapshot() {
				nonSnapshotTimeRange = append(nonSnapshotTimeRange, start, start235959)
				snapshotTimeRange = append(snapshotTimeRange,
					start000000.AddDate(0, 0, 1),
					end)
			} else {
				nonSnapshotTimeRange = append(nonSnapshotTimeRange, start, end)
			}
		} else if end.Before(end235959) {
			// 结束日期未对齐
			if canSnapshot() {
				snapshotTimeRange = append(snapshotTimeRange, start000000.AddDate(0, 0, 1),
					end235959.AddDate(0, 0, -1))
				nonSnapshotTimeRange = append(nonSnapshotTimeRange, start, start235959, end000000, end)
			} else {
				nonSnapshotTimeRange = append(nonSnapshotTimeRange, start, end)
			}
		}
	}
	return
}

func main() {
	type TimePair struct {
		First  time.Time
		Second time.Time
	}
	timePairs := []TimePair{
		{
			time.Date(2019, 12, 31, 0, 0, 0, 0, time.Local),
			time.Date(2019, 12, 31, 23, 59, 59, 0, time.Local),
		},
		{
			time.Date(2019, 12, 31, 0, 0, 0, 0, time.Local),
			time.Date(2020, 12, 31, 23, 59, 59, 0, time.Local),
		},
		{
			time.Date(2019, 12, 31, 0, 0, 0, 0, time.Local),
			time.Date(2030, 1, 2, 23, 59, 58, 0, time.Local),
		},
		{
			time.Date(2019, 12, 31, 0, 0, 0, 0, time.Local),
			time.Date(2019, 12, 31, 23, 59, 58, 0, time.Local),
		},
		{
			time.Date(2019, 12, 31, 0, 0, 1, 0, time.Local),
			time.Date(2019, 12, 31, 23, 59, 59, 0, time.Local),
		},
		{
			time.Date(2019, 12, 31, 0, 0, 1, 0, time.Local),
			time.Date(2019, 12, 31, 23, 59, 58, 0, time.Local),
		},
		{
			time.Date(2019, 12, 31, 1, 0, 0, 0, time.Local),
			time.Date(2020, 1, 1, 1, 59, 58, 0, time.Local),
		},
	}
	for i := 0; i < len(timePairs); i++ {
		snapshotTimeRange, nonSnapshotTimeRange := Time2Range(timePairs[i].First, timePairs[i].Second)
		fmt.Println(i, "快照区间:", snapshotTimeRange, "非快照区间:", nonSnapshotTimeRange)
	}
}
