package plugins

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// ClearCommand define o comando /clear para esvaziar o inventário.
type ClearCommand struct {
	Alvo cmd.Varargs `name:"jogador" optional:"true"`
}

func (c ClearCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, isPlayer := src.(*player.Player)
	if !isPlayer {
		o.Errorf("Este comando só pode ser usado dentro do jogo.")
		return
	}

	alvoStr := strings.TrimSpace(string(c.Alvo))
	var alvoPlayer *player.Player = p

	if alvoStr != "" {
		// Verifica se tem permissão de OP para limpar o inventário de outros
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para limpar o inventário de outros jogadores.")
			return
		}

		encontrado := false
		for entity := range tx.Players() {
			if targetP, ok := entity.(*player.Player); ok {
				if strings.EqualFold(targetP.Name(), alvoStr) {
					alvoPlayer = targetP
					encontrado = true
					break
				}
			}
		}

		if !encontrado {
			p.Messagef("§cJogador '%s' não encontrado ou offline.", alvoStr)
			return
		}
	}

	// Limpa o inventário principal do jogador
	inv := alvoPlayer.Inventory()
	inv.Clear()

	// Opcional: feedback para os envolvidos
	if alvoPlayer.Name() == p.Name() {
		p.Message("§a[Cobblit Engine] Seu inventário foi limpo com sucesso.")
	} else {
		p.Messagef("§a[Cobblit Engine] O inventário de §e%s §afoi limpo.", alvoPlayer.Name())
		alvoPlayer.Message("§e[Cobblit Engine] Seu inventário foi limpo por um administrador.")
	}
}

// RegistrarClearCommand registra o comando /clear no motor.
func RegistrarClearCommand() {
	cmd.Register(cmd.New("clear", "Limpa o inventário do jogador", []string{"limpar"}, ClearCommand{}))
}
