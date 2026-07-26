package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

var (
	carteiras  = make(map[string]float64)
	ecoMutex   sync.Mutex
	arquivoEco = "economy.json"
)

func CarregarEconomia() {
	ecoMutex.Lock()
	defer ecoMutex.Unlock()

	file, err := os.ReadFile(arquivoEco)
	if err != nil {
		return
	}
	_ = json.Unmarshal(file, &carteiras)
}

// SalvarEconomiaInterno assume que o mutex JÁ ESTÁ travado (evita deadlock)
func salvarEconomiaInterno() {
	data, err := json.MarshalIndent(carteiras, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(arquivoEco, data, 0644)
}

func GetSaldo(nome string) float64 {
	ecoMutex.Lock()
	defer ecoMutex.Unlock()
	return carteiras[nome]
}

func AddSaldo(nome string, quantia float64) {
	ecoMutex.Lock()
	carteiras[nome] += quantia
	salvarEconomiaInterno()
	ecoMutex.Unlock()
}

func RemoveSaldo(nome string, quantia float64) bool {
	ecoMutex.Lock()
	defer ecoMutex.Unlock()
	if carteiras[nome] >= quantia {
		carteiras[nome] -= quantia
		salvarEconomiaInterno()
		return true
	}
	return false
}

// Comando /balance ou /money
type MoneyCommand struct {
	Alvo cmd.Varargs `optional:"true" name:"jogador"`
}

func (c MoneyCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("Este comando só pode ser usado por jogadores.")
		return
	}

	alvo := string(c.Alvo)
	if alvo == "" {
		saldo := GetSaldo(p.Name())
		p.Messagef("§a[Economia] Seu saldo atual é de: §e$%.2f", saldo)
	} else {
		saldo := GetSaldo(alvo)
		p.Messagef("§a[Economia] O saldo de §e%s §aé de: §e$%.2f", alvo, saldo)
	}
}

// Comando /pay com Varargs para evitar travamentos de input
type PayCommand struct {
	Destino cmd.Varargs `name:"jogador"`
	Quantia float64     `name:"quantia"`
}

func (c PayCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("Este comando só pode ser usado por jogadores.")
		return
	}

	destino := string(c.Destino)
	if c.Quantia <= 0 {
		p.Message("§cVocê deve enviar uma quantia maior que zero.")
		return
	}

	remetente := p.Name()
	if remetente == destino {
		p.Message("§cVocê não pode enviar dinheiro para si mesmo.")
		return
	}

	if RemoveSaldo(remetente, c.Quantia) {
		AddSaldo(destino, c.Quantia)
		p.Messagef("§a[Economia] Você enviou §e$%.2f §apara §e%s§.", c.Quantia, destino)
	} else {
		p.Message("§cVocê não tem saldo suficiente para essa transação.")
	}
}

// GiveMoneyCommand corrigido para valores altos salvando de verdade no JSON
type GiveMoneyCommand struct {
	Player string  `name:"player"`
	Amount float64 `name:"amount"`
}

func (c GiveMoneyCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if c.Amount <= 0 {
		o.Errorf("A quantia deve ser maior que zero.")
		return
	}

	AddSaldo(c.Player, c.Amount)

	for entity := range tx.Players() {
		if p, ok := entity.(*player.Player); ok {
			if p.Name() == c.Player {
				p.Message(fmt.Sprintf("§aVocê recebeu §e$%.2f §ade um administrador!", c.Amount))
				break
			}
		}
	}

	o.Printf("Você adicionou $%.2f para %s com sucesso.", c.Amount, c.Player)
}

// BalTopCommand dinâmico real baseado no economy.json
type BalTopCommand struct{}

func (BalTopCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("Este comando só pode ser usado dentro do jogo.")
		return
	}

	ecoMutex.Lock()
	type usuarioSaldo struct {
		Nome   string
		Wallet float64
	}
	var ranking []usuarioSaldo
	for k, v := range carteiras {
		ranking = append(ranking, usuarioSaldo{Nome: k, Wallet: v})
	}
	ecoMutex.Unlock()

	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].Wallet > ranking[j].Wallet
	})

	p.Message("§8--- §6Cobblit Economy Top §8---")
	limite := 5
	if len(ranking) < limite {
		limite = len(ranking)
	}

	if limite == 0 {
		p.Message("§7[Nenhum registro de economia encontrado]")
	} else {
		for i := 0; i < limite; i++ {
			p.Messagef("§7%d. §b%s §8- §e$%.2f", i+1, ranking[i].Nome, ranking[i].Wallet)
		}
	}
	p.Message("§8-----------------------------")
}

func RegistrarComandosEconomia() {
	CarregarEconomia()
	cmd.Register(cmd.New("money", "Consulta o seu saldo", []string{"balance", "bal"}, MoneyCommand{}))
	cmd.Register(cmd.New("pay", "Envia dinheiro para outro jogador", []string{}, PayCommand{}))
	cmd.Register(cmd.New("baltop", "Exibe o ranking mais rico do servidor", []string{"topmoney"}, BalTopCommand{}))
	cmd.Register(cmd.New("givemoney", "Adiciona dinheiro para um jogador (Admin)", []string{"addmoney"}, GiveMoneyCommand{}))
}
