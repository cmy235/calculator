package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// func writeResponse(out Output, response http.ResponseWriter, errType bool) {

func checkInputs(action string, response http.ResponseWriter, request *http.Request) {
	var answer float64
	strX, strY, raw := parseValues(request)
	path := "/" + action + "/?" + raw
	x, y, err := convertToFloat(strX, strY)

	if err != nil {
		writeResponse(Output{}, response, true)
		return
	}

	inCache, value := datastore.get(path)
	if inCache {
		value.operation.Cached = true
		writeResponse(value.operation, response, false)
		value.resetExpiration()
		return
	}

	answer, err = calculate(action, x, y)
	if err != nil {
		writeResponse(Output{}, response, true)
		return
	}

	out := Output{action, x, y, answer, false}
	datastore.set(out, path)
	writeResponse(out, response, false)
}

func calculate(action string, x, y float64) (float64, error) {
	var err error

	switch action {
	case "add":
		return x + y, nil
	case "subtract":
		return x - y, nil
	case "multiply":
		return x * y, nil
	case "divide":
		fmt.Println("x, y => ", x, y)
		if y == 0.0 {
			err = errors.New("Error calculating values")
			return 0, err
		}
		return x / y, nil
	default:
		err = errors.New("Error calculating values")
		return 0.0, err
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
