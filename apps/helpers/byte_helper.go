package helpers

import "bytes"

func FromBytesToString(buf []byte) string {
	n := bytes.IndexByte(buf, 0)
	e := bytes.IndexByte(buf,10)
	if (n - e == 1){
		return string(buf[:e])
	}
	return string(buf[:n])
}
