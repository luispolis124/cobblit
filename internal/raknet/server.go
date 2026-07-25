package raknet

import (
	"fmt"
	"net"
)

type Server struct {
	Listener *net.UDPConn
	Address  string
}

func NewServer(address string) *Server {
	return &Server{
		Address: address,
	}
}

func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", s.Address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	s.Listener = conn

	fmt.Printf("[RakNet] Servidor escutando em %s\n", s.Address)

	// Loop principal de escuta de pacotes UDP brutos
	buffer := make([]byte, 2048)
	for {
		n, remoteAddr, err := s.Listener.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("[Erro] Falha ao ler pacote UDP: %v\n", err)
			continue
		}

		go s.handlePacket(remoteAddr, buffer[:n])
	}
}

func (s *Server) handlePacket(addr *net.UDPAddr, data []byte) {
	if len(data) == 0 {
		return
	}
	
	// Identificador preliminar do pacote RakNet
	packetID := data[0]
	fmt.Printf("[RakNet] Pacote recebido de %s | ID: 0x%x | Tamanho: %d bytes\n", addr.String(), packetID, len(data))
}
