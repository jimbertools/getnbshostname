package getnbshostname

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetname(t *testing.T) {
	//name := GetNetbiosNameFromIp("192.168.2.130")
	//log.Println(name)
	//assert.Equal(t, "test", name)

	name := GetNetbiosNameFromIp("127.0.0.1")
	log.Println(name)
	assert.Equal(t, "", name)
}
