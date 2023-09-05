package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pavva91/gin-gorm-rest/config"
	"github.com/pavva91/gin-gorm-rest/server"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPingPongUnsecured(t *testing.T) {
	var cfg config.ServerConfig

	// Get Configs
	err := cleanenv.ReadConfig("./config/config.yml", &cfg)
	if err != nil {
		log.Println(err)
	}

	apiVersion := fmt.Sprintf("/%s/%s", cfg.Server.ApiPath, cfg.Server.ApiVersion)
	// Build our expected body
	body := gin.H{
		"message": "pong",
	}

	// Grab our router
	router := server.NewRouter(cfg)
	server.MapUrls(cfg)

	// Perform a GET request with that handler.
	// w := performRequest(router, "GET", "/api/v1/ping")
	w := performRequest(router, "GET", apiVersion+"/ping")

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	var response map[string]string
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["message"]
	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["message"], value)
}

func TestPingPongSecuredWithoutJWT(t *testing.T) {
	var cfg config.ServerConfig

	// Get Configs
	err := cleanenv.ReadConfig("./config/config.yml", &cfg)
	if err != nil {
		log.Println(err)
	}

	apiVersion := fmt.Sprintf("/%s/%s", cfg.Server.ApiPath, cfg.Server.ApiVersion)
	// Build our expected body
	body := gin.H{
		"error": "request does not contain an access token",
	}

	// Grab our router
	router := server.NewRouter(cfg)
	// server.MapUrls(cfg)

	// Perform a GET request with that handler.
	// w := performRequest(router, "GET", "/api/v1/ping")
	w := performRequest(router, "GET", apiVersion+"/secured/ping")

	// Assert we encoded correctly,
	// the request gives a 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Convert the JSON response to a map
	var response map[string]string
	err = json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["error"]
	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["error"], value)
}
