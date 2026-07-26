package plugins

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

var (
	spawnPos = mgl64.Vec3{0, 100, 0}
	homes    = make(map[string]mgl64.Vec3)
)

func RegistrarComandosMovimento() {
	cmd.Register(cmd.New("spawn", "Teleporta para o spawn do servidor", []string{}, SpawnCommand{}))
	cmd.Register(cmd.New("setspawn", "Define o spawn do servidor", []string{}, SetSpawnCommand{}))
	cmd.Register(cmd.New("home", "Teleporta para a sua casa salva", []string{}, HomeCommand{}))
	cmd.Register(cmd.New("sethome", "Define sua casa atual", []string{}, SetHomeCommand{}))
	cmd.Register(cmd.New("tp", "Teleporta para um jogador", []string{}, TpCommand{}))
}

type SpawnCommand struct{}

func (SpawnCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		p.Teleport(spawnPos)
		o.Printf("§a[Cobblit] Teleportado para o spawn!")
	} else {
		o.Errorf("Este comando só pode ser usado por jogadores.")
	}
}

type SetSpawnCommand struct{}

func (SetSpawnCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		spawnPos = p.Position()
		o.Printf("§a[Cobblit] Spawn do servidor atualizado com sucesso!")
	} else {
		o.Errorf("Este comando só pode ser usado por jogadores.")
	}
}

type HomeCommand struct{}

func (HomeCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		pos, existe := homes[p.Name()]
		if !existe {
			o.Errorf("Você ainda não definiu uma home! Use /sethome.")
			return
		}
		p.Teleport(pos)
		o.Printf("§a[Cobblit] Teleportado para a sua home!")
	} else {
		o.Errorf("Este comando só pode ser usado por jogadores.")
	}
}

type SetHomeCommand struct{}

func (SetHomeCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if p, ok := src.(*player.Player); ok {
		homes[p.Name()] = p.Position()
		o.Printf("§a[Cobblit] Home definida com sucesso na sua posição atual!")
	} else {
		o.Errorf("Este comando só pode ser usado por jogadores.")
	}
}

type TpCommand struct {
	Destino string `cmd:"jogador"`
}

func (c TpCommand) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("Apenas jogadores podem usar este comando.")
		return
	}

	targetFound := false
	for entity := range tx.Entities() {
		if target, ok := entity.(*player.Player); ok {
			if target.Name() == c.Destino {
				p.Teleport(target.Position())
				o.Printf("§a[Cobblit] Teleportado para %s!", target.Name())
				targetFound = true
				break
			}
		}
	}

	if !targetFound {
		o.Errorf("Jogador '%s' não encontrado ou offline.", c.Destino)
	}
}
