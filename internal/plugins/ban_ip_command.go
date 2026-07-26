package plugins

import (
	"encoding/json"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

var (
	ipBansFile   = "ip_bans.json"
	ipBansMutex  sync.Mutex
	ipBansMap    = make(map[string]string) // IP -> Motivo
)

func CarregarIpBans() {
	ipBansMutex.Lock()
	defer ipBansMutex.Unlock()

	data, err := os.ReadFile(ipBansFile)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &ipBansMap)
}

func salvarIpBansInterno() {
	data, err := json.MarshalIndent(ipBansMap, "", "  ")
	if err == nil {
		_ = os.WriteFile(ipBansFile, data, 0644)
	}
}

// BanIpCommand gerencia o comando /ban-ip
type BanIpCommand struct {
	Alvo   cmd.Varargs `name:"jogador_ou_ip"`
	Motivo cmd.Varargs `name:"motivo" optional:"true"`
}

func (c BanIpCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para usar este comando.")
			return
		}
	}

	argsStr := strings.TrimSpace(string(c.Alvo))
	if argsStr == "" {
		o.Errorf("Você precisa especificar um nome de jogador ou IP válido.")
		return
	}

	// Separa o alvo do motivo caso venha tudo junto no Varargs
	parts := strings.SplitN(argsStr, " ", 2)
	alvo := parts[0]
	razao := "IP Banido"
	if len(parts) > 1 {
		razao = parts[1]
	}
	if strings.TrimSpace(string(c.Motivo)) != "" {
		razao = string(c.Motivo)
	}

	CarregarIpBans()

	var ipAlvo string
	// Verifica se é um IP válido diretamente
	if net.ParseIP(alvo) != nil {
		ipAlvo = alvo
	} else {
		// Procura o jogador online pelo nome para pegar o IP dele
		for entity := range tx.Players() {
			if targetP, ok := entity.(*player.Player); ok {
				if strings.EqualFold(targetP.Name(), alvo) {
					// Extrai o IP removendo a porta do endereço de rede
					host, _, err := net.SplitHostPort(targetP.Addr().String())
					if err == nil {
						ipAlvo = host
					} else {
						ipAlvo = targetP.Addr().String()
					}

					// Desconecta o jogador imediatamente
					targetP.Disconnect("§cVocê foi banido do servidor por este IP.")
					break
				}
			}
		}
	}

	if ipAlvo == "" {
		o.Errorf("§cNão foi possível encontrar o IP ou o jogador '%s' está offline.", alvo)
		return
	}

	ipBansMutex.Lock()
	ipBansMap[ipAlvo] = razao
	salvarIpBansInterno()
	ipBansMutex.Unlock()

	o.Printf("§4[Cobblit Engine] O IP %s foi banido com sucesso. Motivo: %s", ipAlvo, razao)
}

// RegistrarBanIpCommand registra o comando /ban-ip no motor.
func RegistrarBanIpCommand() {
	CarregarIpBans()
	cmd.Register(cmd.New("ban-ip", "Bane o endereço IP de um jogador", []string{}, BanIpCommand{}))
}
