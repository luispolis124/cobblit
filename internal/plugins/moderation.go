package plugins

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

const arquivoBanimentos = "bans.json"

// Lista de jogadores banidos armazenada em memória (Nome -> Motivo)
var banidos = map[string]string{}

// CarregarBans lê o arquivo JSON de banimentos ao ligar o servidor
func CarregarBans() {
	dados, err := os.ReadFile(arquivoBanimentos)
	if err != nil {
		SalvarBans()
		return
	}
	_ = json.Unmarshal(dados, &banidos)
}

// SalvarBans grava a lista atualizada de banidos no arquivo JSON
func SalvarBans() {
	dados, err := json.MarshalIndent(banidos, "", "  ")
	if err == nil {
		_ = os.WriteFile(arquivoBanimentos, dados, 0644)
	}
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

	// Adiciona à lista de banidos e salva no JSON
	banidos[alvoStr] = razao
	SalvarBans()

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

	if _, existe := banidos[alvoStr]; !existe {
		o.Errorf("§cO jogador %s não está banido.", alvoStr)
		return
	}

	// Remove do mapa de banidos e atualiza o JSON
	delete(banidos, alvoStr)
	SalvarBans()

	o.Printf("§a[Cobblit Engine] O jogador %s foi desbanido com sucesso.", alvoStr)
}

func RegistrarComandosModeracao() {
	CarregarBans() // Carrega os banimentos salvos ao iniciar o motor
	cmd.Register(cmd.New("kick", "Expulsa um jogador do servidor", []string{}, KickCommand{}))
	cmd.Register(cmd.New("ban", "Bane um jogador do servidor", []string{}, BanCommand{}))
	cmd.Register(cmd.New("unban", "Remove o banimento de um jogador", []string{}, UnbanCommand{}))
}
