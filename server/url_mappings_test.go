package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutesMappings(t *testing.T) {
	apiVersion := "/api/v1"

	assert.EqualValues(t, 0, len(router.Routes()))

	mapUrls(apiVersion)

	routes := router.Routes()

	expectedNumberOfRoutes := 16
	assert.EqualValues(t, expectedNumberOfRoutes, len(routes))

	expectedHttpMethod := http.MethodPost
	expectedUrl := apiVersion + "/token"
	assert.EqualValues(t, expectedHttpMethod, routes[0].Method)
	assert.EqualValues(t, expectedUrl, routes[0].Path)

	expectedHttpMethod = http.MethodPost
	expectedUrl = apiVersion + "/users"
	assert.EqualValues(t, expectedHttpMethod, routes[1].Method)
	assert.EqualValues(t, expectedUrl, routes[1].Path)
}
