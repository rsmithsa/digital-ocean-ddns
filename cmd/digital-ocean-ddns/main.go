package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/publicsuffix"
)

const apiToken = ""

type server struct{}

type domainRecord struct {
	Id       int    `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Data     string `json:"data"`
	Priority *int   `json:"priority"`
	Port     *int   `json:"port"`
	TTL      int    `json:"ttl"`
	Weight   *int   `json:"weight"`
	Flags    *int   `json:"flags"`
	Tag      string `json:"tag"`
}

type singleDomainRecord struct {
	DomainRecord domainRecord `json:"domain_record"`
}

type domainRecords struct {
	DomainRecords []domainRecord `json:"domain_records"`
}

func hasIPChanged(ipAddr string, tokenString string, hostName string) (domainRecord, bool, error) {
	drs := new(domainRecords)

	domainName, err := publicsuffix.EffectiveTLDPlusOne(hostName)
	if err != nil {
		return domainRecord{}, false, fmt.Errorf("unable to parse hostname: %q", err)
	}

	subDomainName := strings.TrimSuffix(hostName, "."+domainName)

	reqString := fmt.Sprintf(doapiv2.doAPIRecords, domainName, hostName)
	req, err := http.NewRequest("GET", reqString, nil)
	if err != nil {
		return domainRecord{}, false, fmt.Errorf("domain records request failed: %q", err)
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return domainRecord{}, false, fmt.Errorf("domain records (GET) failed: %q", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return domainRecord{}, false, fmt.Errorf("non-OK status code returned [ %d ]; %q", resp.StatusCode, resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(drs)

	if err != nil {
		return domainRecord{}, false, fmt.Errorf("decoding JSON response failed: %q", err)
	}

	if len(drs.DomainRecords) == 0 {
		log.Printf("Record not found: [ %q ]", hostName)
		// New record required
		return domainRecord{Type: "A", Name: subDomainName, Data: ipAddr, TTL: 60}, true, nil
	}

	if drs.DomainRecords[0].Data == ipAddr {
		return drs.DomainRecords[0], false, nil
	} else {
		return drs.DomainRecords[0], true, nil
	}
}

func updateRecord(record domainRecord, tokenString string, hostName string, ipAddr string) (domainRecord, error) {
	log.Print("Update")
	domainName, err := publicsuffix.EffectiveTLDPlusOne(hostName)
	if err != nil {
		return domainRecord{}, fmt.Errorf("unable to parse hostname: %q", err)
	}

	record.Data = ipAddr
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(record)

	reqString := fmt.Sprintf(doapiv2.doUpdateRecord, domainName, record.Id)
	req, err := http.NewRequest("PUT", reqString, body)
	if err != nil {
		return domainRecord{}, fmt.Errorf("domain update request failed: %q", err)
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return domainRecord{}, fmt.Errorf("domain record (PUT) failed: %q", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return domainRecord{}, fmt.Errorf("non-OK status code returned [ %d ]; %q", resp.StatusCode, resp.Status)
	}

	dr := new(singleDomainRecord)
	err = json.NewDecoder(resp.Body).Decode(dr)

	if err != nil {
		return domainRecord{}, fmt.Errorf("decoding JSON response failed: %q", err)
	}

	return dr.DomainRecord, nil
}

func createRecord(record domainRecord, tokenString string, hostName string) (domainRecord, error) {
	log.Print("Update")
	domainName, err := publicsuffix.EffectiveTLDPlusOne(hostName)
	if err != nil {
		return domainRecord{}, fmt.Errorf("unable to parse hostname: %q", err)
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(record)

	reqString := fmt.Sprintf(doapiv2.doCreateRecord, domainName)
	req, err := http.NewRequest("POST", reqString, body)
	if err != nil {
		return domainRecord{}, fmt.Errorf("domain create request failed: %q", err)
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return domainRecord{}, fmt.Errorf("domain record (POST) failed: %q", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return domainRecord{}, fmt.Errorf("non-OK status code returned [ %d ]; %q", resp.StatusCode, resp.Status)
	}

	dr := new(singleDomainRecord)
	err = json.NewDecoder(resp.Body).Decode(dr)

	if err != nil {
		return domainRecord{}, fmt.Errorf("decoding JSON response failed: %q", err)
	}

	return dr.DomainRecord, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		q := r.URL.Query()
		hostname := q.Get("hostname")
		myip := q.Get("myip")

		authorization := r.Header.Get("Authorization")

		if hostname == "" || myip == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if authorization == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		hostname = strings.TrimSpace(hostname)

		rec, changed, err := hasIPChanged(myip, apiToken, hostname)

		if err != nil {
			w.Header().Set("Content-Type", "application/txt")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var responseText string
		if changed {
			if rec.Id == 0 {
				createRecord(rec, apiToken, hostname)
			} else {
				updateRecord(rec, apiToken, hostname, myip)
			}
			responseText = fmt.Sprintf("good %s", myip)
		} else {
			responseText = fmt.Sprintf("nochg %s", myip)
		}

		w.Header().Set("Content-Type", "application/txt")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseText))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	s := &server{}
	http.Handle("/nic/update", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
