CREATE TABLE "users" (
    "id" serial NOT NULL,
    "email" character varying(255),
    "password" character varying(255),
    "created_at" timestamp,
    "updated_at" timestamp,
    "deleted_at" timestamp,
    CONSTRAINT "users_email" UNIQUE ("email"),
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
