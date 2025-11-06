import { expect, test } from "@playwright/test";

test.describe("Books flow", () => {
  test("renders list and creates a book", async ({ page }) => {
    const books: any[] = [];

    await page.route("**/books", async (route, request) => {
      if (request.method() === "GET") {
        await route.fulfill({
          status: 200,
          contentType: "application/json",
          body: JSON.stringify({ data: books }),
        });
        return;
      }

      if (request.method() === "POST") {
        const payload = JSON.parse(request.postData() ?? "{}");
        const newBook = {
          id: `test-${books.length + 1}`,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          ...payload,
        };
        books.push(newBook);
        await route.fulfill({
          status: 201,
          contentType: "application/json",
          body: JSON.stringify({ data: newBook }),
        });
        return;
      }

      await route.fallback();
    });

    await page.goto("/");
    await expect(page.getByText(/No books yet/i)).toBeVisible();

    await page.getByRole("link", { name: /add book/i }).click();
    await page.getByLabel("Title").fill("Refactoring");
    await page.getByLabel("Author").fill("Martin Fowler");
    await page.getByLabel("Price").fill("42");
    await page.getByLabel("Stock").fill("10");
    await page.getByLabel("Currency").fill("USD");
    await page.getByRole("button", { name: /save/i }).click();

    await expect(page).toHaveURL(/admin\/books/);
  });
});
