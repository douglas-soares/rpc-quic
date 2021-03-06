package naming

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/douglas-soares/rpc-quick/src/rpc"
)

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

func (n *naming) LookUp(serverName string) (string, error) {
	var result interface{}
	err := n.client.Call(result, "LookUp", serverName)
	fmt.Println(err)
	if err != nil {

		return "", err
	}
	return "result.(string)", nil
}

func (n *naming) bind(serverName string, serverInfo string) {
	// missing handle threads
	servers[serverName] = serverInfo
	log.Println(serverName, "inserted with infos:", serverInfo)
	log.Println(servers)
}

func (n *naming) lookUp(serverName string) string {
	fmt.Println("lookup result:", servers[serverName])
	return servers[serverName]
}
