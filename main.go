package main

import (
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	server = `https://dnsapi.cn/`

	key    []byte
	action string
	act    *Record
)

func init() {
	if num := len(os.Args); num < 4 {
		exitIfErr(errors.New(strconv.Itoa(num)), `at least four parameters are required, but get %s`)
	}
	var filename, reply = strings.Split(strings.ToLower(filepath.Base(os.Args[0])), "_"), Reply{}
	if action = filename[0]; len(filename) > 1 {
		key, _ = hex.DecodeString(filename[1])
	}
	if len(key) == 0 {
		exitIfErr(errors.New(filename[1]), "invalid credentials: %s")
	}
	exitIfErr(request(http.MethodPost, actions.DomainList, strings.NewReader(`login_token=`+string(key)+`&format=json`), &reply))
	if reply.Status.Code != "1" {
		exitIfErr(errors.New(reply.Status.Message), "refused by server: %s")
	}
	act = reply.Find(os.Args[2])
	if len(act.Domain) == 0 {
		exitIfErr(errors.New(os.Args[1]), "cannot resolve such domain: %s")
	}
	act.Record = os.Args[3]
}

func main() {
	switch action {
	case `del`:
		act.Delete()
		break
	default:
		act.Add()
		break
	}
}
