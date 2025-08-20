package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetFaviconIco(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/favicon.ico", GetFaviconIco)
	
	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/favicon.ico", nil)
	
	// Perform the request
	router.ServeHTTP(w, req)
	
	// Check status code
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGreetUser(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", GreetUser)
	
	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	
	// Perform the request
	router.ServeHTTP(w, req)
	
	// Check status code
	assert.Equal(t, http.StatusOK, w.Code)
}