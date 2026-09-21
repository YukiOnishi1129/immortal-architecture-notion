import { relations, sql } from "drizzle-orm";
import {
  boolean,
  check,
  index,
  integer,
  pgTable,
  text,
  timestamp,
  uniqueIndex,
  uuid,
} from "drizzle-orm/pg-core";

// Accounts table
export const accounts = pgTable(
  "accounts",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    email: text("email").notNull(),
    firstName: text("first_name").notNull(),
    lastName: text("last_name").notNull(),
    isActive: boolean("is_active").notNull().default(true),
    lastLoginAt: timestamp("last_login_at", { withTimezone: true }),
    provider: text("provider").notNull(),
    providerAccountId: text("provider_account_id").notNull(),
    thumbnail: text("thumbnail"),
    createdAt: timestamp("created_at", { withTimezone: true })
      .notNull()
      .defaultNow(),
    updatedAt: timestamp("updated_at", { withTimezone: true })
      .notNull()
      .defaultNow(),
  },
  (table) => {
    return {
      emailIdx: uniqueIndex("accounts_email_idx").on(table.email),
      providerIdx: uniqueIndex("accounts_provider_idx").on(
        table.provider,
        table.providerAccountId,
      ),
    };
  },
);

// Templates table
export const templates = pgTable("templates", {
  id: uuid("id").primaryKey().defaultRandom(),
  name: text("name").notNull(),
  ownerId: uuid("owner_id")
    .notNull()
    .references(() => accounts.id),
  updatedAt: timestamp("updated_at", { withTimezone: true })
    .notNull()
    .defaultNow(),
});

// Fields table
export const fields = pgTable(
  "fields",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    templateId: uuid("template_id")
      .notNull()
      .references(() => templates.id, { onDelete: "cascade" }),
    label: text("label").notNull(),
    order: integer("order").notNull(),
    isRequired: boolean("is_required").notNull().default(false),
  },
  (table) => {
    return {
      templateOrderIdx: uniqueIndex("fields_template_order_idx").on(
        table.templateId,
        table.order,
      ),
      orderCheck: check("order_check", sql`${table.order} > 0`),
    };
  },
);

// Notes table
export const notes = pgTable(
  "notes",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    title: text("title").notNull(),
    templateId: uuid("template_id")
      .notNull()
      .references(() => templates.id),
    ownerId: uuid("owner_id")
      .notNull()
      .references(() => accounts.id),
    status: text("status").notNull().default("Draft"),
    createdAt: timestamp("created_at", { withTimezone: true })
      .notNull()
      .defaultNow(),
    updatedAt: timestamp("updated_at", { withTimezone: true })
      .notNull()
      .defaultNow(),
  },
  (table) => {
    return {
      ownerIdx: index("notes_owner_idx").on(table.ownerId),
      templateIdx: index("notes_template_idx").on(table.templateId),
      updatedAtIdx: index("notes_updated_at_idx").on(table.updatedAt.desc()),
    };
  },
);

// Sections table
export const sections = pgTable(
  "sections",
  {
    id: uuid("id").primaryKey().defaultRandom(),
    noteId: uuid("note_id")
      .notNull()
      .references(() => notes.id, { onDelete: "cascade" }),
    fieldId: uuid("field_id")
      .notNull()
      .references(() => fields.id),
    content: text("content").notNull().default(""),
  },
  (table) => {
    return {
      noteFieldIdx: uniqueIndex("sections_note_field_idx").on(
        table.noteId,
        table.fieldId,
      ),
      noteIdx: index("sections_note_idx").on(table.noteId),
      fieldIdx: index("sections_field_idx").on(table.fieldId),
    };
  },
);

// Relations
export const accountsRelations = relations(accounts, ({ many }) => ({
  templates: many(templates),
  notes: many(notes),
}));

export const templatesRelations = relations(templates, ({ one, many }) => ({
  owner: one(accounts, {
    fields: [templates.ownerId],
    references: [accounts.id],
  }),
  fields: many(fields),
  notes: many(notes),
}));

export const fieldsRelations = relations(fields, ({ one, many }) => ({
  template: one(templates, {
    fields: [fields.templateId],
    references: [templates.id],
  }),
  sections: many(sections),
}));

export const notesRelations = relations(notes, ({ one, many }) => ({
  owner: one(accounts, {
    fields: [notes.ownerId],
    references: [accounts.id],
  }),
  template: one(templates, {
    fields: [notes.templateId],
    references: [templates.id],
  }),
  sections: many(sections),
}));

export const sectionsRelations = relations(sections, ({ one }) => ({
  note: one(notes, {
    fields: [sections.noteId],
    references: [notes.id],
  }),
  field: one(fields, {
    fields: [sections.fieldId],
    references: [fields.id],
  }),
}));

// Type exports
export type Account = typeof accounts.$inferSelect;
export type NewAccount = typeof accounts.$inferInsert;

export type Template = typeof templates.$inferSelect;
export type NewTemplate = typeof templates.$inferInsert;

export type Field = typeof fields.$inferSelect;
export type NewField = typeof fields.$inferInsert;

export type Note = typeof notes.$inferSelect;
export type NewNote = typeof notes.$inferInsert;

export type Section = typeof sections.$inferSelect;
export type NewSection = typeof sections.$inferInsert;
