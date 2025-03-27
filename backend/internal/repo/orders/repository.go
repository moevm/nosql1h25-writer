package orders

import "github.com/sv-tools/mongoifc"

type repository struct {
	ordersColl mongoifc.Collection
}

func New(ordersColl mongoifc.Collection) Repo {
	return &repository{ordersColl: ordersColl}
}
