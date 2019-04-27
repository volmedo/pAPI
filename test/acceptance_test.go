package test

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

func iCreateANewPaymentDescribedInJSONAs(jsonPayment *gherkin.DocString) error {
	return godog.ErrPending
}

func iGetAResponse(respStatus string) error {
	return godog.ErrPending
}

func theResponseContainsAPaymentDescribedInJSONAs(jsonPayment *gherkin.DocString) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I create a new payment described in JSON as:$`, iCreateANewPaymentDescribedInJSONAs)
	s.Step(`^I get a "([^"]*)" response$`, iGetAResponse)
	s.Step(`^the response contains a payment described in JSON as:$`, theResponseContainsAPaymentDescribedInJSONAs)
}
