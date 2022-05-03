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
	spnego.NewClient(cl, nil, "")
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

func (h *Handler) GetSynopsis(id string) (string, error) {
	data, err := h.getErrata(id)
	if err != nil {
		return "", err
	}
	/* 	type Rhba struct {
	   		Synopsis string `json:"synopsis"`
	   	}
	   	type Errata struct {
	   		Rhba Rhba `json:"rhba"`
	   	}
	   	type Res struct {
	   		Errata Errata `json:"errata"`
	   	} */

	type res struct {
		Errata struct {
			Rhba struct {
				Synopsis string
			}
		}
	}

	type res2 struct {
		Errata struct {
			Rhsa struct {
				Synopsis string
			}
		}
	}

	var r res
	err = json.Unmarshal(data, &r)
	if err != nil {
		return "", err
	}
	if r == (res{}) {
		var r2 res2
		err = json.Unmarshal(data, &r2)
		if err != nil {
			return "", err
		}
		return r2.Errata.Rhsa.Synopsis, nil
	}
	return r.Errata.Rhba.Synopsis, nil
}
