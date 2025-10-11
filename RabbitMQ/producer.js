// Import the connection manager library for RabbitMQ (handles auto-reconnect)
import { connect } from "amqp-connection-manager";

// Import dotenv to load environment variables from .env file
import dotenv from "dotenv";

// Load variables from .env file into process.env
dotenv.config();

// Get RabbitMQ URL and queue name from environment variables
const RABBITMQ_URL = process.env.RABBITMQ_URL;
const QUEUE_NAME = process.env.QUEUE_NAME;

// Function to send a message to the queue
const sendMessage = async (message) => {
  // 1️⃣ Connect to RabbitMQ with automatic reconnect
  const connection = connect([RABBITMQ_URL]);

  // 2️⃣ Create a channel wrapper (channel abstraction over connection)
  const channelWrapper = connection.createChannel({
    json: false, // we are sending plain text, not JSON
    // Setup function runs when channel is created (or recreated on reconnect)
    setup: async (channel) => {
      // Ensure the queue exists and is durable (survives RabbitMQ restart)
      await channel.assertQueue(QUEUE_NAME, { durable: true });
    },
  });

  // 3️⃣ Send a message to the queue
  try {
    await channelWrapper.sendToQueue(
      QUEUE_NAME,                // queue name
      Buffer.from(message),      // convert message to buffer
      { persistent: true }       // make message persistent (survives RabbitMQ restart)
    );
    console.log(`✅ Sent: ${message}`);
  } catch (err) {
    // Handle any error that occurs while sending message
    console.error("❌ Failed to send message:", err);
  } finally {
    // 4️⃣ Close connection gracefully after a small delay
    // Delay ensures the message is sent before closing
    setTimeout(() => connection.close(), 500);
  }
};

// Example usage of the sendMessage function
sendMessage("Hello from RabbitMQ 🐇 (production-ready)");
