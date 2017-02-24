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
	// RPCParseError means invalid JSON was received by the server.An error occurred on the server while parsing the JSON text.
	RPCParseError = &ErrorObj{Code: -32700, Message: "Parse error"}
	// RPCInvalidRequest means the JSON sent is not a valid Request object.
	RPCInvalidRequest = &ErrorObj{Code: -32600, Message: "Invalid Request"}
	// RPCNotFound means the method does not exist / is not available.
	RPCNotFound = &ErrorObj{Code: -32601, Message: "Method not found"}
	// RPCInvalidParams means Invalid method parameter(s).
	RPCInvalidParams = &ErrorObj{Code: -32602, Message: "Invalid params"}
	// RPCInternalError means Internal JSON-RPC error.
	RPCInternalError = &ErrorObj{Code: -32603, Message: "Internal error"}
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

// ErrorObj ...
type ErrorObj struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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
	val, err := json.Marshal(p)
	return string(val), err
}

// Notification Creates a JSON-RPC 2.0 notification object, return JsonRpc json.
func Notification(method string, args ...interface{}) (string, error) {
	return Request(nil, method, args...)
}

//Batch ...
func Batch(batch ...string) (arrstr string) {
	if len(batch) == 0 {
		return "[]"
	}
	arrstr = "["
	for index := 0; index < len(batch)-1; index++ {
		arrstr += batch[index]
		arrstr += ","
	}
	arrstr += batch[len(batch)-1]
	arrstr += ",]"
	return
}

// ParseReq ...
func ParseReq(msg string) (req *ClientRequest, err error) {
	p := make(map[string]interface{}, 0)
	if err = validateMsg(msg, &p); err != nil {
		return
	}
	req, err = parseReqMap(p)
	return
}

// ParseReqBatch ...
func ParseReqBatch(msg string) (req []*ClientRequest, err error) {
	if msg == "" || len(msg) < 2 {
		err = errors.New("empty message")
		return
	}
	payloads := make([]map[string]interface{}, 0)
	if err = validateMsg(msg, &payloads); err == nil {
		req = make([]*ClientRequest, len(payloads))
		for i, val := range payloads {
			r, _ := parseReqMap(val)
			req[i] = r
		}
	}
	return
}

// ParseRes ...
func ParseRes(msg string) (res *ClientResponse, err error) {
	p := make(map[string]interface{}, 0)
	if err = validateMsg(msg, &p); err != nil {
		return
	}
	return parseResMap(p)
}

// ParseResBatch ...
func ParseResBatch(msg string) (res []*ClientResponse, err error) {
	if msg == "" || len(msg) < 2 {
		err = errors.New("empty message")
		return
	}
	payloads := make([]map[string]interface{}, 0)
	if err = validateMsg(msg, &payloads); err == nil {
		res = make([]*ClientResponse, len(payloads))
		for i, val := range payloads {
			r, _ := parseResMap(val)
			res[i] = r
		}
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
	data, err := json.Marshal(p)
	return string(data), err
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
	data, err := json.Marshal(p)
	return string(data), err
}

// CreateError a JsonRpc error
func CreateError(code int, msg string, data ...interface{}) (obj *ErrorObj) {
	obj = &ErrorObj{Code: code, Message: msg}
	if len(data) > 0 {
		obj.Data = data[0]
	}
	return obj
}

func parseResMap(val map[string]interface{}) (r *ClientResponse, err error) {
	r = &ClientResponse{PlayLoad: &PayloadRes{}}
	if version, ok := val["jsonrpc"]; ok {
		r.PlayLoad.Version = version.(string)
	}
	if ID, ok := val["id"]; ok {
		r.PlayLoad.ID = ID
	}
	if errs, ok := val["error"]; ok {
		r.PlayLoad.Error = &ErrorObj{}
		if err2, ok := errs.(map[string]interface{}); ok {
			r.PlayLoad.Error.Code = int(err2["code"].(float64))
			r.PlayLoad.Error.Message = err2["message"].(string)
			r.PlayLoad.Error.Data = err2["data"]
		}
	}
	if result, ok := val["result"]; ok {
		r.PlayLoad.Result = result
	}
	err = checkResType(r)
	return
}
func parseReqMap(val map[string]interface{}) (r *ClientRequest, err error) {
	r = &ClientRequest{PlayLoad: &PayloadReq{}}
	if version, ok := val["jsonrpc"]; ok {
		r.PlayLoad.Version = version.(string)
	}
	if ID, ok := val["id"]; ok {
		r.PlayLoad.ID = ID
	}
	if method, ok := val["method"]; ok {
		r.PlayLoad.Method = method.(string)
	}
	if params, ok := val["params"]; ok {
		r.PlayLoad.Params = params
	}
	err = checkReqType(r)
	return
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
func validateMsg(msg string, p interface{}) (err error) {
	if msg == "" {
		err = errors.New("empty jsonrpc message")
		return
	}
	err = json.Unmarshal([]byte(msg), p)
	if err != nil {
		err = errors.New("invalid jsonrpc message structures")
	}
	return
}
func checkResType(res *ClientResponse) (err error) {
	p := res.PlayLoad
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
func checkReqType(res *ClientRequest) (err error) {
	p := res.PlayLoad
	if p.Version != jsonRPCVersion {
		res.Type = invalid
		err = errors.New("invalid jsonrpc version")
	} else if p.Method == "" {
		res.Type = invalid
		err = errors.New("invalid jsonrpc object")
	} else if p.ID == nil {
		res.Type = notification
	} else {
		res.Type = request
	}
	return
}
