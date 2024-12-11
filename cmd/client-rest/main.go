package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	lorem "github.com/derektata/lorem/ipsum"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	_time := time.Now().In(time.UTC)
	_deadline := timestamppb.New(_time).AsTime().Format(time.RFC3339)
	// random lorem ipsum text for the description and title
	g := lorem.NewGenerator()
	_title := g.Generate(10)
	_description := g.GenerateParagraphs(1)

	var body string

	// create a new todo
	resp, err := http.Post("http://localhost:5000/v1/todo", "application/json", strings.NewReader(fmt.Sprintf(`
		{
			"todo": {
				"title": "%s",
				"description": "%s",
				"deadline": "%s"
			}
		}
	`, _title, _description, _deadline)))
	if err != nil {
		log.Fatalf("failed to create todo: %v", err)
	}

	// read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Create response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// parse id of the created todo
	var created struct {
		Id string `json:"id"`
	}

	err = json.Unmarshal(bodyBytes, &created)
	if err != nil {
		log.Fatalf("failed to parse response body: %v", err)
	}

	// call the read endpoint
	resp, err = http.Get("http://localhost:5000/v1/todo/" + created.Id)
	if err != nil {
		log.Fatalf("failed to read todo: %v", err)
	}
	bodyBytes, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Read response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// call the update endpoint
	req, _ := http.NewRequest("PUT", "http://localhost:5000/v1/todo/"+created.Id, strings.NewReader(fmt.Sprintf(`
		{
			"todo": {
				"title": "%s updated",
				"description": "%s updated",
				"deadline": "%s"
			}
		}
	`, _title, _description, _deadline)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to update todo: %v", err)
	}
	bodyBytes, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Update response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// call read all todos
	resp, err = http.Get("http://localhost:5000/v1/todo/all")
	if err != nil {
		log.Fatalf("failed to read all todos: %v", err)
	}
	bodyBytes, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to read response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("ReadAll response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// call delete endpoint
	req, _ = http.NewRequest("DELETE", "http://localhost:5000/v1/todo/"+created.Id, nil)
	resp, _ = http.DefaultClient.Do(req)
	bodyBytes, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed to delete todo: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}
