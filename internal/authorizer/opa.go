package authorizer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type OPAAuthorizer struct {
	host              string
	rule              string
	partialEvaluation bool
}

type opaResponse struct {
	Result bool `json:"result"`
}

type opaPayloadInput struct {
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Headers http.Header `json:"headers"`
}

type opaPayload struct {
	Input *opaPayloadInput `json:"input"`
}

func (opaAuthorizer OPAAuthorizer) IsAuthorized(request *Request) bool {
	authPayloadAsJSON, errm := json.Marshal(toOPAPayload(request))
	if errm != nil {
		panic(errm)
	}

	authResponse, errp := http.Post(buildURL(opaAuthorizer), "application/json", bytes.NewBuffer(authPayloadAsJSON))
	if errp != nil {
		panic(errp)
	}

	body, err := ioutil.ReadAll(authResponse.Body)
	if err != nil {
		panic(err.Error())
	}

	var response opaResponse
	if errj := json.Unmarshal(body, &response); errj != nil {
		panic(errj)
	}

	return response.Result
}

func toOPAPayload(request *Request) *opaPayload {
	return &opaPayload{
		Input: &opaPayloadInput{
			Method:  request.Method,
			Path:    request.Path,
			Headers: request.Headers,
		},
	}
}

func NewOPAAuthorizer(opaHost string, rule string, partialEvaluation bool) *OPAAuthorizer {
	return &OPAAuthorizer{host: opaHost, rule: rule, partialEvaluation: partialEvaluation}
}

func buildURL(opaAuthorizer OPAAuthorizer) string {
	return "http://" + opaAuthorizer.host + "/v1/data/" + opaAuthorizer.rule + "?partial=" + strconv.FormatBool(opaAuthorizer.partialEvaluation)
}
