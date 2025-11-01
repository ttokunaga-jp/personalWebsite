import { describe, expect, it } from "vitest";

import { extractTokenFromHash } from "./auth-session";

describe("extractTokenFromHash", () => {
  it("returns token when present in hash fragment", () => {
    expect(extractTokenFromHash("#token=abc123&state=xyz")).toBe("abc123");
  });

  it("ignores leading hashes and whitespace", () => {
    expect(extractTokenFromHash("##token=  secret ")).toBe("secret");
  });

  it("returns null when token parameter missing or blank", () => {
    expect(extractTokenFromHash("#state=123")).toBeNull();
    expect(extractTokenFromHash("#token=")).toBeNull();
    expect(extractTokenFromHash("")).toBeNull();
  });
});
