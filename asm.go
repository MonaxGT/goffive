package goffive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type asm struct {
	conn *http.Client
	url string
	token string
}

type Policies struct {
	Totalitems int      `json:"totalItems"`
	Items      []Policy `json:"items"`
}

type Policy struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	VersionLastChange string   `json:"versionLastChange"`
	Description       string   `json:"description"`
	VirtualServers    []string `json:"virtualServers"`
}

type Signatories struct {
	Items []Signature `json:"items"`
}

type Signature struct {
	SignatureReference struct {
		Name        string `json:"name"`
		SignatureId uint64 `json:"signatureId"`
	} `json:"signatureReference"`
	ID             string `json:"id"`
	Block          bool   `json:"block"`
	Learn          bool   `json:"learn"`
	Enabled        bool   `json:"enabled"`
	Alarm          bool   `json:"alarm"`
	PerformStaging bool   `json:"performStaging"`
}

func (c *asm) query(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-F5-Auth-Token", c.token)
	resp, err := c.conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *asm) Policies() ([]Policy, error) {
	url := fmt.Sprintf("%s/mgmt/tm/asm/policies", c.url)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewBuffer(body))

	var p Policies
	err = dec.Decode(&p)
	if err != nil {
		return nil, err
	}

	var policies []Policy
	for _, v := range p.Items {
		policies = append(policies, v)
	}
	return policies, nil
}

func (c *asm) Signatories(policy string) ([]Signature, error) {
	url := fmt.Sprintf("%s/mgmt/tm/asm/policies/%s/signatures", c.url, policy)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewBuffer(body))

	var s Signatories
	err = dec.Decode(&s)
	if err != nil {
		return nil, err
	}

	var signatories []Signature
	for _, v := range s.Items {
		signatories = append(signatories, v)
	}
	return signatories, nil
}
