package main

import (
  "testing"
  "regexp"
  "net/http/httptest"
)

func TestServerRunning(t *testing.T) {
  request := httptest.NewRequest("http://localhost:8000", "/", nil)
}
