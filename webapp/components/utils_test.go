package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_appendQueryParams(t *testing.T) {
	basePath := "/api/example"
	res := appendQueryParams(basePath)
	assert.Equal(t, basePath, res)
	res = appendQueryParams(basePath, "param1")
	assert.Equal(t, basePath, res)
	res = appendQueryParams(basePath, "key1", "val1")
	assert.Equal(t, basePath+"?key1=val1", res)
	res = appendQueryParams(basePath, "key1", "val1", "key2")
	assert.Equal(t, basePath+"?key1=val1", res)
	res = appendQueryParams(basePath, "key1", "val1", "key2", "val2")
	assert.Equal(t, basePath+"?key1=val1&key2=val2", res)
}
