package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	ApiURL = "https://groupietrackers.herokuapp.com/api"
)

var (
	ArtistsArray    []Artist
	ArtistLocations map[uint][]string
	ArtistDates     map[uint][]string
	InfoMap         map[uint]ArtistInfo
)

type ErrorData struct {
	Code                 int
	Message, Description string
}
type DatesLocation struct {
	Index []DateLocation `json:"index"`
}

type DateLocation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Artist struct {
	Id           uint     `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	FirstAlbum   string   `json:"firstAlbum"`
	CreationDate uint     `json:"creationDate"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}
type Locations struct {
	Id       uint     `json:"id"`
	Location []string `json:"locations"`
	Dates    string   `json:"dates"`
}
type Dates struct {
	Id          uint     `json:"id"`
	DatesValues []string `json:"dates"`
}
type Relation struct {
	Id             uint                `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type ArtistInfo struct {
	TheArtist     Artist
	TheirConcerts []string
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs

		fmt.Print("\nProgram ")
		switch sig {
		case syscall.SIGINT:
			fmt.Print("interrupted.")
		case syscall.SIGTERM:
			fmt.Print("terminated.")
		}
		fmt.Println(" Exiting.")
		os.Exit(0)
	}()

	port := "8080"
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           nil,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	http.HandleFunc("/", Home)
	http.HandleFunc("/detail/", Detail)
	http.HandleFunc("/search", Search)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server starting on http://localhost:" + port)

	InfoMap = make(map[uint]ArtistInfo, len(ArtistsArray))
	ArtistLocations = make(map[uint][]string, len(ArtistsArray))
	ArtistDates = make(map[uint][]string, len(ArtistsArray))

	if errJson := fetchJSONData(ApiURL+"/artists", &ArtistsArray); errJson != nil {
		fmt.Fprintln(os.Stderr, errJson.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, artist := range ArtistsArray {
			var dates Dates
			if err := fetchJSONData(artist.ConcertDates, &dates); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				break
			}
			ArtistDates[artist.Id] = dates.DatesValues
		}
	}()

	go func() {
		defer wg.Done()
		for _, artist := range ArtistsArray {
			var location Locations
			if err := fetchJSONData(artist.Locations, &location); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				break
			}
			ArtistLocations[artist.Id] = location.Location
		}
	}()

	wg.Wait()

	if errSrv := server.ListenAndServe(); errSrv != nil {
		log.Fatal(errSrv)
	}
}
