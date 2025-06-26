package geocode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type nominatimResp []struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func Geocode(ctx context.Context, address string) (lat, lon float64, err error) {
	base := "https://nominatim.openstreetmap.org/search"
	q := url.Values{"q": {address}, "format": {"json"}, "limit": {"1"}}
	req, err := http.NewRequestWithContext(ctx, "GET", base+"?"+q.Encode(), nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "food-delivery-app-server http://localhost:8080")

	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("geocode API returned %d", res.StatusCode)
	}

	var out nominatimResp
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		return
	}
	if len(out) == 0 {
		err = fmt.Errorf("no result for %q", address)
		return
	}

	if lat, err = strconv.ParseFloat(out[0].Lat, 64); err != nil {
		return
	}
	if lon, err = strconv.ParseFloat(out[0].Lon, 64); err != nil {
		return
	}
	return
}
