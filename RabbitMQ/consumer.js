// Import the connection manager library for RabbitMQ (supports auto-reconnect)
import { connect } from "amqp-connection-manager";

// Import dotenv to load environment variables from .env file
import dotenv from "dotenv";

// Load variables from .env file into process.env
dotenv.config();

// Get RabbitMQ URL and queue name from environment variables
const RABBITMQ_URL = process.env.RABBITMQ_URL;
const QUEUE_NAME = process.env.QUEUE_NAME;

// Connect to RabbitMQ with auto-reconnect support
const connection = connect([RABBITMQ_URL]);

// Create a channel wrapper (channel abstraction over connection)
const channelWrapper = connection.createChannel({
  // Setup function runs when channel is created (or re-created on reconnect)
  setup: async (channel) => {
    // Ensure the queue exists and is durable (survives RabbitMQ restarts)
    await channel.assertQueue(QUEUE_NAME, { durable: true });

    // Limit the number of unacknowledged messages per consumer to 1
    // Ensures fair dispatch among multiple consumers
    await channel.prefetch(1);

    // Log to console that consumer is ready
    console.log("👂 Waiting for messages...");

    // Start consuming messages from the queue
    await channel.consume(QUEUE_NAME, async (msg) => {
      // Check if message exists (sometimes msg can be null)
      if (msg !== null) {
        // Convert message buffer to string
        const content = msg.content.toString();

        // Log the received message
        console.log("📩 Received:", content);

        // Simulate some processing (e.g., saving to DB, sending email)
        await new Promise((res) => setTimeout(res, 1000));

        // Acknowledge to RabbitMQ that message was processed successfully
        channel.ack(msg);
      }
    });
  },
});

// Graceful shutdown: handle Ctrl+C or SIGINT
process.on("SIGINT", async () => {
  // Log shutdown message
  console.log("\n🛑 Closing connection...");
  // Close connection to RabbitMQ
  await connection.close();
  // Exit the process
  process.exit(0);
});
