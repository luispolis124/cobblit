package plugins

import (
	"strconv"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type TimeSetCommand struct {
	Tempo cmd.Varargs `name:"tempo"`
}

func (c TimeSetCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	// Restringe o comando apenas para operadores ou console
	if p, ok := src.(*player.Player); ok {
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para usar este comando.")
			return
		}
	}

	arg := strings.ToLower(strings.TrimSpace(string(c.Tempo)))
	if arg == "" {
		o.Errorf("Você precisa especificar um horário (ex: day, night, 0, 1000).")
		return
	}

	var novoTempo int
	var nomeTempo string

	switch arg {
	case "day", "dia":
		novoTempo = 1000
		nomeTempo = "Dia"
	case "night", "noite":
		novoTempo = 13000
		nomeTempo = "Noite"
	case "noon", "meio-dia":
		novoTempo = 6000
		nomeTempo = "Meio-Dia"
	case "midnight", "meia-noite":
		novoTempo = 18000
		nomeTempo = "Meia-Noite"
	default:
		val, err := strconv.Atoi(arg)
		if err != nil {
			o.Errorf("§cHorário inválido. Use: day, night, noon, midnight ou número de ticks.")
			return
		}
		novoTempo = val
		nomeTempo = arg
	}

	// Altera o tempo acessando o mundo através da transação
	tx.World().SetTime(novoTempo)
	o.Printf("§a[Cobblit Engine] O horário do mundo foi alterado para: %s (%d ticks)", nomeTempo, novoTempo)
}

func RegistrarComandosMundo() {
	cmd.Register(cmd.New("timeset", "Altera o horário do mundo", []string{"time"}, TimeSetCommand{}))
}
