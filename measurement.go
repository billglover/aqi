package aqi

import (
	"context"
	"time"
)

// Measurement is an air quality measurement for a single location.
type Measurement struct {
	AQI int
}

// Latest takes a context and a station identifier and returns the
// latest available air quality reading for the station. An error is returned
// if unable to get the latest measurement.
func (c *Client) Latest(ctx context.Context, station string) (Measurement, error) {
	var m Measurement
	req, err := c.NewRequest("GET", "/feed/"+station+"/", nil)
	if err != nil {
		return m, err
	}

	var fr feedResponse
	_, err = c.Do(ctx, req, &fr)
	if err != nil {
		return m, err
	}

	m.AQI = fr.Data.AQI

	return m, nil
}

// FeedResposne is the response returned from the AQI public API. It contains
// location, overal AQI reading and a breakdown of the various air quality
// measures if known.
type feedResponse struct {
	Status string `json:"status"`
	Data   struct {
		AQI          int `json:"aqi"`
		Idx          int `json:"idx"`
		Attributions []struct {
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"attributions"`
		City struct {
			Geo  []float64 `json:"geo"`
			Name string    `json:"name"`
			URL  string    `json:"url"`
		} `json:"city"`
		Dominentpol string `json:"dominentpol"`
		IAQI        struct {
			Co struct {
				V float64 `json:"v"`
			} `json:"co"`
			No2 struct {
				V float64 `json:"v"`
			} `json:"no2"`
			O3 struct {
				V float64 `json:"v"`
			} `json:"o3"`
			Pm10 struct {
				V int `json:"v"`
			} `json:"pm10"`
			Pm25 struct {
				V int `json:"v"`
			} `json:"pm25"`
			So2 struct {
				V float64 `json:"v"`
			} `json:"so2"`
			W struct {
				V float64 `json:"v"`
			} `json:"w"`
			Wg struct {
				V float64 `json:"v"`
			} `json:"wg"`
		} `json:"iaqi"`
		Time struct {
			S  string `json:"s"`
			Tz string `json:"tz"`
			V  int    `json:"v"`
		} `json:"time"`
		Debug struct {
			Sync time.Time `json:"sync"`
		} `json:"debug"`
	} `json:"data"`
}
