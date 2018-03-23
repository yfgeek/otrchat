package main

import(
    "fmt"
    "net"
    "os"
    "encoding/json"
    "../core"
	"log"
	"golang.org/x/crypto/otr"
	"errors"
	crand "crypto/rand"
	"github.com/fatih/color"
	"bufio"
	"strings"
)

type Client struct{
    conn *net.UDPConn
    gkey bool
    userID int
    userName string
    sendMessages chan string
    receiveMessages chan string
    toUserID int
}


const(
	LIST_USER = iota
	NEW_USER
	NEW_MESSAGE
	DELETE_USER
)


var (
	otrPrivKey   otr.PrivateKey
	otrConv      otr.Conversation
	otrSecChange otr.SecurityChange
)



func (c *Client) func_sendMessage(sid int,msg string){

	m:= core.Message{
		Status:sid,
		UserName: c.userName,
		Content: msg,
	}
	str, err := json.Marshal(m)

	if err != nil {
		fmt.Println("json err:", err)
	}

	str= []byte(string(str))
	_, err = c.conn.Write(str)
    checkError(err,"func_sendMessage")
}

func (c *Client) sendMessage() {
    for c.gkey {
        msg := <- c.sendMessages
		m:= core.Message{
			Status:NEW_MESSAGE,
			UserName: c.userName,
			UserID: c.userID,
			Content: msg,
			ToUserID: c.toUserID,
		}

		str, err := json.Marshal(m)

		if err != nil {
			fmt.Println("json err:", err)
		}
		str= []byte(string(str))
		_,err = c.conn.Write(str)
        checkError(err,"sendMessage")
    }

}

func (c *Client) receiveMessage() {
    var buf [8192]byte
    for c.gkey {
        n,err := c.conn.Read(buf[0:])
        checkError(err, "receiveMessage")
        c.receiveMessages <- string(buf[0:n])
    }
    
}

func (c *Client) getMessage() {
    for c.gkey {
		var inputReader *bufio.Reader
		inputReader = bufio.NewReader(os.Stdin)
		msg,err := inputReader.ReadString('\n')
		msg= strings.Trim(msg,"\n")
		checkError(err, "getMessage")
		switch msg {
		case "/quit":
			c.gkey = false
		case "/list":
			c.func_sendMessage(LIST_USER,"")
		case "/talk":
			fmt.Println("Who would you like speak to? (ID)")
			_, err = fmt.Scanln(&c.toUserID)
			fmt.Println("Establishing the secure tunnel over OTR protocol...")
			c.sendMessages <- otr.QueryMessage
		default:
			color.Green("%s [%s] : %s", c.userName,core.NowTime(), msg)
			msgToPeer, err := otrConv.Send([]byte(msg))
			checkError(err,"otr conv sending")
			for _, m := range msgToPeer {
				c.sendMessages <- string(m)
			}
		}
    }
}

func (c *Client) printMessage() {
    //var msg string
    for c.gkey {
        msg := <- c.receiveMessages
        var m core.Message
        json.Unmarshal([]byte(msg),&m)
		switch m.Status {
		case LIST_USER:
			fmt.Println(m.Content)
		case NEW_MESSAGE:
			if c.toUserID==0{
				c.toUserID = m.UserID
			}
			bytes := []byte(m.Content)
			var eee error
			out, encrypted, _, mPeer, eee := otrConv.Receive(bytes)
			checkError(eee,"new chat")
			if len(out) > 0 {
				if !encrypted {
					color.Red("<OTR> Conversation not yet encrypted!")
					color.Red("%s [%s] : %s", m.UserName, m.Time, string(out))
				}else{
					color.Blue("%s [%s] : %s", m.UserName, m.Time, string(out))
				}
			}
			if len(mPeer) > 0 {
				for _, msg := range mPeer {
					c.sendMessages <-string(msg)
				}
			}
		case NEW_USER:
			c.userID = m.UserID
			fmt.Println("Your ID is: ",c.userID)
			fmt.Println(m.Content)
		default:
			fmt.Println(m.UserName,":",m.Content)
		}

    }
}



func GeneratePrivateKey(){
	newKey := new(otr.PrivateKey)
	newKey.Generate(crand.Reader)
	keyBytes := newKey.Serialize(nil)
	rest, ok := otrPrivKey.Parse(keyBytes)
	if !ok {
		fmt.Println(fmt.Errorf("ERROR: Failed to parse private key \n"))
	}
	if len(rest) > 0 {
		fmt.Println(errors.New("ERROR: data remaining after parsing private key"))
	}
	otrConv.PrivateKey = &otrPrivKey
	otrConv.FragmentSize = 1000
	fmt.Println("Fingerprint: ", otrConv.PrivateKey.Fingerprint())
}

func checkError(err error, funcName string){
    if err != nil{
        fmt.Fprintf(os.Stderr,"Fatal error:%s in func:%s",err.Error(), funcName)
        os.Exit(1)
    }
}

func main(){
	log.SetFlags(log.Lshortfile)
	config := &core.Config{}
	config.ReadConfig()

	log.Println("Configruation", fmt.Sprintf(`
Remote Server：%s
Remote Server Port %s
`, config.RemoteAddr,config.ListenAddr))

    service := config.RemoteAddr + config.ListenAddr
    udpAddr, err := net.ResolveUDPAddr("udp4", service)
    checkError(err,"main")

    var c Client
    c.gkey = true
    c.toUserID = 0


    checkError(err,"main")
    fmt.Print("Nickname: ")
    _,err = fmt.Scanln(&c.userName)
    checkError(err,"main")

	GeneratePrivateKey()

	c.sendMessages = make(chan string)
	c.receiveMessages = make(chan string)

    c.conn,err = net.DialUDP("udp",nil,udpAddr)
    checkError(err,"main")
    defer c.conn.Close()

    c.func_sendMessage(NEW_USER, "Welcome ! My friend, " + c.userName)

    go c.printMessage()
    go c.receiveMessage()

    go c.sendMessage()
    c.getMessage()

    c.func_sendMessage(DELETE_USER,c.userName + "left the group")
    fmt.Println("Exited")


    os.Exit(0)
}
