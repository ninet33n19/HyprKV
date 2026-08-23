package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ninet33n19/HyprKV/internal/store"
)

type TestGetRequest struct {
	name           string
	key            string
	expectedStatus int
	expectedBody   string
}

func TestServerGET(t *testing.T) {
	tests := []TestGetRequest{
		{
			name:           "Key Found",
			key:            "foo",
			expectedStatus: http.StatusOK,
			expectedBody:   "bar",
		},
		{
			name:           "Key Not Found",
			key:            "xyz",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "key not found\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := store.NewStore()
			if err := store.Set("foo", "bar"); err != nil {
				t.Fatalf("set up store: %v", err)
			}

			server := NewServer(store)

			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.SetPathValue("key", test.key)

			response := httptest.NewRecorder()

			server.Get(response, request)

			if response.Code != test.expectedStatus {
				t.Errorf(
					"expected status %d, got %d",
					test.expectedStatus,
					response.Code,
				)
			}

			if response.Body.String() != test.expectedBody {
				t.Errorf(
					"expected body %q, got %q",
					test.expectedBody,
					response.Body.String(),
				)
			}
		})
	}
}
