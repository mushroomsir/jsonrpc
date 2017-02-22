package jsonrpc

import (
	"encoding/json"
	"errors"
)

var (
	jsonRPCVersion = "2.0"
	invalid        = "invalid"
	notification   = "notification"
	request        = "request"
	errorType      = "error"
	success        = "success"
)

// ClientRequest represents a JSON-RPC data from client.
type ClientRequest struct {
	Type     string
	PlayLoad *PayloadReq
}

// PayloadReq represents a JSON-RPC request.
type PayloadReq struct {
	Version string `json:"jsonrpc"`
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// Object to pass as request parameter to the method.
	Params interface{} `json:"params,omitempty"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	ID interface{} `json:"id,omitempty"`
}

// ParseReq ...
func ParseReq(msg string) (req *ClientRequest, err error) {
	if msg == "" {
		err = errors.New("empty jsonrpc message")
		return
	}
	p := &PayloadReq{}
	err = json.Unmarshal([]byte(msg), p)
	if err != nil {
		err = errors.New("invalid jsonrpc message structures")
		return
	}
	req = &ClientRequest{PlayLoad: p}
	if p.Version != jsonRPCVersion {
		req.Type = invalid
		err = errors.New("invalid jsonrpc version")
	} else if p.ID == "" {
		req.Type = notification
	} else if p.Method == "" {
		req.Type = request
	} else {
		err = errors.New("invalid jsonrpc method")
	}
	return
}

// Request creates a JSON-RPC 2.0 request object, return JsonRpc json.
// the id must be {String|Integer|nil} type
func Request(id interface{}, method string, args ...interface{}) (result string, err error) {
	if err = validateID(id); err != nil {
		return
	}
	p := &PayloadReq{
		Version: jsonRPCVersion,
		Method:  method,
		ID:      id,
	}
	if len(args) > 0 {
		p.Params = args[0]
	}
	b, err := json.Marshal(p)
	return string(b), err
}

// Notification Creates a JSON-RPC 2.0 notification object, return JsonRpc json.
func Notification(method string, args ...interface{}) (string, error) {
	return Request(nil, method, args...)
}

// ClientResponse represents a JSON-RPC response returned to a client.
type ClientResponse struct {
	Type     string
	PlayLoad *PayloadRes
}

// PayloadRes represents a JSON-RPC request.
type PayloadRes struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *ErrorObj   `json:"error,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

// ParseRes ...
func ParseRes(msg string) (res *ClientResponse, err error) {
	if msg == "" {
		err = errors.New("empty jsonrpc message")
		return
	}
	p := &PayloadRes{}
	err = json.Unmarshal([]byte(msg), p)
	if err != nil {
		err = errors.New("invalid jsonrpc message structures")
		return
	}
	res = &ClientResponse{PlayLoad: p}
	if p.Version != jsonRPCVersion {
		res.Type = invalid
		err = errors.New("invalid jsonrpc version")
	} else if p.Error != nil {
		res.Type = errorType
	} else if p.Result != nil {
		res.Type = success
	} else {
		err = errors.New("invalid jsonrpc object")
	}
	return
}

// Success Creates a JSON-RPC 2.0 success response object, return JsonRpc json.
// The result parameter is required
func Success(id interface{}, result interface{}) (str string, err error) {
	if err = validateID(id); err != nil {
		return
	}
	if result == nil {
		err = errors.New("result parameter is required")
		return
	}
	p := &PayloadRes{
		Version: jsonRPCVersion,
		Result:  result,
		ID:      id,
	}
	json, err := json.Marshal(p)
	return string(json), err
}

// CreateError a JsonRpc error
func CreateError(code int, msg string) *ErrorObj {
	return &ErrorObj{Code: code, Message: msg}
}

//Error Creates a JSON-RPC 2.0 error response object, return JsonRpc json.
func Error(id interface{}, rpcerr *ErrorObj) (str string, err error) {
	if err = validateID(id); err != nil {
		return
	}
	p := &PayloadRes{
		Version: jsonRPCVersion,
		Error:   rpcerr,
		ID:      id,
	}
	json, err := json.Marshal(p)
	return string(json), err
}

// ErrorObj ...
type ErrorObj struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func validateID(id interface{}) (err error) {
	if id != nil {
		switch id.(type) {
		case string:
		case int:
		default:
			err = errors.New("invalid id that MUST contain a String, Number, or NULL value")
		}
	}
	return
}
