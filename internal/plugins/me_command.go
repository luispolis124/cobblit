package plugins

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
)

// MeCommand define o comando /me para emotes e ações no chat global.
type MeCommand struct {
	Acao cmd.Varargs `name:"action"`
}

func (c MeCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	acaoStr := strings.TrimSpace(string(c.Acao))
	if acaoStr == "" {
		o.Errorf("Uso correto: /me <ação>")
		return
	}

	var nomeRemetente string
	if p, ok := src.(*player.Player); ok {
		nomeRemetente = p.Name()
	} else {
		nomeRemetente = "Console"
	}	

	// Formata a mensagem de emote no estilo tradicional do Minecraft
	mensagemEmote := "§d* " + nomeRemetente + " " + acaoStr + " §r"

	// Envia para o chat global do servidor visível para todos os jogadores
	chat.Global.WriteString(mensagemEmote)
}

// RegistrarMeCommand registra o comando /me no motor.
func RegistrarMeCommand() {
	cmd.Register(cmd.New("me", "Executa uma ação no chat", []string{}, MeCommand{}))
}
