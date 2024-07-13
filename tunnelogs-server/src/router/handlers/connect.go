package handlers

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"tunnelogs-server/logger"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	lobbies map[string]*Lobby = make(map[string]*Lobby)
	lock    sync.Mutex
)

type MessageEvent int

const (
	ProducerDisconnected MessageEvent = 1
)

type Lobby struct {
	name     string
	clients  []*websocket.Conn
	producer *websocket.Conn
	started  bool
	lock     sync.Mutex
}

func newLobby(name string) *Lobby {
	return &Lobby{
		name:     name,
		clients:  make([]*websocket.Conn, 0),
		producer: nil,
		started:  false,
		lock:     sync.Mutex{},
	}
}

func (lobby *Lobby) isReady() bool {
	lobby.lock.Lock()
	defer lobby.lock.Unlock()

	return len(lobby.clients) > 0 && lobby.producer != nil
}

func (lobby *Lobby) addProducer(newProducer *websocket.Conn) {
	lobby.lock.Lock()
	defer lobby.lock.Unlock()

	if lobby.producer != nil {
		lobby.producer.WriteMessage(websocket.CloseMessage, nil)
		lobby.producer.Close()
	}

	lobby.producer = newProducer
}

func (lobby *Lobby) addClient(newClient *websocket.Conn) {
	lobby.lock.Lock()
	defer lobby.lock.Unlock()
	lobby.clients = append(lobby.clients, newClient)
}

const BYE_MESSAGE = "bye"

func (lobby *Lobby) Start() {
	lobby.lock.Lock()
	if lobby.started {
		lobby.lock.Unlock()
		return
	}

	lobby.started = true
	lobby.lock.Unlock()

	log := logger.Log.Named(fmt.Sprintf("[Lobby::%s]", lobby.name))

	for {
		func() {
			lobby.lock.Lock()
			defer lobby.lock.Unlock()

			if lobby.producer == nil {
				time.Sleep(time.Second)
				return
			}

			messageType, message, err := lobby.producer.ReadMessage()
			if err != nil {
				log.Error("failed to read message from producer", zap.Error(err))
				lobby.producer.WriteMessage(websocket.CloseMessage, nil)
				lobby.producer.Close()
				lobby.producer = nil
				return
			}

			log.Warn("sending message to clients", zap.Int("Clients", len(lobby.clients)))

			removeClients := make([]int, 0)
			for i, client := range lobby.clients {
				err = client.WriteMessage(messageType, message)
				if err != nil {
					// TODO do better error checking
					log.Error("failed to send message to client", zap.Error(err))
					client.WriteMessage(websocket.CloseMessage, nil)
					client.Close()
					removeClients = append(removeClients, i)
				}
			}

			for _, index := range removeClients {
				lobby.clients = append(lobby.clients[:index], lobby.clients[index+1:]...)
			}
		}()
	}
}

func ClientWSHandler(c echo.Context) error {
	log := logger.Log.Named("[ClientWSHandler]")
	log.Debug("started")

	lobbyName := c.Param("name")
	log.Info("Client connecting to lobby", zap.String("Name", lobbyName))
	if lobbyName == "" {
		log.Error("lobby name cannot be empty")
		return errors.New("lobby name cannot be empty")
	}

	lock.Lock()
	defer lock.Unlock()

	lobby, ok := lobbies[lobbyName]
	if !ok {
		lobby = newLobby(lobbyName)
		lobbies[lobbyName] = lobby
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Error("failed to upgrade to ws", zap.Error(err))
		return err
	}

	log.Info("Adding new client")
	lobby.addClient(conn)

	if lobby.isReady() {
		log.Info("lobby is ready. startting...", zap.String("Name", lobbyName))
		go lobby.Start()
	} else {
		log.Info("lobby is not ready", zap.String("Name", lobbyName))
	}

	log.Debug("finished")
	return nil
}

func ServerWSHandler(c echo.Context) error {
	log := logger.Log.Named("[ServerWSHandler]")
	log.Debug("started")

	lobbyName := c.Param("name")
	log.Info("Server connecting to lobby", zap.String("Name", lobbyName))
	if lobbyName == "" {
		log.Error("lobby name cannot be empty")
		return errors.New("lobby name cannot be empty")
	}

	lock.Lock()
	defer lock.Unlock()

	lobby, ok := lobbies[lobbyName]
	if !ok {
		lobby = newLobby(lobbyName)
		lobbies[lobbyName] = lobby
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Error("failed to upgrade to ws", zap.Error(err))
		return err
	}

	log.Info("Adding New Producer")
	lobby.addProducer(conn)

	if lobby.isReady() {
		log.Info("lobby is ready. startting...", zap.String("Name", lobbyName))
		go lobby.Start()
	} else {
		log.Info("lobby is not ready", zap.String("Name", lobbyName))
	}

	log.Debug("finished")
	return nil
}
