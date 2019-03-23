package main

// TODO
// Fix test cases since refactor
import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandlers(t *testing.T) {
	datastore.new()

	var params = []struct {
		path   string
		output string
		action func(http.ResponseWriter, *http.Request)
	}{
		{
			"/add/?x=4&y=10",
			`{"Action":"add","X":4,"Y":10,"Answer":14,"Cached":false}`,
			add,
		},

		{
			"/subtract/?x=2&y=1",
			`{"Action":"subtract","X":2,"Y":1,"Answer":1,"Cached":false}`,
			subtract,
		},
		{
			"/multiply/?x=-2903425&y=0.0119",
			`{"Action":"multiply","X":-2903425,"Y":0.0119,"Answer":-34550.7575,"Cached":false}`,
			multiply,
		},
		{
			"/divide/?x=-34243&y=10.020",
			`{"Action":"divide","X":-34243,"Y":10.02,"Answer":-3417.4650698602795,"Cached":false}`,
			divide,
		},
	}

	for _, values := range params {

		req, err := http.NewRequest("GET", values.path, nil)
		if err != nil {
			t.Error("Request fail ", err)
		}

		responseRecorder := httptest.NewRecorder()
		handler := http.HandlerFunc(values.action)
		handler.ServeHTTP(responseRecorder, req)

		if status := responseRecorder.Code; status != http.StatusOK {
			t.Errorf("Expected status: %v but received %v", status, http.StatusOK)
		}

		received := responseRecorder.Body.String()
		if received != values.output {
			t.Errorf("Expected %s \n received: %s", values.output, received)
		}

	}

}

func TestExpireOldKeys(t *testing.T) {
	datastore.new()

	var params = []struct {
		toSet        *Output
		inCache      bool
		url          string
		timeSinceSet int
	}{
		{
			&Output{"add", 2, 5, 7, false},
			false,
			"url1",
			-61,
		},
		{
			&Output{"subtract", 9, 0, 9, false},
			false,
			"url2",
			-61,
		},
		{
			&Output{"divide", 10, 2, 5, false},
			true,
			"url3",
			-55,
		},
	}

	for _, values := range params {
		datastore.set(values.toSet, values.url, time.Now().Add(time.Duration(values.timeSinceSet)*time.Second))
		datastore.expireOldKeys()

		result, _ := datastore.get(values.url)
		if result != values.inCache {
			t.Errorf("Expected %v received %v", result, values.inCache)
		}
	}

}

func TestWriteResponse(t *testing.T) {
	var params = []struct {
		result       *Output
		bytesWritten []byte
	}{
		{
			&Output{"divide", 9, 3, 3, true},
			[]byte{123, 34, 65, 99, 116, 105, 111, 110, 34, 58, 34, 100, 105, 118, 105, 100, 101, 34, 44, 34, 88, 34, 58, 57, 44, 34, 89, 34, 58, 51, 44, 34, 65, 110, 115, 119, 101, 114, 34, 58, 51, 44, 34, 67, 97, 99, 104, 101, 100, 34, 58, 116, 114, 117, 101, 125},
		},
		{
			&Output{"add", 2, 2, 4, false},
			[]byte{123, 34, 65, 99, 116, 105, 111, 110, 34, 58, 34, 97, 100, 100, 34, 44, 34, 88, 34, 58, 50, 44, 34, 89, 34, 58, 50, 44, 34, 65, 110, 115, 119, 101, 114, 34, 58, 52, 44, 34, 67, 97, 99, 104, 101, 100, 34, 58, 102, 97, 108, 115, 101, 125},
		},
		{
			&Output{"multiply", 2, 2, 4, false},
			[]byte{123, 34, 65, 99, 116, 105, 111, 110, 34, 58, 34, 109, 117, 108, 116, 105, 112, 108, 121, 34, 44, 34, 88, 34, 58, 50, 44, 34, 89, 34, 58, 50, 44, 34, 65, 110, 115, 119, 101, 114, 34, 58, 52, 44, 34, 67, 97, 99, 104, 101, 100, 34, 58, 102, 97, 108, 115, 101, 125},
		},
	}

	w := httptest.NewRecorder()

	for _, values := range params {
		received := writeResponse(values.result, w)
		if bytes.Equal(values.bytesWritten, received) == false {
			t.Errorf("Expected %v, received %v", values.bytesWritten, received)
		}
	}
}

func TestCompute(t *testing.T) {
	var params = []struct {
		url    string
		action string
	}{
		{
			"/divide/?x=20&y=10",
			"divide",
		},
		{
			"/add/?x=20&y=100",
			"add",
		},
		{
			"/subtract/?x=0.02&y=-8",
			"subtract",
		},
		{
			"/add/?x=44520&y=1239",
			"add",
		},
	}

	datastore.new()
	w := httptest.NewRecorder()

	for _, values := range params {
		req, err := http.NewRequest("GET", values.url, nil)
		if err != nil {
			t.Error(err)
		}
		compute(values.action, w, req)
		in, _ := datastore.get(values.url)

		if !in {
			t.Errorf("Expected %s to be in the cache, was not found", values.url)
		}
	}
}

func TestConvertToFloat(t *testing.T) {
	var params = []struct {
		strX string
		strY string
		x    float64
		y    float64
	}{
		{
			"9.1314",
			"2.0",
			9.1314,
			2.0,
		},
		{
			"1",
			"-0.7",
			1,
			-0.7,
		},
		{
			"2481234012649182376419283461928",
			"-9482395767642679128736498769",
			2481234012649182376419283461928,
			-9482395767642679128736498769,
		},
	}

	for _, values := range params {
		x, y, err := convertToFloat(values.strX, values.strY)
		if x != values.x || y != values.y || err != nil {
			t.Errorf("Expected %e, %e but received %e, %e", values.x, values.y, x, y)
		}
	}
}

func TestCalculate(t *testing.T) {
	var params = []struct {
		action string
		x      float64
		y      float64
		answer float64
	}{
		{
			"add",
			2.2,
			1,
			3.2,
		},
		{
			"subtract",
			92342348243,
			234234243,
			92108114000,
		},
		{
			"multiply",
			-4325,
			0.023,
			-99.475,
		},
	}

	for _, values := range params {
		result := calculate(values.action, values.x, values.y)
		if result != values.answer {
			t.Errorf("Expected %v, received %v", values.answer, result)
		}
	}
}

func TestSet(t *testing.T) {
	datastore.new()

	var params = []struct {
		url    string
		result *Output
	}{
		{
			"/divide/?x=20&y=10",
			&Output{"add", 20, 10, 2, false},
		},
		{
			"/add/?x=20&y=100",
			&Output{"add", 20, 100, 120, false},
		},
		{
			"/subtract/?x=2&y=-8",
			&Output{"subtract", 2, -8, 10, false},
		},
		{
			"/add/?x=44520&y=1239",
			&Output{"add", 44520, 1239, 45759, false},
		},
	}

	for _, values := range params {
		datastore.set(values.result, values.url, time.Now())
		output := datastore.keyVals[values.url]

		if output.operation != values.result {
			t.Errorf("Expected %v, received %v", values.result, output.operation)
		}
	}

	if len(datastore.keyVals) != len(params) {
		t.Errorf("Values missing from cache")
	}
}

func TestGet(t *testing.T) {
	datastore.new()

	var params = []struct {
		url    string
		result *Output
	}{
		{
			"/divide/?x=20&y=10",
			&Output{"add", 20, 10, 2, false},
		},
		{
			"/add/?x=20&y=100",
			&Output{"add", 20, 100, 120, false},
		},
		{
			"/subtract/?x=2&y=-8",
			&Output{"subtract", 2, -8, 10, false},
		},
		{
			"/add/?x=44520&y=1239",
			&Output{"add", 44520, 1239, 45759, false},
		},
	}

	for _, values := range params {
		datastore.set(values.result, values.url, time.Now())
		in, output := datastore.get(values.url)

		if !in || output.operation != values.result {
			t.Errorf("Expected %v, received %v", values.result, output.operation)
		}
	}
}

func TestResetExpiration(t *testing.T) {
	var params = []struct {
		op         *Output
		expireTime time.Time
	}{
		{
			&Output{"add", 20, 10, 2, false},
			time.Now().Add(-time.Second * 60),
		},
		{
			&Output{"add", 20, 100, 120, false},
			time.Now().Add(-time.Second * 10),
		},
		{
			&Output{"subtract", 2, -8, 10, false},
			time.Now().Add(-time.Second * 45),
		},
		{
			&Output{"add", 44520, 1239, 45759, false},
			time.Now().Add(-time.Second * 25),
		},
	}

	for _, values := range params {
		cacheVal := &val{
			values.op,
			values.expireTime,
		}

		cacheVal.resetExpiration()

		if cacheVal.expiration.Unix() < values.expireTime.Unix() {
			t.Errorf("Expected new expiration time to be later than original %v, instead received %v", values.expireTime, cacheVal.expiration)
		}
	}
}
