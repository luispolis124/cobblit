package plugins

import (
	"log"

	"github.com/df-mc/dragonfly/server/player"
)

// RegistrarBoasVindas gerencia a lógica quando um jogador entra
func RegistrarBoasVindas(p *player.Player) {
	log.Printf("[Plugin: BoasVindas] Jogador processado: %s", p.Name())
	p.Message("§a[Cobblit] §fMódulo de Boas-Vindas ativado com sucesso!")
}
