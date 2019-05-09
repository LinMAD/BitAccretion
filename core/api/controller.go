package api

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/LinMAD/BitAccretion/core/assembly"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// GraphStorageKey a key to get value from memory where stored graph
const (
	GraphStorageKey = "graph"
	tagController   = "CONTROLLER"
)

type controller struct {
	api *API
}

var socketUpgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: false, // https://github.com/gorilla/websocket/issues/203 unexpected byte fatal
}

func (c *controller) getTrafficData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "%s", "Incorect HTTP Method")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(assembly.WriteToJSON(c.api.storage.Get(GraphStorageKey))))
}

func (c *controller) getTrafficDataViaWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := socketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("%s: Client (%s) subscribed to web socket", tagController, conn.RemoteAddr())

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("%s %v", tagController, err.Error())
			return
		}

		// Make delay between responses
		time.Sleep(time.Duration(c.api.surveyDelay) * time.Millisecond)

		if string(msg) == "get_traffic" {
			buf := new(bytes.Buffer)
			data := assembly.WriteToJSON(c.api.storage.Get(GraphStorageKey))

			err = binary.Write(buf, binary.LittleEndian, data)
			if err != nil {
				log.Printf("%s %v", tagController, err.Error())
				return
			}

			err = conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			err = conn.Close()
			if err != nil {
				log.Fatalf("%s: Unable to close client err %s", tagController, err.Error())
			}

			log.Printf("%s: Client desconnected from web socket", tagController)
			return
		}
	}
}
