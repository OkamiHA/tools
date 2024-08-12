package rule

import (
	"errors"
	"example_api/internal/constants"
	"example_api/internal/utils"
	"fmt"
	"reflect"
	"strings"
)

type (
	ConditionRule struct {
		ParamId  string      `json:"paramId"`
		Value    interface{} `json:"value"`
		Operator string      `json:"operator"`
	}
)

func BuildModSecurityRuleFromCondition(condition ConditionRule) (string, error) {
	operatorNew, ok := MAP_OPERATOR[condition.Operator]
	if !ok {
		return constants.EmptyString, errors.New(fmt.Sprintf(constants.InvalidOperatorMessage, condition.Operator))
	}
	isNotMatch := false
	if strings.HasPrefix(condition.Operator, "not") {
		isNotMatch = true
	}
	operatorValue, err := BuildModSecurityRuleOperator(operatorNew, condition.Value, isNotMatch)
	if err != nil {
		return constants.EmptyString, err
	}
	return fmt.Sprintf("SecRule %s \"%s\"", condition.ParamId, operatorValue), nil
}

func BuildModSecurityRuleOperator(operator string, value interface{}, isNotMatch bool) (string, error) {
	typeOfValue := reflect.TypeOf(value).Kind()
	if utils.IsNumeric(value) {
		return BuildModSecurityRuleFromNumeric(operator, value)
	}
	switch typeOfValue {
	case reflect.Slice:
		values, ok := value.([]interface{})
		print(ok)
		return BuildModSecurityRuleFromSliceValue(operator, values, isNotMatch)
	case reflect.String:
		valueStr, _ := value.(string)
		return BuildModSecurityRuleFromString(operator, valueStr)
	default:
		return constants.EmptyString, errors.New(fmt.Sprintf(constants.InvalidValueByOperatorMessage, operator))
	}
}

func BuildModSecurityRuleFromSliceValue(operator string, values []interface{}, isNotMatch bool) (string, error) {
	if !utils.IsSliceContainsElement([]string{"@pm", "@streq", "@rx"}, operator) {
		return constants.EmptyString, errors.New(fmt.Sprintf(constants.InvalidOperatorMessage, operator))
	}
	result := constants.EmptyString
	for index, element := range values {
		if index == 0 {
			if isNotMatch {
				result = fmt.Sprintf("!%v", element)
			} else {
				result = fmt.Sprintf("%v", element)
			}
		} else {
			result += fmt.Sprintf(" %v", element)
		}
	}
	return fmt.Sprintf("%s %s", operator, result), nil
}

func BuildModSecurityRuleFromNumeric(operator string, value interface{}) (string, error) {
	allowOperator := []string{"@gt", "@ge", "@lt", "@le", "@eq"}
	if utils.IsSliceContainsElement(allowOperator, operator) {
		return fmt.Sprintf("%s %v", operator, value), nil
	}
	return constants.EmptyString, errors.New("invalid operator")
}

func BuildModSecurityRuleFromString(operator string, value string) (string, error) {
	allowOperator := []string{"@rx", "@pmFromFile", "@pm"}
	if !utils.IsSliceContainsElement(allowOperator, operator) {
		return constants.EmptyString, errors.New(fmt.Sprintf(constants.InvalidOperatorMessage, operator))
	}
	return fmt.Sprintf("%s %s", operator, value), nil
}
