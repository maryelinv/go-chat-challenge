package stooq

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func FetchQuote(code string) (float64, error) {
	url := fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", code)
	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	r := csv.NewReader(res.Body)
	records, err := r.ReadAll()
	if err != nil || len(records) < 2 {
		return 0, errors.New("bad csv")
	}

	closeStr := records[1][6]
	if closeStr == "N/D" {
		return 0, errors.New("no data")
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(closeStr), 64)
	return val, err
}
