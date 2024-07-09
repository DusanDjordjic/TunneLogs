package handlers

import (
	"errors"
	"fmt"
	"sync"
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

type Lobby struct {
	name     string
	client   *websocket.Conn
	producer *websocket.Conn
	started  bool
}

func (lobby *Lobby) isReady() bool {
	return lobby.client != nil && lobby.producer != nil
}

const BYE_MESSAGE = "bye"

func (lobby *Lobby) Start() {
	if lobby.started {
		return
	}
	log := logger.Log.Named(fmt.Sprintf("[Lobby::%s]", lobby.name))

	for {
		messageType, message, err := lobby.producer.ReadMessage()
		if err != nil {
			// TODO do better error checking
			log.Error("failed to read message from producer", zap.Error(err))
			log.Info("sending bye to client")
			err := lobby.client.WriteMessage(websocket.TextMessage, []byte(BYE_MESSAGE))
			if err != nil {
				log.Error("failed to send bye to client", zap.Error(err))
			}

			lobby.client.Close()
			lobby.producer.Close()
			break
		}

		err = lobby.client.WriteMessage(messageType, message)
		if err != nil {
			// TODO do better error checking
			log.Error("failed to send message to client", zap.Error(err))
			log.Info("sending bye to producer")
			err := lobby.producer.WriteMessage(websocket.TextMessage, []byte(BYE_MESSAGE))
			if err != nil {
				log.Error("failed to send bye to server", zap.Error(err))
			}

			lobby.client.Close()
			lobby.producer.Close()
			break
		}
	}

	lock.Lock()
	defer lock.Unlock()
	delete(lobbies, lobby.name)
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
		lobby = &Lobby{
			name:     lobbyName,
			client:   nil,
			producer: nil,
			started:  false,
		}

		lobbies[lobbyName] = lobby
	}

	if lobby.client != nil {
		// TODO support multiple clients
		log.Error("client is already connected")
		return errors.New("client is already connected")
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Error("failed to upgrade to ws", zap.Error(err))
		return err
	}

	lobby.client = conn

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
		lobby = &Lobby{
			name:     lobbyName,
			client:   nil,
			producer: nil,
			started:  false,
		}
		lobbies[lobbyName] = lobby
	}

	if lobby.producer != nil {
		// we only want one producer
		log.Error("producer is already connected")
		return errors.New("producer is already connected")
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Error("failed to upgrade to ws", zap.Error(err))
		return err
	}

	lobby.producer = conn

	if lobby.isReady() {
		log.Info("lobby is ready. startting...", zap.String("Name", lobbyName))
		go lobby.Start()
	} else {
		log.Info("lobby is not ready", zap.String("Name", lobbyName))
	}

	log.Debug("finished")
	return nil
}
