package driver

import (
	"strconv"
	"strings"
)

/*
	结果：
	bytes := [4]byte{1,2,3,4}
	str := convert(bytes[:])
*/
func convertByteToString(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, " ")
}
