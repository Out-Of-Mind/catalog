package structures

import (
	"time"
)

type RequestDataJSON struct {
	GroupName    string `json:"group_name"`
	CategoryName string `json:"category_name"`
	ItemName     string `json:"item_name"`
	JWT          string `json:"jwt"`
	RJWT         string `json:"rjwt"`
}

type RequestJSON struct {
	Action string          `json:"action"`
	JWT    string          `json:"jwt_token"`
	Data   RequestDataJSON `json:"data"`
}

type ResponseDataJSON struct {
	JWT  string `json:"jwt"`
	RJWT string `json:"rjwt"`

	InviteLink string `json:"invite_link"`
	SelectLink string `json:"select_link"`

	Error string `json:"error"`
}

type ResponseJSON struct {
	Succes bool             `json:"success"`
	Data   ResponseDataJSON `json:"data"`
}

type IndexItems struct {
	JWT, RJWT string
	Items     []IndexData
}

type IndexData struct {
	ID           int64
	CategoryName string
	CategoryID   string
	ItemNames    []string
}

type ByID []IndexData

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

type ErrorTemplate struct {
	Show bool
	Text string
}

type LoginData struct {
	CSRFToken string
	Error     ErrorTemplate
}

type RegisterData LoginData

type JWT struct {
	Header    JWTHeader
	Payload   JWTPayload
	Signature JWTSignature
}

type JWTHeader struct {
	Alg  string `json:"alg"`
	Type string `json:"type"`
}

type JWTPayload struct {
	Exp   time.Time `json:"exp"`
	Value int       `json:"value"`
}

type JWTSignature struct {
	Hash string
}

type DashboardData struct {
	JWT, RJWT, UserName string
	Groups              []DashboardGroup
	OwnedGroups         []DashboardOwnedGroup
}

type DashboardGroup struct {
	GroupName        string
	GroupWelcomeLink string
}

type DashboardOwnedGroup struct {
	GroupName        string
	GroupWelcomeLink string
}
