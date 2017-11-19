package init_handler

import (
	"github.com/jnovikov/hackforces/back/libs/storage"
	"github.com/jnovikov/hackforces/back/service_controller/flag_handler"
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
