package scopename

type ScopeName int

const (
	APP = iota
	IF
	FOR
	FUNCTION
)

var scopeNames = map[ScopeName]string{
	APP:      "App",
	IF:       "If",
	FOR:      "For",
	FUNCTION: "Function",
}

func (s ScopeName) String() string {
	return scopeNames[s]
}
func (l ScopeName) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}
