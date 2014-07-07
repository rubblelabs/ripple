package peers

import (
	"bufio"
	"code.google.com/p/goprotobuf/proto"
	"crypto/sha512"
	"fmt"
	"github.com/donovanhide/sslconn"
	"github.com/golang/glog"
	"github.com/rubblelabs/ripple/crypto"
	"github.com/rubblelabs/ripple/peers/protocol"
	"math/big"
	"net"
	"time"
)

type ConnectionStatus int

const (
	peerTimeout   = time.Second * 5
	peerBuffer    = 4096 // 256 * 1024 //256k
	defaultCipher = "ALL:!LOW::!MD5:@STRENGTH"
)

type PeerConnection struct {
	Host    string
	Port    string
	Trusted bool
	Conn    net.Conn
}

type Conn struct {
	Host     string
	Port     string
	Resolved string
	*sslconn.Conn
	conn net.Conn
}

func Listen(m *Manager, port string) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		glog.Fatalln("HandleIncoming:", err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		glog.Fatalln("HandleIncoming:", err)
	}
	defer l.Close()
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			glog.Errorln("HandleIncoming:", err)
			continue
		}
		host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			glog.Errorln("HandleIncoming:", err)
			continue
		}
		glog.Infoln("Incoming Host: %s Port: %s", host, port)
		m.AddPeer(host, port, false, conn)
	}
}

func NewConn(p *PeerConnection) (*Conn, error) {
	c := &Conn{
		Host: p.Host,
		Port: p.Port,
		conn: p.Conn,
	}
	names, err := net.LookupAddr(c.Host)
	if err == nil {
		c.Resolved = names[0]
	} else {
		c.Resolved = "No reverse DNS possible"
	}
	if c.conn == nil {
		var err error
		if c.conn, err = net.DialTimeout("tcp", c.String(), peerTimeout); err != nil {
			return c, fmt.Errorf("Connect: %s", err.Error())
		}
	}
	config := &sslconn.Config{
		CipherList: defaultCipher,
	}
	if c.Conn, err = sslconn.NewConn(c.conn, c.conn, config, false); err != nil {
		return c, fmt.Errorf("Connect: %s", err.Error())
	}
	if err := c.Handshake(); err != nil {
		return c, fmt.Errorf("Connect: %s", err.Error())
	}
	return c, nil
}

func (c *Conn) run(incoming chan protocol.ExtendedMessage, outgoing chan proto.Message) {
	go c.writePump(outgoing)
	c.readPump(incoming)
	close(incoming)
	c.Shutdown()
	c.Free()
	c.conn.Close()
	glog.Errorf("%s Peer: Connection closed", c.String())
}

func (c *Conn) writePump(out chan proto.Message) {
	w := bufio.NewWriterSize(c, peerBuffer)
	encoder := protocol.NewEncoder(w)
	var err error
	for msg := range out {
		if err = encoder.Encode(msg); err != nil {
			glog.Errorf("%s: Peer Write Encode: %s\n", c.String(), err.Error())
			return
		}
		if err = w.Flush(); err != nil {
			glog.Errorf("%s: Peer Write Flush: %s\n", c.String(), err.Error())
			return
		}
	}
	glog.Errorf("%s: Peer Write Pump ended", c.String())
}

func (c *Conn) readPump(in chan protocol.ExtendedMessage) {
	r := bufio.NewReaderSize(c, peerBuffer)
	decoder := protocol.NewDecoder(r)
	for {
		msg, err := decoder.Decode()
		if err != nil {
			glog.Errorf("%s: Peer Read: %s\n", c.String(), err.Error())
			return
		}
		in <- msg
	}
}

func (c *Conn) getSessionCookie() ([]byte, error) {
	hasher := sha512.New()
	if _, err := hasher.Write(c.GetFinishedMessage(1024)); err != nil {
		return nil, fmt.Errorf("Peer GetSessionCookie: %s", err.Error())
	}
	left := big.NewInt(0).SetBytes(hasher.Sum(nil))
	hasher.Reset()
	if _, err := hasher.Write(c.GetPeerFinishedMessage(1024)); err != nil {
		return nil, fmt.Errorf("Peer GetSessionCookie: %s", err.Error())
	}
	right := big.NewInt(0).SetBytes(hasher.Sum(nil))
	return crypto.Sha512Half(left.Xor(left, right).Bytes()), nil
}

func (c *Conn) String() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *PeerConnection) String() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
