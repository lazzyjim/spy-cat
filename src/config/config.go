package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Common struct {
	DB Postgres `json:"db"`
}

func (c Common) Validation() error {
	err := c.DB.Validate()
	if err != nil {
		return fmt.Errorf("section db has an error, err:%s", err)
	}
	return nil
}

func Fetch(path string) (*Common, error) {

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("path does not contains file with configuration, err: " + err.Error())
	}

	co := Common{}
	err = json.Unmarshal(content, &co)
	if err != nil {
		return nil, fmt.Errorf("incorrect configuration file structure, err: " + err.Error())
	}

	err = co.Validation()
	if err != nil {
		return nil, fmt.Errorf("incorrect configuration data, err: " + err.Error())
	}
	return &co, nil
}
