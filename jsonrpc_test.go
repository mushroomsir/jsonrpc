package jsonrpc_test

import "testing"
import "github.com/stretchr/testify/assert"
import "github.com/mushroomsir/jsonrpc"

var (
	jsonRPCVersion = "2.0"
	invalid        = "invalid"
	notification   = "notification"
	request        = "request"
	errorType      = "error"
	success        = "success"
)

func TestProducer(t *testing.T) {

	t.Run("jsonrpc with request that should be", func(t *testing.T) {
		assert := assert.New(t)

		val, err := jsonrpc.Request(123, "update")
		assert.Equal("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":123}", val)

		val, err = jsonrpc.Request("123", "update")
		assert.Equal("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":\"123\"}", val)

		val, err = jsonrpc.Request(true, "update")
		assert.NotNil(err)
		assert.Equal("invalid id that MUST contain a String, Number, or NULL value", err.Error())

		val, err = jsonrpc.Notification("update")
		assert.Equal("{\"jsonrpc\":\"2.0\",\"method\":\"update\"}", val)

		val, err = jsonrpc.Notification("update", 0)
		assert.Equal("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"params\":0}", val)

	})
	t.Run("jsonrpc with response that should be", func(t *testing.T) {
		assert := assert.New(t)

		val, err := jsonrpc.Success("123", nil)
		assert.NotNil(err)
		assert.Equal("result parameter is required", err.Error())

		val, err = jsonrpc.Success("123", "OK")
		assert.Nil(err)
		assert.Equal("{\"jsonrpc\":\"2.0\",\"result\":\"OK\",\"id\":\"123\"}", val)

		val, err = jsonrpc.Success(123, []string{})
		assert.Nil(err)
		assert.Equal("{\"jsonrpc\":\"2.0\",\"result\":[],\"id\":123}", val)

		rpcerr := jsonrpc.CreateError(1, "test")
		val, err = jsonrpc.Error(nil, rpcerr)
		assert.Equal("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":1,\"message\":\"test\"}}", val)
	})

	t.Run("jsonrpc with parsereq that should be", func(t *testing.T) {

		assert := assert.New(t)

		val, err := jsonrpc.ParseReq("")
		assert.Empty(val)
		assert.NotNil(err)
		assert.Equal("empty jsonrpc message", err.Error())

		val, err = jsonrpc.ParseReq("{\"jsonrpc\":\"2.0\",\"method\":\"update\",\"id\":\"123\"}")
		assert.Nil(err)
		assert.Equal("123", val.PlayLoad.ID)
		assert.Equal("update", val.PlayLoad.Method)
		assert.Equal(request, val.Type)
	})
}
