package plugins

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

var (
	bansFile   = "bans.json"
	bansMutex  sync.Mutex
	banListMap = make(map[string]string)
)

// CarregarBans lê o arquivo JSON de banimentos ao ligar o servidor de forma segura
func CarregarBans() {
	bansMutex.Lock()
	defer bansMutex.Unlock()

	dados, err := os.ReadFile(bansFile)
	if err != nil {
		salvarBansInterno()
		return
	}
	_ = json.Unmarshal(dados, &banListMap)
}

// salvarBansInterno grava a lista de banidos assumindo que o mutex JÁ ESTÁ travado
func salvarBansInterno() {
	dados, err := json.MarshalIndent(banListMap, "", "  ")
	if err == nil {
		_ = os.WriteFile(bansFile, dados, 0644)
	}
}

// SalvarBans grava a lista atualizada de banidos com segurança de concorrência
func SalvarBans() {
	bansMutex.Lock()
	defer bansMutex.Unlock()
	salvarBansInterno()
}

// KickCommand gerencia o comando /kick com segurança e Varargs anti-crash
type KickCommand struct {
	Alvo   cmd.Varargs `name:"jogador"`
	Motivo cmd.Varargs `name:"motivo" optional:"true"`
}

func (c KickCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para usar este comando.")
			return
		}
	}

	alvoStr := strings.TrimSpace(string(c.Alvo))
	if alvoStr == "" {
		o.Errorf("Você precisa especificar o nome do jogador para expulsar.")
		return
	}

	razao := strings.TrimSpace(string(c.Motivo))
	if razao == "" {
		razao = "Expulso do servidor"
	}

	o.Printf("§e[Cobblit Engine] Jogador %s expulso. Motivo: %s", alvoStr, razao)
}

// BanCommand gerencia o comando /ban com persistência em JSON
type BanCommand struct {
	Alvo   cmd.Varargs `name:"jogador"`
	Motivo cmd.Varargs `name:"motivo" optional:"true"`
}

func (c BanCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para usar este comando.")
			return
		}
	}

	alvoStr := strings.TrimSpace(string(c.Alvo))
	if alvoStr == "" {
		o.Errorf("Você precisa especificar o nome do jogador para banir.")
		return
	}

	razao := strings.TrimSpace(string(c.Motivo))
	if razao == "" {
		razao = "Banido do servidor"
	}

	bansMutex.Lock()
	banListMap[alvoStr] = razao
	salvarBansInterno()
	bansMutex.Unlock()

	o.Printf("§4[Cobblit Engine] O jogador %s foi banido. Motivo: %s", alvoStr, razao)
}

// UnbanCommand gerencia o comando /unban removendo do JSON
type UnbanCommand struct {
	Alvo cmd.Varargs `name:"jogador"`
}

func (c UnbanCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para usar este comando.")
			return
		}
	}

	alvoStr := strings.TrimSpace(string(c.Alvo))
	if alvoStr == "" {
		o.Errorf("Você precisa especificar o nome do jogador para desbanir.")
		return
	}

	bansMutex.Lock()
	if _, existe := banListMap[alvoStr]; !existe {
		bansMutex.Unlock()
		o.Errorf("§cO jogador %s não está banido.", alvoStr)
		return
	}

	delete(banListMap, alvoStr)
	salvarBansInterno()
	bansMutex.Unlock()

	o.Printf("§a[Cobblit Engine] O jogador %s foi desbanido com sucesso.", alvoStr)
}

func RegistrarComandosModeracao() {
	CarregarBans()
	cmd.Register(cmd.New("kick", "Expulsa um jogador do servidor", []string{}, KickCommand{}))
	cmd.Register(cmd.New("ban", "Bane um jogador do servidor", []string{}, BanCommand{}))
	cmd.Register(cmd.New("unban", "Remove o banimento de um jogador", []string{}, UnbanCommand{}))
}
