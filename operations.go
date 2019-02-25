package main

import (
	"net/http"
	"strconv"
	"time"
)

func compute(action string, response http.ResponseWriter, request *http.Request) {
	var answer float64
	strX, strY, raw := parseValues(request)
	path := "/" + action + "/?" + raw

	x, y, err := convertToFloat(strX, strY)

	if err != nil {
		var errResponse *Output
		if strX == "" || strY == "" {
			errResponse = &Output{"The operation failed.  Did you forget the operands?", x, y, 0, false}
		} else {
			errResponse = &Output{"The operation failed.  Please don't try to do math with letters!", x, y, 0, false}
		}
		_ = writeResponse(errResponse, response)
		return
	}

	inCache, value := datastore.get(path)
	if inCache {
		value.operation.Cached = true
		_ = writeResponse(value.operation, response)
		value.resetExpiration()
		return
	}

	answer = calculate(action, x, y)
	out := &Output{action, x, y, answer, false}
	datastore.set(out, path, time.Now())
	_ = writeResponse(out, response)
}

func calculate(action string, x, y float64) float64 {
	switch action {
	case "add":
		return x + y
	case "subtract":
		return x - y
	case "multiply":
		return x * y
	case "divide":
		return x / y
	default:
		panic("Compuation not found")
	}
}

func convertToFloat(strX, strY string) (float64, float64, error) {
	x, err := strconv.ParseFloat(strX, 64)
	if err != nil {
		return 0, 0, err
	}

	y, err := strconv.ParseFloat(strY, 64)
	if err != nil {
		return 0, 0, err
	}

	return x, y, nil
}
