package controller

// To contain generic structs and handlers for controllers

type JSONResponse struct {
	Message  string      `json:"message"`
	Response interface{} `json:"response"`
}
