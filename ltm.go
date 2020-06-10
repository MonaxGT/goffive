package goffive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ltm struct {
	conn *http.Client
	url string
	token string
}

type Pools struct {
	Items []Pool `json:"items"`
}

type Pool struct {
	AllowNat              string `json:"allowNat"`
	AllowSnat             string `json:"allowSnat"`
	Description           string `json:"description"`
	FullPath              string `json:"fullPath"`
	Generation            int    `json:"generation"`
	IgnorePersistedWeight string `json:"ignorePersistedWeight"`
	IPTosToClient         string `json:"ipTosToClient"`
	IPTosToServer         string `json:"ipTosToServer"`
	Kind                  string `json:"kind"`
	LinkQosToClient       string `json:"linkQosToClient"`
	LinkQosToServer       string `json:"linkQosToServer"`
	LoadBalancingMode     string `json:"loadBalancingMode"`
	MembersReference      struct {
		IsSubcollection bool   `json:"isSubcollection"`
		Link            string `json:"link"`
	} `json:"membersReference"`
	MinActiveMembers       int    `json:"minActiveMembers"`
	MinUpMembers           int    `json:"minUpMembers"`
	MinUpMembersAction     string `json:"minUpMembersAction"`
	MinUpMembersChecking   string `json:"minUpMembersChecking"`
	Name                   string `json:"name"`
	Partition              string `json:"partition"`
	QueueDepthLimit        int    `json:"queueDepthLimit"`
	QueueOnConnectionLimit string `json:"queueOnConnectionLimit"`
	QueueTimeLimit         int    `json:"queueTimeLimit"`
	ReselectTries          int    `json:"reselectTries"`
	SelfLink               string `json:"selfLink"`
	SlowRampTime           int    `json:"slowRampTime"`
}

func (c *ltm) query(url string) ([]byte, error) {
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


func (c *ltm) Common() error {
	url := fmt.Sprintf("%s/mgmt/tm/ltm", c.url)
	body, err := c.query(url)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewBuffer(body))
	var t interface{}
	err = dec.Decode(&t)
	if err != nil {
		return err
	}
	fmt.Println(t)
	return nil
}

func (c *ltm) Pools() ([]Pool, error) {
	url := fmt.Sprintf("%s/mgmt/tm/ltm/pool", c.url)
	body, err := c.query(url)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewBuffer(body))

	var p Pools
	err = dec.Decode(&p)
	if err != nil {
		return nil, err
	}

	var pools []Pool
	for _, v := range p.Items {
		pools = append(pools, v)
	}
	return pools, nil
}
