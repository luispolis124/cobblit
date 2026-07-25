package plugins

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/title"
	"time"
)

const (
	TitleHeader = "§l§bCOBBLIT ENGINE"
	TitleSub    = "§7Bem-vindo ao servidor alpha!"
)

func RegistrarBoasVindas(p *player.Player) {
	nome := p.Name()
	
	// Envia mensagem no chat
	p.Message(fmt.Sprintf("§a[Cobblit] Olá, §e%s§a! Aproveite a experiência de alta performance.", nome))
	
	// Configura o título passando o tempo correto de duração única
	t := title.New(TitleHeader).WithSubtitle(TitleSub).WithDuration(5 * time.Second)
	p.SendTitle(t)
}
