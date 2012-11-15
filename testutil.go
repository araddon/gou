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

// take two floats, compare, need to be within .1%
func CloseEnuf(a, b float64) bool {
	c := a / b
	if c > .999 && c < 1.001 {
		return true
	}
	return false
}
func StartTest() {
	lastTest = time.Now()
}

func StopCheck() {
	t := time.Now()
	if lastTest.Add(time.Second*1).UnixNano() < t.UnixNano() {
		Log(INFO, "Stopping Test")
		finished <- true
		stopper()
		os.Exit(1)
	}
}

// Simple Fetch Wrapper, given a url it returns bytes
func Fetch(url string) (ret []byte, err error) {
	resp, err := http.Get(url)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
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
func PostJson(url, body string) (ret string, err error, resp *http.Response) {
	//Post(url string, bodyType string, body io.Reader) 
	Debug(url)
	buf := bytes.NewBufferString(body)
	resp, err = http.Post(url, "application/json", buf)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
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

// issues http delete an application/json to url with body
// ie:   type = application/json
func DeleteJson(url, body string) (ret string, err error, resp *http.Response) {
	//Post(url string, bodyType string, body io.Reader) 
	Debug(url)
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("DELETE", url, buf)
	if err != nil {
		return
	}
	//
	req.Header.Add("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req) //(url, "application/json", buf)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
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

// posts a www-form encoded form to url with body
func PostForm(url, body string) (ret string, err error, resp *http.Response) {
	//Post(url string, bodyType string, body io.Reader) 
	buf := bytes.NewBufferString(body)
	resp, err = http.Post(url, "application/x-www-form-urlencoded", buf)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		Log(WARN, url, body, err.Error())
		return "", err, resp
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err, resp
	}

	return string(data), nil, resp
}
