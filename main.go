package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type WeatherData struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Rain struct {
		OneH float64 `json:"1h"`
	} `json:"rain"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}
type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Country   string  `json:"country"`
	State     string  `json:"state"`
}

func KelvinToCelc(temp float64) float64 {
	return temp - 273.15
}

func MeterSecToKmHour(speed float64) float64 {
	return speed * 3.6
}

func main() {
	//place to insert apiKey
	apiKey := ""

	cityName := flag.String("name", "warsaw", "cityName")
	flag.Parse()

	geoLink := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&appid=%s", *cityName, apiKey)
	geoResp, err := http.Get(geoLink)
	if err != nil {
		panic(err)
	}
	geoBody, err := io.ReadAll(geoResp.Body)

	var location []Location
	err = json.Unmarshal(geoBody, &location)
	if err != nil {
		panic(err)
	}
	finalLoction := location[0]

	lat := strconv.FormatFloat(finalLoction.Latitude, 'f', -1, 64)
	lon := strconv.FormatFloat(finalLoction.Longitude, 'f', -1, 64)

	link := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", lat, lon, apiKey)

	resp, err := http.Get(link)
	if err != nil {
		panic("Error with getting data from link")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Country: %s\n", weatherData.Sys.Country)
	fmt.Printf("City: %s\n", weatherData.Name)
	fmt.Printf("Temp: %0.2f*C\n", KelvinToCelc(weatherData.Main.Temp))
	fmt.Printf("Feels like: %0.2f*C\n", KelvinToCelc(weatherData.Main.FeelsLike))
	fmt.Printf("Speed of wind: %0.2fkm/h\n", MeterSecToKmHour(weatherData.Wind.Speed))
	fmt.Printf("Rain in 1h: %0.2fmm\n", weatherData.Rain.OneH)

}
