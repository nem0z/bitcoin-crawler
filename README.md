## Component

1. Crawler
This is the core component of the application, basically an engine which start new connections and handle their result.
There is only 1 per instance

2. Node
A Node describe a node in the Bitcoin network and keep track of some information about it:
```go
type Node struct {
	Timestamp time.Time `json:"timestamp"`
	Info      *Info     `json:"info"`
	Addr      *Addr     `json:"addr"`
	Ping      bool      `json:"ping"`
}
```

3. Peer
A peer is bascially a connection to a Bitcoin node, it's used to communicate with this node
```go
type Peer struct {
	ip        string
	port      int
	conn      net.Conn
	handlers  Handlers
	Info      *Info
	Addrs     []*Addr
	PingNonce []byte
	PingAt    time.Time
	PongAt    time.Time
	queue     chan *message.Message
	onClose   chan *Node
}
```

## Workflow

1. Main

- Init database
- Load known peers
- Start crawlers with known peers
- Load OK peers from database
- Wait for interupt (SIGINT or SIGTERM)

2. On Interupt

- Export nodes to `/expot/nodes.json`
- Save / Update nodes in DB

3. Create a crawler

- Create a struct crawler
- Register all the provided addr to the crawler
- Start the ResultHandler -> Basically update the status of a Node in map nodes
- Start N worker to handle new nodes
- Start moniroting -> Basically a go routine which log the state of the map of nodes every 10s

4. Worker
- Verify if the Node is already known (if yes stop here)
- Create + setup the new Peer
- Start the Peer
- Send 10 GetAddr message to the Peer
- After 60s close the connection (should be reworked to use context)
- Add the peer (node) to the map of known nodes

5. Create a new Peer
- Create a tcp connection to the IPv4 or IPv6 address
- Init the struct
  
6. Start a Peer
- Send Version message
- Send Verack message
- Send Ping message
- Send GetAddr message
- Start handler for differents messages