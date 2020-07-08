package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter the location")
	in.Scan()
	var aaa string = in.Text()

	uurl := fmt.Sprintf("https://geocoder.api.here.com/search/6.2/geocode.json?languages=en-US&maxresults=1&searchtext=%s&app_id=DG9nbJeGww0CjNyc3cet&app_code=2JPpm7n9_AWwiOUasdlb4g", aaa)
	//	fmt.Println(uurl)
	links := []string{
		uurl,
	}

	checkUrls(links)
}

//var thamee string = "Hiiii"

func checkUrls(urls []string) {
	c := make(chan string)
	var wg sync.WaitGroup

	for _, link := range urls {
		wg.Add(1)
		go checkUrl(link, c, &wg)
	}

	go func() {

		wg.Wait()
		close(c)
	}()

	for msg := range c {
		fmt.Println(msg)
	}
}

func checkUrl(url string, c chan string, wg *sync.WaitGroup) {
	defer (*wg).Done()
	resp, err := http.Get(url)

	if err != nil {
		c <- "We could not reach:" + url
	} else {
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		// Convert response body to string
		_ = string(bodyBytes)
		//fmt.Println("API Response as String:\n" + bodyString)

		// Convert response body to Todo struct
		var todoStruct Todo
		json.Unmarshal(bodyBytes, &todoStruct)
		//fmt.Printf("API Response as struct %+v\n", todoStruct)

		var lat float64 = (todoStruct.Response.View[0].Result[0].Location.NavigationPosition[0].Latitude)
		var long float64 = (todoStruct.Response.View[0].Result[0].Location.NavigationPosition[0].Longitude)
		var s1 string = fmt.Sprintf("%f", lat)
		var s2 string = fmt.Sprintf("%f", long)

		m1 := func() {
			SendLoc(s1, s2)
		}
		m2 := func() {
			SendLocPetrol(s1, s2)
		}
		m3 := func() {
			SendLocShopping(s1, s2)
		}

		Parallelize(m1, m2, m3)

	}

}
func Parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}

func SendLocShopping(lat string, long string) {
	var lotlong string = lat + "%2C" + long
	var locurl string = fmt.Sprintf("https://places.demo.api.here.com/places/v1/browse?at=%s&cat=shopping&app_id=xGBkLAFbpw2Q1rR5X9wj&app_code=jOOYUM_Bu2I9BIBDd8KEQg", lotlong)
	resp, err := http.Get(locurl)
	if err != nil {
		fmt.Println("Error is generated")
	} else {
		defer resp.Body.Close()
		bodyByt, _ := ioutil.ReadAll(resp.Body)

		var shopping EatDrink
		json.Unmarshal(bodyByt, &shopping)
		fmt.Println("For Shopping Category :")
		for i := range shopping.Results.Items[0:3] {
			fmt.Printf("     %s Shopping is %d mts away from your location\n", shopping.Results.Items[i].Title, shopping.Results.Items[i].Distance)
		}

	}

}

func SendLocPetrol(lat string, long string) {
	var lotlong string = lat + "%2C" + long
	var locurl string = fmt.Sprintf("https://places.demo.api.here.com/places/v1/browse?at=%s&cat=petrol-station&app_id=xGBkLAFbpw2Q1rR5X9wj&app_code=jOOYUM_Bu2I9BIBDd8KEQg", lotlong)
	resp, err := http.Get(locurl)
	if err != nil {
		fmt.Println("Error is generated")
	} else {
		defer resp.Body.Close()
		bodyByt, _ := ioutil.ReadAll(resp.Body)

		var pstation EatDrink
		json.Unmarshal(bodyByt, &pstation)
		fmt.Println("For Petrol Station Category :")
		for i := range pstation.Results.Items[0:3] {
			fmt.Printf("     %s Petrol station is %d mts away from your location\n", pstation.Results.Items[i].Title, pstation.Results.Items[i].Distance)
		}

	}

}

func SendLoc(lat string, long string) {
	var lotlong string = lat + "%2C" + long
	var locurl string = fmt.Sprintf("https://places.demo.api.here.com/places/v1/browse?at=%s&cat=eat-drink&app_id=xGBkLAFbpw2Q1rR5X9wj&app_code=jOOYUM_Bu2I9BIBDd8KEQg", lotlong)
	resp, err := http.Get(locurl)
	if err != nil {
		fmt.Println("Error is generated")
	} else {
		defer resp.Body.Close()
		bodyByt, _ := ioutil.ReadAll(resp.Body)

		// Convert response body to string
		//var bodyStr string = string(bodyByt)

		var eat EatDrink
		json.Unmarshal(bodyByt, &eat)
		fmt.Println("For Eat Drink Category :")
		//fmt.println()
		for i := range eat.Results.Items[0:3] {
			fmt.Printf("     %s Restaurant is %d mts away from your location\n", eat.Results.Items[i].Title, eat.Results.Items[i].Distance)
		}
		//fmt.print(eat.Results.Items[0].Title)
		//fmt.Println(eat.Results.Items[0].Distance)
		//fmt.Println(eat.Results.Items[0].Title)
		//fmt.Printf("API Response as struct %+v\n", todoStruct)

		//Needed -*******
		//fmt.Println(eat.Results.Items[0:3]) ********

		// fmt.Println("Testing********************:\n" + thamee)
		// intSlice := []int{1, 2, 3, 4, 5}
		// fmt.Printf("Slice: %v\n", intSlice)

		// first := intSlice[0:3]
		// fmt.Printf("First element: %d\n", first)

	}

}

type Todo struct {
	Response struct {
		MetaInfo struct {
			Timestamp string `json:"Timestamp"`
		} `json:"MetaInfo"`
		View []struct {
			Type   string `json:"_type"`
			ViewID int    `json:"ViewId"`
			Result []struct {
				Relevance    float64 `json:"Relevance"`
				MatchLevel   string  `json:"MatchLevel"`
				MatchQuality struct {
					City float64 `json:"City"`
				} `json:"MatchQuality"`
				Location struct {
					LocationID      string `json:"LocationId"`
					LocationType    string `json:"LocationType"`
					DisplayPosition struct {
						Latitude  float64 `json:"Latitude"`
						Longitude float64 `json:"Longitude"`
					} `json:"DisplayPosition"`
					NavigationPosition []struct {
						Latitude  float64 `json:"Latitude"`
						Longitude float64 `json:"Longitude"`
					} `json:"NavigationPosition"`
					MapView struct {
						TopLeft struct {
							Latitude  float64 `json:"Latitude"`
							Longitude float64 `json:"Longitude"`
						} `json:"TopLeft"`
						BottomRight struct {
							Latitude  float64 `json:"Latitude"`
							Longitude float64 `json:"Longitude"`
						} `json:"BottomRight"`
					} `json:"MapView"`
					Address struct {
						Label          string `json:"Label"`
						Country        string `json:"Country"`
						State          string `json:"State"`
						County         string `json:"County"`
						City           string `json:"City"`
						PostalCode     string `json:"PostalCode"`
						AdditionalData []struct {
							Value string `json:"value"`
							Key   string `json:"key"`
						} `json:"AdditionalData"`
					} `json:"Address"`
				} `json:"Location"`
			} `json:"Result"`
		} `json:"View"`
	} `json:"Response"`
}

type EatDrink struct {
	Results struct {
		Next  string `json:"next"`
		Items []struct {
			Position      []float64 `json:"position"`
			Distance      int       `json:"distance"`
			Title         string    `json:"title"`
			AverageRating int       `json:"averageRating"`
			Category      struct {
				ID     string `json:"id"`
				Title  string `json:"title"`
				Href   string `json:"href"`
				Type   string `json:"type"`
				System string `json:"system"`
			} `json:"category"`
			Icon     string        `json:"icon"`
			Vicinity string        `json:"vicinity"`
			Having   []interface{} `json:"having"`
			Type     string        `json:"type"`
			Href     string        `json:"href"`
			Tags     []struct {
				ID    string `json:"id"`
				Title string `json:"title"`
				Group string `json:"group"`
			} `json:"tags,omitempty"`
			ID           string `json:"id"`
			OpeningHours struct {
				Text       string `json:"text"`
				Label      string `json:"label"`
				IsOpen     bool   `json:"isOpen"`
				Structured []struct {
					Start      string `json:"start"`
					Duration   string `json:"duration"`
					Recurrence string `json:"recurrence"`
				} `json:"structured"`
			} `json:"openingHours,omitempty"`
			ChainIds         []string `json:"chainIds,omitempty"`
			AlternativeNames []struct {
				Name     string `json:"name"`
				Language string `json:"language"`
			} `json:"alternativeNames,omitempty"`
		} `json:"items"`
	} `json:"results"`
	Search struct {
		Context struct {
			Location struct {
				Position []float64 `json:"position"`
				Address  struct {
					Text        string `json:"text"`
					House       string `json:"house"`
					Street      string `json:"street"`
					PostalCode  string `json:"postalCode"`
					District    string `json:"district"`
					City        string `json:"city"`
					County      string `json:"county"`
					StateCode   string `json:"stateCode"`
					Country     string `json:"country"`
					CountryCode string `json:"countryCode"`
				} `json:"address"`
			} `json:"location"`
			Type string `json:"type"`
			Href string `json:"href"`
		} `json:"context"`
	} `json:"search"`
}
