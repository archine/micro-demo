package vo

type OrderDetail struct {
	Username string   `json:"username"`
	Mobile   string   `json:"mobile"`
	Goods    []string `json:"goods"`
}
