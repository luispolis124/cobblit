package plugins

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// CobblitCommand define a estrutura do comando administrativo com Varargs anti-crash
type CobblitCommand struct {
	Sub cmd.Varargs `optional:"true" name:"subcomando"`
}

// Run executa a lógica protegida e interativa do motor
func (c CobblitCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	// Verifica se é um jogador e se ele tem privilégios de operador
	if p, ok := src.(*player.Player); ok {
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para usar este comando administrativo.")
			return
		}
	}

	sub := strings.ToLower(strings.TrimSpace(string(c.Sub)))

	if sub == "" || sub == "help" {
		o.Printf("§e=== COMANDOS DO COBBLIT ENGINE ===")
		o.Printf("§b/cobblit status §7- Mostra o status do motor e sistemas")
		o.Printf("§b/cobblit info   §7- Mostra informações detalhadas da sessão")
		o.Printf("§b/cobblit reload §7- Recarrega as configurações de OPs e Bans")
		return
	}

	switch sub {
	case "status":
		o.Printf("§a[Cobblit Engine] Status do Sistema: §eONLINE §a(Alta Performance Ativa)")
		o.Printf("§b[Cobblit Engine] Desenvolvido em Go com arquitetura concorrente.")
	case "info":
		o.Printf("§e=== COBBLIT ENGINE INFO ===")
		if p, ok := src.(*player.Player); ok {
			o.Printf("§7Operador Conectado: §f%s", p.Name())
			pos := p.Position()
			o.Printf("§7Posição Atual: §fX: %.1f, Y: %.1f, Z: %.1f", pos.X(), pos.Y(), pos.Z())
		} else {
			o.Printf("§7Executado via: §fConsole do Servidor")
		}
		o.Printf("§7Versão do Motor: §fAlpha 1.0.0 (Dragonfly v0.11)")
	case "reload":
		// Recarrega os arquivos JSON em tempo de execução
		CarregarOps()
		CarregarBans()
		o.Printf("§a[Cobblit Engine] Arquivos 'ops.json' e 'bans.json' recarregados com sucesso!")
	default:
		o.Errorf("§c[Cobblit Admin] Subcomando desconhecido: %s. Use /cobblit help", sub)
	}
}

// RegistrarComandos registra o comando principal no registrador global
func RegistrarComandos() {
	cmd.Register(cmd.New("cobblit", "Gerencia o motor Cobblit Engine", []string{"cb"}, CobblitCommand{}))
}
