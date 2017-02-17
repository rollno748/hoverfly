package matching

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func TestEmptyTemplateShouldMatchOnAnyRequest(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	store := RequestTemplateStore{}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"sdv": []string{"ascd"},
		},
	}
	result, _ := store.GetResponse(r, false, simulation)

	Expect(result.Body).To(Equal("test-body"))
}

func TestTemplateShouldMatchOnBody(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Body: StringToPointer("body"),
		},
		Response: models.ResponseDetails{
			Body: "body",
		},
	})

	store := RequestTemplateStore{}

	r := models.RequestDetails{
		Body: "body",
	}
	result, err := store.GetResponse(r, false, simulation)
	Expect(err).To(BeNil())

	Expect(result.Body).To(Equal("body"))
}

func TestReturnResponseWhenAllHeadersMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: headers,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	store := RequestTemplateStore{}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
			"header2": []string{"val2"},
		},
	}

	result, _ := store.GetResponse(r, false, simulation)

	Expect(result.Body).To(Equal("test-body"))
}

func TestReturnNilWhenOneHeaderNotPresentInRequest(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: headers,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	store := RequestTemplateStore{}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
		},
	}

	result, _ := store.GetResponse(r, false, simulation)

	Expect(result).To(BeNil())
}

func TestReturnNilWhenOneHeaderValueDifferent(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: headers,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	store := RequestTemplateStore{}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1"},
			"header2": []string{"different"},
		},
	}
	result, _ := store.GetResponse(r, false, simulation)

	Expect(result).To(BeNil())
}

func TestReturnResponseWithMultiValuedHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: headers,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	store := RequestTemplateStore{}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Body:        "test-body",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result, _ := store.GetResponse(r, false, simulation)

	Expect(result.Body).To(Equal("test-body"))
}

func TestReturnNilWithDifferentMultiValuedHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: headers,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	store := RequestTemplateStore{}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-differnet"},
			"header2": []string{"val2"},
		},
	}

	result, _ := store.GetResponse(r, false, simulation)

	Expect(result).To(BeNil())
}

func TestEndpointMatchWithHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers:     headers,
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	store := RequestTemplateStore{}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query:       "q=test",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}
	result, _ := store.GetResponse(r, false, simulation)

	Expect(result.Body).To(Equal("test-body"))
}

func TestEndpointMismatchWithHeadersReturnsNil(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": []string{"val1-a", "val1-b"},
		"header2": []string{"val2"},
	}

	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers:     headers,
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: models.ResponseDetails{
			Body: "test-body",
		},
	})

	store := RequestTemplateStore{}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://testhost.com",
		Path:        "/a/1",
		Query:       "q=different",
		Headers: map[string][]string{
			"header1": []string{"val1-a", "val1-b"},
			"header2": []string{"val2"},
		},
	}

	result, _ := store.GetResponse(r, false, simulation)

	Expect(result).To(BeNil())
}

func TestAbleToMatchAnEmptyPathInAReasonableWay(t *testing.T) {
	RegisterTestingT(t)

	response := models.ResponseDetails{
		Body: "test-body",
	}
	destination := "testhost.com"
	method := "GET"
	path := ""
	query := "q=test"
	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Destination: &destination,
			Path:        &path,
			Method:      &method,
			Query:       &query,
		},
		Response: response,
	})

	store := RequestTemplateStore{}

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Query:       "q=test",
	}
	result, _ := store.GetResponse(r, false, simulation)

	Expect(result.Body).To(Equal("test-body"))

	r = models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query:       "q=test",
	}

	result, _ = store.GetResponse(r, false, simulation)

	Expect(result).To(BeNil())
}

func TestRequestTemplateResponsePairCanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	method := "POST"

	requestTemplateResponsePair := models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Method: &method,
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	pairView := requestTemplateResponsePair.ConvertToV1RequestResponsePairView()

	Expect(pairView.Request.RequestType).To(Equal(StringToPointer("template")))
	Expect(pairView.Request.Method).To(Equal(StringToPointer("POST")))
	Expect(pairView.Request.Destination).To(BeNil())
	Expect(pairView.Request.Path).To(BeNil())
	Expect(pairView.Request.Scheme).To(BeNil())
	Expect(pairView.Request.Query).To(BeNil())

	Expect(pairView.Response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnDestinationAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Destination: StringToPointer("*.com"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	store := RequestTemplateStore{}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	response, err := store.GetResponse(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnPathAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Path: StringToPointer("/api/*"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	store := RequestTemplateStore{}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	response, err := store.GetResponse(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnMethodAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Method: StringToPointer("*T"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	store := RequestTemplateStore{}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	response, err := store.GetResponse(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnSchemeAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Scheme: StringToPointer("H*"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	store := RequestTemplateStore{}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Scheme:      "http",
		Path:        "/api/1",
	}

	response, err := store.GetResponse(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnQueryAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Query: StringToPointer("q=*"),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	store := RequestTemplateStore{}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Query:       "q=anything-i-want",
	}

	response, err := store.GetResponse(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnBodyndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Body: StringToPointer(`{"json": "object", "key": *}`),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	store := RequestTemplateStore{}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Body:        `{"json": "object", "key": "value"}`,
	}

	response, err := store.GetResponse(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestTemplatesCanUseGlobsOnBodyAndNotMatchWhenTheBodyIsWrong(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Body: StringToPointer(`{"json": "object", "key": *}`),
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	store := RequestTemplateStore{}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Body:        `[{"json": "objects", "key": "value"}]`,
	}

	_, err := store.GetResponse(request, false, simulation)
	Expect(err).ToNot(BeNil())
}

func TestTemplatesCanUseGlobsOnHeadersAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.Templates = append(simulation.Templates, models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: map[string][]string{
				"unique-header": []string{"*"},
			},
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	})

	store := RequestTemplateStore{}

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Headers: map[string][]string{
			"unique-header": []string{"totally-unique"},
		},
	}

	response, err := store.GetResponse(request, false, simulation)
	Expect(err).To(BeNil())

	Expect(response.Body).To(Equal("template matched"))
}

func TestRequestTemplateResponsePair_ConvertToRequestResponsePairView_CanBeConvertedToARequestResponsePairView_WhileIncomplete(t *testing.T) {
	RegisterTestingT(t)

	method := "POST"

	requestTemplateResponsePair := models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Method: &method,
		},
		Response: models.ResponseDetails{
			Body: "template matched",
		},
	}

	pairView := requestTemplateResponsePair.ConvertToRequestResponsePairView()

	Expect(pairView.Request.RequestType).To(Equal(StringToPointer("template")))
	Expect(pairView.Request.Method).To(Equal(StringToPointer("POST")))
	Expect(pairView.Request.Destination).To(BeNil())
	Expect(pairView.Request.Path).To(BeNil())
	Expect(pairView.Request.Scheme).To(BeNil())
	Expect(pairView.Request.Query).To(BeNil())

	Expect(pairView.Response.Body).To(Equal("template matched"))
}