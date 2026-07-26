package plugins

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/title"
)

// RegistrarBoasVindas envia uma mensagem e um título épico na tela quando o jogador entra
func RegistrarBoasVindas(p *player.Player) {
	nome := p.Name()
	
	// Envia um título em destaque no meio da tela
	// Parâmetros: Título principal, Subtítulo, FadeIn, Duração, FadeOut (em ticks)
	t := title.New("§b§lCOBBLIT ENGINE").WithSubtitle("§eBem-vindo(a), " + nome + "!")
	p.SendTitle(t)

	// Mensagem tradicional no chat com as instruções
	p.Message("§8[§bCobblit Engine§8] §fSeja bem-vindo(a) ao servidor, §e" + nome + "§f!")
	p.Message("§7Digite §b/cobblit help §7para ver os comandos disponíveis.")
}
