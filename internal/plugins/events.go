package plugins

import (
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/google/uuid"
)

func RegistrarOuvintesEventos() {
	chat.Global.Subscribe(ChatListener{})
}

type ChatListener struct{}

func (ChatListener) Message(v ...any) {
	// Processa as mensagens globais do chat
}

func (ChatListener) UUID() uuid.UUID {
	return uuid.New()
}
