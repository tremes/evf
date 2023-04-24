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

func TestSynopsis(t *testing.T) {
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
			name:       "errata with rhsa synopsis",
			writerJson: rhsa,
			exp:        "Moderate: OpenShift Container Platform 4.12.5 security update",
		},
		{
			name:       "errata with rhba synopsis",
			writerJson: rhba,
			exp:        "Moderate: OpenShift Container Platform 4.12.5 security update",
		},
		{
			name:       "empty errata response",
			writerJson: `{"errata": {}}`,
			exp:        "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
		})
	}
}
