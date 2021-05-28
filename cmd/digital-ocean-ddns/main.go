package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rsmithsa/digital-ocean-ddns/internal/doapiv2"

	"golang.org/x/net/publicsuffix"
)

const apiToken = ""

type server struct{}

func hasIPChanged(ipAddr string, token string, hostName string) (doapiv2.DomainRecord, bool, error) {
	domainName, err := publicsuffix.EffectiveTLDPlusOne(hostName)
	if err != nil {
		return doapiv2.DomainRecord{}, false, fmt.Errorf("unable to parse hostname: %q", err)
	}

	subDomainName := strings.TrimSuffix(hostName, "."+domainName)

	result, err := doapiv2.GetDomainRecordsByNameAndType(domainName, token, "A", hostName)
	if err != nil {
		return doapiv2.DomainRecord{}, false, err
	}

	if len(result) == 0 {
		log.Printf("Record not found: [ %q ]", hostName)
		// New record required
		return doapiv2.DomainRecord{Type: "A", Name: subDomainName, Data: ipAddr, TTL: 60}, true, nil
	}

	if result[0].Data == ipAddr {
		return result[0], false, nil
	} else {
		return result[0], true, nil
	}
}

func updateRecord(record doapiv2.DomainRecord, token string, hostName string, ipAddr string) (doapiv2.DomainRecord, error) {
	log.Printf("Updating record: [ %q ]", hostName)
	domainName, err := publicsuffix.EffectiveTLDPlusOne(hostName)
	if err != nil {
		return doapiv2.DomainRecord{}, fmt.Errorf("unable to parse hostname: %q", err)
	}

	record.Data = ipAddr
	return doapiv2.UpdateDomainRecord(domainName, token, record)
}

func createRecord(record doapiv2.DomainRecord, token string, hostName string) (doapiv2.DomainRecord, error) {
	log.Printf("Creating record: [ %q ]", hostName)
	domainName, err := publicsuffix.EffectiveTLDPlusOne(hostName)
	if err != nil {
		return doapiv2.DomainRecord{}, fmt.Errorf("unable to parse hostname: %q", err)
	}

	return doapiv2.CreateDomainRecord(domainName, token, record)
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
