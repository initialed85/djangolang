import createClient from "openapi-fetch";
import { createHooks } from "swr-openapi";
import type * as Schema from "./schema";

export const djangolangApi = createClient<Schema.paths>({
  baseUrl: "http://localhost:7070",
});

export const { use: useDjangolang } = createHooks(
  djangolangApi,
  "djangolang-api",
);
