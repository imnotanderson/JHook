package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

const SUCCESS_MARK string = "build success:"

var htm []byte

func main() {
	var err error
	htm, err = ioutil.ReadFile("test.htm")
	if checkErr(err) {
		return
	}
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./"))))
	http.HandleFunc("/", Handle)
	http.ListenAndServe(":4000", nil)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	args := r.FormValue("args")
	if args != "" {
		downloadUrl := GoBuild(args)
		w.Write([]byte(downloadUrl))
	} else {
		w.Write(htm)
	}
}

func GoBuild(args string) string {
	pCmd := exec.Command(fmt.Sprintf("build_%s.sh", args))
	out, err := pCmd.CombinedOutput()
	if checkErr(err) {
		return "Err:" + err.Error()
	}
	str := string(out)
	if str[:len(SUCCESS_MARK)] == SUCCESS_MARK {
		return str
	}
	str = str[:strings.Index(str, SUCCESS_MARK)]
	return str
}

func checkErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
