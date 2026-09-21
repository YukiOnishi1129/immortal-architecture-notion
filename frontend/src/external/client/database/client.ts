import { drizzle, type NodePgDatabase } from "drizzle-orm/node-postgres";
import { Pool } from "pg";
import * as schema from "./schema";

if (!process.env.DATABASE_URL) {
  throw new Error("DATABASE_URL environment variable is not set");
}

// Singleton pattern for database connection
const globalForDb = globalThis as unknown as {
  pool: Pool | undefined;
  drizzle: NodePgDatabase<typeof schema> | undefined;
};

if (!globalForDb.pool) {
  globalForDb.pool = new Pool({
    connectionString: process.env.DATABASE_URL,
  });
}

if (!globalForDb.drizzle) {
  globalForDb.drizzle = drizzle(globalForDb.pool, { schema });
}

export const db = globalForDb.drizzle;

export type Database = typeof db;
