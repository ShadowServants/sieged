package round_handler

import (
	"hackforces/libs/storage"
	"hackforces/service_controller/flag_handler"
	"hackforces/libs/statusstorage"
	"hackforces/libs/round_storage"
)

type HandlerFactory struct {
	roundH *RoundHandler
}

func NewHandlerFactory() *HandlerFactory{
	rhf := new(HandlerFactory)
	rhf.roundH = NewRoundHandler()
	return rhf
}

func (roundF *HandlerFactory) SetIpStorage(st storage.Storage) *HandlerFactory{
	roundF.roundH.IpStorage = st
	return roundF
}

func (roundF *HandlerFactory) SetTeamStorage(st storage.Storage) *HandlerFactory {
	roundF.roundH.TeamStorage = st
	return roundF
}

func (roundF *HandlerFactory) SetPointStorage(st storage.Storage) *HandlerFactory {
	roundF.roundH.Points = &flaghandler.PointsStorage{St: st}
	return roundF
}

func (roundF *HandlerFactory) SetStatusStorage(st storage.Storage) *HandlerFactory {
	roundF.roundH.St = statusstorage.NewStatusStorage(st)
	return roundF
}

func (roundF *HandlerFactory) SetRoundStorage(st storage.Storage) *HandlerFactory {
	roundF.roundH.Rounds = &round_storage.RoundStorage{St: st}
	return roundF
}

func (roundF *HandlerFactory) GetHandler() *RoundHandler{
	return roundF.roundH
}