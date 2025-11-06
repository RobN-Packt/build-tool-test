import type { SQSEvent, SQSHandler } from "aws-lambda";
import { z } from "zod";

const purchaseSchema = z.object({
  bookId: z.string().uuid(),
  quantity: z.number().int().positive(),
  customerId: z.string().min(1).optional(),
});

export const handler: SQSHandler = async (event: SQSEvent) => {
  for (const record of event.Records) {
    const body = JSON.parse(record.body);
    const validation = purchaseSchema.safeParse(body);
    if (!validation.success) {
      console.error("Invalid purchase payload", validation.error.format());
      continue;
    }
    const purchase = validation.data;
    console.info("Processing purchase", {
      bookId: purchase.bookId,
      quantity: purchase.quantity,
      customerId: purchase.customerId,
      messageId: record.messageId,
    });
    // In a real system we'd update inventory or call the API here.
  }
};
