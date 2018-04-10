package init_handler

import (
	"hackforces/service_controller/flag_handler"
	"hackforces/libs/storage"
)

type InitHandlerFactory struct {
	initHandler InitHandler
}

func (ihf *InitHandlerFactory) SetTeamStorage(st storage.Storage) *InitHandlerFactory {
	ihf.initHandler.TeamStorage = st
	return ihf
}

func (ihf *InitHandlerFactory) SetPointStorage(st storage.Storage) *InitHandlerFactory {
	ihf.initHandler.Ps = &flaghandler.PointsStorage{St: st}
	return ihf
}
