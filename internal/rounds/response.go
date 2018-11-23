package rounds

import (
	"sieged/internal/team"
)

type Response struct {
	Responses []team.Status `json:"responses"`
}
