package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	actions = struct{ CreateRecord, DeleteRecord, RecordList, DomainList string }{"Record.Create", "Record.Remove", "Record.List", "Domain.List"}
)

type (
	//Record info about current operation
	Record struct {
		//Domain, SubDomain, Record is provided by os.Args.
		Domain, SubDomain, Record,
		TTL string //Default 600.
		Type string //Only TXT.
	}

	//RequestBody Request body for DNSPod API
	RequestBody struct {
		//for replace/delete accurate field.
		replacements map[string]string

		//LoginToken provided by DNSPod. You can get this string via `https://console.dnspod.cn/account/token`.
		LoginToken string `json:"login_token"`
		//Format always json.
		Format string `json:"format"`
		//Default is 600. Acutally, method will use the value returned by server.
		TTL string `json:"ttl"`
		//Domain is the root of sub domain.
		Domain string `json:"domain,omitempty"`
		//SubDomain is the full domain with domain and prefix.
		SubDomain string `json:"sub_domain,omitempty"`
		//RecordType is 'A', 'CNAME', 'TXT', etc... (must be upper)
		RecordType string `json:"record_type,omitempty"`
		//Record is the value of subdomain (just for del and add).
		Record string `json:"value,omitempty"`
		//RecordID is the value of record matched by server (just for del action).
		RecordID string `json:"record_id,omitempty"`
		//RecordLineID is a value defined by DNSPod (this program will fill it by '0').
		RecordLineID string `json:"record_line_id,omitempty"`
		//Offset is page navi.
		Offset string `json:"offset,omitempty"`
		//Length will tell server how many records we'd like to get by once request (default & max are both 100).
		Length string `json:"length,omitempty"`
	}
)

//NewRequestBody Create a standard request body
func (act *Record) NewRequestBody(mustKV ...interface{}) *RequestBody {
	if len(act.TTL) == 0 {
		act.TTL = "600"
	}
	var temp = make(map[string]string, len(mustKV)/2+len(mustKV)%2)
	for index := range mustKV {
		if obj, yes := mustKV[index].(map[string]string); yes {
			for k, v := range obj {
				temp[k] = v
			}
		} else if str, yes := mustKV[index].(string); yes {
			temp[str] = ""
		}
	}
	return &RequestBody{
		replacements: temp,
		LoginToken:   string(key),
		Format:       `json`,
		TTL:          act.TTL,
		Domain:       act.Domain,
		SubDomain:    act.SubDomain,
		RecordType:   act.Type,
		Record:       act.Record,
	}
}

//GetReader Create a simple reader with current request body values
func (body *RequestBody) GetReader() io.Reader {
	body.RecordType = strings.ToUpper(body.RecordType)
	values, mappings := url.Values{}, map[string]string{}
	bin, _ := json.Marshal(body)
	json.Unmarshal(bin, &mappings)
	for k, v := range mappings {
		values.Set(k, v)
	}
	for k, v := range body.replacements {
		if len(v) == 0 {
			values.Del(k)
		} else {
			values.Set(k, v)
		}
	}
	fmt.Println(values.Encode())
	return strings.NewReader(values.Encode())
}

//CreateAction Create a standard record info with sub domain input
//
//Notice: please use `reply.Find` method to validate action, `Find` method will connect DNSPod service to check your domain is a hosted input.
//
//* If your domain is _dnsauth.go.var.ink the domain given by arguments might be 'go.var.ink'. Acutally your domain is 'var.ink' instead of 'go.var.ink', so you will hit an error about 'go.var.ink' is not your domain.
//
//* This method will find all your domains and match your sub domain to makes your domain valid for DNSPod request to avoid such problem.
func CreateAction(subdomain, record string) (act *Record) {
	return &Record{"", subdomain, record, "600", "TXT"}
}

//Delete a record which type and value matched with provider.
//
//Notice: this method will find the accurate record and send delete request.
func (act *Record) Delete() {
	for reply, index, records := struct {
		Status Status `json:"status"`
	}{}, 0, act.ListRecords(); index < len(records); index++ {
		body := act.NewRequestBody(map[string]string{"record_id": records[index]}, "ttl", "sub_domain").GetReader()
		if exitIfErr(request(http.MethodPost, actions.DeleteRecord, body, &reply)); reply.Status.Code != "1" {
			os.Stderr.WriteString(fmt.Sprintf("refused by server: %s\r\n", reply.Status.Message))
		}
	}
}

//ListRecords will fetch all records matched with Type and sub domain.
func (act *Record) ListRecords() (records []string) {
	var reply = Reply{}
	for i, body := 0, act.NewRequestBody(); len(reply.Records) == 100 || len(reply.Status.Code) == 0; i++ {
		body.Offset, body.Length, body.TTL = strconv.Itoa(i), `100`, ""
		exitIfErr(request(http.MethodPost, actions.RecordList, body.GetReader(), &reply))
		records = append(records, reply.GetRecords()[act.Record]...)
	}
	return
}

//Add a record which named after sub domain.
func (act *Record) Add() {
	body, reply := act.NewRequestBody(map[string]string{"record_line_id": "0"}).GetReader(), Reply{}
	if exitIfErr(request(http.MethodPost, actions.CreateRecord, body, &reply)); reply.Status.Code != "1" {
		exitIfErr(errors.New(reply.Status.Message), "refused by server: %s")
	}
}
