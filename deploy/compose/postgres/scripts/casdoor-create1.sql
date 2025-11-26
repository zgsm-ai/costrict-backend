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
-- Sequence structure for casbin_api_rule_id_seq
-- ----------------------------
CREATE SEQUENCE IF NOT EXISTS "public"."casbin_api_rule_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for casbin_rule_id_seq
-- ----------------------------
CREATE SEQUENCE IF NOT EXISTS "public"."casbin_rule_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for casbin_user_rule_id_seq
-- ----------------------------
CREATE SEQUENCE IF NOT EXISTS "public"."casbin_user_rule_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for permission_rule_id_seq
-- ----------------------------
CREATE SEQUENCE IF NOT EXISTS "public"."permission_rule_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for record_id_seq
-- ----------------------------
CREATE SEQUENCE IF NOT EXISTS "public"."record_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Table structure for adapter
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."adapter" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "table" varchar(100) COLLATE "pg_catalog"."default",
  "use_same_db" bool,
  "type" varchar(100) COLLATE "pg_catalog"."default",
  "database_type" varchar(100) COLLATE "pg_catalog"."default",
  "host" varchar(100) COLLATE "pg_catalog"."default",
  "port" int4,
  "user" varchar(100) COLLATE "pg_catalog"."default",
  "password" varchar(150) COLLATE "pg_catalog"."default",
  "database" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for application
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."application" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "logo" varchar(200) COLLATE "pg_catalog"."default",
  "homepage_url" varchar(100) COLLATE "pg_catalog"."default",
  "description" varchar(100) COLLATE "pg_catalog"."default",
  "organization" varchar(100) COLLATE "pg_catalog"."default",
  "cert" varchar(100) COLLATE "pg_catalog"."default",
  "default_group" varchar(100) COLLATE "pg_catalog"."default",
  "header_html" text COLLATE "pg_catalog"."default",
  "enable_password" bool,
  "enable_sign_up" bool,
  "enable_signin_session" bool,
  "enable_auto_signin" bool,
  "enable_code_signin" bool,
  "enable_saml_compress" bool,
  "enable_saml_c14n10" bool,
  "enable_saml_post_binding" bool,
  "use_email_as_saml_name_id" bool,
  "enable_web_authn" bool,
  "enable_link_with_email" bool,
  "org_choice_mode" varchar(255) COLLATE "pg_catalog"."default",
  "saml_reply_url" varchar(500) COLLATE "pg_catalog"."default",
  "providers" text COLLATE "pg_catalog"."default",
  "signin_methods" varchar(2000) COLLATE "pg_catalog"."default",
  "signup_items" varchar(3000) COLLATE "pg_catalog"."default",
  "signin_items" text COLLATE "pg_catalog"."default",
  "grant_types" varchar(1000) COLLATE "pg_catalog"."default",
  "tags" text COLLATE "pg_catalog"."default",
  "saml_attributes" varchar(1000) COLLATE "pg_catalog"."default",
  "is_shared" bool,
  "ip_restriction" varchar(255) COLLATE "pg_catalog"."default",
  "client_id" varchar(100) COLLATE "pg_catalog"."default",
  "client_secret" varchar(100) COLLATE "pg_catalog"."default",
  "redirect_uris" varchar(1000) COLLATE "pg_catalog"."default",
  "forced_redirect_origin" varchar(100) COLLATE "pg_catalog"."default",
  "token_format" varchar(100) COLLATE "pg_catalog"."default",
  "token_signing_method" varchar(100) COLLATE "pg_catalog"."default",
  "token_fields" varchar(1000) COLLATE "pg_catalog"."default",
  "expire_in_hours" int4,
  "refresh_expire_in_hours" int4,
  "signup_url" varchar(200) COLLATE "pg_catalog"."default",
  "signin_url" varchar(200) COLLATE "pg_catalog"."default",
  "forget_url" varchar(200) COLLATE "pg_catalog"."default",
  "affiliation_url" varchar(100) COLLATE "pg_catalog"."default",
  "ip_whitelist" varchar(200) COLLATE "pg_catalog"."default",
  "terms_of_use" varchar(100) COLLATE "pg_catalog"."default",
  "signup_html" text COLLATE "pg_catalog"."default",
  "signin_html" text COLLATE "pg_catalog"."default",
  "theme_data" json,
  "footer_html" text COLLATE "pg_catalog"."default",
  "form_css" text COLLATE "pg_catalog"."default",
  "form_css_mobile" text COLLATE "pg_catalog"."default",
  "form_offset" int4,
  "form_side_html" text COLLATE "pg_catalog"."default",
  "form_background_url" varchar(200) COLLATE "pg_catalog"."default",
  "form_background_url_mobile" varchar(200) COLLATE "pg_catalog"."default",
  "failed_signin_limit" int4,
  "failed_signin_frozen_time" int4
)
;

-- ----------------------------
-- Table structure for casbin_api_rule
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."casbin_api_rule" (
  "id" int8 NOT NULL DEFAULT nextval('casbin_api_rule_id_seq'::regclass),
  "ptype" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v0" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v1" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v2" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v3" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v4" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v5" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."casbin_rule" (
  "id" int8 NOT NULL DEFAULT nextval('casbin_rule_id_seq'::regclass),
  "ptype" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v0" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v1" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v2" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v3" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v4" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v5" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;

-- ----------------------------
-- Table structure for casbin_user_rule
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."casbin_user_rule" (
  "id" int8 NOT NULL DEFAULT nextval('casbin_user_rule_id_seq'::regclass),
  "ptype" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v0" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v1" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v2" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v3" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v4" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v5" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;

-- ----------------------------
-- Table structure for cert
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."cert" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "scope" varchar(100) COLLATE "pg_catalog"."default",
  "type" varchar(100) COLLATE "pg_catalog"."default",
  "crypto_algorithm" varchar(100) COLLATE "pg_catalog"."default",
  "bit_size" int4,
  "expire_in_years" int4,
  "certificate" text COLLATE "pg_catalog"."default",
  "private_key" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for enforcer
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."enforcer" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "updated_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "description" varchar(100) COLLATE "pg_catalog"."default",
  "model" varchar(100) COLLATE "pg_catalog"."default",
  "adapter" varchar(100) COLLATE "pg_catalog"."default",
  "enforcer" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for group
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."group" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "updated_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "manager" varchar(100) COLLATE "pg_catalog"."default",
  "contact_email" varchar(100) COLLATE "pg_catalog"."default",
  "type" varchar(100) COLLATE "pg_catalog"."default",
  "parent_id" varchar(100) COLLATE "pg_catalog"."default",
  "is_top_group" bool,
  "title" varchar(255) COLLATE "pg_catalog"."default",
  "key" varchar(255) COLLATE "pg_catalog"."default",
  "children" text COLLATE "pg_catalog"."default",
  "is_enabled" bool
)
;

-- ----------------------------
-- Table structure for invitation
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."invitation" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "updated_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "code" varchar(100) COLLATE "pg_catalog"."default",
  "is_regexp" bool,
  "quota" int4,
  "used_count" int4,
  "application" varchar(100) COLLATE "pg_catalog"."default",
  "username" varchar(100) COLLATE "pg_catalog"."default",
  "email" varchar(100) COLLATE "pg_catalog"."default",
  "phone" varchar(100) COLLATE "pg_catalog"."default",
  "signup_group" varchar(100) COLLATE "pg_catalog"."default",
  "default_code" varchar(100) COLLATE "pg_catalog"."default",
  "state" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for ldap
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."ldap" (
  "id" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "owner" varchar(100) COLLATE "pg_catalog"."default",
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "server_name" varchar(100) COLLATE "pg_catalog"."default",
  "host" varchar(100) COLLATE "pg_catalog"."default",
  "port" int4,
  "enable_ssl" bool,
  "allow_self_signed_cert" bool,
  "username" varchar(100) COLLATE "pg_catalog"."default",
  "password" varchar(100) COLLATE "pg_catalog"."default",
  "base_dn" varchar(100) COLLATE "pg_catalog"."default",
  "filter" varchar(200) COLLATE "pg_catalog"."default",
  "filter_fields" varchar(100) COLLATE "pg_catalog"."default",
  "default_group" varchar(100) COLLATE "pg_catalog"."default",
  "password_type" varchar(100) COLLATE "pg_catalog"."default",
  "auto_sync" int4,
  "last_sync" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for model
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."model" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "description" varchar(100) COLLATE "pg_catalog"."default",
  "model_text" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for organization
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."organization" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "website_url" varchar(100) COLLATE "pg_catalog"."default",
  "logo" varchar(200) COLLATE "pg_catalog"."default",
  "logo_dark" varchar(200) COLLATE "pg_catalog"."default",
  "favicon" varchar(200) COLLATE "pg_catalog"."default",
  "has_privilege_consent" bool,
  "password_type" varchar(100) COLLATE "pg_catalog"."default",
  "password_salt" varchar(100) COLLATE "pg_catalog"."default",
  "password_options" varchar(100) COLLATE "pg_catalog"."default",
  "password_obfuscator_type" varchar(100) COLLATE "pg_catalog"."default",
  "password_obfuscator_key" varchar(100) COLLATE "pg_catalog"."default",
  "password_expire_days" int4,
  "country_codes" text COLLATE "pg_catalog"."default",
  "default_avatar" varchar(200) COLLATE "pg_catalog"."default",
  "default_application" varchar(100) COLLATE "pg_catalog"."default",
  "user_types" text COLLATE "pg_catalog"."default",
  "tags" text COLLATE "pg_catalog"."default",
  "languages" varchar(255) COLLATE "pg_catalog"."default",
  "theme_data" json,
  "master_password" varchar(200) COLLATE "pg_catalog"."default",
  "default_password" varchar(200) COLLATE "pg_catalog"."default",
  "master_verification_code" varchar(100) COLLATE "pg_catalog"."default",
  "ip_whitelist" varchar(200) COLLATE "pg_catalog"."default",
  "init_score" int4,
  "enable_soft_deletion" bool,
  "is_profile_public" bool,
  "use_email_as_username" bool,
  "enable_tour" bool,
  "ip_restriction" varchar(255) COLLATE "pg_catalog"."default",
  "nav_items" varchar(1000) COLLATE "pg_catalog"."default",
  "widget_items" varchar(1000) COLLATE "pg_catalog"."default",
  "mfa_items" varchar(300) COLLATE "pg_catalog"."default",
  "account_items" varchar(5000) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for payment
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."payment" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "provider" varchar(100) COLLATE "pg_catalog"."default",
  "type" varchar(100) COLLATE "pg_catalog"."default",
  "product_name" varchar(100) COLLATE "pg_catalog"."default",
  "product_display_name" varchar(100) COLLATE "pg_catalog"."default",
  "detail" varchar(255) COLLATE "pg_catalog"."default",
  "tag" varchar(100) COLLATE "pg_catalog"."default",
  "currency" varchar(100) COLLATE "pg_catalog"."default",
  "price" float8,
  "return_url" varchar(1000) COLLATE "pg_catalog"."default",
  "is_recharge" bool,
  "user" varchar(100) COLLATE "pg_catalog"."default",
  "person_name" varchar(100) COLLATE "pg_catalog"."default",
  "person_id_card" varchar(100) COLLATE "pg_catalog"."default",
  "person_email" varchar(100) COLLATE "pg_catalog"."default",
  "person_phone" varchar(100) COLLATE "pg_catalog"."default",
  "invoice_type" varchar(100) COLLATE "pg_catalog"."default",
  "invoice_title" varchar(100) COLLATE "pg_catalog"."default",
  "invoice_tax_id" varchar(100) COLLATE "pg_catalog"."default",
  "invoice_remark" varchar(100) COLLATE "pg_catalog"."default",
  "invoice_url" varchar(255) COLLATE "pg_catalog"."default",
  "out_order_id" varchar(100) COLLATE "pg_catalog"."default",
  "pay_url" varchar(2000) COLLATE "pg_catalog"."default",
  "success_url" varchar(2000) COLLATE "pg_catalog"."default",
  "state" varchar(100) COLLATE "pg_catalog"."default",
  "message" varchar(2000) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for permission
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."permission" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "description" varchar(100) COLLATE "pg_catalog"."default",
  "users" text COLLATE "pg_catalog"."default",
  "groups" text COLLATE "pg_catalog"."default",
  "roles" text COLLATE "pg_catalog"."default",
  "domains" text COLLATE "pg_catalog"."default",
  "model" varchar(100) COLLATE "pg_catalog"."default",
  "adapter" varchar(100) COLLATE "pg_catalog"."default",
  "resource_type" varchar(100) COLLATE "pg_catalog"."default",
  "resources" text COLLATE "pg_catalog"."default",
  "actions" text COLLATE "pg_catalog"."default",
  "effect" varchar(100) COLLATE "pg_catalog"."default",
  "is_enabled" bool,
  "submitter" varchar(100) COLLATE "pg_catalog"."default",
  "approver" varchar(100) COLLATE "pg_catalog"."default",
  "approve_time" varchar(100) COLLATE "pg_catalog"."default",
  "state" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for permission_rule
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."permission_rule" (
  "id" int8 NOT NULL DEFAULT nextval('permission_rule_id_seq'::regclass),
  "ptype" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v0" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v1" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v2" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v3" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v4" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "v5" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying
)
;

-- ----------------------------
-- Table structure for plan
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."plan" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "description" varchar(100) COLLATE "pg_catalog"."default",
  "price" float8,
  "currency" varchar(100) COLLATE "pg_catalog"."default",
  "period" varchar(100) COLLATE "pg_catalog"."default",
  "product" varchar(100) COLLATE "pg_catalog"."default",
  "payment_providers" varchar(100) COLLATE "pg_catalog"."default",
  "is_enabled" bool,
  "role" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for pricing
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."pricing" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "description" varchar(100) COLLATE "pg_catalog"."default",
  "plans" text COLLATE "pg_catalog"."default",
  "is_enabled" bool,
  "trial_duration" int4,
  "application" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for product
-- ----------------------------
CREATE TABLE IF NOT EXISTS "public"."product" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "image" varchar(100) COLLATE "pg_catalog"."default",
  "detail" varchar(1000) COLLATE "pg_catalog"."default",
  "description" varchar(200) COLLATE "pg_catalog"."default",
  "tag" varchar(100) COLLATE "pg_catalog"."default",
  "currency" varchar(100) COLLATE "pg_catalog"."default",
  "price" float8,
  "quantity" int4,
  "sold" int4,
  "is_recharge" bool,
  "providers" varchar(255) COLLATE "pg_catalog"."default",
  "return_url" varchar(1000) COLLATE "pg_catalog"."default",
  "state" varchar(100) COLLATE "pg_catalog"."default"
)
;
