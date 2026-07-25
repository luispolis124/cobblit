package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"

	"github.com/luispolis124/cobblit/internal/players"
	"github.com/luispolis124/cobblit/internal/plugins"
	"github.com/luispolis124/cobblit/internal/world"
)

func main() {
	log.Println("[Cobblit Engine] Inicializando núcleos de simulação...")
	time.Sleep(1 * time.Second)

	gameWorld := world.NewWorld("CobblitAlpha")
	_ = gameWorld

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	conf, err := server.DefaultConfig().Config(logger)
	if err != nil {
		log.Fatalf("[Cobblit Engine] Erro ao carregar config: %v", err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	// Registra os comandos personalizados do Cobblit Engine
	plugins.RegistrarComandos()

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	go func() {
		for pInstance := range srv.Accept() {
			players.RegistrarEntrada(pInstance)
			plugins.RegistrarBoasVindas(pInstance)

			go func(pl *player.Player) {
				gameWorld.StreamChunks(0, 0)
			}(pInstance)
		}
	}()

	go func() {
		log.Println("--- Cobblit Engine Iniciada com Sucesso ---")
		srv.Listen()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("[Cobblit Engine] Sinal recebido (%v). Desativando conexões...", sig)

	if err := srv.Close(); err != nil {
		log.Printf("[Cobblit Engine] Erro ao fechar: %v", err)
	} else {
		log.Println("[Cobblit Engine] Servidor encerrado com segurança.")
	}
}
