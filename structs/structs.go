package structs

type Request struct {
	Expression string `json:"expression"`
}
type ResponseOK struct {
	Id int `json:"id"`
}
type ResponseERROR struct {
	Error string `json:"error"`
}
type ResponseResult struct {
	Id     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}
type ResponseExpression struct {
	Expression ResponseResult `json:"expression"`
}
