package team

type Status struct {
	TeamId        int    `json:"team_id"`
	StatusMessage string `json:"status_message"`
	Status        string `json:"status"`
	Points        Score  `json:"points"`
}


