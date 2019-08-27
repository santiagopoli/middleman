package authorizer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type OPAAuthorizer struct {
	opaURL string
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

func (opa OPAAuthorizer) IsAuthorized(request *Request) bool {
	authPayloadAsJSON, errm := json.Marshal(toOPAPayload(request))
	if errm != nil {
		panic(errm)
	}

	authResponse, errp := http.Post(opa.opaURL, "application/json", bytes.NewBuffer(authPayloadAsJSON))
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

func NewOPAAuthorizer(opaURL string) *OPAAuthorizer {
	return &OPAAuthorizer{opaURL: opaURL}
}
