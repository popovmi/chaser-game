package main

func (srv *server) initTickers() {
	go func() {
		for {
			select {
			case <-srv.fpsTicker.C:
				srv.broadcastState()
			case <-srv.quit:
				srv.fpsTicker.Stop()
				return
			}
		}
	}()
}
