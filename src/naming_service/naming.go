package naming

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/douglas-soares/rpc-quick/src/rpc"
)

type namingRequest struct {
	ServerName string
	Addr       string
}

type namingResult struct {
	Addr string
}

type naming struct {
	addr   string
	server rpc.Server
	client rpc.Client
}

var servers map[string]string

func NewNamingService(addr string) *naming {
	servers = make(map[string]string)
	server := rpc.NewServer()

	return &naming{
		addr:   addr,
		server: server,
	}
}

func (n *naming) ListenAndServe(tlsConfig *tls.Config) error {
	n.server.Register("Bind", n.bind)
	n.server.Register("LookUp", n.lookUp)
	return n.server.ListenAndServe(n.addr, tlsConfig, nil)
}

func (n *naming) StartClient(tlsConfig *tls.Config) {
	n.client = rpc.NewClient(n.addr, tlsConfig, nil)
}

func (n *naming) Bind(serverName string, serverAddr string) error {
	return n.client.Call("Bind", namingRequest{
		ServerName: serverName,
		Addr:       serverAddr,
	}, nil)
}

func (n *naming) LookUp(serverName string) (string, error) {
	var result namingResult
	err := n.client.Call("LookUp", serverName, &result)
	if err != nil {
		return "", err
	}
	return result.Addr, nil
}

func (n *naming) bind(req namingRequest) {
	// missing handle threads
	servers[req.ServerName] = req.Addr
	log.Println(req.ServerName, "inserted with infos:", req.Addr)
	log.Println(servers)
}

func (n *naming) lookUp(serverName string) namingResult {
	fmt.Println("lookup result:", servers[serverName])
	return namingResult{Addr: servers[serverName]}
}
