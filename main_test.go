package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const jsonFilePath = "fixtures/test.json"

func TestHandlerOK(t *testing.T) {
	assert := assert.New(t)

	jsonData, err := ioutil.ReadFile(jsonFilePath)
	assert.NoError(err)

	for _, method := range allMethods {
		spy := &writerSpy{}
		logger := log.New(spy, "", log.LstdFlags)
		handler := jsonFileWithCORSOptionsHandlerFunc(jsonFilePath, logger)

		req := httptest.NewRequest(method, "/foo", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		header := resp.Header()
		for name, value := range headers {
			assert.Equal(value, header.Get(name))
		}
		if method != http.MethodOptions {
			assert.Equal("application/json", header.Get("Content-Type"))
			if method != http.MethodHead {
				assert.Equal(string(jsonData), resp.Body.String(), method)
			}
		}
		assert.Equal(http.StatusOK, resp.Code)
		assert.Equal(1, spy.calls)
	}
}

func TestHandlerResponseNotFoundWhenFileDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	spy := &writerSpy{}
	logger := log.New(spy, "", log.LstdFlags)
	handler := jsonFileWithCORSOptionsHandlerFunc("/tmp/foo", logger)

	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(http.StatusNotFound, resp.Code)
	assert.Equal(1, spy.calls)
}

type writerSpy struct {
	mu    sync.Mutex
	calls int
}

func (l *writerSpy) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.calls++
	return 0, nil
}
