ALTER TABLE "accounts" RENAME COLUMN "name" TO "first_name";--> statement-breakpoint
ALTER TABLE "accounts" ADD COLUMN "last_name" text NOT NULL;--> statement-breakpoint
ALTER TABLE "accounts" ADD COLUMN "is_active" boolean DEFAULT true NOT NULL;--> statement-breakpoint
ALTER TABLE "accounts" ADD COLUMN "last_login_at" timestamp with time zone;--> statement-breakpoint
ALTER TABLE "accounts" ADD COLUMN "provider" text NOT NULL;--> statement-breakpoint
ALTER TABLE "accounts" ADD COLUMN "provider_account_id" text NOT NULL;--> statement-breakpoint
ALTER TABLE "accounts" ADD COLUMN "thumbnail" text;--> statement-breakpoint
ALTER TABLE "accounts" ADD COLUMN "updated_at" timestamp with time zone DEFAULT now() NOT NULL;--> statement-breakpoint
CREATE UNIQUE INDEX "accounts_provider_idx" ON "accounts" USING btree ("provider","provider_account_id");--> statement-breakpoint
ALTER TABLE "accounts" DROP COLUMN "auth_id";--> statement-breakpoint
ALTER TABLE "fields" ADD CONSTRAINT "order_check" CHECK ("fields"."order" > 0);