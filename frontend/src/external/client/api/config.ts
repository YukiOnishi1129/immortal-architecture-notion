import { AccountsApi, NotesApi, TemplatesApi } from "./generated/apis";
import { Configuration } from "./generated/runtime";

const defaultBaseUrl = "http://localhost:8080";

const basePath = (process.env.API_BASE_URL || defaultBaseUrl).replace(
  /\/+$/,
  "",
);

const configuration = new Configuration({
  basePath,
});

export const accountsApiClient = new AccountsApi(configuration);
export const notesApiClient = new NotesApi(configuration);
export const templatesApiClient = new TemplatesApi(configuration);
