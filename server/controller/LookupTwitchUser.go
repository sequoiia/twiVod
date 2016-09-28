package controller

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"github.com/sequoiia/twiVod/models"
	"fmt"
	"encoding/json"
)

var HttpClient *http.Client

func LookupTwitchUser(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	log.Printf("Looking up %s\n", args["TwitchUsername"])

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/users/%s?on_site=1", args["TwitchUsername"]), nil); if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Server error."))
		log.Println("Something went wrong during a lookup request for a Twitch user.")
	}

	req.Header.Set("Client-ID", "twiVod - github.com/sequoiia/twivod")

	resp, err := HttpClient.Do(req)

	defer resp.Body.Close()

	var TwitchUser models.TwitchUser

	err = json.NewDecoder(resp.Body).Decode(&TwitchUser)

	w.Write([]byte(fmt.Sprintf("%s", TwitchUser)))
}

func LookupTwitchUsers(w http.ResponseWriter, r *http.Request) {
	var limit, offset, query string
	status, value := setQueryVariables(r, "limit"); if status {
		limit = value
	} else {
		limit = "8"
	}

	status, value = setQueryVariables(r, "offset"); if status {
		offset = value
	} else {
		offset = "0"
	}

	status, value = setQueryVariables(r, "query"); if status {
		query = value
	}

	if query == "" {
		w.WriteHeader(400)
		w.Write([]byte("No query provided."))
	} else {
		log.Printf("Searching for users named %s", query)

		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/search/channels?q=%s&offset=%s&limit=%s", query, offset, limit), nil); if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Server error."))
			log.Println("Something went wrong during a lookup request for a Twitch user.")
		}

		req.Header.Set("Client-ID", "twiVod - github.com/sequoiia/twivod")

		resp, err := HttpClient.Do(req)

		defer resp.Body.Close()

		var TwitchUsers models.TwitchUsersSearch

		err = json.NewDecoder(resp.Body).Decode(&TwitchUsers)

		w.Write([]byte(fmt.Sprintf("%s", TwitchUsers)))
	}
}

func LookupTwitchUserVods(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)

	var limit, offset string
	status, value := setQueryVariables(r, "limit"); if status {
		limit = value
	} else {
		limit = "8"
	}

	status, value = setQueryVariables(r, "offset"); if status {
		offset = value
	} else {
		offset = "0"
	}

	log.Printf("Looking up VODS from %s\n", args["TwitchUsername"])

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/channels/%s/videos?limit=%s&offset=%s&broadcasts=true&on_site=1", args["TwitchUsername"], limit, offset), nil); if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Server error."))
		log.Println("Something went wrong during a lookup request for a Twitch user.")
	}

	req.Header.Set("Client-ID", "twiVod - github.com/sequoiia/twivod")

	resp, err := HttpClient.Do(req)

	defer resp.Body.Close()

	var TwitchUserVods models.TwitchUserVods

	err = json.NewDecoder(resp.Body).Decode(&TwitchUserVods)

	w.Write([]byte(fmt.Sprintf("%s", TwitchUserVods)))
}

func setQueryVariables(r *http.Request, valueName string) (bool, string) {
	if r.URL.Query()[valueName] != nil {
		return true, r.URL.Query()[valueName][0]
	} else {
		return false, ""
	}
}