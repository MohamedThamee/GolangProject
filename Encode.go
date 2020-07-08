package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}
type Author struct {
	Firstname  string `json:"firstname"`
	Secondname string `json:"secondname"`
}

var datas TotalStuct

var books []Book

func main() {
	r := mux.NewRouter()
	books = append(books, Book{ID: "1", Isbn: "789", Title: "asdf", Author: &Author{Firstname: "Thamee", Secondname: "Muthari"}})
	books = append(books, Book{ID: "2", Isbn: "789", Title: "asdf", Author: &Author{Firstname: "Sdfgf", Secondname: "asfgf"}})
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/{id}", getLocationRes).Methods("GET")

	log.Println("Server is ready to handle requests at", 5000)
	log.Fatal(http.ListenAndServe(":5000", r))

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
	json.NewEncoder(w).Encode(&Book{})

}

type User1 struct {
	EatDrinkData NeededArr
}
type User2 struct {
	PetrolData NeededArr
}
type User3 struct {
	ShoppingData NeededArr
}

var bodyString1 NeededArr
var bodyString2 NeededArr
var bodyString3 NeededArr

func getLocationRes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	// for _, item := range books {
	// 	if item.ID == params["id"] {
	// 		json.NewEncoder(w).Encode(item)
	// 		return
	// 	}
	// }
	// json.NewEncoder(w).Encode(&Book{})

	uurl := fmt.Sprintf("https://geocoder.api.here.com/search/6.2/geocode.json?languages=en-US&maxresults=1&searchtext=%s&app_id=DG9nbJeGww0CjNyc3cet&app_code=2JPpm7n9_AWwiOUasdlb4g", params["id"])
	//fmt.Println(uurl)

	resp, err := http.Get(uurl)

	if err != nil {
		//c <- "We could not reach:" + url
		fmt.Println("Error Occured:\n")

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
		fmt.Printf("API Response as struct %s %s\n", s1, s2)
		// json.NewEncoder(w).Encode(s1)
		// json.NewEncoder(w).Encode(s2)

		//var ret string = SendLoc(s1, s2)
		//fmt.Printf("Return Value", ret)

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
		fmt.Print("Return Value", bodyString1)
		//u := User{petrolData: bodyString3, EatDrinkData: bodyString1, shoppingData: bodyString2}
		u1 := User1{EatDrinkData: bodyString1}
		u2 := User2{PetrolData: bodyString3}
		u3 := User3{ShoppingData: bodyString2}
		//datas = append(datas, NeededArr{petrolData: bodyString3})
		//datas := TotalStuct{Eat: bodyString3, Petrol: bodyString1, Shopping: bodyString2}
		//datas = append(datas, TotalStuct{Eat: bodyString3, Petrol: bodyString1, Shopping: bodyString2})
		json.NewEncoder(w).Encode(u1)
		json.NewEncoder(w).Encode(u2)
		json.NewEncoder(w).Encode(u3)

		// json.NewEncoder(w).Encode(bodyString2)
		// json.NewEncoder(w).Encode(bodyString3)

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

		bodyString1 = eat.Results.Items[0:3]
		//js, _ := json.Marshal(EatDrink{})

		//content, err1 := json.Marshal(bodyByt)
		//fmt.Print("content", bodyString)
		//return "bodyString"

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
			fmt.Printf(" %s Shopping is %d mts away from your location\n", shopping.Results.Items[i].Title, shopping.Results.Items[i].Distance)
		}

		bodyString2 = shopping.Results.Items[0:3]

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

		bodyString3 = pstation.Results.Items[0:3]

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

type NeededArr []struct {
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
}

type TotalStuct struct {
	Eat []struct {
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
			Text       string      `json:"text"`
			Label      string      `json:"label"`
			IsOpen     bool        `json:"isOpen"`
			Structured interface{} `json:"structured"`
		} `json:"openingHours"`
		ChainIds []string `json:"chainIds,omitempty"`
	} `json:"Eat"`
	Petrol []struct {
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
			Text       string      `json:"text"`
			Label      string      `json:"label"`
			IsOpen     bool        `json:"isOpen"`
			Structured interface{} `json:"structured"`
		} `json:"openingHours"`
		ChainIds []string `json:"chainIds,omitempty"`
	} `json:"petrol"`
	Shopping []struct {
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
			Text       string      `json:"text"`
			Label      string      `json:"label"`
			IsOpen     bool        `json:"isOpen"`
			Structured interface{} `json:"structured"`
		} `json:"openingHours"`
		ChainIds []string `json:"chainIds,omitempty"`
	} `json:"shopping"`
}
