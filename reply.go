package main

import "strings"

type (
	//Status defined Op Code & Msg about Code.
	Status struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	//Reply Server reply
	Reply struct {
		Status Status `json:"status"`

		//Records all records
		Records []struct {
			ID      string `json:"id"`      //id used for access
			Name    string `json:"name"`    //subdomain
			Type    string `json:"type"`    //type info
			Value   string `json:"value"`   //record value
			Enabled string `json:"enabled"` //true / false
		} `json:"records"`

		//Domains all domains
		Domains []struct {
			Status string `json:"status"` //enable / others
			Name   string `json:"name"`   //root domain
			TTL    string `json:"ttl"`    //default ttl value
		} `json:"domains"`
	}
)

//Find match domain from subdomain.
//please use this method to validate action, this method will connect DNSPod service to check your domain is a hosted input.
func (reply *Reply) Find(subdomain string) (act *Record) {
	for _, d := range reply.Domains {
		if d.Status == "enable" {
			if strings.Index(subdomain, "."+d.Name) != -1 {
				return &Record{d.Name, subdomain[:strings.LastIndex(subdomain, "."+d.Name)], "", d.TTL, "TXT"}
			}
		}
	}
	return
}

//GetRecords get all records (filtered by Type) and trans to mappings.
func (reply *Reply) GetRecords() map[string][]string {
	var records = make(map[string][]string, 0)
	for _, record := range reply.Records {
		records[record.Value] = append(records[record.Value], record.ID)
	}
	return records
}
