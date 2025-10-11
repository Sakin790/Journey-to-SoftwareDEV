import { connect } from "amqp-connection-manager";

// RabbitMQ connection info
const RABBITMQ_URL = "amqp://guest:guest@localhost:5672";
const QUEUE_NAME = "terminal_queue";

// Connect to RabbitMQ with auto-reconnect
const connection = connect([RABBITMQ_URL]);

// Create a channel
const channelWrapper = connection.createChannel({
  setup: async (channel) => {
    // Ensure durable queue exists
    await channel.assertQueue(QUEUE_NAME, { durable: true });
  },
});

console.log("Type something to send to RabbitMQ and press Enter:");

// Listen to terminal input
process.stdin.on("data", async (data) => {
  const message = data.toString().trim();
  if (!message) return; // ignore empty input

  try {
    // Send persistent message to queue
    await channelWrapper.sendToQueue(QUEUE_NAME, Buffer.from(message), { persistent: true });
    console.log(`✅ Sent: ${message}`);
  } catch (err) {
    console.error("❌ Failed to send message:", err);
  }
});
