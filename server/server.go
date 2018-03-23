package main

import(
    "fmt"
    "net"
    "os"
    "encoding/json"
    "../core"
    "github.com/phayes/freeport"
	"log"
)

type Packet struct{
	message  string
	toUserID int
	toIP     *net.UDPAddr
}

type Server struct{
    conn     *net.UDPConn
    messages chan core.Message
    clients  map [int]Client
}

type Client struct{
    userID int
    userName string
    userAddr *net.UDPAddr
}

const(
	LIST_USER = iota
	NEW_USER
	NEW_MESSAGE
	DELETE_USER
)

var userInitialID = 1

func (s *Server) handleMessage(){
    var buf [8192]byte
    n, addr, err := s.conn.ReadFromUDP(buf[0:])
    if err != nil{
        return
    }
    msg := buf[0:n]
    m := s.analyzeMessage(msg)
    m.Time = core.NowTime()
	fmt.Println( "From", m.UserID, "To", m.ToUserID, "Message：",m.Content)
	switch m.Status{
		case LIST_USER:
			msg:= fmt.Sprintf("\n")
			for _, v := range s.clients {
				msg += fmt.Sprintf("%d -->> %s \n",v.userID,v.userName)
			}
			m.Content = msg
			m.FromIP = addr
			m.ToIP = addr
			s.messages <- m
		case NEW_USER:
            var c Client
            c.userAddr = addr
            c.userID = userInitialID
			userInitialID++
            c.userName = m.UserName
            s.clients[c.userID] = c

            m.UserID = c.userID
            m.ToUserID = -1
			m.FromIP = addr
			m.ToIP = addr

            s.messages <- m
        case NEW_MESSAGE:
			m.FromIP = addr
			search, _ := s.clients[m.ToUserID]
			m.ToIP = search.userAddr
			s.messages <- m
        case DELETE_USER:
			m.FromIP = addr
			m.ToIP = addr
			m.ToUserID = -1
			s.messages <- m

		default:
            fmt.Println("Cannot read the message:", string(msg))
    }
}

func (s *Server) analyzeMessage(msg []byte) (m core.Message) {
    json.Unmarshal(msg, &m)
    return
}


func (s *Server) sendMessage() {
    for{
        msg := <- s.messages
		str, _ := json.Marshal(msg)
		if msg.ToUserID>0{
			s.conn.WriteToUDP([]byte(str), msg.ToIP)
		}else{
			s.conn.WriteToUDP([]byte(str), msg.FromIP)
		}
    }

}

func checkError(err error){
    if err != nil{
        fmt.Fprintf(os.Stderr,"Fatal error:%s",err.Error())
        os.Exit(1)
    }
}

func main(){
	log.SetFlags(log.Llongfile)
	// Generate a random server port
	port, err := freeport.GetFreePort()
	if err != nil {
		// 随机端口失败就采用 7448
		port = 1200
	}
	// Default config
	//k:=core.RandPassword()
	config := &core.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
	}
	config.ReadConfig()
	config.SaveConfig()

	log.Println("Configruation", fmt.Sprintf(`
Server Port：
%s
	`, config.ListenAddr))

	if err != nil {
		log.Fatalln(err)
	}
	udpAddr, err := net.ResolveUDPAddr("udp4",config.ListenAddr)

	if err != nil {
		log.Fatalln(err)
	}

    var s Server
    s.messages = make(chan core.Message,20)
    s.clients =make(map[int]Client,0)

    s.conn,err = net.ListenUDP("udp",udpAddr)
    checkError(err)

    go s.sendMessage()

    for{
        s.handleMessage()
    }
}
