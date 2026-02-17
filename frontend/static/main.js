const peerConnection = new RTCPeerConnection({
  iceServers: [{ urls: "stun:stun.1.google.com:19302" }],
});

// example for a data channel
const dataChannel = peerConnection.createDataChannel("dataChannel");

dataChannel.onmessage = (event) => console.log("Received:", event.data);

dataChannel.onopen = () => dataChannel.send("Hello from the browser");

const offer = await peerConnection.createOffer();

await peerConnection.setLocalDescription(offer);

let iceCandidates = [];

peerConnection.onconnectionstatechange = (event) =>
  console.log("connection state: ", event);

peerConnection.onicecandidate = async (event) => {
  if (event.candidate) {
    iceCandidates.push(event.candidate.toJSON());
  } else {
    // ALL candidates gathered
    console.log("SDP offer ready for signaling:");

    const offerData = {
      sdp: peerConnection.localDescription.sdp,
      type: peerConnection.localDescription.type,
      candidates: iceCandidates,
    };

    console.dir(offerData);

    // send offerData to your go pion application via your signaling server
    let answer = await sendToServer(offerData);

    const remoteDesc = new RTCSessionDescription(answer.answer);

    await peerConnection.setRemoteDescription(remoteDesc);
  }
};

// 1. listen for the moment go starts sending the track
peerConnection.ontrack = (event) => {
  // 2. hotwire the incoming stream directory to your html video tag
  const videoElement = document.getElementById("my-video");
  videoElement.srcObject = event.streams[0];
};

const sendToServer = async (data) => {
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
