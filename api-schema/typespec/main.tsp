import "@typespec/http";
import "@typespec/openapi3";
import "./models/common.tsp";
import "./models/account.tsp";
import "./models/template.tsp";
import "./models/note.tsp";
import "./routes/accounts.tsp";
import "./routes/templates.tsp";
import "./routes/notes.tsp";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

@service({
  title: "Mini Notion API",
})
@server("http://localhost:8080", "Development server")
@server("https://api.mini-notion.com", "Production server")
namespace MiniNotion;
