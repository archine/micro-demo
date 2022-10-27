package vo

type OrderList struct {
	No     string  `json:"no"`
	Price  float64 `json:"price"`
	Status int     `json:"status"`
}

type OrderDetail struct {
	Username string   `json:"username"`
	Mobile   string   `json:"mobile"`
	Goods    []string `json:"goods"`
}
