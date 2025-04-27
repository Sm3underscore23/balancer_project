package service

type CheckerService interface {
	CheckerOnce()
	CheckerWithTicker()
}
