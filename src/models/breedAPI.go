package models

import (
	"encoding/json"
	"net/http"
	"spy-cat/src/db"
	"time"
)

const URL = "https://api.thecatapi.com/v1/breeds"

func ValidateBreed(breedName string) (bool, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(URL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var breeds []db.Breed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return false, err
	}

	for _, breed := range breeds {
		if breed.Name == breedName {
			return true, nil
		}
	}
	return false, nil
}
