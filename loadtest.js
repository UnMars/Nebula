import ws from 'k6/ws';
import { check } from 'k6';
import { Trend, Rate } from 'k6/metrics';
import { fail } from 'k6';


// === Metrics ===
const broadcastLatency = new Trend('broadcast_latency');
const connectionErrors = new Rate('connection_errors');

// === Configuration ===
export const options = {
  stages: [
    { duration: '30s', target: 100 },  // Ramp up to 100
    { duration: '1m', target: 500 },   // Ramp up to 500
    { duration: '1m', target: 1000 },  // Ramp up to 1000
    { duration: '1m', target: 1000 },  // Stay at 1000
    { duration: '30s', target: 0 },    // Ramp down
  ],
  thresholds: {
    'broadcast_latency': [{
        threshold : 'p(95)<500',
        abortOnFail : true
    }], // Latency should be < 500ms
    'connection_errors': [{
        threshold : 'rate<0.01',
        abortOnFail : true
    }], // Errors < 1%
  },
};

const ROOM = "general";
const BASE_URL = __ENV.SERVER_URL || "ws://localhost:8080/ws";

export default function () {
  const username = `user_${__VU}`;
  const url = `${BASE_URL}?room=${ROOM}&username=${username}`;
  const params = { tags: { my_tag: 'nebula_ws' } };

  const res = ws.connect(url, params, function (socket) {
    
    socket.on('open', function open() {
      // Send a message every 1s
      socket.setInterval(function timeout() {
        const payload = {
          type: "message",
          room: ROOM,
          sender: username,
          content: `Message from ${username}`,
          sendAt: Date.now() // Timestamp for latency
        };
        
        try {
          socket.send(JSON.stringify(payload));
        } catch (e) {
          console.error(`VU ${__VU} send error: ${e}`);
        }
      }, 1000); 
    });

    socket.on('message', function (msg) {
      try {
        const data = JSON.parse(msg);
        // Calculate latency if sendAt is present
        if (data.sendAt) {
          const latency = Date.now() - data.sendAt;
          broadcastLatency.add(latency);
        }
      } catch (e) {
        // Ignore parsing errors
      }
    });

    socket.on('error', function (e) {
      connectionErrors.add(1);
      console.error(`VU ${__VU} error: ${e.error()}`);
      fail(`WS error: ${e.error()}`);
    });

    socket.on('close', function () {
      console.log(`VU ${__VU} closed with code ${socket.closeCode} reason=${socket.closeReason}`);
    });

    // Keep connection alive for 30s per iteration
    socket.setTimeout(function () {
      socket.close();
    }, 30000);
  });

  check(res, { 'status is 101': (r) => r && r.status === 101 });
}