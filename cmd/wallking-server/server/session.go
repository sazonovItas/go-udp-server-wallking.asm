package server

import (
	"asm-game/server/internal/game/player"
	"fmt"
	"log"
	"sync"
	"time"

	convertBytes "asm-game/server/internal/convertbytes"
)

type Session struct {
	// Sync for different thread
	sync.Mutex

	// Container for players
	CntPl       int
	SessionTime int32
	Players     map[string]*player.Player
}

func (s *Session) NewPlayer(addr string, data []byte) {
	pl, ok := s.Players[addr]
	ok = ok && !(pl == nil)
	if !ok {
		pl = player.New(addr, data)
		s.Lock()
		s.CntPl++
		s.Players[addr] = pl
		pl.SessionUpTime = s.SessionTime
		s.Unlock()
		fmt.Printf("New Player:\n%s\n", pl.String())
	} else {
		fmt.Printf("New name %s for player:\n%s\n", string(data), pl.String())
		s.Lock()
		pl.Info.Name = string(data)
		pl.Uptime = time.Now()
		pl.SessionUpTime = s.SessionTime
		s.Unlock()
	}
}

func (s *Session) UpdatePlayer(addr string, data []byte) {
	pl, ok := s.Players[addr]
	if !ok {
		return
	}

	s.Lock()
	_, ok = convertBytes.ByteSliceToT[int32](data[:4])
	// && pl.SessionUpTime <= uptime
	if ok {
		pl.Info.Update(data[4:])
		pl.Uptime = time.Now()
	}
	s.Unlock()
	log.Printf("%s", s.String())
}

func (s *Session) ExitPlayer(addr string) {
	_, ok := s.Players[addr]
	if !ok {
		return
	}

	s.Lock()
	s.CntPl--
	delete(s.Players, addr)
	s.Unlock()
}

func (s *Session) String() string {
	s.Lock()
	defer s.Unlock()
	var str string = fmt.Sprintf("Cnt players %d\n", s.CntPl)

	for _, v := range s.Players {
		if v == nil {
			continue
		}

		str += fmt.Sprintf("%s\n", v.String())
	}

	return str
}
