package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/oauth2/google"
)

var (
	project  string
	token    string
	location string
	verbose  bool

	client *http.Client
)

const (
	apiEndpoint = "https://reachability.googleapis.com/v1beta1"

	reachable = "REACHABLE"
)

func main() {
	ctx := context.Background()

	flag.StringVar(&project, "project", "", "")
	flag.StringVar(&token, "token", "", "")
	flag.StringVar(&location, "location", "", "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.Parse()

	if project == "" {
		log.Fatalln("Please provide a project name")
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

	testIDs, err := listTests(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(testIDs))
	for _, tt := range testIDs {
		go func(test string) {
			defer wg.Done()

			r, err := rerunTest(test)
			if err != nil {
				log.Fatalf("Error when triggering rerun for %v: %v", test, err)
			}
			printResult(r)
		}(tt)
	}
	wg.Wait()
}

func printResult(r resourceResponse) {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "%v\t", r.Name)
	fmt.Fprintf(buf, "%v\t", r.ReachabilityDetails.Result)

	out := os.Stdout
	if r.ReachabilityDetails.Result != reachable {
		out = os.Stderr
	}
	fmt.Fprintln(out, buf.String())
}

func rerunTest(resource string) (resourceResponse, error) {
	var r resourceResponse

	req, _ := http.NewRequest("POST", apiEndpoint+"/"+resource+":rerun", nil)
	res, err := client.Do(req)
	if err := handleHTTPError(res, err); err != nil {
		return r, err
	}

	var op opResponse
	if err := unmarshalBody(res.Body, &r); err != nil {
		return r, err
	}

	// Poll the long running operation results.
	for {
		req, _ = http.NewRequest("GET", apiEndpoint+"/"+r.Name, nil)
		res, err = client.Do(req)
		if err := unmarshalBody(res.Body, &op); err != nil {
			return r, err
		}
		if op.Done {
			return op.Resource, nil
		}
		// TODO(jbd): Don't loop indefinitely.
		time.Sleep(20 * time.Millisecond)
	}
}

func listTests(ctx context.Context) ([]string, error) {
	req, _ := http.NewRequest("GET", apiEndpoint+"/projects/"+project+"/locations/"+location+"/connectivityTests", nil)
	res, err := client.Do(req)
	if handleHTTPError(res, err) != nil {
		return nil, err
	}

	var r listResponse
	if err := unmarshalBody(res.Body, &r); err != nil {
		return nil, err
	}
	tests := make([]string, 0, len(r.Resources))
	for _, rr := range r.Resources {
		tests = append(tests, rr.Name)
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

func unmarshalBody(r io.Reader, dst interface{}) error {
	all, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Printf("%s", all)
	}
	return json.Unmarshal(all, dst)
}

type listResponse struct {
	Resources []resourceResponse `json:"resources"`
}

type opResponse struct {
	Name     string           `json:"name"`
	Resource resourceResponse `json:"response"`
	Done     bool             `json:"done"`
}

type resourceResponse struct {
	Name                string                      `json:"name"`
	UpdateTime          time.Time                   `json:"updateTime"`
	ReachabilityDetails reachabilityDetailsResponse `json:"reachabilityDetails"`
}

type reachabilityDetailsResponse struct {
	Result     string    `json:"result"`
	VerifyTime time.Time `json:"verifyTime"`
}
