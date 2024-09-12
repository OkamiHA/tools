package golang

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMergeTwoAggregateResults(t *testing.T) {
	firstAggregate := map[string]JSONObject{
		"level 1": {"age": 15, "name": "hello"},
		"level 3": {"age": 16, "name": "world"},
	}
	secondAggregate := map[string]JSONObject{
		"level 1": {"location": "VN", "sex": "male"},
		"level 2": {"location": "EN", "sex": "female"},
	}

	result := MergeTwoAggregateResults(firstAggregate, secondAggregate)
	j, _ := json.Marshal(result)
	fmt.Println(j)
}
