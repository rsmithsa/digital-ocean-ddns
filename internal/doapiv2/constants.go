package doapiv2

const doAPIRoot string = "https://api.digitalocean.com/v2"

const doAPIRecords string = doAPIRoot + "/domains/%s/records?type=A&name=%s"
const doUpdateRecord string = doAPIRoot + "/domains/%s/records/%d"
const doCreateRecord string = doAPIRoot + "/domains/%s/records"
