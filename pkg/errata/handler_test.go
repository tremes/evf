package errata

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

func Test_Synopsis(t *testing.T) {
	rhsa := `{
	"errata": {
		"rhsa": {
			"synopsis": "Moderate: OpenShift Container Platform 4.12.5 security update"
		}
	}
}`
	rhba := `{
	"errata": {
		"rhba": {
			"synopsis": "Moderate: OpenShift Container Platform 4.12.5 security update"
		}
	}
}`
	tests := []struct {
		name       string
		writerJson string
		exp        string
	}{
		{
			name:       "rhsa",
			writerJson: rhsa,
			exp:        "Moderate: OpenShift Container Platform 4.12.5 security update",
		},
		{
			name:       "rhba",
			writerJson: rhba,
			exp:        "Moderate: OpenShift Container Platform 4.12.5 security update",
		},
	}

	for _, test := range tests {
		httpServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
			if request.URL.Path == "/rest/api/2/search/test.json" {
				_, err := writer.Write([]byte(test.writerJson))
				if err != nil {
					assert.NoError(t, err)
				}
			}
		}))
		defer httpServer.Close()

		cl := client.NewClientWithPassword("", "", "", &config.Config{})
		handler := &Handler{
			url: httpServer.URL,
			cli: spnego.NewClient(cl, nil, ""),
		}

		synopsis, err := handler.Synopsis(context.Background(), "rest/api/2/search/test")
		assert.NoError(t, err)
		assert.Equal(t, test.exp, synopsis)
	}
}
