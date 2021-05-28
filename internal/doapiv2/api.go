package doapiv2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAllDomainRecords(domainName string, token string) ([]DomainRecord, error) {
	reqString := fmt.Sprintf(domainRecords, domainName)

	return getDomainRecordsInternal(reqString, token)
}

func GetDomainRecordsByNameAndType(domainName string, token string, recordType string, hostName string) ([]DomainRecord, error) {
	reqString := fmt.Sprintf(domainRecordsFilter, domainName, recordType, hostName)

	return getDomainRecordsInternal(reqString, token)
}

func getDomainRecordsInternal(url string, token string) ([]DomainRecord, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []DomainRecord{}, fmt.Errorf("domain records request failed: %q", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return []DomainRecord{}, fmt.Errorf("domain records (GET) failed: %q", err)
	}

	if resp.StatusCode != 200 {
		return []DomainRecord{}, fmt.Errorf("non-OK status code returned [ %d ]; %q", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()

	result := new(domainRecordsResult)
	err = json.NewDecoder(resp.Body).Decode(result)

	if err != nil {
		return []DomainRecord{}, fmt.Errorf("decoding JSON response failed: %q", err)
	}

	return result.DomainRecords, nil
}

func CreateDomainRecord(domainName string, token string, record DomainRecord) (DomainRecord, error) {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(record)

	url := fmt.Sprintf(domainsCreateRecord, domainName)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return DomainRecord{}, fmt.Errorf("domain create request failed: %q", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return DomainRecord{}, fmt.Errorf("domain record (POST) failed: %q", err)
	}

	if resp.StatusCode != 200 {
		return DomainRecord{}, fmt.Errorf("non-OK status code returned [ %d ]; %q", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()

	dr := new(domainRecordResult)
	err = json.NewDecoder(resp.Body).Decode(dr)

	if err != nil {
		return DomainRecord{}, fmt.Errorf("decoding JSON response failed: %q", err)
	}

	return dr.DomainRecord, nil
}

func UpdateDomainRecord(domainName string, token string, record DomainRecord) (DomainRecord, error) {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(record)

	url := fmt.Sprintf(domainsUpdateRecord, domainName, record.Id)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return DomainRecord{}, fmt.Errorf("domain update request failed: %q", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return DomainRecord{}, fmt.Errorf("domain record (PUT) failed: %q", err)
	}

	if resp.StatusCode != 200 {
		return DomainRecord{}, fmt.Errorf("non-OK status code returned [ %d ]; %q", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()

	dr := new(domainRecordResult)
	err = json.NewDecoder(resp.Body).Decode(dr)

	if err != nil {
		return DomainRecord{}, fmt.Errorf("decoding JSON response failed: %q", err)
	}

	return dr.DomainRecord, nil
}
