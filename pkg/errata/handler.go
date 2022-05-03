package errata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

type errataResponse struct {
	Errata errata `json:"errata"`
}
type errata struct {
	// these are two types (in fact the same) with different naming only
	Rhba rhba `json:"rhba,omitempty"`
	Rhsa rhba `json:"rhsa,omitempty"`
}

type rhba struct {
	Synopsis string `json:"synopsis"`
}

type Handler struct {
	url string
	cli *spnego.Client
}

func New(url, krb5ConfFile, username, realm, password string) (*Handler, error) {
	cfg, err := config.Load(krb5ConfFile)
	if err != nil {
		return nil, err
	}

	cl := client.NewClientWithPassword(username, realm, password, cfg)
	if err != nil {
		return nil, err
	}
	return &Handler{
		url: url,
		cli: spnego.NewClient(cl, nil, ""),
	}, nil
}

func (h *Handler) getErrata(id string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s.json", h.url, id)
	//TODO use request with context
	r, _ := http.NewRequest("GET", url, nil)
	res, err := h.cli.Do(r)
	if err != nil {
		return nil, err
	}
	//TODO use limit reader
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return data, nil
}

func (h *Handler) Synopsis(id string) (string, error) {
	data, err := h.getErrata(id)
	if err != nil {
		return "", err
	}
	var er errataResponse
	err = json.Unmarshal(data, &er)
	if err != nil {
		return "", err
	}
	if er.Errata.Rhba.Synopsis != "" {
		return er.Errata.Rhba.Synopsis, nil
	}
	return er.Errata.Rhsa.Synopsis, nil
}
