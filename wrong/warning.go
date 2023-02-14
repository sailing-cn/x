package wrong

type Warning struct {
	s string
}

func (w Warning) Error() string {
	return w.s
}

func (w Warning) Message() string {
	return w.s
}
func Warn(text string) error {
	return &Warning{text}
}
