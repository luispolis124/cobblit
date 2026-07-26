package plugins

import (
	"runtime"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

var startTime = time.Now()

// StatusCommand define o comando /status para monitoramento de desempenho.
type StatusCommand struct{}

func (StatusCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Converte bytes para Megabytes (MB)
	allocMB := float64(m.Alloc) / 1024 / 1024
	sysMB := float64(m.Sys) / 1024 / 1024
	numGoroutines := runtime.NumGoroutine()
	uptime := time.Since(startTime).Truncate(time.Second)

	p, isPlayer := src.(*player.Player)

	if isPlayer {
		p.Message("§8--- §b[Cobblit Engine] Status do Servidor §8---")
		p.Messagef("§7Tempo ligado (Uptime): §e%v", uptime)
		p.Messagef("§7Goroutines ativas: §a%d", numGoroutines)
		p.Messagef("§7Memória Alocada: §e%.2f MB", allocMB)
		p.Messagef("§7Memória do Sistema (Sys): §e%.2f MB", sysMB)
		p.Message("§8-------------------------------------------")
	} else {
		o.Printf("[Cobblit Engine] Status - Uptime: %v | Goroutines: %d | Memória Alocada: %.2f MB | Sistema: %.2f MB", uptime, numGoroutines, allocMB, sysMB)
	}
}

// GcCommand define o comando /gc para limpeza de lixo de memória.
type GcCommand struct{}

func (GcCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	var mBefore runtime.MemStats
	runtime.ReadMemStats(&mBefore)

	// Executa a coleta de lixo forçada do Go
	runtime.GC()

	var mAfter runtime.MemStats
	runtime.ReadMemStats(&mAfter)

	liberadoMB := float64(mBefore.Alloc-mAfter.Alloc) / 1024 / 1024
	if liberadoMB < 0 {
		liberadoMB = 0
	}

	p, isPlayer := src.(*player.Player)

	if isPlayer {
		p.Messagef("§a[Cobblit Engine] Limpeza de memória (GC) concluída! Memória liberada: §e%.2f MB", liberadoMB)
	} else {
		o.Printf("[Cobblit Engine] Limpeza de memória (GC) concluída! Memória liberada: %.2f MB", liberadoMB)
	}
}

// RegistrarComandosStatus registra os comandos /status e /gc no motor.
func RegistrarComandosStatus() {
	cmd.Register(cmd.New("status", "Exibe o uso de memória e desempenho do motor", []string{"tps", "perf"}, StatusCommand{}))
	cmd.Register(cmd.New("gc", "Força a limpeza de lixo de memória do servidor", []string{"clearmem"}, GcCommand{}))
}
