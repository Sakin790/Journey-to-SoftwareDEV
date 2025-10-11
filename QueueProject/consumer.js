import { connect } from "amqp-connection-manager";

// RabbitMQ connection info
const RABBITMQ_URL = "amqp://guest:guest@localhost:5672";
const QUEUE_NAME = "terminal_queue";

// Connect to RabbitMQ with auto-reconnect
const connection = connect([RABBITMQ_URL]);

const channelWrapper = connection.createChannel({
  setup: async (channel) => {
    // Ensure durable queue exists
    await channel.assertQueue(QUEUE_NAME, { durable: true });

    // Prefetch 1 message at a time for fair dispatch
    await channel.prefetch(1);

    console.log("👂 Waiting for messages...");

    // Consume messages
    await channel.consume(QUEUE_NAME, async (msg) => {
      if (msg !== null) {
        const content = msg.content.toString();
        console.log("📩 Received:", content);

        // Simulate processing
        await new Promise((res) => setTimeout(res, 500));

        // Acknowledge message
        channel.ack(msg);
      }
    });
  },
});

// Graceful shutdown
process.on("SIGINT", async () => {
  console.log("\n🛑 Closing connection...");
  await connection.close();
  process.exit(0);
});
