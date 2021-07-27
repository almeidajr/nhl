package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Roster struct {
	Person struct {
		ID       int    `json:"id"`
		FullName string `json:"fullName"`
		Link     string `json:"link"`
	} `json:"person"`
	JerseyNumber string `json:"jerseyNumber"`
	Position     struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Abbreviation string `json:"abbreviation"`
	} `json:"position"`
}

type RostersResponse struct {
	Rosters []Roster `json:"roster"`
}

func GetAllRosters(teamID int) ([]Roster, error) {
	requestUrl := fmt.Sprintf("%s/teams/%d/roster", baseUrl, teamID)
	res, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response RostersResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Rosters, nil
}
