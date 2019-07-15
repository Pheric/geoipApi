package geoipApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GeoipApi string
type QueryResponse struct {
	Message string
	Ip string
	Hostname string
	ContinentCode string
	ContinentName string
	CountryCode2 string
	CountryCode3 string
	CountryName string
	CountryCapital string
	StateProv string
	District string
	City string
	Zipcode string
	Latitude string
	Longitude string
	IsEu bool
	CallingCode string
	CountryTld string
	Languages string
	CountryFlag string
	Isp string
	ConnectionType string
	Organization string
	GeonameId int
	Currency struct {
		Code string
		Name string
		Symbol string
	}
	TimeZone struct {
		Name string
		Offset string
		CurrentTime string
		CurrentTimeUnix string
		IsDst string
		DstSavings string
	}
}
type Response struct {
	Error error
	QueryResponse *QueryResponse
}

func (g GeoipApi) FetchInfo(rchan chan Response, ipAddr string) {
	var response Response
	defer func(r *Response) {
		rchan <- *r
	}(&response)
	qr := new(QueryResponse)

	req, err := http.NewRequest("GET", "https://api.ipgeolocation.io/ipgeo", nil)
	if err != nil {
		response.Error = fmt.Errorf("error creating request for ip \"%s\": %v", ipAddr, err)
		return
	}

	q := req.URL.Query()
	q.Add("apiKey", string(g))
	q.Add("ip", ipAddr)
	req.URL.RawQuery = q.Encode()

	client := http.Client{}
	r, err := client.Do(req)
	if err != nil {
		response.Error = fmt.Errorf("error fetching data for ip \"%s\": %v", ipAddr, err)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error = fmt.Errorf("error reading response for ip \"%s\": %v", ipAddr, err)
		return
	}
	_ = r.Body.Close()

	if err := json.Unmarshal(data, qr); err != nil {
		response.Error = fmt.Errorf("error unmarshalling response for ip \"%s\": %v", ipAddr, err)
		return
	}
	response.QueryResponse = qr
}