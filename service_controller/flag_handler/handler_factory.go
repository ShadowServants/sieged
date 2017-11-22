package flaghandler

import (
	"hackforces/libs/storage"
	"hackforces/service_controller/flag_handler/flagstorage"
	"hackforces/libs/statusstorage"
)

type FlagHandlerFactory struct {
	flaghandler *FlagHandler
}

func NewFlagHandlerFactory() *FlagHandlerFactory {
	ff := new(FlagHandlerFactory)
	ff.flaghandler = NewFlagHandler()
	return ff
}

func (ff *FlagHandlerFactory) SetRoundStorage(st storage.Storage) *FlagHandlerFactory {

	ff.flaghandler.RoundSt = &RoundStorage{st}
	return ff
}

func (ff *FlagHandlerFactory) SetPointStorage(st storage.Storage) *FlagHandlerFactory {
	ff.flaghandler.Points = &PointsStorage{st}
	return ff
}

func (ff *FlagHandlerFactory) SetFlagStorage(st storage.Storage) *FlagHandlerFactory {
	ff.flaghandler.Flags = flagstorage.NewFlagStorage(st)
	return ff
}

func (ff *FlagHandlerFactory) SetStatusStorage(st storage.Storage) *FlagHandlerFactory {
	ff.flaghandler.StatusStorage = statusstorage.NewStatusStorage(st)
	return ff
}

func (ff *FlagHandlerFactory) SetTeamFlagsSet(ks storage.KeySet) *FlagHandlerFactory {
	ff.flaghandler.TeamFlagsSet = ks
	return ff
}

func (ff *FlagHandlerFactory) GetFlagHandler() *FlagHandler {
	return ff.flaghandler
}

