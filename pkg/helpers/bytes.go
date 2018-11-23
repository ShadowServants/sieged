package helpers

func FromBytesToString(buf []byte, index int) string {
	return string(buf[:index])
}
