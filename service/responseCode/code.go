package responseCode

var code map[string]int

func init() {
	code = map[string]int{
		"authErr":  10001,
		"dbErr":    10002,
		"loginErr": 10003,
	}
}
