package plugins

import (
	"log"
	"sync"

	"github.com/df-mc/dragonfly/server/world"
)

var (
	mundosAtivos = make(map[string]*world.World)
	worldMutex   sync.Mutex
)

func CarregarOuCriarMundo(nome string) *world.World {
	worldMutex.Lock()
	defer worldMutex.Unlock()

	if m, existe := mundosAtivos[nome]; existe {
		return m
	}

	// Cria o mundo utilizando o construtor correto do Dragonfly atual
	novoMundo := world.New()
	mundosAtivos[nome] = novoMundo

	log.Printf("[Cobblit Worlds] Mundo '%s' inicializado com sucesso!", nome)
	return novoMundo
}
