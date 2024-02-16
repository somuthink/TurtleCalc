package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func ExecutuionTime(text string) int {
	var exec_time int
	syms := []string{"+", "-", "*", "/"}

	costs := make([]int, 4)
	res, _ := glob_db.db.Query("SELECT (exec_time) FROM `operations`")
	i := 0
	for res.Next() {
		var val int
		res.Scan(&val)
		costs[i] = val
		i++
	}

	if rune(text[0]) == rune('-') {
		text = text[1:]
	}

	for i, cost := range costs {
		exec_time += cost * strings.Count(text, syms[i])
	}

	return exec_time
}

func splitOperations(id int, expression string) ([]Operation, error) {
	var operations []Operation

	if len(expression) == 0 {
		return operations, errors.New("expression is empty")
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

	return operations, nil
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

func createGroups(id int, prob_slice []string) ([]Operation, error) {
	var res []Operation
	strike := false
	for i := 1; i < len(prob_slice)-1; i += 2 {
		sym := prob_slice[i]

		if sym == "*" || sym == "/" {
			if strike {
				res[len(res)-1].text += sym + prob_slice[i+1]
			} else {
				starts_with_minus := ""
				if i-2 >= 0 {
					if prob_slice[i-2] == "-" {
						starts_with_minus = "-"
					}
				}
				res = append(res, Operation{
					id:   id,
					text: fmt.Sprintf("%s%s%s%s", starts_with_minus, prob_slice[i-1], sym, prob_slice[i+1]),
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

	resSplit, err := splitOperations(id, string_simpl_prob)
	if err != nil {
		return nil, err
	}

	res = append(res, resSplit...)

	return res, nil
}

func parseProblem(prob string) ([]string, error) {
	separatorPattern := "([0-9]+)|([*+/\\-])"

	regex := regexp.MustCompile(separatorPattern)

	matches := regex.FindAllStringSubmatch(prob, -1)

	if len(matches) == 0 {
		return nil, errors.New("no matches found")
	}

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

	return result, nil
}
