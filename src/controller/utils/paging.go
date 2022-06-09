package utils

import (
	"strconv"
)

func GetTotalPage(count string) int {

	c, _ := strconv.Atoi(count)

	if count == "" {
		return 1
	}

	if maxPageCnt := c % 10; maxPageCnt == 0 {
		return c / 10
	} else {
		return (c / 10) + 1
	}
}

func GetOffset(pS string, pN string) string {
	pageSize, _ := strconv.Atoi(pS)
	pageNo, _ := strconv.Atoi(pN)

	offset := strconv.Itoa((pageNo - 1) * pageSize)
	if pageNo == 1 {
		offset = "0"
	}
	return offset
}
