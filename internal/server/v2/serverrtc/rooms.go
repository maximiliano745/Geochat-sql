package v2

import (
	"log"
	"math/rand" // Importa el paquete "math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Participant struct {
	Host bool
	Conn *websocket.Conn
}

type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant
}

func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant)
}

func (r *RoomMap) Get(roomID string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Map[roomID]
}

func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock() // Usamos Lock en lugar de RLock, ya que estamos modificando el mapa
	defer r.Mutex.Unlock()

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	var letters = []rune("abcdefghijkmnopqrstABCDEFGHIJKMNOPQRST134567890")
	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
	}
	roomId := string(b)
	r.Map[roomId] = []Participant{}

	return roomId
}

func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	p := Participant{host, conn}

	log.Println("Insert dentro del Room con RoomID: ", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)

}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	delete(r.Map, roomID)
}
