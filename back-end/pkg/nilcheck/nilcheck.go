package nilcheck

func NotNilBool(newVal *bool, oldVal bool) bool {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func NotNilUint(newVal *uint, oldVal uint) uint {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func NotNilFloat32(newVal *float32, oldVal float32) float32 {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func NotNilString(newVal *string, oldVal string) string {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}
