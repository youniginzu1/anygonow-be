package api

import (
	"net/http"
	"testing"

	"github.com/aqaurius6666/urlbuilder"
	"github.com/stretchr/testify/assert"
)

var (
	USER_ID     = "72f048be-f355-4521-b89f-958c997b4d48"
	HANDYMAN_ID = "74dab9ef-7f39-417a-804e-37532e58886c"
	CATEGORY_ID = "21fe7ae9-ea14-42ae-bc7c-14592071c77a"
	ZIPCODE     = "100"
	BASE_URL    = "http://localhost:50001/api"
	BASE_PORT   = 50001
	BASE_HOST   = "localhost"
)
var (
	BASE = urlbuilder.From(urlbuilder.UrlBuilder{
		Port:     BASE_PORT,
		Host:     BASE_HOST,
		BasePath: "/api",
		Scheme:   "http",
	})
	BASE_GUEST    = urlbuilder.From(*BASE)
	BASE_HANDYMAN = urlbuilder.From(*BASE).WithQueryParam("userId", HANDYMAN_ID).WithQueryParam("role", "HANDYMAN")
	BASE_CUSTOMER = urlbuilder.From(*BASE).WithQueryParam("userId", USER_ID).WithQueryParam("role", "CUSTOMER")
	BASE_ADMIN    = urlbuilder.From(*BASE).WithQueryParam("userId", USER_ID).WithQueryParam("role", "ADMIN")
)

type TestCase struct {
	Name               string
	Url                string
	Body               []byte
	ExpectedError      error
	ExpectedStatusCode int
	Method             string
}

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		http.Get("http://localhost:9091/api/businesses?categoryId=29c6a5eb-eb52-42cf-90aa-d410045903d6")
	}
}

func GetTestCase_OrderApi() []TestCase {
	return []TestCase{
		{
			Name: "Test get /orders ok",
			Url: BASE_HANDYMAN.
				WithPath("/orders").Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /orders/projects ok",
			Url: BASE_CUSTOMER.
				WithPath("/orders/projects").Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
	}
}

func GetTestCase_ContactApi() []TestCase {
	return []TestCase{
		{
			Name: "Test get /contacts/states ok",
			Url: BASE_GUEST.
				WithPath("/contacts/states").Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /contacts/:id ok",
			Url: BASE_HANDYMAN.
				WithPath("/contacts/:id").WithPathParam("id", HANDYMAN_ID).Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /contacts/:id ok 2",
			Url: BASE_CUSTOMER.
				WithPath("/contacts/:id").WithPathParam("id", USER_ID).Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /contacts/:id ok 2",
			Url: BASE_HANDYMAN.
				WithPath("/contacts/:id").WithPathParam("id", USER_ID).Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
	}
}

func TestApi(t *testing.T) {
	testCase := append(append(append(append(GetTestCase_OrderApi(),
		GetTestCase_BusinessApi()...),
		GetTestCase_AdminApi()...),
		GetTestCase_UserApi()...),
		GetTestCase_ContactApi()...)
	for _, c := range testCase {
		t.Run(c.Name, func(t *testing.T) {
			if c.Method == http.MethodGet {
				resp, err := http.Get(c.Url)
				if c.ExpectedError != nil {
					assert.EqualError(t, err, c.ExpectedError.Error())
				} else {
					assert.Nil(t, err, err)
				}
				assert.Equal(t, c.ExpectedStatusCode, resp.StatusCode)
			}
		})
	}
}

func GetTestCase_AdminApi() []TestCase {
	return []TestCase{
		{
			Name: "Test get /admin/categories ok",
			Url: BASE_ADMIN.
				WithPath("/admin/categories").Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /admin/groups ok",
			Url: BASE_ADMIN.
				WithPath("/admin/groups").Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /admin/categories unauthorized",
			Url: BASE_CUSTOMER.
				WithPath("/admin/categories").Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusUnauthorized,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /admin/users unauthorized",
			Url: BASE_CUSTOMER.
				WithPath("/admin/users").Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusUnauthorized,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /admin/users unauthorized",
			Url: BASE_HANDYMAN.
				WithPath("/admin/users").Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusUnauthorized,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /admin/users",
			Url: BASE_ADMIN.
				WithPath("/admin/users").Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /admin/businesses",
			Url: BASE_ADMIN.
				WithPath("/admin/businesses").Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
	}
}

func GetTestCase_BusinessApi() []TestCase {
	return []TestCase{
		{
			Name:               "Test get /businesses",
			Url:                BASE_GUEST.WithPath("/businesses").Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /businesses with category",
			Url: BASE_GUEST.WithPath("/businesses").
				WithQueryParam("categoryId", CATEGORY_ID).
				Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /businesses with category, zipcode",
			Url: BASE_GUEST.WithPath("/businesses").
				WithQueryParam("categoryId", CATEGORY_ID).
				WithQueryParam("zipcode", ZIPCODE).
				Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /businesses with zipcode",
			Url: BASE_GUEST.WithPath("/businesses").
				WithQueryParam("zipcode", ZIPCODE).
				Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /businesses/:id",
			Url: BASE_GUEST.WithPath("/businesses/:id").
				WithPathParam("id", HANDYMAN_ID).Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /businesses/:id/services",
			Url: BASE_GUEST.WithPath("/businesses/:id/services").
				WithPathParam("id", HANDYMAN_ID).Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /businesses/:id/rating",
			Url: BASE_GUEST.WithPath("/businesses/:id/rating").
				WithPathParam("id", HANDYMAN_ID).Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /businesses/:id/feedbacks",
			Url: BASE_GUEST.WithPath("/businesses/:id/feedbacks").
				WithPathParam("id", HANDYMAN_ID).Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /businesses/interest",
			Url: BASE_GUEST.WithPath("/businesses/interest").
				Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /businesses with zipcode",
			Url: BASE_GUEST.WithPath("/businesses").WithPath("/businesses").
				WithQueryParam("zipcode", ZIPCODE).
				Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
		{
			Name: "Test get /businesses with zipcode",
			Url: BASE_GUEST.WithPath("/businesses").WithPath("/businesses").
				WithQueryParam("zipcode", ZIPCODE).
				Build(),
			ExpectedError:      nil,
			Body:               nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
	}
}

func GetTestCase_UserApi() []TestCase {
	return []TestCase{
		{
			Name:               "Test get /users/:id",
			Url:                BASE_CUSTOMER.WithPath("/users/:id").WithPathParam("id", USER_ID).Build(),
			Body:               nil,
			ExpectedError:      nil,
			ExpectedStatusCode: http.StatusOK,
			Method:             http.MethodGet,
		},
	}
}
