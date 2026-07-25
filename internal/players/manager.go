package players

import (
	"log"

	"github.com/df-mc/dragonfly/server/player"
)

// Session armazena dados personalizados de cada entidade/jogador no Cobblit
type Session struct {
	Player *player.Player
	Admin  bool
}

// Gerenciador central de jogadores ativos
var ActivePlayers = make(map[string]*Session)

// RegistrarEntrada processa o jogador assim que ele entra no servidor
func RegistrarEntrada(p *player.Player) {
	name := p.Name()
	log.Printf("[Cobblit Players] Jogador conectado: %s (UUID: %s)", name, p.UUID().String())

	// Adiciona à lista de sessões ativas
	ActivePlayers[name] = &Session{
		Player: p,
		Admin:  false, // Aqui você pode definir regras de cargo futuramente
	}

	// Mensagem de boas-vindas personalizada do motor
	p.Message("§b[Cobblit Engine] §fSessão inicializada com sucesso. Bom jogo!")
}

// RegistrarSaida limpa os dados quando o jogador sai
func RegistrarSaida(p *player.Player) {
	name := p.Name()
	log.Printf("[Cobblit Players] Jogador desconectado: %s", name)
	delete(ActivePlayers, name)
}
