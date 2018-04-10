package round_handler




type ChainChecker interface {
	SetNext(next ChainChecker) ChainChecker
	Process(teamId int,round int) TeamResponse
}

