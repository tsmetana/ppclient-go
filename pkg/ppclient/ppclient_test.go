package ppclient

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var fakeResultBody = []ppRelease{
	{
		Shortname: "openshift-1.0",
		Phase:     "Unsupported",
	},
	{
		Shortname: "openshift-4.9",
		Phase:     "Maintenance",
	},
	{
		Shortname: "openshift-4.11",
		Phase:     "Planning / Development / Testing",
	},
}

func newPpClientWithEndpoint(endpoint string) PpClient {
	return &client{
		client: &http.Client{
			Timeout: time.Duration(10 * time.Second),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		endpoint: endpoint,
	}
}

func TestGetReleases(t *testing.T) {
	var fakeResult PpReleaseList
	for _, r := range fakeResultBody {
		fakeResult = append(fakeResult, NewPpRelease(r.Shortname, r.Phase))
	}
	var testCases = []struct {
		desc      string
		expResult []PpRelease
		resp      string
		expError  bool
		shortname string
	}{
		{
			desc:      "sucessful request / response",
			expResult: fakeResult,
			shortname: "openshift",
			expError:  false,
		},
		{
			desc:      "empty response, JSON decoding fails",
			expResult: nil,
			shortname: "",
			expError:  true,
		},
		{
			desc:      "error status returned, should fail",
			expResult: nil,
			shortname: "badrequest",
			expError:  true,
		},
	}
	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		shortnameParam := r.URL.Query().Get("shortname")
		if shortnameParam == "openshift" {
			fakeResponseBody, err := json.Marshal(fakeResultBody)
			if err != nil {
				t.Fatalf("fatal error creating fake response: %v", err)
			}
			w.Write(fakeResponseBody)
		} else if shortnameParam == "badrequest" {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
		}
		// otherwise: return OK, but with empty body
	}))
	defer testServer.Close()

	client := newPpClientWithEndpoint(testServer.URL)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			releases, err := client.GetReleases(tc.shortname)
			if err != nil && !tc.expError {
				t.Errorf("expected no error but got: %v", err)
			}
			if len(releases) != len(tc.expResult) {
				t.Errorf("expected %d releases, got %d", len(releases), len(tc.expResult))
			}
		})
	}
}
