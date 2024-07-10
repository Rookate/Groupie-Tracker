package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"common"
)

var VarArtists []Artist

var VarDateLocation DatesLocation

func fetchJSONData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	t, errTmpl := template.ParseFiles(tmpl)
	if errTmpl != nil {
		http.Error(w, "Error parsing template "+tmpl, http.StatusInternalServerError)
		return
	}

	if errExec := t.Execute(w, data); errExec != nil {
		fmt.Fprintln(os.Stderr, errExec.Error())
		return
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorPage(w, http.StatusNotFound, "Page "+r.URL.Path[1:]+" Not Found.", "The page you requested was not found.")
		return
	}

	if len(r.URL.Query()) > 0 {
		ErrorPage(w, http.StatusBadRequest, "Invalid URL", "The URL you requested is not valid.")
		return
	}

	groupieResp, errGGet := http.Get(ApiURL)
	if errGGet != nil {
		fmt.Fprintln(os.Stderr, errGGet.Error())
		ErrorPage(w, http.StatusInternalServerError, "Failed to fetch Groupie API", "There was a problem fetching "+ApiURL+".")
		return
	}
	defer groupieResp.Body.Close()

	if groupieResp.StatusCode != http.StatusOK {
		ErrorPage(w, http.StatusInternalServerError, "Failed to fetch Groupie API", "There was a problem fetching "+ApiURL+". Page returned code "+strconv.Itoa(groupieResp.StatusCode))
		return
	}

	renderTemplate(w, "templates/index.html", ArtistsArray)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	artistId := r.URL.Path[len("/detail/"):]
	id, errAtoi := strconv.Atoi(artistId)

	if errAtoi != nil {
		fmt.Fprintln(os.Stderr, errAtoi.Error())
		ErrorPage(w, http.StatusInternalServerError, "Not an ID", "Word "+artistId+" is not a valid ID.")
		return
	}
	if id <= 0 || id > len(ArtistsArray) {
		ErrorPage(w, http.StatusInternalServerError, "ID not valid", "ID "+artistId+" does not exist.")
		return
	}

	artist := ArtistsArray[id-1]

	if _, present := InfoMap[artist.Id]; !present {
		var wg sync.WaitGroup
		var location Locations
		var relation Relation

		wg.Add(2)
		go func() {
			defer wg.Done()
			if errLJson := fetchJSONData(artist.Locations, &location); errLJson != nil {
				fmt.Fprintln(os.Stderr, errLJson.Error())
				ErrorPage(w, http.StatusInternalServerError, "Failed to fetch artist's Locations info.", "There was an error fetching "+artist.Locations+".")
			}
		}()

		go func() {
			defer wg.Done()
			if errRJson := fetchJSONData(artist.Relations, &relation); errRJson != nil {
				fmt.Fprintln(os.Stderr, errRJson.Error())
				ErrorPage(w, http.StatusInternalServerError, "Failed to fetch artist's Relations info.", "There was an error fetching "+artist.Relations+".")
			}
		}()

		wg.Wait()

		var concerts []string
		for _, loc := range location.Location {
			formattedLoc := loc

			formattedLoc = strings.Replace(formattedLoc, "-", ", ", -1)
			formattedLoc = strings.Replace(formattedLoc, "_", " ", -1)
			formattedLoc = common.Capitalize(formattedLoc)
			formattedLoc = strings.Replace(formattedLoc, "Usa", "USA", 1)
			formattedLoc = strings.Replace(formattedLoc, "Uk", "UK", 1)

			concerts = append(concerts, formattedLoc+": "+strings.Join(relation.DatesLocations[loc], ", "))
		}

		InfoMap[artist.Id] = ArtistInfo{
			TheArtist:     artist,
			TheirConcerts: concerts,
		}
	}

	renderTemplate(w, "templates/detail.html", InfoMap[artist.Id])
}

func Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	results := SearchArtists(query, VarArtists, VarDateLocation)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func SearchArtists(query string, artists []Artist, datesLocation DatesLocation) []map[string]interface{} {
	query = strings.ToLower(query)
	var filteredArtists []map[string]interface{}

	for _, artist := range ArtistsArray {
		found := false
		if strings.Contains(strings.ToLower(artist.Name), query) ||
			strings.Contains(strings.ToLower(strings.Join(artist.Members, " ")), query) ||
			strings.Contains(strings.ToLower(artist.FirstAlbum), query) ||
			strings.Contains(fmt.Sprint(artist.CreationDate), query) {
			filteredArtists = append(filteredArtists, createArtistResult(artist))
			continue
		}

		for _, loc := range ArtistLocations[artist.Id] {
			if strings.Contains(loc, strings.ToLower(query)) {
				filteredArtists = append(filteredArtists, createArtistResult(artist))
				found = true
				break
			}
		}
		if found {
			continue
		}

		for _, date := range ArtistDates[artist.Id] {
			if strings.Contains(date, query) {
				filteredArtists = append(filteredArtists, createArtistResult(artist))
				break
			}
		}
	}

	return filteredArtists
}

func ErrorPage(w http.ResponseWriter, code int, msg, desc string) {
	renderTemplate(w, "templates/error.html", ErrorData{
		Code:        code,
		Message:     msg,
		Description: desc,
	})
}

func createArtistResult(artist Artist) map[string]interface{} {
	return map[string]interface{}{
		"type":  "artist/band",
		"id":    artist.Id,
		"name":  artist.Name,
		"image": artist.Image,
	}
}

/*
		if _, present := ArtistLocations[artist.Id]; !present {
			var location Locations
			if err := fetchJSONData(artist.Locations, &location); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				break
			}
			ArtistLocations[artist.Id] = location.Location
		}

		for _, loc := range ArtistLocations[artist.Id] {
			if strings.Contains(loc, strings.ToLower(query)) {
				filteredArtists = append(filteredArtists, createArtistResult(artist))
				found = true
				break
			}
		}
		if found {
			found = false
			continue
		}

		if _, present := ArtistDates[artist.Id]; !present {
			var dates Dates
			if err := fetchJSONData(artist.ConcertDates, &dates); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				break
			}
			ArtistDates[artist.Id] = dates.DatesValues
		}
			for _, date := range ArtistDates[artist.Id] {
		if strings.Contains(date, query) {
			filteredArtists = append(filteredArtists, createArtistResult(artist))
			break
		}
	}


*/
