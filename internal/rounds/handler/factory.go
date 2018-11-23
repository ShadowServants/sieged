package handler

import (
	"sieged/internal/rounds"
	"sieged/internal/team/score"
	"sieged/internal/team/status"
	"sieged/pkg/storage"
)

type Factory struct {
	roundH *rounds.Handler
}

func NewFactory() *Factory {
	rhf := new(Factory)
	rhf.roundH = rounds.NewHandler()
	return rhf
}

func (roundF *Factory) SetIpStorage(st storage.Storage) *Factory {
	roundF.roundH.IpStorage = st
	return roundF
}

func (roundF *Factory) SetTeamStorage(st storage.Storage) *Factory {
	roundF.roundH.TeamStorage = st
	return roundF
}

func (roundF *Factory) SetPointStorage(st storage.Storage) *Factory {
	roundF.roundH.ScoreStorage = &score.Storage{St: st}
	return roundF
}

func (roundF *Factory) SetStatusStorage(st storage.Storage) *Factory {
	roundF.roundH.St = status.NewStorage(st)
	return roundF
}

func (roundF *Factory) SetRoundStorage(st storage.Storage) *Factory {
	roundF.roundH.Rounds = &rounds.Storage{St: st}
	return roundF
}

func (roundF *Factory) GetHandler() *rounds.Handler {
	return roundF.roundH
}
