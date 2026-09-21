package usecase_test

// b2i converts bool to int for gomock Times().
func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// strPtr helper for optional string pointers.
func strPtr(s string) *string { return &s }
