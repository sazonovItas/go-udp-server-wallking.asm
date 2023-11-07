package server

import (
	"asm-game/server/game/player"
	"fmt"
	"sync"
)

type Session struct {
	// Sync for different thread
	sync.Mutex

	// Container for players
	CntPl   int
	Players map[string]*player.Player
}

func (s *Session) NewPlayer(addr string, data []byte) {
	pl, ok := s.Players[addr]
	ok = ok && !(pl == nil)
	if !ok {
		pl = player.New(addr, data)
		s.Lock()
		s.CntPl++
		s.Players[addr] = pl
		s.Unlock()
		fmt.Printf("New Player:\n%s\n", pl.String())
	} else {
		fmt.Printf("New name %s for player:\n%s\n", string(data), pl.String())
		s.Lock()
		pl.Info.Name = string(data)
		s.Unlock()
	}
}

func (s *Session) UpdatePlayer(addr string, data []byte) {
	pl, ok := s.Players[addr]
	ok = ok && !(pl == nil)
	if !ok {
		return
	}

	s.Lock()
	pl.Info.Update(data)
	s.Unlock()
}

func (s *Session) ExitPlayer(addr string, data []byte) {
	pl, ok := s.Players[addr]
	ok = ok && !(pl == nil)
	if !ok {
		return
	}

	s.Lock()
	s.CntPl--
	s.Players[addr] = nil
	s.Unlock()
}
