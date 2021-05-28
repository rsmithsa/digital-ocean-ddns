package doapiv2

const DOAPIRoot string = "https://api.digitalocean.com/v2"

const DODomainRecordsFilter string = DOAPIRoot + "/domains/%s/records?type=%s&name=%s"
const DODomainsUpdateRecord string = DOAPIRoot + "/domains/%s/records/%d"
const DODomainsCreateRecord string = DOAPIRoot + "/domains/%s/records"
