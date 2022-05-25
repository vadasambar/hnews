package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	appsv1 "github.com/vadasambar/hnews/api/v1"
)

var (
	// scoreRegex is the regex for conditions ">[number-here]", ">=[number-here]", "<[number-here]",
	// "<=[number-here]", "!=[number-here]" and "=[number-here]"
	scoreRegex = regexp.MustCompile("^[[:space:]]*(>|>=|=|!=|<|<=)[[:space:]]*([[:digit:]]+[[:space:]]*$)")
)

// EvalCond takes a value and a condition
// and evaluates the condition on the value
// e.g., value = 5, cond = "<10"
// evalCond would do 5 < 10 => returns true
func EvalCond(value int, cond appsv1.Comparison) bool {
	// https: //play.golang.com/p/B8ZgghEBK4k

	result := scoreRegex.FindAllStringSubmatch(string(cond), -1)
	if len(result[0]) < 3 {
		return false
	}

	comparisonOperator := strings.TrimSpace(result[0][1])
	condValue, err := strconv.Atoi(strings.TrimSpace(result[0][2]))
	if err != nil {
		fmt.Println("err", err)
		return false
	}
	switch comparisonOperator {
	case ">":
		return value > condValue
	case ">=":
		return value >= condValue
	case "<":
		return value < condValue
	case "<=":
		return value <= condValue
	case "=":
		return value == condValue
	case "!=":
		return value != condValue
	}

	return false
}
