// Package controller
/**
  @author: zhenxiangjin
  @update_date: 2024/1/22
  @note:
*/

package controller

import (
	"coolrun-backend-go/common"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

// @BasePath /api/v1

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许任何来源
	}} // 默认配置或自定义配置

// Client 结构体用于保存客户端连接
type Client struct {
	Connection      *websocket.Conn `json:"connection"`
	Room            *Room           `json:"room"`
	UserOpenID      string          `json:"user_open_id"`
	UserDisplayName string          `json:"user_display_name"`
	// 其他所需属性
}

// Room 结构体用于表示虚拟房间
type Room struct {
	CreatedAt          time.Time        `json:"created_at"`
	DeletedAt          time.Time        `json:"deleted_at"`
	ID                 int              `json:"id"`
	Name               string           `json:"name"`
	Owner              string           `json:"owner"`
	Capacity           int              `json:"capacity"`
	CurrentNum         int              `json:"current_num"`
	Clients            map[*Client]bool `json:"clients"`
	Broadcast          chan []byte      `json:"broadcast"`            // 用于广播消息到房间内的所有客户端
	LastBroadcastCache []byte           `json:"last_broadcast_cache"` // 保存最后一次广播的消息
}

type PersonalInformation struct {
	Username  string  `json:"username"`
	UserID    string  `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
	Distance  float64 `json:"distance"`
}

var Rooms = make(map[string]*Room)

// CreateRoom
// @Summary 创建房间
// @Description {Websocket} 创建房间,如果创建成功则自动加入该房间.
// @Accept application/json
// @Produce application/json
// @Param room_id query int true "Room ID"
// @Param room_name query string true "Room Name"
// @Param creator query int true "User ID(NOT OpenID)"
// @Success 200 {object} BroadcastMessage "成功"
// @Failure 400 {string} BroadcastMessage "请求错误,错误信息在属性 Message 里,websocket 会自动断联."
// @Router /run/room/create [get]
func CreateRoom(c *gin.Context) {
	// create room by user id
	// TODO(zhenxiang) random generate room ID or by sequence. 解决随机命名冲突问题.
	roomID := c.Query("room_id")
	roomName := c.Query("room_name")
	creator := c.Query("creator")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.SysError(err.Error())
		return
	}
	// check input information
	if roomID == "" || roomName == "" || creator == "" {
		err := SendMessage(conn, fmt.Sprintf("Error: Input information is insuffient!"))
		err = conn.Close()
		if err != nil {
			common.SysError(err.Error())
		}
		return
	}
	// check room
	_, exists := Rooms[roomID]
	if exists {
		err := SendMessage(conn, fmt.Sprintf("Error: Room %s is already exist!", roomID))
		err = conn.Close()
		if err != nil {
			common.SysError(err.Error())
		}
		return
	}
	// create room and save room in rooms
	// TODO(zhenxiang) add user id & user name information into client struct
	client := &Client{Connection: conn}
	clients := make(map[*Client]bool, 10)
	initBroadcastInfo := BroadcastMessage{
		RoomName:           roomName,
		Capacity:           common.VIPRoomCapacity,
		CurrentUserNumber:  0,
		ClientsInformation: make(map[string]PersonalInformation),
	}
	initBroadcastMessage, err := json.Marshal(initBroadcastInfo)
	if err != nil {
		common.SysError(err.Error())
	}
	intRoomID, err := strconv.Atoi(roomID)
	if err != nil {
		common.SysError(err.Error())
		err = SendMessage(conn, "Error: RoomID Convert error.")
		err := conn.Close()
		if err != nil {
			common.SysError(err.Error())
		}
		return
	}
	newRoom := Room{
		ID:                 intRoomID,
		Name:               roomName,
		Owner:              creator,
		Capacity:           common.VIPRoomCapacity, // TODO(zhenxiang) remove hardcode here
		CurrentNum:         0,
		Clients:            clients,
		Broadcast:          make(chan []byte),
		LastBroadcastCache: initBroadcastMessage,
	}
	newRoom.AddClient(client) // TODO(zhenxiang) 解耦加入部分,使用重定向代替.
	client.Room = &newRoom
	Rooms[strconv.Itoa(newRoom.ID)] = &newRoom
}

func SendMessage(conn *websocket.Conn, message interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Write the message to the WebSocket connection.
	err = conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}

	return nil
}

func (room *Room) AddClient(client *Client) {
	// TODO(zhenxiang) 检查房间是否已满,而且检查是否有权限进入房间.
	room.Clients[client] = true
	client.Room = room
	room.CurrentNum += 1

	go readFromClient(client)

	// 订阅房间的广播通道
	go func() {
		for message := range room.Broadcast {
			// 发送消息到客户端
			err := client.Connection.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				// 处理错误，可能需要从房间移除客户端
				common.SysError(err.Error())
			}
		}
	}()

	// Broadcast, new user
	room.broadcastMessage(fmt.Sprintf("User:%s has been joined this room!", client.UserDisplayName))
}

// clientInputMessage Encapsulated message from client
type clientInputMessage struct {
	PersonalInformation
}

func readFromClient(client *Client) {
	for {
		_, message, err := client.Connection.ReadMessage()
		if err != nil {
			common.SysError(err.Error())
			// TODO(zhenxiang) 处理错误，可能需要从房间移除客户端
			break
		}
		// encode message
		var clientMessage clientInputMessage
		err = json.Unmarshal(message, &clientMessage)
		if err != nil {
			common.SysError(err.Error())
			err := client.Connection.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			if err != nil {
				return
			}
		}
		// 处理接收到的消息，并将其添加到房间的汇总数据中
		client.Room.processClientMessage(clientMessage)
	}
}

type BroadcastMessage struct {
	RoomName           string                         `json:"room_name"`
	Capacity           int                            `json:"capacity"`
	CurrentUserNumber  int                            `json:"current_user_number"`
	ClientsInformation map[string]PersonalInformation `json:"clients_information"`
	Message            string                         `json:"message"`
}

func (room *Room) processClientMessage(message clientInputMessage) {
	// 更新房间内的消息汇总数据
	var broadcastMessage BroadcastMessage
	err := json.Unmarshal(room.LastBroadcastCache, &broadcastMessage)
	if err != nil {
		common.SysError(err.Error())
	}
	// remove sent message
	broadcastMessage.Message = ""
	// Merge client input message
	// load broadcast info from cache
	broadcastMessage.ClientsInformation[message.UserID] = message.PersonalInformation
	// 定期或按需将汇总数据广播到房间内的所有客户端
	room.broadcast(broadcastMessage)
}

// JoinRoom
// @Summary 加入房间
// @Accept application/json
// @Produce application/json
// @Param user_name query int true "User Display Name"
// @Param user_id query int true "User ID(NOT OpenID)"
// @Param room_id path int true "Room ID"
// @Success 200 {object} BroadcastMessage "成功"
// @Failure 400 {string} BroadcastMessage "请求错误,错误信息在属性 Message 里,websocket 会自动断联."
// @Router /run/room/join/{room_id} [get]
func JoinRoom(c *gin.Context) {
	// grab user info
	userName := c.Query("user_name")
	userID := c.Query("user_id")
	roomID := c.Param("room_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.SysError(err.Error())
		return
	}
	// check input information
	if roomID == "" || userName == "" || userID == "" {
		err := SendMessage(conn, fmt.Sprintf("Error: Input information is insuffient!"))
		err = conn.Close()
		if err != nil {
			common.SysError(err.Error())
		}
		return
	}
	intRoomID, err := strconv.Atoi(roomID)
	if err != nil {
		common.SysError(err.Error())
		err := SendMessage(conn, fmt.Sprintf("Error: Convert Room ID error."))
		err = conn.Close()
		if err != nil {
			common.SysError(err.Error())
		}
		return
	}

	// 在此处处理新连接，将其加入到虚拟房间中
	client := &Client{Connection: conn, UserOpenID: userID, UserDisplayName: userName}
	joinRoom(client, intRoomID)
}

func joinRoom(client *Client, roomID int) {
	// check if room exist
	room, exists := Rooms[strconv.Itoa(roomID)]
	common.SysLog(fmt.Sprintf("Room %d is %v exists.", roomID, exists))
	// test: check room info
	//fmt.Printf("Room info: %v", *room)
	// end test
	if !exists {
		err := client.Connection.WriteMessage(websocket.TextMessage, []byte("Room not exist."))
		if err != nil {
			common.SysError(err.Error())
		}
		err = client.Connection.Close()
		if err != nil {
			common.SysError(err.Error())
		}
		return
	}
	room.AddClient(client)
}

func (room *Room) broadcastMessage(message string) {
	var broadcastMessage BroadcastMessage
	err := json.Unmarshal(room.LastBroadcastCache, &broadcastMessage)
	if err != nil {
		common.SysError(err.Error())
	}
	broadcastMessage.CurrentUserNumber = room.CurrentNum
	broadcastMessage.Message = message
	room.broadcast(broadcastMessage)
}

func (room *Room) broadcast(message BroadcastMessage) {
	broadcastMessage, err := json.Marshal(message)
	if err != nil {
		common.SysError(err.Error())
		return
	}

	bufferCurrentNum := 0
	for client := range room.Clients {
		// check if a client is still alive.
		if client.Connection == nil {
			delete(room.Clients, client)
			continue
		}
		bufferCurrentNum += 1
		room.Broadcast <- broadcastMessage
		room.LastBroadcastCache = broadcastMessage
	}
	room.CurrentNum = bufferCurrentNum
}
