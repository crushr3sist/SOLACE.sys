class SOLACE_sys {
  constructor() {
    this.peerConnection = new RTCPeerConnection({
      iceServers: [{ urls: "stun:stun.1.google.com:19302" }],
    });
    this.iceCandidates = [];
    this.onStateChangeCallback = null;
    this.onNewMessageCallback = null;

    this.dataChannel = this.peerConnection.createDataChannel("dataChannel");
    this.dataChannel.onopen = () =>
      this.dataChannel.send("Hello from the browser");

    this.peerConnection.onicecandidate = async (event) => {
      if (event.candidate) {
        this.iceCandidates.push(event.candidate.toJSON());
      } else {
        // ALL candidates gathered
        console.log("SDP offer ready for signaling:");

        const offerData = {
          sdp: this.peerConnection.localDescription.sdp,
          type: this.peerConnection.localDescription.type,
          candidates: this.iceCandidates,
        };

        console.dir(offerData);

        // send offerData to your go pion application via your signaling server
        let answer = await this.sendToServer(offerData);

        const remoteDesc = new RTCSessionDescription(answer.answer);

        await this.peerConnection.setRemoteDescription(remoteDesc);
      }
    };

    this.dataChannel.onmessage = (event) => {
      const newMessage = event.data;
      console.log(newMessage);
      if (this.onNewMessageCallback) {
        this.onNewMessageCallback(newMessage);
      }
    };

    this.peerConnection.onconnectionstatechange = () => {
      const newState = this.peerConnection.connectionState;
      if (this.onStateChangeCallback) {
        this.onStateChangeCallback(newState);
      }
    };
  }

  async init() {
    const offer = await this.peerConnection.createOffer();
    await this.peerConnection.setLocalDescription(offer);
  }

  onStateChange(callback) {
    this.onStateChangeCallback = callback;
  }
  onMessageChange(callback) {
    this.onNewMessageCallback = callback;
  }

  get() {
    let message = "";
    return message;
  }

  send(message) {
    this.dataChannel.send(message);
  }

  sendToServer = async (data) => {
    try {
      const response = await fetch("/webrtc/offer", {
        method: "POST",
        headers: {
          "Content-Type": "application/json; charset=UTF-8",
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        throw new Error("HTTP error! status: " + response.status);
      }
      const responseData = await response.json();
      console.log("success: ", responseData);
      return responseData;
    } catch (error) {
      console.error("Error: ", error);
    }
  };
}
const connectionInstance = new SOLACE_sys();

let state = document.getElementById("state");
let sendButton = document.getElementById("send_button");
let messageInput = document.getElementById("message");
let messagesList = document.getElementById("messages");

connectionInstance.onStateChange((newState) => {
  console.log("State changed to: ", newState);

  state.innerText = newState;

  if (newState === "connected") {
    state.style.color = "green";
  } else if (newState === "failed" || newState === "disconnected") {
    state.style.color = "red";
  } else {
    state.style.color = "orange"; // connecting, new, etc.
  }
});

connectionInstance.onMessageChange((newMessage) => {
  // Create the new list item
  let newListElem = document.createElement("li");
  newListElem.textContent = newMessage;

  // Append it to the list (preserves previous items)
  messagesList.appendChild(newListElem);
});

sendButton.addEventListener("click", () => {
  console.log(messageInput.value);
  connectionInstance.send(messageInput.value);
});

connectionInstance.init();

// // 1. listen for the moment go starts sending the track
// peerConnection.ontrack = (event) => {
//   // 2. hotwire the incoming stream directory to your html video tag
//   const videoElement = document.getElementById("my-video");
//   videoElement.srcObject = event.streams[0];
// };
