# Cobblit Engine - Arquitetura e Diretrizes de Alta Performance

Este documento descreve os pilares arquiteturais do **Cobblit Engine**, inspirados em modelos de alta concorrência em Go:

1. **Estrutura de Rotinas (Goroutines):** 
   Utilização de `go func()` para lidar com aceitação de clientes, simulação de mundos e eventos simultâneos sem bloquear a thread principal do motor.

2. **Uso como Biblioteca:** 
   O núcleo do servidor atua como um pacote importável, permitindo que o `main.go` funcione como ponto flexível de configuração e inicialização de regras customizadas.

3. **Gerenciamento de Pacotes:** 
   Gerenciamento estrito de dependências utilizando módulos nativos do Go (`go mod init` e `go get`) para isolamento e reprodutibilidade.

4. **Comunicação Assíncrona:** 
   Emprego de canais (`channels`) para a troca segura e desacoplada de dados entre threads em tarefas como chat, streaming de chunks e eventos de plugins.
