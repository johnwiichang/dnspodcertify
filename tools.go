package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

//request send a request to DNSPod API and check whether error occurred or not.
func request(method, action string, body io.Reader, entity ...interface{}) (error, string) {
	req, _ := http.NewRequest(method, server+action, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err, "cannot send request to dnspod: %s"
	}
	if resp.StatusCode > 399 {
		return errors.New(resp.Status), "invalid status: %s"
	}
	if len(entity) > 0 {
		if resp.Body == nil {
			return errors.New("no content"), "invalid body: %s"
		}
		defer resp.Body.Close()
		bin, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err, "error occurred during body reader: %s"
		}
		return json.Unmarshal(bin, entity[0]), "error occurred during body convert: %s"
	}
	return nil, ""
}

//print & exit if error occurred.
//you can use this way to exit when critical error occurred.
func exitIfErr(err error, format ...string) {
	if err != nil {
		if len(format) == 0 {
			format = []string{"%s"}
		}
		os.Stderr.WriteString(fmt.Sprintf(format[0], err.Error()) + "\r\n")
		os.Exit(1)
	}
}
