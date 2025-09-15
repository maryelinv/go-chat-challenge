package main

import "strings"

func isStockCmd(s string) bool {
	return strings.HasPrefix(strings.TrimSpace(s), "/stock=")
}

func parseStockCode(s string) string {
	code := strings.TrimPrefix(strings.TrimSpace(s), "/stock=")
	return strings.ToLower(strings.TrimSpace(code))
}
