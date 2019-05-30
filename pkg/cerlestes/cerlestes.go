package cerlestes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	moonPriceUrl = "https://ore.cerlestes.de/data/overview/moon-10000002.json"
	orePriceUrl  = "https://ore.cerlestes.de/data/overview/ore-10000002.json"
)

func PriceList() (map[int64]float64, error) {
	var prices struct {
		Date     time.Time
		Overview string
		Prices   map[string]struct {
			Bp98 float64
		}
	}

	ps := map[int64]float64{}

	data, err := fetch(moonPriceUrl)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &prices); err != nil {
		return nil, err
	}

	for k, v := range prices.Prices {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, err
		}

		ps[id] = v.Bp98 / 100
	}

	data, err = fetch(orePriceUrl)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &prices); err != nil {
		return nil, err
	}

	for k, v := range prices.Prices {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, err
		}

		ps[id] = v.Bp98 / 100
	}

	return ps, nil
}

func fetch(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
