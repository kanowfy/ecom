package handlers

import (
	"strconv"
	"strings"
)

func parseDate(date string, part_idx int) uint32 {
	parts := strings.Split(date, "/")
	part, _ := strconv.ParseUint(parts[part_idx], 10, 32)
	return uint32(part)
}

func getMonthPart(date string) uint32 {
	return parseDate(date, 0)
}

func getYearPart(date string) uint32 {
	return parseDate(date, 1)
}
