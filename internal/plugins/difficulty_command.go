package plugins

import (
	"strconv"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// DifficultyCommand gerencia o comando /difficulty
type DifficultyCommand struct {
	Nivel cmd.Varargs `name:"nivel"`
}

func (c DifficultyCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para alterar a dificuldade.")
			return
		}
	}

	argStr := strings.TrimSpace(string(c.Nivel))
	if argStr == "" {
		o.Errorf("Uso correto: /difficulty <0-3 | peaceful, easy, normal, hard>")
		return
	}

	var diff int
	argLower := strings.ToLower(argStr)

	switch argLower {
	case "0", "peaceful", "p":
		diff = 0
	case "1", "easy", "e":
		diff = 1
	case "2", "normal", "n":
		diff = 2
	case "3", "hard", "h":
		diff = 3
	default:
		if val, err := strconv.Atoi(argLower); err == nil && val >= 0 && val <= 3 {
			diff = val
		} else {
			o.Errorf("§cNível de dificuldade inválido. Use 0 (Peaceful), 1 (Easy), 2 (Normal) ou 3 (Hard).")
			return
		}
	}

	// Atribui diretamente as variáveis globais de dificuldade do Dragonfly
	var d world.Difficulty
	switch diff {
	case 0:
		d = world.DifficultyPeaceful
	case 1:
		d = world.DifficultyEasy
	case 2:
		d = world.DifficultyNormal
	case 3:
		d = world.DifficultyHard
	}

	// Aplica a dificuldade no mundo
	tx.World().SetDifficulty(d)

	nomesDificuldade := []string{"Pacífico", "Fácil", "Normal", "Difícil"}
	nomeAtual := nomesDificuldade[diff]

	o.Printf("§a[Cobblit Engine] A dificuldade do mundo foi alterada para: §e%s (%d)", nomeAtual, diff)
}

// RegistrarDifficultyCommand registra o comando /difficulty no motor.
func RegistrarDifficultyCommand() {
	cmd.Register(cmd.New("difficulty", "Define a dificuldade do jogo", []string{"diff"}, DifficultyCommand{}))
}
