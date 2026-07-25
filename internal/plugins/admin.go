package plugins

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// CobblitCommand define a estrutura do comando administrativo
type CobblitCommand struct {
	Sub string `optional:"true" name:"subcomando"`
}

// Run executa a lógica do comando quando digitado (incluindo o tx *world.Tx exigido pela API)
func (c CobblitCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("Este comando só pode ser usado por jogadores no jogo.")
		return
	}

	if c.Sub == "" || c.Sub == "help" {
		p.Message("§e=== COMANDOS DO COBBLIT ===")
		p.Message("§b/cobblit status §7- Mostra o status do motor")
		p.Message("§b/cobblit info   §7- Mostra informações da sessão")
		return
	}

	switch c.Sub {
	case "status":
		p.Message("§a[Cobblit Engine] Status do Sistema: §eONLINE §a(Alta Performance Ativa)")
		p.Message("§b[Cobblit Engine] Desenvolvido em Go com arquitetura concorrente.")
	case "info":
		p.Message("§e=== COBBLIT ENGINE INFO ===")
		p.Message(fmt.Sprintf("§7Jogador Conectado: §f%s", p.Name()))
		pos := p.Position()
		p.Message(fmt.Sprintf("§7Posição Atual: §fX: %.1f, Y: %.1f, Z: %.1f", pos.X(), pos.Y(), pos.Z()))
		p.Message("§7Versão do Motor: §fAlpha 1.0.0")
	default:
		p.Message(fmt.Sprintf("§c[Cobblit Admin] Subcomando desconhecido: %s. Use /cobblit help", c.Sub))
	}
}

// RegistrarComandos registra o comando no registrador global do Dragonfly
func RegistrarComandos() {
	cmd.Register(cmd.New("cobblit", "Gerencia o motor Cobblit Engine", []string{"cb"}, CobblitCommand{}))
}
