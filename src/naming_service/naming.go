package naming

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/douglas-soares/rpc-quick/src/rpc"
)

type NamingResult struct {
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
	return n.server.ListenAndServe(n.addr, tlsConfig)
}

func (n *naming) StartClient(tlsConfig *tls.Config) {
	n.client = rpc.NewClient(n.addr, tlsConfig)
}

func (n *naming) Bind(serverName string, serverAddr string) error {
	return n.client.Call(nil, "Bind", serverName, serverAddr)
}

func (n *naming) LookUp(serverName string) (NamingResult, error) {
	var result NamingResult
	err := n.client.Call(&result, "LookUp", serverName)
	if err != nil {
		return NamingResult{}, err
	}
	return result, nil
}

func (n *naming) bind(serverName string, serverInfo string) {
	// missing handle threads
	servers[serverName] = serverInfo
	log.Println(serverName, "inserted with infos:", serverInfo)
	log.Println(servers)
}

func (n *naming) lookUp(serverName string) NamingResult {
	fmt.Println("lookup result:", servers[serverName])
	return NamingResult{Addr: servers[serverName]}
}
