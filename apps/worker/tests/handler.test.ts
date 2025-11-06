import { describe, expect, it, vi } from "vitest";

import { handler } from "../src/handler";

describe("handler", () => {
  it("processes valid messages", async () => {
    const info = vi.spyOn(console, "info").mockImplementation(() => undefined);
    const error = vi.spyOn(console, "error").mockImplementation(() => undefined);

    await handler({
      Records: [
        {
          messageId: "1",
          body: JSON.stringify({ bookId: "3fa85f64-5717-4562-b3fc-2c963f66afa6", quantity: 2, customerId: "c1" }),
        } as any,
      ],
    } as any);

    expect(info).toHaveBeenCalledWith("Processing purchase", expect.objectContaining({ quantity: 2 }));
    expect(error).not.toHaveBeenCalled();
    info.mockRestore();
    error.mockRestore();
  });

  it("skips invalid messages", async () => {
    const error = vi.spyOn(console, "error").mockImplementation(() => undefined);

    await handler({
      Records: [
        {
          messageId: "1",
          body: JSON.stringify({ bookId: "not-uuid", quantity: 0 }),
        } as any,
      ],
    } as any);

    expect(error).toHaveBeenCalled();
    error.mockRestore();
  });
});
