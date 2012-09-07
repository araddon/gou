package gou

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	finished chan bool
	lastTest time.Time = time.Now()
	stopper  func()    = func() {}
)

// Use this in combo with StopCheck() for test functions that must start
// processes such as 
func SetStopper(f func()) {
	stopper = f
}

// dumb simple assert for testing, printing
//    Assert(len(items) == 9, t, "Should be 9 but was %d", len(items))
func Assert(is bool, t *testing.T, format string, args ...interface{}) {
	if is == false {
		if logger == nil {
			logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
		}
		DoLog(3, ERROR, fmt.Sprintf(format, args...), logger)
		t.Fail()
	}
}

func StopCheck() {
	t := time.Now()
	if lastTest.Add(time.Second*1).UnixNano() < t.UnixNano() {
		Log(INFO, "Stopping Test")
		finished <- true
		stopper()
	}
	lastTest = t
}

// Simple Fetch Wrapper, given a url it returns bytes
func Fetch(url string) (ret []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		Log(WARN, err.Error())
		return
	}
	ret, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

// posts an application/json to url with body
// ie:   type = application/json
func PostJson(url, body string) (ret string, err error) {
	//Post(url string, bodyType string, body io.Reader) 
	Debug(url)
	buf := bytes.NewBufferString(body)
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		Log(WARN, err.Error())
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// posts a www-form encoded form to url with body
func PostForm(url, body string) (ret string, err error, resp *http.Response) {
	//Post(url string, bodyType string, body io.Reader) 
	buf := bytes.NewBufferString(body)
	resp, err = http.Post(url, "application/x-www-form-urlencoded", buf)
	if err != nil {
		Log(WARN, err.Error())
		return "", err, resp
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err, resp
	}

	return string(data), nil, resp
}
