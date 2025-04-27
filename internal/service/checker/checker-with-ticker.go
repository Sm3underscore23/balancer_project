package checker

import (
	"time"
)

func (p *poolChecker) CheckerWithTicker() {
	t := time.NewTicker(time.Second * 3)
	done := make(chan struct{}, 1)
	for {
		select {
		case <-t.C:
			p.CheckerOnce()
			done <- struct{}{}
		}
		<- done
	}
}
