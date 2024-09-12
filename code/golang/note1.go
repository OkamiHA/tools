package golang

import "errors"

type JSONObject = map[string]interface{}
type (
	AggregateResult struct {
		Id   string
		Data JSONObject
	}
)

func MakeMapAggregateResult(data []JSONObject) (map[string]JSONObject, error) {
	result := make(map[string]JSONObject)
	for _, item := range data {
		idValue, ok := item["_id"]
		if !ok {
			return nil, errors.New("field _id not found")
		}
		result[idValue.(string)] = item
	}
	return result, nil
}
func MergeTwoAggregateResults(firstAggregate, secondAggregate map[string]JSONObject) map[string]JSONObject {
	if len(firstAggregate) == 0 {
		return secondAggregate
	}
	if len(secondAggregate) == 0 {
		return firstAggregate
	}
	result := make(map[string]JSONObject)
	for key, firstAggregateValue := range firstAggregate {
		secondAggregateValue, ok := secondAggregate[key]
		if ok {
			result[key] = MergeTwoJSONObject(firstAggregateValue, secondAggregateValue)
			delete(firstAggregate, key)
			delete(secondAggregate, key)
		}
	}
	for key, firstAggregateValue := range firstAggregate {
		result[key] = firstAggregateValue
	}
	for key, secondAggregateValue := range secondAggregate {
		result[key] = secondAggregateValue
	}
	return result
}

func MergeTwoJSONObject(firstObject, secondObject JSONObject) JSONObject {
	for key, value := range secondObject {
		firstObject[key] = value
	}
	return firstObject
}
