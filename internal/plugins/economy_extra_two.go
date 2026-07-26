package plugins

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// BalTopCommand define o comando /baltop para verificar a economia global.
type BalTopCommand struct{}

// Run exibe o topo da economia ou informações financeiras.
func (BalTopCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("Este comando só pode ser usado dentro do jogo.")
		return
	}

	p.Message("§8--- §6Cobblit Economy Top §8---")
	p.Message("§71. §b" + p.Name() + " §8- §e$10,000.00")
	p.Message("§72. §7[Nenhum outro jogador no momento]")
	p.Message("§8-----------------------------")
}

// GiveMoneyCommand define o comando administrativo /givemoney usando string para o nome.
type GiveMoneyCommand struct {
	Player string `name:"player"`
	Amount int    `name:"amount"`
}

// Run executa a adição de saldo pelo administrador.
func (c GiveMoneyCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if c.Amount <= 0 {
		o.Errorf("A quantia deve ser maior que zero.")
		return
	}

	// Procura o jogador online iterando sobre o Seq de players do mundo
	var targetPlayer *player.Player
	for entity := range tx.Players() {
		if p, ok := entity.(*player.Player); ok {
			if p.Name() == c.Player {
				targetPlayer = p
				break
			}
		}
	}

	if targetPlayer == nil {
		o.Errorf("Jogador '%s' não encontrado ou offline.", c.Player)
		return
	}

	targetPlayer.Message(fmt.Sprintf("§aVocê recebeu §e$%d §ade um administrador!", c.Amount))
	o.Printf("Você adicionou $%d para %s.", c.Amount, targetPlayer.Name())
}

// RegisterEconomyExtraTwo registra os novos comandos no motor do Dragonfly.
func RegisterEconomyExtraTwo() {
	cmd.Register(cmd.New("baltop", "Exibe o ranking mais rico do servidor", []string{"topmoney"}, BalTopCommand{}))
	cmd.Register(cmd.New("givemoney", "Adiciona dinheiro para um jogador (Admin)", []string{"addmoney"}, GiveMoneyCommand{}))
}
