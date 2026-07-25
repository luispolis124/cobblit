package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"

	"github.com/luispolis124/cobblit/internal/players"
)

func main() {
	log.Println("[Cobblit Engine] Inicializando núcleos de simulação...")
	time.Sleep(1 * time.Second)

	srv := server.New()
	srv.CloseOnProgramEnd()

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	// Loop principal de aceitação de jogadores com o novo sistema
	go func() {
		for p := range srv.Accept() {
			// Registra o player no sistema modular
			players.RegistrarEntrada(p)

			// Opcional: Se quiser lidar com a saída do jogador em background
			go func(pl *player.Player) {
				// Gerenciamento de eventos de saída futuramente se necessário
			}(p)
		}
	}()

	log.Println("--- Cobblit Engine Iniciada com Sucesso ---")
	srv.Listen()

	// Gerenciamento de fechamento seguro via Ctrl + C
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
