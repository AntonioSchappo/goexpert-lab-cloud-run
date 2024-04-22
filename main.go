package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	Location struct {
		Cep      string `json:"cep"`
		Location string `json:"localidade"`
	}

	WeatherApiResp struct {
		Location struct {
			Name string `json:"name"`
		} `json:"location"`
		Current Current `json:"current"`
	}

	Current struct {
		Temp_C float64 `json:"temp_c"`
		Temp_F float64 `json:"temp_f"`
	}
)

const (
	weatherApiKey = "58550850a19a4f7fb6a132328242104"
)

func main() {
	http.HandleFunc("GET /{cep}", handle)
	http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	if valid := validCep(cep); !valid {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	respCep, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		if respCep.StatusCode == http.StatusNotFound {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		http.Error(w, "can not get location", http.StatusInternalServerError)
		return
	}

	defer respCep.Body.Close()
	var l Location
	err = json.NewDecoder(respCep.Body).Decode(&l)
	if err != nil {
		http.Error(w, "can not decode location", http.StatusInternalServerError)
		return
	}

	respWeather, err := http.Get("http://api.weatherapi.com/v1/current.json?q=" + l.Location + "&key=" + weatherApiKey)
	if respWeather.StatusCode != http.StatusOK || err != nil {
		fmt.Println("http://api.weatherapi.com/v1/current.json?q=" + l.Location + "&key=" + weatherApiKey)

		http.Error(w, "can not get weather", respWeather.StatusCode)
		return
	}
	defer respWeather.Body.Close()

	var weather WeatherApiResp
	err = json.NewDecoder(respWeather.Body).Decode(&weather)
	currentTemp, err := getCurrentTemp(weather)
	if err != nil {
		http.Error(w, "can not decode weather", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentTemp)
}

func validCep(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	for _, c := range cep {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// func getLocation(cep string) (string, error) {
// 	resp, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()
// 	return
// }

func getCurrentTemp(w WeatherApiResp) (Current, error) {
	// w.Current.Temp_K = w.Current.Temp_C + 273.15
	return w.Current, nil
}
