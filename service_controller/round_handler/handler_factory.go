package round_handler

import (
	"hackforces/libs/storage"
	"hackforces/service_controller/flag_handler"
	"hackforces/libs/statusstorage"
)

type HandlerFactory struct {
	round_h *RoundHandler
}

func NewHandlerFactory() *HandlerFactory{
	rhf := new(HandlerFactory)
	rhf.round_h = NewRoundHandler()
	return rhf
}

func (round_f *HandlerFactory) SetIpStorage(st storage.Storage) *HandlerFactory{
	round_f.round_h.IpStorage = st
	return round_f
}

func (round_f *HandlerFactory) SetTeamStorage(st storage.Storage) *HandlerFactory {
	round_f.round_h.TeamStorage = st
	return round_f
}

func (round_f *HandlerFactory) SetPointStorage(st storage.Storage) *HandlerFactory {
	round_f.round_h.Points = &flaghandler.PointsStorage{st}
	return round_f
}

func (round_f *HandlerFactory) SetStatusStorage(st storage.Storage) *HandlerFactory {
	round_f.round_h.St = statusstorage.NewStatusStorage(st)
	return round_f
}

func (round_f *HandlerFactory) SetRoundStorage(st storage.Storage) *HandlerFactory {
	round_f.round_h.Rounds = &flaghandler.RoundStorage{st}
	return round_f
}

func (round_f *HandlerFactory) GetHandler() *RoundHandler{
	return round_f.round_h
}