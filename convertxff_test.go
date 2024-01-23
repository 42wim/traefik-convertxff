package traefik_convertxff_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	convertxff "github.com/42wim/traefik-convertxff"
)

func TestConvertXFF_ServeHTTP(t *testing.T) {
	testCases := []struct {
		desc          string
		header        string
		expectedValue string
		expectedCode  int
	}{
		{
			desc:          "ipv4 and bracketed ipv6",
			header:        "1.2.3.4, [2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
			expectedValue: "1.2.3.4,2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expectedCode:  http.StatusOK,
		},
		{
			desc:          "ipv4 only",
			header:        "1.2.3.4",
			expectedValue: "1.2.3.4",
			expectedCode:  http.StatusOK,
		},
		{
			desc:          "bracketed ipv6 only",
			header:        "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]",
			expectedValue: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expectedCode:  http.StatusOK,
		},
		{
			desc:          "empty header value",
			header:        "",
			expectedValue: "",
			expectedCode:  http.StatusOK,
		},
		{
			desc:          "bracketed ipv6 and ipv4",
			header:        "[2001:0db8:85a3:0000:0000:8a2e:0370:7334],1.2.3.4",
			expectedValue: "2001:0db8:85a3:0000:0000:8a2e:0370:7334,1.2.3.4",
			expectedCode:  http.StatusOK,
		},
		{
			desc:          "bracketed, normal ipv4-mapped ipv6 and ipv4",
			header:        "[::ffff:2a5b:3cde], ::ffff:2a6b:3cde, 1.2.3.4",
			expectedValue: "42.91.60.222,42.107.60.222,1.2.3.4",
			expectedCode:  http.StatusOK,
		},
	}
	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			fat, _ := convertxff.New(context.Background(), next, &convertxff.Config{}, "TestConvertXFF")
			r := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			r.Header.Set(convertxff.XFF, test.header)
			w := httptest.NewRecorder()
			fat.ServeHTTP(w, r)
			if test.expectedCode != w.Code {
				t.Errorf("Expexted code: %d; Received: %d", test.expectedCode, w.Code)
			}
			value := r.Header.Get(convertxff.XFF)
			if test.expectedValue != value {
				t.Errorf("Expexted new header value: %s; Received value: %s", test.expectedValue, r.Header.Get(convertxff.XFF))
			}
		})
	}
}
