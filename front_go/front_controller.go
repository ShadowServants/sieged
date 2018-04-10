package front_go

import (
	"time"
	"log"
	"hackforces/libs/round_storage"
)

type FrontController struct {
	RoundStorage *round_storage.RoundStorage
	RoundTime int
	CurrentRound int
	GameEnabled bool
	ServiceMap map[string]string
	//TeamsList map[int]
	ticker *time.Ticker
}



func (fr *FrontController) RoundJobsDispatch() {
	for t := range fr.ticker.C {
		fr.CurrentRound++
		log.Println("New round started at ",t)
		fr.RoundJobsHandler()
	}
}

func (fr *FrontController) RoundJobsHandler() {

}

func (fr *FrontController) StartsRounds() {
	fr.ticker = time.NewTicker(time.Millisecond * time.Duration(fr.RoundTime))
}

func (fr *FrontController) PauseRounds() {
	fr.ticker.Stop()
}


