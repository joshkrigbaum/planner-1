package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//Define structures to receive weather forecast from JSON
type current struct {
	Time                 uint    `json:"time"`                 //	1453402675,
	Summary              string  `json:"summary"`              //	"Rain",
	Icon                 string  `json:"icon"`                 //	"rain",
	NearestStormDistance uint    `json:"nearestStormDistance"` //	0,
	PrecipIntensity      float64 `json:"precipIntensity"`      //	0.1685,
	PrecipIntensityError float64 `json:"precipIntensityError"` //	0.0067,
	PrecipProbability    float64 `json:"precipProbability"`    //	1,
	PrecipType           string  `json:"precipType"`           //	"rain",
	Temperature          float64 `json:"temperature"`          //	48.71,
	ApparentTemperature  float64 `json:"apparentTemperature"`  //	46.93,
	Dewpoint             float64 `json:"dewPoint"`             //	47.7,
	Humidity             float64 `json:"humidity"`             //	0.96,
	WindSpeed            float64 `json:"windSpeed"`            //	4.64,
	WindBearing          int     `json:"windBearing"`          //	186,
	Visibility           float64 `json:"visibility"`           //	4.3,
	CloudCover           float64 `json:"cloudCover"`           //	0.73,
	Pressure             float64 `json:"pressure"`             //	1009.7,
	Ozone                float64 `json:"ozone"`                //	328.35
}

type dailyData struct {
	Time                          uint64  `json:"time"`        //	1453402675,
	Summary                       string  `json:"summary"`     //	"Rain",
	Icon                          string  `json:"icon"`        //	"rain",
	SunriseTime                   uint    `json:"sunriseTime"` //	1453391560,
	SunsetTime                    uint    `json:"sunsetTime"`  //	1453424361,
	MoonPhase                     float64 `json:"moonPhase"`   //	0.43
	PrecipIntensity               float64 `json:"precipIntensity"`
	PrecipitationIntensityMax     float64 `json:"precipIntensityMax"`
	PrecipitationIntensityMaxTime float64 `json:"precipIntensityMaxTime"`
	PrecipProbability             float64 `json:"precipProbability"`           //	1,
	PrecipType                    string  `json:"precipType"`                  //	"rain",
	TemperatureHigh               float64 `json:"temperatureHigh"`             //	41.42,
	TemperatureHighTime           uint    `json:"temperatureHighTime"`         //	1453417200
	TemperatureLow                float64 `json:"temperatureLow"`              //	41.42,
	TemperatureLowTime            uint    `json:"temperatureLowTime"`          //	1453417200
	ApparentTemperatureHigh       float64 `json:"apparentTemperatureHigh"`     //	46.93,
	ApparentTemperatureHighTime   float64 `json:"apparentTemperatureHighTime"` //	46.93,
	ApparentTemperatureLow        float64 `json:"apparentTemperatureLow"`      //	46.93,
	ApparentTemperatureLowTime    float64 `json:"apparentTemperatureLowTime"`  //	46.93,
	Dewpoint                      float64 `json:"dewPoint"`                    //	47.7,
	Humidity                      float64 `json:"humidity"`                    //	0.96,
	Pressure                      float64 `json:"pressure"`
	WindSpeed                     float64 `json:"windSpeed"` //	4.64,
	WindGust                      float64 `json:"windGust"`
	WindGustTime                  float64 `json:"windGustTime"`
	WindBearing                   int     `json:"windBearing"` //	186,
	CloudCover                    float64 `json:"cloudCover"`
	UVIndex                       float64 `json:"uvIndex"`
	UVIndexTime                   float64 `json:"uvIndexTime"`
	Visibility                    float64 `json:"visibility"`                 //	4.3,
	Ozone                         float64 `json:"ozone"`                      //	328.35
	TemperatureMin                float64 `json:"temperatureMin"`             //	41.42,
	TemperatureMinTime            uint    `json:"temperatureMinTime"`         //	1453417200
	TemperatureMax                float64 `json:"temperatureMax"`             //	41.42,
	TemperatureMaxTime            uint    `json:"temperatureMaxTime"`         //	1453417200
	ApparentTemperatureMin        float64 `json:"apparentTemperatureMin"`     //	46.93,
	ApparentTemperatureMinTime    float64 `json:"apparentTemperatureMinTime"` //	46.93,
	ApparentTemperatureMax        float64 `json:"apparentTemperatureMax"`     //	46.93,
	ApparentTemperatureMaxTime    float64 `json:"apparentTemperatureMaxTime"` //	46.93,
}

type daily struct {
	Summary string      `json:"summary"` //	"Rain for the hour.",
	Icon    string      `json:"icon"`    //	"rain",
	Data    []dailyData `json:"data"`
}

type alert struct {
	Title       string `json:"title"`       //	"Flood Watch for Mason, WA",
	Time        uint   `json:"time"`        //	1453375020,
	Expires     uint   `json:"expires"`     //	1453407300,
	Description string `json:"description"` //	"...FLOOD WATCH...\n",
	URL         string `json:"uri"`         //	"http:/..."
}

type darkskyForecast struct {
	Latitude  float64 `json:"latitude"`  //	40.47780682531368,
	Longitude float64 `json:"longitude"` //	-86.93875375799722,
	Timezone  string  `json:"timezone"`  //	"America/Indiana/Indianapolis",
	Current   current `json:"currently"`
	Daily     daily
	Alerts    []alert
	Offset    int `json:"offset"` //	-4
} // End of receiving structure for weather forecast

type wotdType struct {
	Word      string
	Pronounce string
	POS       string
	Defs      []string
}

type sound struct {
	wave string `xml: "wav"	json:	"wave"`
	wpr  string `xml:	"wpr"	json:	"wpr"`
}

type entry struct {
	ew        string `xml: "ew"	json: "word"`
	subj      string `xml: "subj"	json: "subject"`
	syllables string `xml: "hw"	json: "syllables"`
	sound     string `xml: "sound"	json:	"sound"`
	pronounce string `xml:	"pr"	json:	"pronounce"`
	pos       string `xml: "fl"	json:	"pos"`
}

// Define receiving structure for WOTD XML
type wotdFormat struct {
	entryList string   `xml: "entry_list"	json: "entryList"`
	Word      string   `xml: "ew" json: "word"`
	Pronounce string   `xml: "pr" json: "pronounce"`
	POS       string   `xml: "fl" json: "pos"`
	Defs      []string `xml: "dt" json: "def"`
}

//Define structure to receive WOTD from XML
type entryList struct {
	XMLName xml.Name `xml:"entry_list"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Entry   struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
		Ew   struct {
			Text string `xml:",chardata"`
		} `xml:"ew"`
		Subj struct {
			Text string `xml:",chardata"`
		} `xml:"subj"`
		Hw struct {
			Text string `xml:",chardata"`
		} `xml:"hw"`
		Sound struct {
			Text string `xml:",chardata"`
			Wav  struct {
				Text string `xml:",chardata"`
			} `xml:"wav"`
			Wpr struct {
				Text string `xml:",chardata"`
			} `xml:"wpr"`
		} `xml:"sound"`
		Pr struct {
			Text string `xml:",chardata"`
		} `xml:"pr"`
		Fl struct {
			Text string `xml:",chardata"`
		} `xml:"fl"`
		In []struct {
			Text string `xml:",chardata"`
			If   struct {
				Text string `xml:",chardata"`
			} `xml:"if"`
		} `xml:"in"`
		Et struct {
			Text string `xml:",chardata"`
			It   []struct {
				Text string `xml:",chardata"`
			} `xml:"it"`
		} `xml:"et"`
		Def struct {
			Text string `xml:",chardata"`
			Vt   struct {
				Text string `xml:",chardata"`
			} `xml:"vt"`
			Date struct {
				Text string `xml:",chardata"`
			} `xml:"date"`
			Sn []struct {
				Text string `xml:",chardata"`
			} `xml:"sn"`
			Dt []struct {
				Text string `xml:",chardata"`
				Sx   struct {
					Text string `xml:",chardata"`
					Sxn  struct {
						Text string `xml:",chardata"`
					} `xml:"sxn"`
				} `xml:"sx"`
				Vi struct {
					Text string `xml:",chardata"`
					It   struct {
						Text string `xml:",chardata"`
					} `xml:"it"`
				} `xml:"vi"`
			} `xml:"dt"`
		} `xml:"def"`
		Uro []struct {
			Text string `xml:",chardata"`
			Ure  struct {
				Text string `xml:",chardata"`
			} `xml:"ure"`
			Sound struct {
				Text string `xml:",chardata"`
				Wav  struct {
					Text string `xml:",chardata"`
				} `xml:"wav"`
				Wpr struct {
					Text string `xml:",chardata"`
				} `xml:"wpr"`
			} `xml:"sound"`
			Pr struct {
				Text string `xml:",chardata"`
			} `xml:"pr"`
			Fl struct {
				Text string `xml:",chardata"`
			} `xml:"fl"`
		} `xml:"uro"`
	} `xml:"entry"`
}

//Define structures to receive configuration from JSON
type configStruct struct {
	Debug                 bool
	DarkSkyKey            string
	Latitude              string
	Longitude             string
	Excludes              string
	WeatherURL            string
	WeatherReloadInterval int
	QotdURL               string
	QotdReloadInterval    int
	WotdURL               string
	WotdReloadInterval    int
	PhotosDir             string
	CSSDirectory          string
	PhotosReloadInterval  int
	TimeCheckInterval     int
	HTMLFile              string
	PhotoDir              string
	PhotoReloadInterval   int
	LogFile               string
	MWrss                 string
	MWurl                 string
	MWkey                 string
} // End of receiving structure for configuration

var forecast darkskyForecast

//var HTMLFile string

func main() {
	log.Println("\n  INFO: Starting Planner Application.\n")

	log.Println("  INFO: Loading Configuration from json/config.json.\n")
	config := getConfig()
	//displayConfig(config)

	log.Println("  INFO: Calling startWeather()")
	go startWeather(config)
	time.Sleep(10 * time.Second)

	log.Println("  INFO: Calling startWOTD()")
	go startWOTD(config)
	time.Sleep(10 * time.Second)

	log.Println("  INFO: Calling startPhotos()")
	go startPhotos(config)
	select {}
}

func startWeather(config configStruct) {
	// Initial Weather load on startup
	log.Println("  INFO: Initial Weather() Load")
	getWeather(config)

	//==================================
	// Repeat Weather load every weatherdReloadInterval
	ticker := time.NewTicker(time.Hour * time.Duration(config.WeatherReloadInterval))
	for range ticker.C {
		log.Println("  INFO: Periodic Weather() Load")
		getWeather(config)
	}
	log.Printf("\n  INFO: *** Error: Exit on range ticker in function startWeather(). ***\n\n")
}

func startWOTD(config configStruct) {
	// Initial WOTD load on startup
	log.Println("  INFO: Initial WOTD() Load")
	getWOTD(config)

	//==================================
	// Repeat WOTD load every wotdReloadInterval
	ticker := time.NewTicker(time.Hour * time.Duration(config.WotdReloadInterval))
	for range ticker.C {
		log.Println("  INFO: Periodic WOTD() Load")
		getWOTD(config)
	}
	log.Printf("\n  INFO: *** Error: Exit on range ticker in function startWOTD(). ***\n\n")
}

func startPhotos(config configStruct) {
	// Initial Photos load on startup
	log.Println("  INFO: Initial Photos() Load")
	getPhotos(config)

	//==================================
	// Repeat WOTD load every wotdReloadInterval
	ticker := time.NewTicker(time.Minute * time.Duration(config.PhotoReloadInterval))
	for range ticker.C {
		log.Println("  INFO: Periodic Photos() Load")
		getPhotos(config)
	}
	log.Printf("\n  INFO: *** Error: Exit on range ticker in function startPhotos(). ***\n\n")
}

//*************************************************************

func getPhotos(config configStruct) {
	cssBytes, err := ioutil.ReadFile(config.CSSDirectory)
	if err != nil {
		log.Fatalln("ReadFile failed w/ err", err)
	}
	css := string(cssBytes)

	rand.Seed(time.Now().Unix())

	deck, err := ioutil.ReadDir(config.PhotoDir)
	if err != nil {
		log.Println("  INFO: ReadDir error:", err)
	}

	index := rand.Intn(len(deck))
	photo := deck[index].Name()

	startStr := "background: url(../photos/"
	stopStr := ") no-repeat center center fixed"
	start := strings.Index(css, startStr)
	stop := strings.Index(css, stopStr) + len(stopStr)
	oldStr := css[start:stop]
	newStr := startStr + photo + stopStr
	css = strings.Replace(css, oldStr, newStr, 1)

	cssFile := []byte(css)
	ioutil.WriteFile(config.CSSDirectory, cssFile, 0644)
}

func getWeather(config configStruct) {
	htmlBytes, err := ioutil.ReadFile(config.HTMLFile)
	if err != nil {
		log.Fatalln("ReadFile failed w/ err", err)
	}
	html := string(htmlBytes)

	darkskyURL := config.WeatherURL + config.DarkSkyKey + "/" + config.Latitude + "," + config.Longitude + "?" + config.Excludes
	forecast = getForecast(darkskyURL)
	forecast.Daily.Data = forecast.Daily.Data[:3]

	startStr := "<span id=\"currentTemp\">"
	stopStr := " &#8457"
	valueStr := string(truncate(forecast.Current.Temperature, 0))
	start := strings.Index(html, startStr)
	stop := strings.Index(html, stopStr) + len(stopStr)
	oldStr := html[start:stop]
	newStr := startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"currentHumidity\">"
	stopStr = " %</span>"
	valueStr = string(truncate(forecast.Current.Humidity*100, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"currentWindSpeed\">"
	stopStr = " mph</span>"
	valueStr = string(truncate(forecast.Current.WindSpeed, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"currentVisibility\">"
	stopStr = " mi.</span>"
	valueStr = string(truncate(forecast.Current.Visibility, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<h2><span id=\"day1\">"
	stopStr = "<!--d1--></span></h2>"
	valueStr = getWeekday(forecast.Daily.Data[0].Time)
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<span id=\"lowTemp1\">"
	stopStr = " &#8457;<!--1--></span>"
	valueStr = string(truncate(forecast.Daily.Data[0].TemperatureLow, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<span id=\"highTemp1\">"
	stopStr = " &#8457;<!--2--></span>"
	valueStr = string(truncate(forecast.Daily.Data[0].TemperatureHigh, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"humidity1\">"
	stopStr = " %<!--1--></span>"
	valueStr = string(truncate(forecast.Daily.Data[0].Humidity*100, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"windspeed1\">"
	stopStr = " mph<!--1--></span>"
	valueStr = string(truncate(forecast.Daily.Data[0].WindSpeed, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"visibility1\">"
	stopStr = " mi.<!--1--></span>"
	valueStr = string(truncate(forecast.Daily.Data[0].Visibility, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<h2><span id=\"day2\">"
	stopStr = "<!--d2--></span></h2>"
	valueStr = getWeekday(forecast.Daily.Data[1].Time)
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<span id=\"lowTemp2\">"
	stopStr = " &#8457;<!--3--></span>"
	valueStr = string(truncate(forecast.Daily.Data[1].TemperatureLow, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<span id=\"highTemp2\">"
	stopStr = " &#8457;<!--4--></span>"
	valueStr = string(truncate(forecast.Daily.Data[1].TemperatureHigh, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"humidity2\">"
	stopStr = " %<!--2--></span>"
	valueStr = string(truncate(forecast.Daily.Data[1].Humidity*100, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"windspeed2\">"
	stopStr = " mph<!--2--></span>"
	valueStr = string(truncate(forecast.Daily.Data[1].WindSpeed, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"visibility2\">"
	stopStr = " mi.<!--2--></span>"
	valueStr = string(truncate(forecast.Daily.Data[1].Visibility, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<h2><span id=\"day3\">"
	stopStr = "<!--d3--></span></h2>"
	valueStr = getWeekday(forecast.Daily.Data[2].Time)
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<span id=\"lowTemp3\">"
	stopStr = " &#8457;<!--5--></span>"
	valueStr = string(truncate(forecast.Daily.Data[2].TemperatureLow, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<span id=\"highTemp3\">"
	stopStr = " &#8457;<!--6--></span>"
	valueStr = string(truncate(forecast.Daily.Data[2].TemperatureHigh, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"humidity3\">"
	stopStr = " %<!--3--></span>"
	valueStr = string(truncate(forecast.Daily.Data[2].Humidity*100, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"windspeed3\">"
	stopStr = " mph<!--3--></span>"
	valueStr = string(truncate(forecast.Daily.Data[2].WindSpeed, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<br> <span id=\"visibility3\">"
	stopStr = " mi.<!--3--></span>"
	valueStr = string(truncate(forecast.Daily.Data[2].Visibility, 0))
	start = strings.Index(html, startStr)
	stop = strings.Index(html, stopStr) + len(stopStr)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	htmlFile := []byte(html)
	ioutil.WriteFile(config.HTMLFile, htmlFile, 0644)

	log.Println("  INFO: Finished getWeather()\n")
}

func getConfig() configStruct {
	// Read config.json file and assign values to struct config ===================================
	var config configStruct
	configFile, err := ioutil.ReadFile("json/config.json")
	if err != nil {
		log.Printf("  INFO: File error reading json/config.json: %v\n", err)
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(configFile), &config)
	if err != nil {
		log.Fatalln("  FATAL: Error unmarshaling json/config.json:", err)
	}

	return config
}

func getForecast(darkskyURL string) darkskyForecast {
	var forecast darkskyForecast

	_, err := os.Stat("json/darksky.json")
	if !os.IsNotExist(err) {
		err := os.Remove("json/darksky.json")
		if err != nil {
			log.Println("  FATAL: Error removing json/darksky.json.")
			log.Println("   FATAL: Exiting program.")
			os.Exit(1)
		} else {
			log.Println("  INFO: json/darksky.json has been deleted.")
		}
	} else {
		log.Println("  INFO: json/darksky.json does not exist.")
	}
	log.Println("  INFO: json/darksky.json has been loaded.")

	data, err := http.Get(darkskyURL)
	if err != nil {
		log.Println("    FATAL: Error reading ", darkskyURL)
		log.Println("    FATAL: Exiting program.")
		os.Exit(1)
	}

	// Convert raw data to []bytes.
	dataBYTES, err := ioutil.ReadAll(data.Body)
	data.Body.Close()
	if err != nil {
		log.Printf("  FATAL: Error reading body of %s: %v", darkskyURL, err)
		log.Println("  FATAL: Exiting program.")
		os.Exit(1)
	}

	_, err = os.Stat("json/darksky.json")
	if !os.IsNotExist(err) {
		err := os.Remove("json/darksky.json")
		if err != nil {
			log.Println("  FATAL: Error removing json/darksky.json.")
			log.Println("  FATAL: Exiting program.")
			os.Exit(1)
		}
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, dataBYTES, "", "    ")
	if err != nil {
		log.Println("  INFO: Error pretty printing JSON")
	}

	darksky, err := os.OpenFile("json/darksky.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("  INFO: Error opening 'json/darksky.json':", err)
	}
	defer darksky.Close()

	darksky.WriteString(prettyJSON.String())

	// Start of unmarshal
	weatherData, err := ioutil.ReadFile("json/darksky.json")
	err = json.Unmarshal(weatherData, &forecast)
	if err != nil {
		log.Println("  FATAL: Error unmarshaling json/config.json:", err)
		log.Println("  FATAL: Exiting program.")
		os.Exit(1)
	}

	log.Println("  INFO: Finished getForecastData()")
	return forecast
}

func getWOTD(config configStruct) {
	rssURL := config.MWrss
	data, err := http.Get(rssURL)
	if err != nil {
		log.Println("  INFO: Error on http.Get(rssURL):", err)
	}
	dataBYTES, err := ioutil.ReadAll(data.Body)
	if err != nil {
		fmt.Println("  INFO: Error on ioutil.ReadAll(data.Body):", err)
	}
	data.Body.Close()
	rss := string(dataBYTES)

	word := extract(rss, "<![CDATA[", "]]>")
	wotdURL := config.MWurl + word + "?key=" + config.MWkey
	data, err = http.Get(wotdURL)
	if err != nil {
		log.Println("  INFO: Error on http.Get(wotdURL):", err)
	}
	dataBYTES, err = ioutil.ReadAll(data.Body)
	if err != nil {
		log.Println("  INFO: Error on ioutil.ReadAll(data.Body):", err)
	}
	data.Body.Close()
	var def1 entryList
	err = xml.Unmarshal(dataBYTES, &def1)

	var wotdInfo wotdType

	wotdInfo.Word = def1.Entry.ID
	wotdInfo.Pronounce = def1.Entry.Pr.Text
	wotdInfo.POS = def1.Entry.Fl.Text
	numdefs := len(def1.Entry.Def.Dt)

	x := 0
	for x < numdefs {
		wotdInfo.Defs = append(wotdInfo.Defs, string(def1.Entry.Def.Dt[x].Text))
		x++
	}

	htmlBytes, err := ioutil.ReadFile("planner.html")
	if err != nil {
		log.Fatalln("ReadFile failed w/ err", err)
	}
	html := string(htmlBytes)

	startStr := "<span id=\"word\">"
	stopStr := ":&nbsp;<!--w1--></span>"
	valueStr := string(wotdInfo.Word)
	start := strings.Index(html, startStr)
	//fmt.Println("1. start =", start)
	stop := strings.Index(html, stopStr) + len(stopStr)
	//fmt.Println("1. stop =", stop)
	oldStr := html[start:stop]
	newStr := startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<span id=\"pronounce\">[&nbsp;"
	stopStr = "&nbsp;]<!--w2--></span>"
	valueStr = "&nbsp;" + string(wotdInfo.Pronounce)
	start = strings.Index(html, startStr)
	//fmt.Println("2. start =", start)
	stop = strings.Index(html, stopStr) + len(stopStr)
	//fmt.Println("2. stop =", stop)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	startStr = "<span id=\"pos\">"
	stopStr = "<!--w3--></span>"
	valueStr = "&nbsp;" + string(wotdInfo.POS)
	start = strings.Index(html, startStr)
	//fmt.Println("3. start =", start)
	stop = strings.Index(html, stopStr) + len(stopStr)
	//fmt.Println("3. stop =", stop)
	oldStr = html[start:stop]
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	valueStr = ""
	d := 0
	for d < numdefs {
		startStr = "<span id=\"defs\">"
		stopStr = "<!--w4--></span>"
		//extract(src string, startStr string, stopStr string) string
		cleanerdef := erase(string(wotdInfo.Defs[d]), ":")
		valueStr = valueStr + "&nbsp;&nbsp;&nbsp;Definition " + strconv.Itoa(d+1) + ") &nbsp;" + cleanerdef + "<br>"
		//fmt.Printf("ValueStr = %s\n", valueStr)
		start = strings.Index(html, startStr)
		//fmt.Println("4. start =", start)
		stop = strings.Index(html, stopStr) + len(stopStr)
		//fmt.Println("4. stop =", stop)
		oldStr = html[start:stop]
		//tempStr := ""
		d++
	}

	//fmt.Println("tempStr =", tempStr)
	newStr = startStr + valueStr + stopStr
	html = strings.Replace(html, oldStr, newStr, 1)

	//fmt.Println("html =", html)

	htmlFile := []byte(html)
	ioutil.WriteFile("planner.html", htmlFile, 0644)

	log.Println("  INFO: Finished getWOTD()\n")
}

func truncate(x interface{}, p int) string {
	//fmt.Println("x =", x)
	fmtStr := "%." + strconv.Itoa(p) + "f"
	xstring := fmt.Sprintf(fmtStr, x)
	return xstring
}

func convertToInt(d float64) string {
	xstring := fmt.Sprintf("%0.f", d*100)
	return xstring
}

func getWeekday(UnixTime uint64) string {
	//UnixTime := forecast.Daily.Data[0].Time
	timeStr := strconv.FormatUint(UnixTime, 10)
	reportedTime, _ := strconv.ParseInt(timeStr, 10, 64)
	tm := time.Unix(reportedTime, 0)
	//tm := getTime(timeStr)
	day0 := fmt.Sprintf("%v", tm.Weekday())
	return day0
}

func getTime(reportedStr string) time.Time {
	reportedTime, _ := strconv.ParseInt(reportedStr, 10, 64)
	tm := time.Unix(reportedTime, 0)
	return tm
}

func extract(src string, startStr string, stopStr string) string {
	start := strings.Index(src, startStr) + len(startStr)
	if start == -1 {
		return "NotFound"
	}

	stop := strings.Index(src, stopStr)
	if stop == -1 {
		return "NotFound"
	}

	found := src[start:stop]
	return found
}

func erase(src string, ch string) string {
	if len(ch) > 1 || len(ch) == 0 {
		return "erase() failed on ch"
	}
	p := 0
	for p < len(src) {
		if string(src[p]) == ch {
			src = src[:p] + src[p+1:]
		}
		p++
	}
	return src
}

func findOriginal(src string, startStr string, stopStr string, newStr string, occurance int) string {
	n := 0
	var start int
	var stop int

	for n < occurance {
		start := strings.Index(src, startStr)
		if start == -1 {
			return "NotFound"
		}
		stop = strings.Index(src, stopStr) + len(stopStr)
		if stop == -1 {
			return "NotFound"
		}
		n++
	}

	OriginalStr := src[start:stop]
	return OriginalStr
}
