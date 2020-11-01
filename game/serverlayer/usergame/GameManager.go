package game

import (
	gamedomain "github.com/app/game/applayer/game/domain"
	usergamedomain "github.com/app/game/serverlayer/usergame/domain"
	"github.com/patrickmn/go-cache"
	"time"
)

type GameManager struct {
	cacheByGameID *cache.Cache
	cacheByUserID *cache.Cache
}

func InitManager(cacheRate time.Duration) usergamedomain.IGameManager {
	cacheByGameID := cache.New(5*time.Minute, cacheRate)
	cacheByUserID := cache.New(5*time.Minute, cacheRate)
	return &GameManager{
		cacheByGameID: cacheByGameID,
		cacheByUserID: cacheByUserID,
	}
}

func (g *GameManager) AddGame(game gamedomain.IGameEnvironment) {
	gid := game.GetGameID()
	uid := game.GetBelongingUserID()
	_, exists := g.cacheByGameID.Get(gid)
	if exists {
		return
	}

	g.cacheByGameID.Add(gid, game, cache.DefaultExpiration)

	_cur, ok := g.cacheByUserID.Get(uid)
	if ok {
		cur := _cur.([]gamedomain.IGameEnvironment)
		cur = append(cur, game)
		g.cacheByUserID.Set(uid, cur, cache.DefaultExpiration)

		return
	}

	toAdd := []gamedomain.IGameEnvironment{
		game,
	}

	g.cacheByUserID.Add(uid, toAdd, cache.DefaultExpiration)
}

func (g *GameManager) GetGameByGameID(id string) (gamedomain.IGameEnvironment, bool) {
	r, b := g.cacheByGameID.Get(id)
	if b {
		return r.(gamedomain.IGameEnvironment), b
	}
	return nil, b
}

func (g *GameManager) GetGamesByUserID(id string) (gamedomain.IGameEnvironments, bool) {
	r, b := g.cacheByUserID.Get(id)
	if b {
		return r.(gamedomain.IGameEnvironments), b
	}
	return nil, b
}

func (g *GameManager) RemoveGameByGameID(id string) {
	//defer fmt.Println(g.cacheByUserID.Items())
	//defer fmt.Println(g.cacheByGameID.Items())

	_game, ok := g.cacheByGameID.Get(id)
	if !ok {
		return
	}
	game := _game.(gamedomain.IGameEnvironment)
	uid := game.GetBelongingUserID()

	g.cacheByGameID.Delete(id)
	_gs, ok := g.cacheByUserID.Get(uid)
	if !ok {
		return
	}

	gs := _gs.([]gamedomain.IGameEnvironment)

	if len(gs) <= 1 {
		g.cacheByUserID.Delete(uid)
		return
	}
	for in := 0; in < len(gs); in++ {
		if gs[in].GetGameID() == id {
			gs = append(gs[:in], gs[in+1:]...)
			in--
		}
	}

	g.cacheByUserID.Set(uid, gs, cache.DefaultExpiration)

}

func (g *GameManager) RemoveGamesByUserID(id string) {
	_gs, ok := g.cacheByUserID.Get(id)
	if !ok {
		return
	}
	gs := _gs.([]gamedomain.IGameEnvironment)
	for _, v := range gs {
		g.cacheByGameID.Delete(v.GetGameID())
	}

	g.cacheByUserID.Delete(id)
}
