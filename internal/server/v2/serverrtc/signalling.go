package v2

import (
	//"log"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap

func CreateRoomRequestHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	roomID := AllRooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}

	log.Println(AllRooms.Map)

	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

type brodcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var brodcast = make(chan brodcastMsg)

func brodcaster() {
	for {
		msg := <-brodcast

		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil {
					log.Fatal(err)
					client.Conn.Close()
				}
			}

		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func JoinRoomRequestHandle(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]
	//fmt.Println("aca en JoinRoomRequestHandle con id: ", roomID)

	if !ok {
		log.Println("roomID missing in URL parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	go brodcaster()

	for {
		var msg brodcastMsg
		err := ws.ReadJSON(&msg.Message)

		if err != nil {
			log.Fatal("error: ", err)
		}

		msg.Client = ws
		msg.RoomID = roomID[0]

		brodcast <- msg

	}
}
