package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"time"

	"golang.org/x/oauth2/google"
)

var (
	project  string
	token    string
	location string
	testIDs  []string

	client *http.Client
)

const apiEndpoint = "https://reachability.googleapis.com/v1beta1"

func main() {
	ctx := context.Background()
	flag.StringVar(&project, "project", "", "")
	flag.StringVar(&token, "token", "", "")
	flag.StringVar(&location, "location", "", "")
	flag.Parse()

	if project == "" {
		log.Fatalln("Please provide a project name.")
	}
	if location == "" {
		location = "global"
	}

	// List all matching tests and rerun.
	var err error
	client, err = google.DefaultClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if len(testIDs) == 0 {
		testIDs, err = listTests(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, tt := range testIDs {
		if err := rerunTest(tt); err != nil {
			log.Fatalf("Error when triggering rerun for %v: %v", tt, err)
		}
	}
}

func rerunTest(id string) error {
	req, _ := http.NewRequest("POST", apiEndpoint+"/projects/"+project+"/locations/"+location+"/connectivityTests/"+id+":rerun", nil)
	res, err := client.Do(req)
	return handleHTTPError(res, err)
}

func listTests(ctx context.Context) ([]string, error) {
	req, _ := http.NewRequest("GET", apiEndpoint+"/projects/"+project+"/locations/"+location+"/connectivityTests", nil)
	res, err := client.Do(req)
	if handleHTTPError(res, err) != nil {
		return nil, err
	}
	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var r listResponse
	if err := json.Unmarshal(all, &r); err != nil {
		return nil, err
	}
	tests := make([]string, 0, len(r.Resources))
	for _, rr := range r.Resources {
		tests = append(tests, path.Base(rr.Name))
	}
	return tests, nil
}

func handleHTTPError(res *http.Response, err error) error {
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error handling response, status code = %v", res.StatusCode)
	}
	return nil
}

type listResponse struct {
	Resources []resourceResponse `json:"resources"`
}

type resourceResponse struct {
	Name                string                      `json:"name"`
	UpdateTime          time.Time                   `json:"updateTime"`
	ReachabilityDetails reachabilityDetailsResponse `json:"reachabilityDetails"`
}

type reachabilityDetailsResponse struct {
	Result string `json:"result"`
}
