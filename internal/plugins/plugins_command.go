package plugins

import (
	"fmt"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

// Lista global de subsistemas nativos do Cobblit Engine
var modulosNativos = []string{
	"CobblitCore",
	"EconomySystem",
	"WorldManager",
	"ProtectionEngine",
	"OperatorSystem",
}

func RegistrarModulo(nome string) {
	for _, m := range modulosNativos {
		if m == nome {
			return
		}
	}
	modulosNativos = append(modulosNativos, nome)
}

type PluginsCommand struct{}

func (c PluginsCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	var listaFormatada []string

	// Adiciona nativos em verde
	for _, mod := range modulosNativos {
		listaFormatada = append(listaFormatada, fmt.Sprintf("§a%s", mod))
	}

	// Adiciona os plugins externos .so (Verde se ativo, Vermelho se falhou ao carregar, estilo PocketMine)
	for _, ext := range PluginsExternos {
		if ext.Ativo {
			listaFormatada = append(listaFormatada, fmt.Sprintf("§a%s", ext.Nome))
		} else {
			listaFormatada = append(listaFormatada, fmt.Sprintf("§c%s (Erro)", ext.Nome))
		}
	}

	total := len(listaFormatada)
	o.Printf("§7Plugins (%d): §r%s", total, strings.Join(listaFormatada, "§7, §r"))
}

func RegistrarPluginsCommand() {
	cmd.Register(cmd.New("plugins", "Lista os plugins/módulos rodando no servidor", []string{"pl"}, PluginsCommand{}))
}
