package request

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type (
	AttackStaticRequest struct {
		Search             string
		Companies          []string
		Segmentations      []string
		MethodAttacks      []string
		Distributors       []string
		DateFrom           *int64
		DateTo             *int64
		TotalAttackFrom    *int64
		TotalAttackTo      *int64
		HighestBytesFrom   *int64
		HighestBytesTo     *int64
		HighestPacketsFrom *int64
		HighestPacketsTo   *int64
		DurationFrom       *int64
		DurationTo         *int64
		TotalIpFrom        *int64
		TotalIpTo          *int64
		BpsValue           *int64
		PpsValue           *int64
		DurationBaseline   *int64
	}

	AttackStaticPipeline struct {
		FirstMatch  bson.M
		SecondMatch bson.M
		Group       bson.M
	}
)

func BuildStaticAttackFilterRangeByField(arrayFilter []bson.M, fieldName string, isCheckSize bool, from, to *int64) ([]bson.M, error) {
	if from == nil && to == nil {
		return arrayFilter, nil
	}
	if from != nil && to != nil {
		var filter bson.M
		if isCheckSize {
			filter = bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$gte": bson.A{bson.M{"$size": fmt.Sprintf("$%s", fieldName)}, from}},
						bson.M{"$lte": bson.A{bson.M{"$size": fmt.Sprintf("$%s", fieldName)}, to}},
					},
				},
			}
		} else {
			filter = bson.M{fieldName: bson.M{"$gte": from, "$lte": to}}
		}
		arrayFilter = append(arrayFilter, filter)
	}
	return nil, fmt.Errorf("value of key %s is invalid", fieldName)
}

func BuildStaticAttackFilterByValueInput(fieldName string, value *int64) bson.M {
	if value == nil {
		return nil
	}
	return bson.M{fieldName: bson.M{"$gte": value}}

}

func BuildStaticAttackGroupCondtion(fieldName string, value *int64) bson.M {
	if value == nil {
		return bson.M{"$sum": 1}
	}
	return bson.M{"$sum": bson.M{"$cond": bson.A{
		bson.M{"$gte": bson.A{fmt.Sprintf("$%s", fieldName), value}}, 1, 0,
	}}}
}

func BuildStaticAttackFromSearch(queryString string, isCustomerFilter bool) bson.M {
	if len(queryString) == 0 {
		return nil
	}
	if isCustomerFilter {
		return bson.M{
			"$or": bson.A{
				bson.M{"customer_name": bson.M{"$regex": queryString, "$options": "i"}},
				bson.M{"industry.name": bson.M{"$regex": queryString, "$options": "i"}},
				bson.M{"distributor.name": bson.M{"$regex": queryString, "$options": "i"}},
			},
		}
	}
	return bson.M{"attacktype": bson.M{"$regex": queryString, "$options": "i"}}
}

func BuildFilterFromRequest(req AttackStaticRequest, groupByField string) (bson.M, []bson.M, error) {
	customerFilter := BuildStaticAttackFromSearch(req.Search, true)
	pipelineOverall, err := BuildFilterPipelineFromRequest(req, groupByField, true)
	if err != nil {
		return nil, nil, err
	}
	pipelineStaticDetail, err := BuildFilterPipelineFromRequest(req, groupByField, false)
	if err != nil {
		return nil, nil, err
	}
	

}

func BuildFilterPipelineFromRequest(req AttackStaticRequest, groupByField string, isOverallPipeline bool) ([]bson.M, error) {
	firstMatchFilter, err := BuildFirstMatchFilter(req, isOverallPipeline)
	var pipeline []bson.M
	if err != nil {
		return nil, err
	}
	pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": firstMatchFilter}})
	if isOverallPipeline {
		pipeline = append(pipeline, BuildGroupByWithoutCondition(req, groupByField))
	} else {
		pipeline = append(pipeline, BuildGroupByWithCondition(req, groupByField))
	}
	secondMatchFilter, err := BuildSecondMatchFilter(req, isOverallPipeline)
	if err != nil {
		return nil, err
	}
	if len(secondMatchFilter) != 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"$and": secondMatchFilter}})
	} 
	return pipeline
}

func BuildFirstMatchFilter(req AttackStaticRequest, isIgnoreFilterValue bool) ([]bson.M, error) {
	var err error
	firstMatchFilters:=  []bson.M{{"company_id": bson.M{"$exists": true}}}
	atktypeFilter := BuildStaticAttackFromSearch(req.Search, false)
	if atktypeFilter != nil {
		firstMatchFilters = append(firstMatchFilters, atktypeFilter)
	}
	firstMatchFilters, err = BuildStaticAttackFilterRangeByField(firstMatchFilters, "last_attack_time", false, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	filterFromArray := BuildStaticAttackFilterFromArrayValue(req)
	if len(filterFromArray) != 0 {
		firstMatchFilters = append(firstMatchFilters, filterFromArray...)
	}
	if isIgnoreFilterValue {
		return firstMatchFilters, nil
	}
	arrayOrFilter := make([]bson.M, 0)
	orFilter := BuildStaticAttackFilterByValueInput("highest_bytes", req.BpsValue)
	if orFilter != nil {
		arrayOrFilter = append(arrayOrFilter, orFilter)
	}
	orFilter = BuildStaticAttackFilterByValueInput("highest_packets", req.BpsValue)
	if orFilter != nil {
		arrayOrFilter = append(arrayOrFilter, orFilter)
	}
	orFilter = BuildStaticAttackFilterByValueInput("duration", req.BpsValue)
	if orFilter != nil {
		arrayOrFilter = append(arrayOrFilter, orFilter)
	}
	if len(arrayOrFilter) != 0 {
		firstMatchFilters = append(firstMatchFilters, bson.M{"$or": arrayOrFilter})
	}
	return firstMatchFilters, nil
}

func BuildSecondMatchFilter(req AttackStaticRequest, isStaticOverall bool) ([]bson.M, error) {
	var err error
	var secondMatchFilters []bson.M
	if isStaticOverall {
		secondMatchFilters, err = BuildStaticAttackFilterRangeByField(secondMatchFilters, "total_ip", true, req.TotalIpFrom, req.TotalIpTo)
		if err != nil {
			return nil, err
		}
		secondMatchFilters, err = BuildStaticAttackFilterRangeByField(secondMatchFilters, "total_attack", false, req.TotalIpFrom, req.TotalIpTo)
		if err != nil {
			return nil, err
		}
		return secondMatchFilters, nil
	}
	secondMatchFilters, err = BuildStaticAttackFilterRangeByField(secondMatchFilters, "highest_bytes", false, req.HighestBytesFrom, req.HighestBytesTo)
	if err != nil {
		return nil, err
	}
	secondMatchFilters, err = BuildStaticAttackFilterRangeByField(secondMatchFilters, "highest_packets", false, req.HighestPacketsFrom, req.HighestPacketsTo)
	if err != nil {
		return nil, err
	}
	secondMatchFilters, err = BuildStaticAttackFilterRangeByField(secondMatchFilters, "duration", false, req.DurationFrom, req.DurationTo)
	if err != nil {
		return nil, err
	}
	return secondMatchFilters, nil
}

func BuildGroupByWithCondition(req AttackStaticRequest, groupByField string) bson.M {
	groupBy := bson.M{}
	groupBy["highest_bytes"] = BuildStaticAttackGroupCondtion("highest_bytes", req.BpsValue)
	groupBy["highest_packets"] = BuildStaticAttackGroupCondtion("highest_packets", req.PpsValue)
	groupBy["duration"] = BuildStaticAttackGroupCondtion("duration", req.DurationBaseline)
	switch groupByField {
	case "company_id":
		groupBy["_id"] = "company_id"
	case "segmentation_id":
		groupBy["_id"] = "segmentation_id"
	case "distributor_id":
		groupBy["_id"] = "distributor_id"
	default:
		groupBy["_id"] = "method_attack"
	}
	return groupBy
}
func BuildGroupByWithoutCondition(req AttackStaticRequest, groupByField string) bson.M {
	groupBy := bson.M{"$ips": bson.M{"$addToSet": "$ip"}}
	switch groupByField {
	case "company_id":
		groupBy["_id"] = "company_id"
		groupBy["distributor_ids"] = bson.M{"$addToSet": "$distributor_id"}
		groupBy["segmentation_ids"] = bson.M{"$addToSet": "$segmentation_id"}
		groupBy["method_attacks"] = bson.M{"$addToSet": "$method_attack"}
	case "segmentation_id":
		groupBy["_id"] = "segmentation_id"
		groupBy["distributor_ids"] = bson.M{"$addToSet": "$distributor_id"}
		groupBy["company_ids"] = bson.M{"$addToSet": "$company_id"}
		groupBy["method_attacks"] = bson.M{"$addToSet": "$method_attack"}
	case "distributor_id":
		groupBy["_id"] = "distributor_id"
		groupBy["segmentation_ids"] = bson.M{"$addToSet": "$segmentation_id"}
		groupBy["company_ids"] = bson.M{"$addToSet": "$company_id"}
		groupBy["method_attacks"] = bson.M{"$addToSet": "$method_attack"}
	default:
		groupBy["_id"] = "method_attack"
		groupBy["distributor_ids"] = bson.M{"$addToSet": "$distributor_id"}
		groupBy["segmentation_ids"] = bson.M{"$addToSet": "$segmentation_id"}
		groupBy["company_ids"] = bson.M{"$addToSet": "$company_id"}
	}
	return groupBy
}

func BuildStaticAttackFilterFromArrayValue(req AttackStaticRequest) []bson.M {
	var filters []bson.M
	if len(req.Companies) != 0 {
		filters = append(filters, bson.M{"$company_id": bson.M{"$in": req.Companies}})
	}
	if len(req.Distributors) != 0 {
		filters = append(filters, bson.M{"$distributor_id": bson.M{"$in": req.Distributors}})
	}
	if len(req.Segmentations) != 0 {
		filters = append(filters, bson.M{"$segmentation_id": bson.M{"$in": req.Segmentations}})
	}
	if len(req.MethodAttacks) != 0 {
		filters = append(filters, bson.M{"$attacktype": bson.M{"$in": req.MethodAttacks}})
	}
	return filters
}
