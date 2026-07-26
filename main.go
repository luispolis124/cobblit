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
)

func main() {
	log.Println("[Cobblit Engine] Inicializando núcleos de simulação...")
	time.Sleep(1 * time.Second)

	cfg := plugins.CarregarOuCriarConfig()

	// Carrega os plugins externos compatíveis com c-shared no Android
	plugins.CarregarPluginsSo()

	_ = plugins.CarregarOuCriarMundo("world")
	_ = plugins.CarregarOuCriarMundo("nether")

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	conf, err := server.DefaultConfig().Config(logger)
	if err != nil {
		log.Fatalf("[Cobblit Engine] Erro ao carregar config base: %v", err)
	}

	conf.Name = cfg.Motd
	conf.MaxPlayers = cfg.MaxPlayers

	srv := conf.New()
	srv.CloseOnProgramEnd()

	plugins.RegistrarComandos()
	plugins.RegistrarComandosOperator()
	plugins.RegistrarComandosModeracao()
	plugins.RegistrarComandosMundo()
	plugins.RegistrarComandosEconomia()
	plugins.RegistrarComandosMovimento()
	plugins.RegisterAdminCommands()  
	plugins.RegistrarOuvintesEventos()
	plugins.RegisterGameRuleCommand() 
	plugins.RegistrarComandosStatus() 
	plugins.RegistrarBanListCommand() 
	plugins.RegistrarMeCommand() 
	plugins.RegistrarClearCommand() 
	plugins.RegistrarBanIpCommand() 
	plugins.RegistrarDifficultyCommand() 
	plugins.RegistrarWhitelistCommand() 
	plugins.RegistrarPluginsCommand()

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	go func() {
		for pInstance := range srv.Accept() {
			var _ *player.Player = pInstance

			log.Printf("[Cobblit Network] Jogador conectado: %s (UUID: %s) IP: %s", pInstance.Name(), pInstance.UUID(), pInstance.Addr())

			players.RegistrarEntrada(pInstance)
			plugins.RegistrarBoasVindas(pInstance)

			nome := pInstance.Name()
			if plugins.GetSaldo(nome) == 0 && cfg.MoedaInicial > 0 {
				plugins.AddSaldo(nome, cfg.MoedaInicial)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		for range ticker.C {
			chat.Global.WriteString("§l§b[Cobblit Engine]§r §7Servidor rodando com alta performance e estabilidade!")
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
