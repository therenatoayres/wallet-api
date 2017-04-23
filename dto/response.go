package dto

//Response dto to send as http response paylod
type Response struct {
	Rates []Rate     `json:"rates"`
	Fails []Currency `json:"fails"`
}
