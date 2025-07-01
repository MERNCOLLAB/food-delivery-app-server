package utils

func SafeString(s *string, def string) string {
	if s != nil {
		return *s
	}
	return def
}
