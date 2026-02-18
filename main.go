package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v3"
	"github.com/pion/webrtc/v4"
)

//go:embed frontend/*
var frontend_dir embed.FS

//go:embed frontend/static/*
var static_dir embed.FS

var peerConnections sync.Map

func main() {

	// create a views engine
	engine := html.NewFileSystem(http.FS(frontend_dir), ".html")

	// start a new fiber instance
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},           // Specify allowed origins
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"}, // Specify allowed methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},                         // Specify allowed headers
		AllowCredentials: true,                                                                 // Needed if you use cookies or sessions
	}))

	// embed the static files
	static_files := fs.FS(static_dir)

	// use those static files
	app.Use("/", static.New("", static.Config{FS: static_files}))

	// serve the index.html page
	app.Get("/", func(c fiber.Ctx) error {

		return c.Render("frontend/index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Post("/webrtc/offer", func(c fiber.Ctx) error {

		sdpOffer := new(webrtc.SessionDescription)

		if err := c.Bind().Body(sdpOffer); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		// log.Printf("Received handshake: sdp=%s, type=%s\n", sdpOffer.SDP, sdpOffer.Type)

		// 1. prepare configuration (stun/turn)
		config := webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{URLs: []string{"stun:stun.1.google.com:19302"}},
			},
		}

		// 2. create a new PeerConnection
		peerConnection, err := webrtc.NewPeerConnection(config)

		if err != nil {
			return err

		}

		peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
			dc.OnMessage(func(msg webrtc.DataChannelMessage) {
				fmt.Println("Received from client:", string(msg.Data))
				sendString := "pong:)" + string(msg.Data)
				dc.SendText(sendString)
			})
		})

		peerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
			log.Printf("Peer Connection State has changed: %s\n", state.String())
			switch state {
			case webrtc.PeerConnectionStateClosed:
				{
					err := peerConnection.GracefulClose()
					if err != nil {
						log.Printf("%s error ocurred", err)
					}
					peerConnections.Delete(peerConnection.ID())
					log.Printf("Peer (%s) Connection State has changed: %s\n", peerConnection.ID(), state.String())

				}
			case webrtc.PeerConnectionStateConnected:
				{
					log.Printf("Peer (%s) Connection State has changed: %s\n", peerConnection.ID(), state.String())

				}
			case webrtc.PeerConnectionStateConnecting:
				{
					log.Printf("Peer (%s) Connection State has changed: %s\n", peerConnection.ID(), state.String())

				}
			case webrtc.PeerConnectionStateFailed:
				{
					log.Printf("Peer (%s) Connection State has changed: %s\n", peerConnection.ID(), state.String())

				}
			case webrtc.PeerConnectionStateDisconnected:
				{
					log.Printf("Peer (%s) Connection State has changed: %s\n", peerConnection.ID(), state.String())

				}

			}

		})

		// 3. set the remote SessionDescription
		err = peerConnection.SetRemoteDescription(*sdpOffer)
		if err != nil {
			return err

		}

		// 4. create an answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			return err

		}

		// 5. sets the LocalDescription and starts gather ICE candidates
		err = peerConnection.SetLocalDescription(answer)

		if err != nil {
			return err

		}

		gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
		select {
		case <-gatherComplete:
			log.Println("ICE gathering completed")
		case <-time.After(10 * time.Second):
			log.Println("ICE gathering timed out")
			return err
		}

		peerConnections.Store(peerConnection.ID(), peerConnection)

		data := map[string]interface{}{
			"status":  200,
			"message": "handshake received, sending back answer",
			"peerId":  peerConnection.ID(),
			"answer":  peerConnection.LocalDescription(),
		}

		return c.JSON(data)
	})

	fmt.Println("listening on http://localhost:3000/")

	// start the server
	log.Fatal(app.Listen(":3000"))

}
