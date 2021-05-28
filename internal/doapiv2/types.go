package doapiv2

type DomainRecord struct {
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

type domainRecordResult struct {
	DomainRecord DomainRecord `json:"domain_record"`
}

type domainRecordsResult struct {
	DomainRecords []DomainRecord `json:"domain_records"`
}
