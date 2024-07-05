package config

import (
	"errors"
	"fmt"
)

const DefaultPostgresPort = 5432
const DefaultPostgresSchema = "public"

type Postgres struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DB       string `json:"name_db"`
	Username string `json:"username"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
	Settings struct {
		MaximumPoolSize   int `json:"maximumPoolSize"`
		ConnectionTimeout int `json:"connectionTimeout"`
	} `json:"settings"`
}

type CatsApi struct {
	Host string `json:"host"`
}

func (r *CatsApi) Validate() error {
	if r.Host == "" {
		return errors.New("property `host` is required")
	}
	return nil
}

func (r *Postgres) Validate() error {

	if r.Port == 0 {
		r.Port = DefaultPostgresPort
	}

	if len(r.Schema) == 0 {
		r.Schema = DefaultPostgresSchema
	}

	if len(r.DB) == 0 {
		return errors.New("property `name_db` is required")
	}

	if len(r.Username) == 0 {
		return errors.New("property `username` is required")
	}

	if len(r.Password) == 0 {
		return errors.New("property `password` is required")
	}

	return nil
}

func (r *Postgres) ConnectionSource() string {
	return fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?search_path=%v&sslmode=disable",
		r.Username,
		r.Password,
		r.Host,
		r.Port,
		r.DB,
		r.Schema,
	)
}
