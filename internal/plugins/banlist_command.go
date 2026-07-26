package plugins

import (
	"sort"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// BanListCommand define o comando /banlist para consultar jogadores banidos.
type BanListCommand struct {
	Tipo string `name:"type" optional:"true"`
}

func (c BanListCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	CarregarBans()

	bansMutex.Lock()
	var nomes []string
	for nome := range banListMap {
		nomes = append(nomes, nome)
	}
	bansMutex.Unlock()

	sort.Strings(nomes)
	total := len(nomes)

	p, isPlayer := src.(*player.Player)

	tipoConsulta := strings.ToLower(c.Tipo)
	if tipoConsulta == "ips" {
		if isPlayer {
			p.Messagef("§aHá §e%d §aIPs banidos no momento.", total)
		} else {
			o.Printf("Há %d IPs banidos no momento.", total)
		}
	} else {
		if isPlayer {
			p.Messagef("§aHá §e%d §ajogadores banidos no momento.", total)
		} else {
			o.Printf("Há %d jogadores banidos no momento.", total)
		}
	}

	if total > 0 {
		listaStr := strings.Join(nomes, ", ")
		if isPlayer {
			p.Message("§7" + listaStr)
		} else {
			o.Printf("%s", listaStr)
		}
	} else {
		if isPlayer {
			p.Message("§7[Nenhum registro de banimento encontrado]")
		} else {
			o.Printf("[Nenhum registro de banimento encontrado]")
		}
	}
}

// RegistrarBanListCommand registra o comando /banlist no motor.
func RegistrarBanListCommand() {
	CarregarBans()
	cmd.Register(cmd.New("banlist", "Exibe a lista de jogadores banidos", []string{"banned"}, BanListCommand{}))
}
