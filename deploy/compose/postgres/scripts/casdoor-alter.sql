/*
 Navicat Premium Dump SQL

 Source Server         : 测试集群
 Source Server Type    : PostgreSQL
 Source Server Version : 170005 (170005)
 Source Host           : localhost:5432
 Source Catalog        : casdoor
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 170005 (170005)
 File Encoding         : 65001

 Date: 05/08/2025 19:06:11
*/


-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."casbin_api_rule_id_seq"
OWNED BY "public"."casbin_api_rule"."id";
SELECT setval('"public"."casbin_api_rule_id_seq"', 69, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."casbin_rule_id_seq"
OWNED BY "public"."casbin_rule"."id";
SELECT setval('"public"."casbin_rule_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."casbin_user_rule_id_seq"
OWNED BY "public"."casbin_user_rule"."id";
SELECT setval('"public"."casbin_user_rule_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."permission_rule_id_seq"
OWNED BY "public"."permission_rule"."id";
SELECT setval('"public"."permission_rule_id_seq"', 3, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."record_id_seq"
OWNED BY "public"."record"."id";
SELECT setval('"public"."record_id_seq"', 287, true);

-- ----------------------------
-- Primary Key structure for table adapter
-- ----------------------------
ALTER TABLE "public"."adapter" ADD CONSTRAINT "adapter_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table application
-- ----------------------------
ALTER TABLE "public"."application" ADD CONSTRAINT "application_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Indexes structure for table casbin_api_rule
-- ----------------------------
CREATE INDEX "IDX_casbin_api_rule_ptype" ON "public"."casbin_api_rule" USING btree (
  "ptype" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_api_rule_v0" ON "public"."casbin_api_rule" USING btree (
  "v0" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_api_rule_v1" ON "public"."casbin_api_rule" USING btree (
  "v1" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_api_rule_v2" ON "public"."casbin_api_rule" USING btree (
  "v2" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_api_rule_v3" ON "public"."casbin_api_rule" USING btree (
  "v3" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_api_rule_v4" ON "public"."casbin_api_rule" USING btree (
  "v4" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_api_rule_v5" ON "public"."casbin_api_rule" USING btree (
  "v5" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table casbin_api_rule
-- ----------------------------
ALTER TABLE "public"."casbin_api_rule" ADD CONSTRAINT "casbin_api_rule_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table casbin_rule
-- ----------------------------
CREATE INDEX "IDX_casbin_rule_ptype" ON "public"."casbin_rule" USING btree (
  "ptype" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_rule_v0" ON "public"."casbin_rule" USING btree (
  "v0" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_rule_v1" ON "public"."casbin_rule" USING btree (
  "v1" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_rule_v2" ON "public"."casbin_rule" USING btree (
  "v2" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_rule_v3" ON "public"."casbin_rule" USING btree (
  "v3" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_rule_v4" ON "public"."casbin_rule" USING btree (
  "v4" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_rule_v5" ON "public"."casbin_rule" USING btree (
  "v5" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table casbin_rule
-- ----------------------------
ALTER TABLE "public"."casbin_rule" ADD CONSTRAINT "casbin_rule_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table casbin_user_rule
-- ----------------------------
CREATE INDEX "IDX_casbin_user_rule_ptype" ON "public"."casbin_user_rule" USING btree (
  "ptype" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_user_rule_v0" ON "public"."casbin_user_rule" USING btree (
  "v0" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_user_rule_v1" ON "public"."casbin_user_rule" USING btree (
  "v1" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_user_rule_v2" ON "public"."casbin_user_rule" USING btree (
  "v2" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_user_rule_v3" ON "public"."casbin_user_rule" USING btree (
  "v3" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_user_rule_v4" ON "public"."casbin_user_rule" USING btree (
  "v4" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_casbin_user_rule_v5" ON "public"."casbin_user_rule" USING btree (
  "v5" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table casbin_user_rule
-- ----------------------------
ALTER TABLE "public"."casbin_user_rule" ADD CONSTRAINT "casbin_user_rule_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table cert
-- ----------------------------
ALTER TABLE "public"."cert" ADD CONSTRAINT "cert_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table enforcer
-- ----------------------------
ALTER TABLE "public"."enforcer" ADD CONSTRAINT "enforcer_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Indexes structure for table group
-- ----------------------------
CREATE UNIQUE INDEX "UQE_group_name" ON "public"."group" USING btree (
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table group
-- ----------------------------
ALTER TABLE "public"."group" ADD CONSTRAINT "group_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Indexes structure for table invitation
-- ----------------------------
CREATE INDEX "IDX_invitation_code" ON "public"."invitation" USING btree (
  "code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table invitation
-- ----------------------------
ALTER TABLE "public"."invitation" ADD CONSTRAINT "invitation_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table ldap
-- ----------------------------
ALTER TABLE "public"."ldap" ADD CONSTRAINT "ldap_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table model
-- ----------------------------
ALTER TABLE "public"."model" ADD CONSTRAINT "model_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table organization
-- ----------------------------
ALTER TABLE "public"."organization" ADD CONSTRAINT "organization_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table payment
-- ----------------------------
ALTER TABLE "public"."payment" ADD CONSTRAINT "payment_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table permission
-- ----------------------------
ALTER TABLE "public"."permission" ADD CONSTRAINT "permission_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Indexes structure for table permission_rule
-- ----------------------------
CREATE INDEX "IDX_permission_rule_ptype" ON "public"."permission_rule" USING btree (
  "ptype" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_permission_rule_v0" ON "public"."permission_rule" USING btree (
  "v0" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_permission_rule_v1" ON "public"."permission_rule" USING btree (
  "v1" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_permission_rule_v2" ON "public"."permission_rule" USING btree (
  "v2" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_permission_rule_v3" ON "public"."permission_rule" USING btree (
  "v3" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_permission_rule_v4" ON "public"."permission_rule" USING btree (
  "v4" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_permission_rule_v5" ON "public"."permission_rule" USING btree (
  "v5" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table permission_rule
-- ----------------------------
ALTER TABLE "public"."permission_rule" ADD CONSTRAINT "permission_rule_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table plan
-- ----------------------------
ALTER TABLE "public"."plan" ADD CONSTRAINT "plan_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table pricing
-- ----------------------------
ALTER TABLE "public"."pricing" ADD CONSTRAINT "pricing_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table product
-- ----------------------------
ALTER TABLE "public"."product" ADD CONSTRAINT "product_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Indexes structure for table provider
-- ----------------------------
CREATE UNIQUE INDEX "UQE_provider_name" ON "public"."provider" USING btree (
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table provider
-- ----------------------------
ALTER TABLE "public"."provider" ADD CONSTRAINT "provider_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Indexes structure for table radius_accounting
-- ----------------------------
CREATE INDEX "IDX_radius_accounting_acct_session_id" ON "public"."radius_accounting" USING btree (
  "acct_session_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_radius_accounting_acct_start_time" ON "public"."radius_accounting" USING btree (
  "acct_start_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_radius_accounting_acct_stop_time" ON "public"."radius_accounting" USING btree (
  "acct_stop_time" "pg_catalog"."timestamp_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_radius_accounting_username" ON "public"."radius_accounting" USING btree (
  "username" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table radius_accounting
-- ----------------------------
ALTER TABLE "public"."radius_accounting" ADD CONSTRAINT "radius_accounting_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Indexes structure for table record
-- ----------------------------
CREATE INDEX "IDX_record_name" ON "public"."record" USING btree (
  "name" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_record_owner" ON "public"."record" USING btree (
  "owner" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table record
-- ----------------------------
ALTER TABLE "public"."record" ADD CONSTRAINT "record_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table resource
-- ----------------------------
ALTER TABLE "public"."resource" ADD CONSTRAINT "resource_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table role
-- ----------------------------
ALTER TABLE "public"."role" ADD CONSTRAINT "role_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table session
-- ----------------------------
ALTER TABLE "public"."session" ADD CONSTRAINT "session_pkey" PRIMARY KEY ("owner", "name", "application");

-- ----------------------------
-- Primary Key structure for table subscription
-- ----------------------------
ALTER TABLE "public"."subscription" ADD CONSTRAINT "subscription_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table syncer
-- ----------------------------
ALTER TABLE "public"."syncer" ADD CONSTRAINT "syncer_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Indexes structure for table token
-- ----------------------------
CREATE INDEX "IDX_token_access_token_hash" ON "public"."token" USING btree (
  "access_token_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_token_code" ON "public"."token" USING btree (
  "code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_token_refresh_token_hash" ON "public"."token" USING btree (
  "refresh_token_hash" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table token
-- ----------------------------
ALTER TABLE "public"."token" ADD CONSTRAINT "token_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table transaction
-- ----------------------------
ALTER TABLE "public"."transaction" ADD CONSTRAINT "transaction_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Indexes structure for table user
-- ----------------------------
CREATE INDEX "IDX_user_created_time" ON "public"."user" USING btree (
  "created_time" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_user_email" ON "public"."user" USING btree (
  "email" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_user_external_id" ON "public"."user" USING btree (
  "external_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_user_id" ON "public"."user" USING btree (
  "id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_user_id_card" ON "public"."user" USING btree (
  "id_card" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_user_invitation" ON "public"."user" USING btree (
  "invitation" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_user_invitation_code" ON "public"."user" USING btree (
  "invitation_code" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_user_phone" ON "public"."user" USING btree (
  "phone" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "IDX_user_universal_id" ON "public"."user" USING btree (
  "universal_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table user
-- ----------------------------
ALTER TABLE "public"."user" ADD CONSTRAINT "user_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Primary Key structure for table user_identity_binding
-- ----------------------------
ALTER TABLE "public"."user_identity_binding" ADD CONSTRAINT "user_identity_binding_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table verification_record
-- ----------------------------
CREATE INDEX "IDX_verification_record_receiver" ON "public"."verification_record" USING btree (
  "receiver" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table verification_record
-- ----------------------------
ALTER TABLE "public"."verification_record" ADD CONSTRAINT "verification_record_pkey" PRIMARY KEY ("owner", "name");

-- ----------------------------
-- Indexes structure for table webhook
-- ----------------------------
CREATE INDEX "IDX_webhook_organization" ON "public"."webhook" USING btree (
  "organization" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table webhook
-- ----------------------------
ALTER TABLE "public"."webhook" ADD CONSTRAINT "webhook_pkey" PRIMARY KEY ("owner", "name");
