package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	server := httptest.NewServer(newHandler(NewStore()))
	defer server.Close()

	response, err := http.Get(server.URL + "/health")
	if err != nil { t.Fatal(err) }
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.StatusCode)
	}
}

func TestCreateTask(t *testing.T) {
	server := httptest.NewServer(newHandler(NewStore()))
	defer server.Close()

	response, err := http.Post(server.URL+"/tasks", "application/json", strings.NewReader(`{"title":"Learn Go"}`))
	if err != nil { t.Fatal(err) }
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", response.StatusCode)
	}
}

func TestCreateTaskRejectsEmptyTitle(t *testing.T) {
	server := httptest.NewServer(newHandler(NewStore()))
	defer server.Close()

	response, err := http.Post(server.URL+"/tasks", "application/json", strings.NewReader(`{"title":"  "}`))
	if err != nil { t.Fatal(err) }
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.StatusCode)
	}
}

func TestDeleteMissingTask(t *testing.T) {
	server := httptest.NewServer(newHandler(NewStore()))
	defer server.Close()

	request, _ := http.NewRequest(http.MethodDelete, server.URL+"/tasks/999", nil)
	response, err := http.DefaultClient.Do(request)
	if err != nil { t.Fatal(err) }
	defer response.Body.Close()
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.StatusCode)
	}
}
