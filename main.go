package main

import ( 
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

func main() {
  http.HandleFunc("/", weatherHandler)
  http.ListenAndServe(":5000", nil)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	count := r.FormValue("count")

	body, err := getWeatherResponseBody(count)

  if err != nil {
    panic(err)
  }

  openWeather := OpenWeather{}
  err = json.Unmarshal(body, &openWeather)
  if err != nil {
    panic(err)
  }

//console print
for i := range openWeather.List  {
    fmt.Printf("\n\nList[%d]: %v", i, 
      openWeather.List[i])
    fmt.Printf("\nWeather in %s, %s is %.2f",
      openWeather.List[i].Name, openWeather.List[i].Country.CountryName,
    openWeather.List[i].Temperature.NormalisedCurrentTemp())
  }

//prints to webserver
  for i := range openWeather.List {
      fmt.Fprintf(w, "\nWeather in %s is %.2f", 
         openWeather.List[i].Name, openWeather.List[i].Country.CountryName,
 openWeather.List[i].Temperature.NormalisedCurrentTemp())
	}
}

//prints 1 blob
  // fmt.Printf("\nopenWeather: %v", openWeather)
	// fmt.Printf("\nList[0]: %v", openWeather.List[0])
	// fmt.Printf("\nName: %s", openWeather.List[0].Name)
	// fmt.Printf("\nCurrentTemp: %.2f",
	// 	openWeather.List[0].Temperature.NormalisedCurrentTemp())

func (t TemperatureDetails) NormalisedCurrentTemp() float64 {
  return t.CurrentTemp - 273.15
}

type OpenWeather struct {
  List []City `json:"list"`
}

type City struct {
  Temperature TemperatureDetails `json:"main"`
  Name    string  `json:"name"`
  Country CountryDetails `json:"sys"`
}

type CountryDetails struct {
  CountryName string  `json:"country"`
}

type TemperatureDetails struct {
  CurrentTemp float64 `json:"temp"`
  MaxTemp     float64 `json:"temp_max"`
}

func getWeatherResponseBody(count string) ([]byte, error) {
  url:= fmt.Sprintf("http://api.openweathermap.org/data/2.5/find?appid=0a12b8f2f0dd011ed6085cb995ff61b4&lat=-37.81&lon=144.96&cnt=%s", count)

  	resp, err := http.Get(url)
    if err != nil {
		  fmt.Printf("Error getting weather: %v", err)
		  return []byte(""), err
	}
  //defer for memory leakage
  defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
	  fmt.Printf("Error reading weather body: %v", err)
	  return []byte(""), err
	}
	
  return body, nil

}