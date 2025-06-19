package stats

import "github.com/samber/lo"

type GraphIn struct {
	X       string
	Y       string
	AggType Aggregation
}

type Aggregation string

const (
	AggregationCount Aggregation = "count"
	AggregationSum   Aggregation = "sum"
	AggregationMax   Aggregation = "max"
	AggregationMin   Aggregation = "min"
	AggregationAvg   Aggregation = "avg"
)

var allowedToBeAsX = map[string]struct{}{
	"user_id":             {},
	"user_system_role":    {},
	"user_active":         {},
	"user_created_at":     {},
	"order_id":            {},
	"order_active":        {},
	"order_freelancer_id": {},
	"order_client_id":     {},
	"order_created_at":    {},
}

func CanBeX(x string) bool {
	_, ok := allowedToBeAsX[x]
	return ok
}

var allowedAggregationsForY = map[string][]Aggregation{
	"count":                  {AggregationCount},
	"user_reviews_count":     {AggregationAvg, AggregationMin, AggregationMax},
	"user_balance":           {AggregationAvg, AggregationMin, AggregationMax, AggregationSum},
	"user_client_rating":     {AggregationAvg, AggregationMin, AggregationMax},
	"user_freelancer_rating": {AggregationAvg, AggregationMin, AggregationMax},
	"order_completion_time":  {AggregationAvg, AggregationMin, AggregationMax, AggregationSum},
	"order_cost":             {AggregationAvg, AggregationMin, AggregationMax, AggregationSum},
}

func AllowedAggregationsForY(x string) ([]Aggregation, bool) {
	value, exists := allowedAggregationsForY[x]
	return value, exists
}

func ValidRequest(x, y string, agg Aggregation) bool {
	value, exists := AllowedAggregationsForY(y)
	return CanBeX(x) && exists && lo.Contains(value, agg)
}

type GraphOut struct {
	X string
	Y float64
}
