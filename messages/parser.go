package messages

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/msgbox/relay/structs"
	"regexp"
	"strings"
)

// Parses a Message from the queue

func Parse(data []byte) *structs.Message {
	msg := &structs.Message{}
	proto.Unmarshal(data, msg)

	return msg
}

func ParseName(s string) string {
	account := strings.Split(s, ":")
	return account[0]
}

// Because I have no idea how regexp in go work
// nothing would compile so this shit pile ends up
// in here
func ParseBox(s string) string {
	reg, _ := regexp.Compile(":(.*?)@")
	re := reg.FindString(s)

	re2 := strings.SplitAfter(re, ":")
	if len(re2) < 2 {
		return ""
	}

	re3 := strings.Split(re2[1], "@")
	if len(re3) < 1 {
		return ""
	}

	return re3[0]
}
