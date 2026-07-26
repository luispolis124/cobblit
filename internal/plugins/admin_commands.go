package plugins

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// ListCommand define o comando /list para ver os jogadores online.
type ListCommand struct{}

func (ListCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	var names []string
	for entity := range tx.Players() {
		if p, ok := entity.(*player.Player); ok {
			names = append(names, p.Name())
		}
	}

	o.Printf("Jogadores online (%d): %s", len(names), strings.Join(names, ", "))
}

// RegisterAdminCommands registra o comando administrativo no motor.
func RegisterAdminCommands() {
	cmd.Register(cmd.New("list", "Exibe todos os jogadores online", []string{"online"}, ListCommand{}))
}
