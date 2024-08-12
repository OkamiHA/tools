package rule

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildModSecurityRuleFromCondition(t *testing.T) {
	condition := ConditionRule{
		ParamId:  "country",
		Operator: "not in",
		Value:    []interface{}{"Vietnam", "China"},
	}
	result, err := BuildModSecurityRuleFromCondition(condition)
	assert.Nil(t, err)
	fmt.Println(result)
	assert.NotNil(t, result)
}
