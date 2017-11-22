package init_handler

import (
	"hackforces/service_controller/flag_handler"
	"hackforces/libs/storage"
)

type InitHandlerFactory struct{
	init_handler InitHandler
}

func (ihf *InitHandlerFactory) SetTeamStorage(st storage.Storage) *InitHandlerFactory {
	ihf.init_handler.TeamStorage = st
	return ihf
}

func (ihf *InitHandlerFactory) SetPointStorage(st storage.Storage) *InitHandlerFactory {
	ihf.init_handler.Ps = &flaghandler.PointsStorage{st};
	return ihf
}

func NewInitHandlerFactory() *InitHandlerFactory{
	ihf := new(InitHandlerFactory)
	return ihf
}
