package authorizer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
	Host    string              `json:"host"`
	Method  string              `json:"method"`
	Path    []string            `json:"path"`
	Headers map[string][]string `json:"headers"`
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
			Host:    request.Host,
			Method:  request.Method,
			Path:    splitPath(request.Path),
			Headers: request.Headers,
		},
	}
}

func buildURL(opaAuthorizer OPAAuthorizer) string {
	return "http://" + opaAuthorizer.host + "/v1/data/" + opaAuthorizer.rule + "?partial=" + strconv.FormatBool(opaAuthorizer.partialEvaluation)
}

func splitPath(path string) []string {
	return strings.Split(path, "/")[1:]
}

func NewOPAAuthorizer(opaHost string, rule string, partialEvaluation bool) *OPAAuthorizer {
	return &OPAAuthorizer{host: opaHost, rule: rule, partialEvaluation: partialEvaluation}
}
