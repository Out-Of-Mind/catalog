package structures

type RequestDataJSON struct {
	CategoryName string `json:"category_name"`
	ItemName string `json:"item_name"`
	RJWT string `json:"rjwt"`
}

type RequestJSON struct {
	Action string `json:"action"`
	JWT string `json:"jwt_token"`
	Data RequestDataJSON `json:"data"`
}

type ResponseDataJSON struct {
	JWT string `json:"jwt"`
	RJWT string `json:"rjwt"`

	Error string `json:"error"`
}

type ResponseJSON struct {
	Succes bool `json:"success"`
	Data ResponseDataJSON`json:"data"`
}

type IndexItems struct {
	Items []IndexData
}

type IndexData struct {
	CategoryName string
	ItemNames []string
}