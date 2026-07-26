package plugins

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

func RegistrarComandosExtras() {
	cmd.Register(cmd.New("ping", "Mostra sua latência atual no servidor", []string{}, PingCommand{}))
	cmd.Register(cmd.New("suicide", "Retorna ao spawn sacrificando seu personagem", []string{"kill"}, SuicideCommand{}))
	cmd.Register(cmd.New("gamemode", "Altera o modo de jogo", []string{"gm"}, GamemodeCommand{}))
}

// --- /ping ---
type PingCommand struct{}

func (PingCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("Apenas jogadores podem usar este comando.")
		return
	}

	latencia := p.Latency()
	o.Printf("§b[Cobblit] Sua latência atual (ping) é de: §e%v", latencia)
}

// --- /suicide ---
type SuicideCommand struct{}

func (SuicideCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("Apenas jogadores podem usar este comando.")
		return
	}

	p.Teleport(spawnPos)
	o.Printf("§c[Cobblit] Você utilizou o comando de resgate e retornou ao spawn!")
}

// --- /gamemode [modo] ---
type GamemodeCommand struct {
	Modo string `cmd:"modo"`
}

func (c GamemodeCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("Apenas jogadores podem usar este comando.")
		return
	}

	var novoGm world.GameMode
	modoStr := strings.ToLower(c.Modo)

	switch modoStr {
	case "0", "survival", "sobrevivencia":
		novoGm = world.GameModeSurvival
	case "1", "creative", "criativo":
		novoGm = world.GameModeCreative
	case "2", "adventure", "aventura":
		novoGm = world.GameModeAdventure
	case "3", "spectator", "espectador":
		novoGm = world.GameModeSpectator
	default:
		o.Errorf("Modo de jogo inválido! Use: survival, creative, adventure ou spectator.")
		return
	}

	p.SetGameMode(novoGm)
	o.Printf("§a[Cobblit] Seu modo de jogo foi alterado com sucesso!")
}
