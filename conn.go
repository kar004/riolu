// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/kar004/websocket"
	"log"
	"net/http"
	"time"
	"fmt"
    "encoding/json"
    s "strings"
)

type datos struct {
    ID int      `json:"id"`
    Nombre string `json:"nombre"`
    Edad string `json:"edad"`
    Accion string `json:"accion"`
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

}
// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		h.broadcast <- message
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			res := &datos{}
    		json.Unmarshal([]byte(message), &res)
    		insert:=s.Contains(res.Accion,"insertar")
    		update:=s.Contains(res.Accion,"modificar")
    		delete:=s.Contains(res.Accion,"eliminar")
    		/*fmt.Println(res.ID)
    		fmt.Println(res.Nombre)
    		fmt.Println(res.Edad)
    		fmt.Println(res.Accion)*/
			///////////////////////////////////////////////////////////////////
    		if update==true{
    		modificar(res.ID, res.Nombre, res.Edad)
			consulta(res.ID)
			X:=row["nombre"].(string)
			Y:=row["edad"].(string)
			mensaje:=([]byte("se ha modificado a "+X+" de "+Y+" años con exito"))
			if err := c.write(websocket.TextMessage, mensaje); err != nil {
				return
			}	
    		}else if insert==true {
    		max:=max()
    		insertar(max,res.Nombre,res.Edad)
			consulta(max)
			X:=row["nombre"].(string)
			Y:=row["edad"].(string)
			mensaje:=([]byte("se ha insertado a "+X+" de "+Y+" años con exito"))
			if err := c.write(websocket.TextMessage, mensaje); err != nil {
				return
			}
    		}else if delete==true {
    		consulta(res.ID)
			X:=row["nombre"].(string)
			Y:=row["edad"].(string)
			eliminar(res.ID)
			mensaje:=([]byte("se ha eliminado a "+X+" de "+Y+" años con exito"))
			if err := c.write(websocket.TextMessage, mensaje); err != nil {
				return
			}
    		}else{
    			fmt.Println("no se encontro indicacion")
    		}
    		///////////////////////////////////////////////////////////////////		
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
	c.readPump()

}