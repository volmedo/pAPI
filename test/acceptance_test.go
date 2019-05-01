package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

const (
	serverURL = "http://localhost"
	apiRoot   = serverURL + "/v1/"
)

var statusCodesAndStrings = map[int]string{
	http.StatusCreated: "Created",
}

type Client struct {
	httpClient *http.Client
	lastResp   *http.Response
}

func newClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: time.Second * 10},
	}
}

func (c *Client) doRequest(method, url, body string) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}

	c.lastResp, err = c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error reading response: %s", err)
	}

	return nil
}

func (c *Client) lastResponse() *http.Response {
	return c.lastResp
}

func (c *Client) iCreateANewPaymentDescribedInJSONAs(jsonPayment *gherkin.DocString) error {
	err := c.doRequest("POST", apiRoot+"payments", jsonPayment.Content)
	if err != nil {
		return fmt.Errorf("Error doing request: %s", err)
	}

	return nil
}

func (c *Client) iGetAResponse(respStatus string) error {
	resp := c.lastResponse()
	if resp == nil {
		return errors.New("Nil response")
	}

	status, ok := statusCodesAndStrings[resp.StatusCode]
	if !ok || !strings.EqualFold(respStatus, status) {
		return fmt.Errorf("Expected status code %s but got %s", status, respStatus)
	}

	return nil
}

func (c *Client) theResponseContainsAPaymentDescribedInJSONAs(jsonPayment *gherkin.DocString) error {
	resp := c.lastResponse()
	defer resp.Body.Close()
	var gotPayment interface{}
	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&gotPayment)
	if err != nil {
		return fmt.Errorf("Invalid JSON response: %s", err)
	}

	var expectedPayment interface{}
	decoder = json.NewDecoder(bytes.NewBuffer([]byte(jsonPayment.Content)))
	err = decoder.Decode(&expectedPayment)
	if err != nil {
		return fmt.Errorf("Invalid JSON string in test specification: %s", err)
	}

	if !reflect.DeepEqual(gotPayment, expectedPayment) {
		return errors.New("Payment data in the response don't match expected payment data")
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	client := newClient()

	s.Step(`^I create a new payment described in JSON as:$`, client.iCreateANewPaymentDescribedInJSONAs)
	s.Step(`^I get a "([^"]*)" response$`, client.iGetAResponse)
	s.Step(`^the response contains a payment described in JSON as:$`, client.theResponseContainsAPaymentDescribedInJSONAs)
}
