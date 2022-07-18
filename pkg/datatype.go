package pkg

func BooleanToInt(data bool) int64 {
	if data == false {
		return 0
	}

	return 1
}

func IntToBoolean(data int64) bool {
	if data == 0 {
		return false
	}

	return true
}
