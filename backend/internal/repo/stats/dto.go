package stats

var fieldInfo = map[string]struct {
	Coll    string // Коллекция-источник
	Path    string // Название поля
	IsArray bool   // Если массив, то будет использоваться представление $size
}{
	// ==== X-поля ====
	"user_id": {
		Coll: "users",
		Path: "_id",
	},
	"user_system_role": {
		Coll: "users",
		Path: "systemRole",
	},
	"user_active": {
		Coll: "users",
		Path: "active",
	},
	"user_created_at": {
		Coll: "users",
		Path: "createdAt",
	},
	"order_id": {
		Coll: "orders",
		Path: "_id",
	},
	"order_active": {
		Coll: "orders",
		Path: "active",
	},
	"order_freelancer_id": {
		Coll: "orders",
		Path: "freelancerId",
	},
	"order_client_id": {
		Coll: "orders",
		Path: "clientId",
	},
	"order_created_at": {
		Coll: "orders",
		Path: "createdAt",
	},

	// ==== Y-поля ====
	"count": { // $sum: 1, Path не используется
		Coll: "any",
		Path: "",
	},
	"user_balance": {
		Coll: "users",
		Path: "balance",
	},
	"user_client_rating": {
		Coll: "users",
		Path: "client.rating",
	},
	"user_freelancer_rating": {
		Coll: "users",
		Path: "freelancer.rating",
	},
	"order_completion_time": {
		Coll: "orders",
		Path: "completionTime",
	},
	"order_cost": {
		Coll: "orders",
		Path: "cost",
	},
	"order_responses_count": {
		Coll:    "orders",
		Path:    "responses",
		IsArray: true,
	},
}
