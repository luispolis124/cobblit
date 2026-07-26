package plugins

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

var (
	whitelistFile   = "whitelist.json"
	whitelistMutex  sync.Mutex
	whitelistAtiva  = false
	whitelistJogos  = make(map[string]bool) // Nome do jogador -> true
)

// CarregarWhitelist lê o arquivo JSON e o estado da whitelist
type WhitelistCommand struct {
	Acao   string      `name:"on|off|list|add|remove|reload"`
	Jogador cmd.Varargs `name:"jogador" optional:"true"`
}

func CarregarWhitelist() {
	whitelistMutex.Lock()
	defer whitelistMutex.Unlock()

	data, err := os.ReadFile(whitelistFile)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &whitelistJogos)
}

func salvarWhitelistInterno() {
	data, err := json.MarshalIndent(whitelistJogos, "", "  ")
	if err == nil {
		_ = os.WriteFile(whitelistFile, data, 0644)
	}
}

func SalvarWhitelist() {
	whitelistMutex.Lock()
	defer whitelistMutex.Unlock()
	salvarWhitelistInterno()
}

func (c WhitelistCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para gerenciar a whitelist.")
			return
		}
	}

	acao := strings.ToLower(strings.TrimSpace(c.Acao))
	alvo := strings.TrimSpace(string(c.Jogador))

	CarregarWhitelist()

	switch acao {
	case "on":
		whitelistAtiva = true
		o.Printf("§a[Cobblit Engine] A whitelist foi ativada.")

	case "off":
		whitelistAtiva = false
		o.Printf("§c[Cobblit Engine] A whitelist foi desativada.")

	case "reload":
		CarregarWhitelist()
		o.Printf("§a[Cobblit Engine] Whitelist recarregada com sucesso.")

	case "list":
		whitelistMutex.Lock()
		var nomes []string
		for nome := range whitelistJogos {
			nomes = append(nomes, nome)
		}
		whitelistMutex.Unlock()

		sort.Strings(nomes)
		total := len(nomes)
		o.Printf("§aHá %d jogadores na whitelist: §7%s", total, strings.Join(nomes, ", "))

	case "add":
		if alvo == "" {
			o.Errorf("Uso correto: /whitelist add <jogador>")
			return
		}
		whitelistMutex.Lock()
		whitelistJogos[strings.ToLower(alvo)] = true
		salvarWhitelistInterno()
		whitelistMutex.Unlock()

		o.Printf("§a[Cobblit Engine] O jogador §e%s §afoi adicionado à whitelist.", alvo)

	case "remove":
		if alvo == "" {
			o.Errorf("Uso correto: /whitelist remove <jogador>")
			return
		}
		whitelistMutex.Lock()
		delete(whitelistJogos, strings.ToLower(alvo))
		salvarWhitelistInterno()
		whitelistMutex.Unlock()

		o.Printf("§e[Cobblit Engine] O jogador §c%s §foi removido da whitelist.", alvo)

		// Desconecta o jogador se ele estiver online e a whitelist estiver ativa
		if whitelistAtiva {
			for entity := range tx.Players() {
				if targetP, ok := entity.(*player.Player); ok {
					if strings.EqualFold(targetP.Name(), alvo) {
						targetP.Disconnect("§cVocê não está na whitelist deste servidor!")
						break
					}
				}
			}
		}

	default:
		o.Errorf("Comandos disponíveis: /whitelist <on|off|list|reload|add|remove>")
	}
}

// RegistrarWhitelistCommand registra o comando /whitelist no motor.
func RegistrarWhitelistCommand() {
	CarregarWhitelist()
	cmd.Register(cmd.New("whitelist", "Gerencia a lista de permissões do servidor", []string{"wl"}, WhitelistCommand{}))
}
