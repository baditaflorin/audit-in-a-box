import { describe, expect, it } from "vitest";
import { trimSlash } from "../../api/client";

describe("trimSlash", () => {
  it("normalizes backend urls", () => {
    expect(trimSlash("http://localhost:25342///")).toBe(
      "http://localhost:25342",
    );
  });
});
