package main

func (srv *server) initTickers() {
	go func() {
		for {
			select {
			case <-srv.rateTicker.C:
				srv.game.Tick()
				srv.broadcastState()
			case <-srv.quit:
				srv.rateTicker.Stop()
				return
			}
		}
	}()
}
