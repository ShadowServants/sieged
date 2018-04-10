package flaghandler

import (
	"hackforces/libs/round_storage"
	"hackforces/libs/statusstorage"
	"hackforces/libs/storage"
	"hackforces/service_controller/flag_handler/flagstorage"
)

type Factory struct {
	handler *FlagHandler
}

func NewFlagHandlerFactory() *Factory {
	ff := new(Factory)
	ff.handler = NewFlagHandler()
	return ff
}

func (ff *Factory) SetRoundStorage(st storage.Storage) *Factory {

	ff.handler.RoundSt = &round_storage.RoundStorage{st}
	return ff
}

func (ff *Factory) SetPointStorage(st storage.Storage) *Factory {
	ff.handler.Points = &PointsStorage{st}
	return ff
}

func (ff *Factory) SetFlagStorage(st storage.Storage) *Factory {
	ff.handler.Flags = flagstorage.NewFlagStorage(st)
	return ff
}

func (ff *Factory) SetStatusStorage(st storage.Storage) *Factory {
	ff.handler.StatusStorage = statusstorage.NewStatusStorage(st)
	return ff
}

func (ff *Factory) SetTeamFlagsSet(ks storage.KeySet) *Factory {
	ff.handler.TeamFlagsSet = ks
	return ff
}

func (ff *Factory) SetTeamNum(teamNum int) *Factory {
	ff.handler.TeamNum = teamNum
	return ff
}

func (ff *Factory) SetRoundDelta(delta int) *Factory {
	ff.handler.RoundDelta = delta
	return ff
}

func (ff *Factory) GetFlagHandler() *FlagHandler {
	return ff.handler
}
