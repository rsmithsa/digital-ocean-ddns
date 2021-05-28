package doapiv2

const apiRoot string = "https://api.digitalocean.com/v2"

const domainRecords string = apiRoot + "/domains/%s/records"
const domainRecordsFilter string = apiRoot + "/domains/%s/records?type=%s&name=%s"

const domainsUpdateRecord string = apiRoot + "/domains/%s/records/%d"
const domainsCreateRecord string = apiRoot + "/domains/%s/records"
