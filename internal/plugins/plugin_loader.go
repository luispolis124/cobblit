package plugins

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

type PluginInfo struct {
	Nome  string
	Ativo bool
}

var PluginsExternos []PluginInfo

func CarregarPluginsSo() {
	pastaPlugins := "plugins"

	if err := os.MkdirAll(pastaPlugins, 0755); err != nil {
		log.Printf("[Cobblit Engine] Erro ao criar diretório de plugins: %v", err)
		return
	}

	arquivos, err := os.ReadDir(pastaPlugins)
	if err != nil {
		log.Printf("[Cobblit Engine] Erro ao ler a pasta de plugins: %v", err)
		return
	}

	for _, arquivo := range arquivos {
		if arquivo.IsDir() || !strings.HasSuffix(arquivo.Name(), ".so") {
			continue
		}

		caminhoSo := filepath.Join(pastaPlugins, arquivo.Name())
		cPath := C.CString(caminhoSo)
		defer C.free(unsafe.Pointer(cPath))

		// Abre a biblioteca .so usando dlopen do C
		handle := C.dlopen(cPath, C.RTLD_LAZY)
		nomePlugin := strings.TrimSuffix(arquivo.Name(), ".so")

		if handle == nil {
			log.Printf("[Cobblit Engine] Erro ao carregar o plugin %s via CGO", arquivo.Name())
			PluginsExternos = append(PluginsExternos, PluginInfo{
				Nome:  nomePlugin,
				Ativo: false,
			})
			continue
		}

		// Procura pela função exportada OnEnable
		symName := C.CString("OnEnable")
		defer C.free(unsafe.Pointer(symName))

		symbol := C.dlsym(handle, symName)
		if symbol != nil {
			// Converte o ponteiro para uma função do Go e executa
			onEnableFn := *(*func())(unsafe.Pointer(&symbol))
			onEnableFn()
		}

		// Registra diretamente pelo gerenciador de módulos para evitar duplicidade na listagem
		RegistrarModulo(nomePlugin)

		log.Printf("[Cobblit Engine] Plugin externo reconhecido e carregado: %s", nomePlugin)
	}
}
