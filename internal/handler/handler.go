// Package handler 提供HTTP 接口处理层。
package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// parseId 解析路径参数 id,必须为正整数。
func parseId(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

// parseDateRange 解析 start_date/end_date 查询参数(格式 2006-01-02)。
// 返回值直接用于查询:过滤条件为 pubtime >= start 且 pubtime < end;
// end_date 会加 1 天,把"闭区间日期"转成开区间时间戳。
// 任一参数非法时返回错误;参数为空表示不过滤该侧。
func parseDateRange(startStr, endStr string) (time.Time, time.Time, error) {
	var start, end time.Time
	var err error

	if startStr != "" {
		if start, err = time.ParseInLocation("2006-01-02", startStr, time.Local); err != nil {
			return start, end, err
		}
	}
	if endStr != "" {
		if end, err = time.ParseInLocation("2006-01-02", endStr, time.Local); err != nil {
			return start, end, err
		}
		end = end.AddDate(0, 0, 1) // 闭区间转开区间:含 end_date 当天
	}
	return start, end, nil
}
