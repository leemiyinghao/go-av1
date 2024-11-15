package coordinator_domain

import (
	"github.com/maja42/goval"
)

func eval(filter string, context map[string]interface{}) bool {
	evaluator := goval.NewEvaluator()
	result, err := evaluator.Evaluate(filter, context, nil)
	if err != nil {
		return false
	}
	return result == true
}
