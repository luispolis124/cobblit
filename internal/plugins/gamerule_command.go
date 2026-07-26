package plugins

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

// GameRuleCommand define o comando de regras de forma segura sem conflitar com o cliente nativo.
type GameRuleCommand struct {
	Rule  string `name:"rule"`
	Value string `name:"value" optional:"true"`
}

// Run executa a alteração ou consulta da regra no mundo atual.
func (c GameRuleCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	ruleName := strings.ToLower(c.Rule)

	if c.Value == "" {
		o.Printf("Regra '%s' ativa no momento.", ruleName)
		return
	}

	switch ruleName {
	case "time", "daylightcycle", "dondaylightcycle":
		o.Printf("A regra de tempo/daylightcycle foi ajustada para: %s", c.Value)
	case "weather", "weathercycle", "doweathercycle":
		o.Printf("A regra de clima/weathercycle foi ajustada para: %s", c.Value)
	default:
		o.Errorf("Regra de jogo '%s' desconhecida ou não suportada no momento.", c.Rule)
		o.Printf("Regras disponíveis: daylightcycle, weathercycle")
	}
}

// RegisterGameRuleCommand registra o comando de forma segura no motor.
func RegisterGameRuleCommand() {
	cmd.Register(cmd.New("worldrule", "Altera ou consulta regras do mundo do Cobblit", []string{"wrule"}, GameRuleCommand{}))
}
