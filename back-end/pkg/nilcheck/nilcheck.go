package nilcheck

func IfNotNilBool(newVal *bool, oldVal bool) bool {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func IfNotNilFloat32(newVal *float32, oldVal float32) float32 {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func IfNotNilString(newVal *string, oldVal string) string {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}
