package internal

import (
	"fmt"
	"regexp"
	"strings"
)

func splitOperations(id int, expression string) []Operation {
	var operations []Operation

	if len(expression) == 0 {
		return operations
	}

	re := regexp.MustCompile(`\-?\d+[\+\-]\d+`)

	matches := re.FindAllString(expression, -1)

	if len(matches) == 0 {
		operations = append(operations, Operation{id: id, text: expression + "+0"})
	}

	if len(matches)%2 != 0 {
		expression = expression + "+0"
		matches = re.FindAllString(expression, -1)
	}

	for _, match := range matches {
		operations = append(operations, Operation{id: id, text: match})
	}

	return operations
}

func containsMinus(s string) bool {
	for _, char := range s {
		if char == '-' {
			return true
		}
	}
	return false
}

func removeSubstrings(str string, substrs []Operation) string {
	for _, elem := range substrs {
		substr := elem.text
		for strings.Contains(str, substr) {
			index := strings.Index(str, substr)
			str = str[:index] + str[index+len(substr):]
		}
	}
	return str
}

func createGroups(id int, prob_slice []string) []Operation {
	var res []Operation
	strike := false
	for i := 1; i < len(prob_slice)-1; i += 2 {
		sym := prob_slice[i]

		if sym == "*" || sym == "/" {
			if strike {
				res[len(res)-1].text += sym + prob_slice[i+1]
			} else {
				res = append(res, Operation{
					id:   id,
					text: fmt.Sprintf("%s%s%s", prob_slice[i-1], sym, prob_slice[i+1]),
				})
				strike = true
			}
		} else {
			strike = false
		}

	}

	string_simpl_prob := removeSubstrings(strings.Join(prob_slice, ""), res)

	re := regexp.MustCompile(`[\+\-]+`)
	string_simpl_prob = re.ReplaceAllStringFunc(string_simpl_prob, func(match string) string {
		if containsMinus(match) {
			return "-"
		}
		return "+"
	})

	res = append(res, splitOperations(id, string_simpl_prob)...)

	return res
}

func parseProblem(prob string) []string {
	separatorPattern := "([0-9]+)|([*+/\\-])"

	regex := regexp.MustCompile(separatorPattern)

	matches := regex.FindAllStringSubmatch(prob, -1)

	result := make([]string, 0)

	for _, match := range matches {
		if match[1] != "" {
			// Operand
			result = append(result, match[1])
		} else {
			// Operator
			result = append(result, match[2])
		}
	}

	return result
}
