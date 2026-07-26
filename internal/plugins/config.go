package plugins

import (
	"encoding/json"
	"os"
)

// ServidorConfig define a estrutura do arquivo config.json
type ServidorConfig struct {
	Motd           string `json:"motd"`
	MaxPlayers     int    `json:"max_players"`
	ViewDistance   int    `json:"view_distance"`
	OnlineMode     bool   `json:"online_mode"`
	MoedaInicial   float64 `json:"moeda_inicial"`
}

var ConfigGlobal = ServidorConfig{
	Motd:         "Cobblit Engine Alpha",
	MaxPlayers:   100,
	ViewDistance: 8,
	OnlineMode:   true,
	MoedaInicial: 100.0, // Dinheiro padrão que novos jogadores ganham
}

const arquivoConfig = "config.json"

// CarregarOuCriarConfig lê o arquivo JSON ou gera um padrão se não existir
func CarregarOuCriarConfig() ServidorConfig {
	file, err := os.ReadFile(arquivoConfig)
	if err != nil {
		// Se o arquivo não existe, salva o padrão
		SalvarConfig(ConfigGlobal)
		return ConfigGlobal
	}

	_ = json.Unmarshal(file, &ConfigGlobal)
	return ConfigGlobal
}

func SalvarConfig(cfg ServidorConfig) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(arquivoConfig, data, 0644)
}
