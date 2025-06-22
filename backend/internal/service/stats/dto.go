package stats

import (
	"github.com/samber/lo"

	"github.com/moevm/nosql1h25-writer/backend/internal/entity"
)

type GraphIn struct {
	X       string
	Y       string
	AggType entity.Aggregation
}

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

func canBeX(x string) bool {
	_, ok := allowedToBeAsX[x]
	return ok
}

var allowedAggregationsForY = map[string][]entity.Aggregation{
	"count":                  {entity.AggregationCount},
	"user_balance":           {entity.AggregationAvg, entity.AggregationMin, entity.AggregationMax, entity.AggregationSum},
	"user_client_rating":     {entity.AggregationAvg, entity.AggregationMin, entity.AggregationMax},
	"user_freelancer_rating": {entity.AggregationAvg, entity.AggregationMin, entity.AggregationMax},
	"order_completion_time":  {entity.AggregationAvg, entity.AggregationMin, entity.AggregationMax, entity.AggregationSum},
	"order_cost":             {entity.AggregationAvg, entity.AggregationMin, entity.AggregationMax, entity.AggregationSum},
	"order_responses_count":  {entity.AggregationAvg, entity.AggregationMin, entity.AggregationMax, entity.AggregationSum},
}

func validRequest(x, y string, agg entity.Aggregation) bool {
	value, exists := allowedAggregationsForY[y]
	return canBeX(x) && exists && lo.Contains(value, agg)
}

type GraphOut struct {
	X string
	Y float64
}
