package session

import (
	"asm-game/server/internal/game/player"
	"fmt"
	"log"
	"sync"
	"time"

	convertBytes "asm-game/server/internal/convertbytes"
)

type Session struct {
	sync.Mutex

	// Container for players
	CntPl       int
	SessionTime int32
	Players     map[string]*player.Player
}

func (s *Session) NewPlayer(addr string, data []byte) {
	s.Lock()
	defer s.Unlock()
	pl, ok := s.Players[addr]
	if !ok {
		pl = player.New(addr, data)
		s.CntPl++
		s.Players[addr] = pl
		pl.SessionUpTime = s.SessionTime
		log.Printf("New Player:\n%s\n", pl.String())
	} else {
		log.Printf("New name %s for player:\n%s\n", string(data), pl.String())
		pl.Info.Name = string(data)
		pl.Uptime = time.Now()
		pl.SessionUpTime = s.SessionTime
	}
}

func (s *Session) UpdatePlayer(addr string, data []byte) {
	pl, ok := s.Players[addr]
	if !ok {
		return
	}

	uptime, ok := convertBytes.ByteSliceToT[int32](data[:4])
	s.Lock()
	if ok {
		pl.Info.Update(data[4:])
		pl.Uptime = time.Now()
		pl.SessionUpTime = uptime
		pl.Updated = true
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
		if v.Updated {
			str += fmt.Sprintf("%s\n", v.String())
		}
	}

	return str
}
