package test

type TestErrorMessage string

const (
	TestUnexpectedValue TestErrorMessage = "expected %v, but got %v"
	TestMismatchCompare TestErrorMessage = "mismatch (-want +got):\n%s"
)

func (m TestErrorMessage) String() string {
	return string(m)
}
