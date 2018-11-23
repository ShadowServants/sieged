package flaghandler

import (
	"sieged/internal/flags"
	"sieged/internal/rounds"
	"sieged/internal/team/score"
	"sieged/internal/team/status"

	"sieged/pkg/storage"
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

	ff.handler.RoundSt = &rounds.Storage{St: st}
	return ff
}

func (ff *Factory) SetPointStorage(st storage.Storage) *Factory {
	ff.handler.Points = &score.Storage{St: st}
	return ff
}

func (ff *Factory) SetFlagStorage(st storage.Storage) *Factory {
	ff.handler.Flags = flags.NewStorage(st)
	return ff
}

func (ff *Factory) SetStatusStorage(st storage.Storage) *Factory {
	ff.handler.StatusStorage = status.NewStorage(st)
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
