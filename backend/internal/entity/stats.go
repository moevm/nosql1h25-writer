package entity

type Point struct {
	X string  `bson:"x"`
	Y float64 `bson:"value"`
}

type Aggregation string

const (
	AggregationCount Aggregation = "count"
	AggregationSum   Aggregation = "sum"
	AggregationMax   Aggregation = "max"
	AggregationMin   Aggregation = "min"
	AggregationAvg   Aggregation = "avg"
)
