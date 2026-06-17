/*
 Navicat Premium Dump SQL

 Source Server         : pgsql临时
 Source Server Type    : PostgreSQL
 Source Server Version : 150017 (150017)
 Source Host           : 192.168.79.130:35432
 Source Catalog        : casdoor
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 150017 (150017)
 File Encoding         : 65001

 Date: 17/03/2026 10:29:08
*/

\c casdoor;

-- ----------------------------
-- Sequence structure for casbin_api_rule_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."casbin_api_rule_id_seq";
CREATE SEQUENCE "public"."casbin_api_rule_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for casbin_rule_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."casbin_rule_id_seq";
CREATE SEQUENCE "public"."casbin_rule_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for casbin_user_rule_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."casbin_user_rule_id_seq";
CREATE SEQUENCE "public"."casbin_user_rule_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for permission_rule_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."permission_rule_id_seq";
CREATE SEQUENCE "public"."permission_rule_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for record_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."record_id_seq";
CREATE SEQUENCE "public"."record_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Table structure for adapter
-- ----------------------------
DROP TABLE IF EXISTS "public"."adapter";
CREATE TABLE "public"."adapter" (
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
-- Records of adapter
-- ----------------------------
INSERT INTO "public"."adapter" VALUES ('built-in', 'api-adapter-built-in', '2026-03-17T01:59:49Z', 'casbin_api_rule', 'f', '', '', '', 0, '', '', '');
INSERT INTO "public"."adapter" VALUES ('built-in', 'user-adapter-built-in', '2026-03-17T01:59:49Z', 'casbin_user_rule', 'f', '', '', '', 0, '', '', '');

-- ----------------------------
-- Table structure for application
-- ----------------------------
DROP TABLE IF EXISTS "public"."application";
CREATE TABLE "public"."application" (
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
-- Records of application
-- ----------------------------
INSERT INTO "public"."application" VALUES ('admin', 'app-built-in', '2026-03-17T01:59:46Z', 'Casdoor', 'https://cdn.casbin.org/img/casdoor-logo_1185x256.png', 'https://casdoor.org', '', 'built-in', 'cert-built-in', '', '', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', '', '', '[{"owner":"","name":"provider_captcha_default","canSignUp":false,"canSignIn":false,"canUnlink":false,"countryCodes":null,"prompted":false,"signupGroup":"","rule":"None","provider":null}]', '[{"name":"Password","displayName":"Password","rule":"All"},{"name":"Verification code","displayName":"Verification code","rule":"All"},{"name":"WebAuthn","displayName":"WebAuthn","rule":"None"},{"name":"Face ID","displayName":"Face ID","rule":"None"}]', '[{"name":"ID","visible":false,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"Random"},{"name":"Username","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Display name","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Password","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Confirm password","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Email","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"Normal"},{"name":"Phone","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Agreement","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"}]', 'null', 'null', '[]', 'null', 'f', '', '8dd1828f0bc708f64365', 'ab4af8f8748bce6756a3bc0c42ba41b716448fe7', '[]', '', 'JWT', '', '[]', 168, 0, '', '', '', '', '', '', '', '', NULL, '', '', '', 2, '', '', '', 0, 0);
INSERT INTO "public"."application" VALUES ('admin', 'costrict_login', '2026-03-17T10:02:21+08:00', 'CoStrict', 'https://cdn.casbin.org/img/casdoor-logo_1185x256.png', '', '', 'costrict_group', 'cert-built-in', '', '', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', '', '', '[{"owner":"","name":"provider_captcha_default","canSignUp":false,"canSignIn":false,"canUnlink":false,"countryCodes":null,"prompted":false,"signupGroup":"","rule":"","provider":null}]', '[{"name":"Password","displayName":"Password","rule":"All"},{"name":"Verification code","displayName":"Verification code","rule":"All"},{"name":"WebAuthn","displayName":"WebAuthn","rule":"None"},{"name":"Face ID","displayName":"Face ID","rule":"None"}]', '[{"name":"ID","visible":false,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"Random"},{"name":"Username","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Display name","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Password","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Confirm password","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Email","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"Normal"},{"name":"Phone","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Agreement","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Signup button","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Providers","visible":true,"required":true,"prompted":false,"type":"","customCss":".provider-img {\n width: 30px;\n margin: 5px;\n }\n .provider-big-img {\n margin-bottom: 10px;\n }\n ","label":"","placeholder":"","options":null,"regex":"","rule":"small"}]', '[{"name":"Back button","visible":true,"label":"","customCss":".back-button {\n      top: 65px;\n      left: 15px;\n      position: absolute;\n}\n.back-inner-button{}","placeholder":"","rule":"None","isCustom":false},{"name":"Languages","visible":true,"label":"","customCss":".login-languages {\n    top: 55px;\n    right: 5px;\n    position: absolute;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Logo","visible":true,"label":"","customCss":".login-logo-box {}","placeholder":"","rule":"None","isCustom":false},{"name":"Signin methods","visible":true,"label":"","customCss":".signin-methods {}","placeholder":"","rule":"None","isCustom":false},{"name":"Username","visible":true,"label":"","customCss":".login-username {}\n.login-username-input{}","placeholder":"","rule":"None","isCustom":false},{"name":"Password","visible":true,"label":"","customCss":".login-password {}\n.login-password-input{}","placeholder":"","rule":"None","isCustom":false},{"name":"Agreement","visible":true,"label":"","customCss":".login-agreement {}","placeholder":"","rule":"None","isCustom":false},{"name":"Forgot password?","visible":true,"label":"","customCss":".login-forget-password {\n    display: inline-flex;\n    justify-content: space-between;\n    width: 320px;\n    margin-bottom: 25px;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Login button","visible":true,"label":"","customCss":".login-button-box {\n    margin-bottom: 5px;\n}\n.login-button {\n    width: 100%;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Signup link","visible":true,"label":"","customCss":".login-signup-link {\n    margin-bottom: 24px;\n    display: flex;\n    justify-content: end;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Providers","visible":true,"label":"","customCss":".provider-img {\n      width: 30px;\n      margin: 5px;\n}\n.provider-big-img {\n      margin-bottom: 10px;\n}","placeholder":"","rule":"small","isCustom":false}]', '["authorization_code","password","client_credentials","token","id_token","refresh_token"]', '[]', 'null', 'f', '', '9e2fc5d4fbcd52ef4f6f', 'ab5d8ba28b0e6c0d6e971247cdc1deb269c9eea3', '[".*/oidc-auth/api/v1/plugin/login/callback",".*/oidc-auth/api/v1/manager/bind/account/callback",".*/oidc-auth/api/v1/manager/login/callback"]', '', 'JWT', '', '[]', 168, 168, '', '', '', '', '', '', '', '', NULL, '', '', '', 2, '', '', '', 5, 15);

-- ----------------------------
-- Table structure for casbin_api_rule
-- ----------------------------
DROP TABLE IF EXISTS "public"."casbin_api_rule";
CREATE TABLE "public"."casbin_api_rule" (
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
-- Records of casbin_api_rule
-- ----------------------------
INSERT INTO "public"."casbin_api_rule" VALUES (1, 'p', 'built-in', '*', '*', '*', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (2, 'p', 'app', '*', '*', '*', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (3, 'p', '*', '*', 'POST', '/api/signup', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (4, 'p', '*', '*', 'GET', '/api/get-email-and-phone', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (5, 'p', '*', '*', 'POST', '/api/login', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (6, 'p', '*', '*', 'GET', '/api/get-app-login', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (7, 'p', '*', '*', 'POST', '/api/logout', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (8, 'p', '*', '*', 'GET', '/api/logout', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (9, 'p', '*', '*', 'POST', '/api/callback', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (10, 'p', '*', '*', 'POST', '/api/device-auth', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (11, 'p', '*', '*', 'GET', '/api/get-account', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (12, 'p', '*', '*', 'GET', '/api/userinfo', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (13, 'p', '*', '*', 'GET', '/api/user', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (14, 'p', '*', '*', 'GET', '/api/health', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (15, 'p', '*', '*', '*', '/api/webhook', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (16, 'p', '*', '*', 'GET', '/api/get-qrcode', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (17, 'p', '*', '*', 'GET', '/api/get-webhook-event', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (18, 'p', '*', '*', 'GET', '/api/get-captcha-status', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (19, 'p', '*', '*', '*', '/api/login/oauth', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (20, 'p', '*', '*', 'GET', '/api/get-application', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (21, 'p', '*', '*', 'GET', '/api/get-organization-applications', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (22, 'p', '*', '*', 'GET', '/api/get-user', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (23, 'p', '*', '*', 'GET', '/api/get-user-application', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (24, 'p', '*', '*', 'GET', '/api/get-resources', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (25, 'p', '*', '*', 'GET', '/api/get-records', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (26, 'p', '*', '*', 'GET', '/api/get-product', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (27, 'p', '*', '*', 'POST', '/api/buy-product', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (28, 'p', '*', '*', 'GET', '/api/get-payment', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (29, 'p', '*', '*', 'POST', '/api/update-payment', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (30, 'p', '*', '*', 'POST', '/api/invoice-payment', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (31, 'p', '*', '*', 'POST', '/api/notify-payment', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (32, 'p', '*', '*', 'POST', '/api/unlink', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (33, 'p', '*', '*', 'POST', '/api/set-password', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (34, 'p', '*', '*', 'POST', '/api/send-verification-code', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (35, 'p', '*', '*', 'GET', '/api/get-captcha', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (36, 'p', '*', '*', 'POST', '/api/verify-captcha', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (37, 'p', '*', '*', 'POST', '/api/verify-code', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (38, 'p', '*', '*', 'POST', '/api/reset-email-or-phone', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (39, 'p', '*', '*', 'POST', '/api/upload-resource', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (40, 'p', '*', '*', 'GET', '/.well-known/openid-configuration', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (41, 'p', '*', '*', 'GET', '/.well-known/webfinger', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (42, 'p', '*', '*', '*', '/.well-known/jwks', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (43, 'p', '*', '*', 'GET', '/api/get-saml-login', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (44, 'p', '*', '*', 'POST', '/api/acs', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (45, 'p', '*', '*', 'GET', '/api/saml/metadata', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (46, 'p', '*', '*', '*', '/api/saml/redirect', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (47, 'p', '*', '*', '*', '/cas', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (48, 'p', '*', '*', '*', '/scim', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (49, 'p', '*', '*', '*', '/api/webauthn', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (50, 'p', '*', '*', 'GET', '/api/get-release', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (51, 'p', '*', '*', 'GET', '/api/get-default-application', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (52, 'p', '*', '*', 'GET', '/api/get-prometheus-info', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (53, 'p', '*', '*', '*', '/api/metrics', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (54, 'p', '*', '*', 'GET', '/api/get-pricing', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (55, 'p', '*', '*', 'GET', '/api/get-plan', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (56, 'p', '*', '*', 'GET', '/api/get-subscription', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (57, 'p', '*', '*', 'GET', '/api/get-provider', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (58, 'p', '*', '*', 'GET', '/api/get-organization-names', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (59, 'p', '*', '*', 'GET', '/api/get-all-objects', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (60, 'p', '*', '*', 'GET', '/api/get-all-actions', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (61, 'p', '*', '*', 'GET', '/api/get-all-roles', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (62, 'p', '*', '*', 'GET', '/api/run-casbin-command', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (63, 'p', '*', '*', 'POST', '/api/refresh-engines', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (64, 'p', '*', '*', 'GET', '/api/get-invitation-info', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (65, 'p', '*', '*', 'GET', '/api/faceid-signin-begin', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (66, 'p', '*', '*', 'POST', '/api/identity/merge', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (67, 'p', '*', '*', 'GET', '/api/identity/info', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (68, 'p', '*', '*', 'POST', '/api/identity/bind', '*', '*');
INSERT INTO "public"."casbin_api_rule" VALUES (69, 'p', '*', '*', 'POST', '/api/identity/unbind', '*', '*');

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS "public"."casbin_rule";
CREATE TABLE "public"."casbin_rule" (
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
-- Records of casbin_rule
-- ----------------------------

-- ----------------------------
-- Table structure for casbin_user_rule
-- ----------------------------
DROP TABLE IF EXISTS "public"."casbin_user_rule";
CREATE TABLE "public"."casbin_user_rule" (
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
-- Records of casbin_user_rule
-- ----------------------------

-- ----------------------------
-- Table structure for cert
-- ----------------------------
DROP TABLE IF EXISTS "public"."cert";
CREATE TABLE "public"."cert" (
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
-- Records of cert
-- ----------------------------
INSERT INTO "public"."cert" VALUES ('admin', 'cert-built-in', '2026-03-17T01:59:46Z', 'Built-in Cert', 'JWT', 'x509', 'RS256', 4096, 20, '-----BEGIN CERTIFICATE-----
MIIE3TCCAsWgAwIBAgIDAeJAMA0GCSqGSIb3DQEBCwUAMCgxDjAMBgNVBAoTBWFk
bWluMRYwFAYDVQQDEw1jZXJ0LWJ1aWx0LWluMB4XDTI2MDMxNzAxNTk0OVoXDTQ2
MDMxNzAxNTk0OVowKDEOMAwGA1UEChMFYWRtaW4xFjAUBgNVBAMTDWNlcnQtYnVp
bHQtaW4wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCyuQTjmKxnIQmP
19iSRU7U6o2BfB5uGPpi+3XkNd5+iuB68yfLaDfTU++t6sfKMHc/wEZtK7kE0P82
DlnO/PqeN3qLXZ+zciPmYBD34J/75dBJcfphaQqAVD2kbM59tniYqkmYgIf+w/cM
AuPxXWItF+R9EmFbQZ0walK2/8UZfTZTB1aFOeSn88LnWRQDw17RD1zhQJ3jui/i
53QnK1soRO0cwlaGK1RXS8CwnX3eMql5+pJsnD1kxOPFOAdJ/PdR0hqfQHAp3hjD
V9PhER0FbXXQYJ76De6Jzq6Gn/dj2RkQrNQ+bc0jU58j8qt4LMD/iz7Bq/LrjRHm
GTSCv4E653d61+Xqoh9JWPwsER4MpGT6yu4RmXAgGCopKezDJ6n9qts2bgr4HLsu
yzeLfi6LHxUmX2Dyao3ETeZoSK7sdkdwSP8Eb3ND+/U1iVI5D4MtIhP1H1V5dH7I
bhjC2St88gbFx1iaaEHBhQIcVFZS6g2u2jOLkibaeCjgmoX8PtFQH1u/ZQQnSKSo
BUvbgzVoxO50Vs9hDEgJLWf59iCvpwgHkdWvEYufqV4CcRb6a7OS4ik25KAbomAx
k+1bT2vm+S3a7uDr9lQGgdaw6ivUizgKGTtx5pD1zKhjFFJ5mF/yhxqwVVy8cg9p
o1/OZX95oVkX4wbY6+iu2/ykKV67yQIDAQABoxAwDjAMBgNVHRMBAf8EAjAAMA0G
CSqGSIb3DQEBCwUAA4ICAQA5KKRBlhkzgxsAosDB8cym4NLbiYNKgN7hzS4OJH1X
HuKKcgwwTISW6CrWLh21lQOGQUNHW7Gs7IKriaKV4tSgWFV36o3lZ7ctSjn6sGtX
/dMs73ZkiB6NWMhNcIubuDlTzmSxWG2Mget6jXLxyvfL5LPrIDxeOkCSvNSEuk5b
fHcb1gBAx0CEbVC+HmpP6gBsb4eBnSq+preBK6CZf6lgIEH236+FbQqDlUxfhhYl
1A8I0zNQr/Zwxzh3J12Npx0d31bQN8qOvcBo+gVj3nuaVd7vMi3ziynT4lQk2p3j
dhoYof+CUVtSUQxDYih9wX5sjkubkPQhqTwfUMny4u7dZpDoJEukMrVcSoJf9IGS
wmLNfM+Pm0ECl+l2PSjRJYNOymj/xgr6JPk20qRZc2GQHtr2wzdZEKjLAUNGzntN
NTjEZNvatDZ+QcRlaWL4HhDp7L02y9p12Y65OcL5acjpz9cniWnIyVK7NTicFtgc
2uBiKu0B7bVSMlJBbA5cvKyjDMlzKsoTVE4rxYFD81gulQYoiL9xV2dLikGdISMW
7XMb/EWlF0Z/OrO/xBlD2gJ6oJApSUt/2NXAT+Snsaw7DrncrskjQdG5PapUH+YI
whezXJM/xkxLhJMav0ByjkwD5NOcKAMq5stqy4Y0R1N4eHZUS3UVWv8JvUaNNpGY
YQ==
-----END CERTIFICATE-----
', '-----BEGIN RSA PRIVATE KEY-----
MIIJKAIBAAKCAgEAsrkE45isZyEJj9fYkkVO1OqNgXwebhj6Yvt15DXeforgevMn
y2g301PvrerHyjB3P8BGbSu5BND/Ng5Zzvz6njd6i12fs3Ij5mAQ9+Cf++XQSXH6
YWkKgFQ9pGzOfbZ4mKpJmICH/sP3DALj8V1iLRfkfRJhW0GdMGpStv/FGX02UwdW
hTnkp/PC51kUA8Ne0Q9c4UCd47ov4ud0JytbKETtHMJWhitUV0vAsJ193jKpefqS
bJw9ZMTjxTgHSfz3UdIan0BwKd4Yw1fT4REdBW110GCe+g3uic6uhp/3Y9kZEKzU
Pm3NI1OfI/KreCzA/4s+wavy640R5hk0gr+BOud3etfl6qIfSVj8LBEeDKRk+sru
EZlwIBgqKSnswyep/arbNm4K+By7Lss3i34uix8VJl9g8mqNxE3maEiu7HZHcEj/
BG9zQ/v1NYlSOQ+DLSIT9R9VeXR+yG4YwtkrfPIGxcdYmmhBwYUCHFRWUuoNrtoz
i5Im2ngo4JqF/D7RUB9bv2UEJ0ikqAVL24M1aMTudFbPYQxICS1n+fYgr6cIB5HV
rxGLn6leAnEW+muzkuIpNuSgG6JgMZPtW09r5vkt2u7g6/ZUBoHWsOor1Is4Chk7
ceaQ9cyoYxRSeZhf8ocasFVcvHIPaaNfzmV/eaFZF+MG2Ovortv8pCleu8kCAwEA
AQKCAgA2ZBp5zZgCVsqbLTxq347rO2rPUl/Mb+ubvh2V01MgPdGteVzOCaybsvrJ
eh9NiExyaWjFhr3rGlPhNqLcXJ9FEWtJCLg2m3EThn8DgSgPOFSYVYi2poocngaY
iMia97Jh+8kkEVHSFf+Udzi51wqu4IuYCMT6SdC10FXV/acVZ8blo+EP8JuHBad2
UuTFdGrhwwQoRUsH6Dz0mC4AUbCPY9KpovdEI+GqI5LD7VpGCjMLvVdcmHwSs90B
w1hsiwEVtJSHjMvV4hK1OIByCYxEf5iQTydywFStOephbEOU7pgzO/Xn/fClZKiU
/yByFWSzIMR7moH2b9tXpKoXdIK2US0rROXi8ZDALkTMjD4A+nH9P+2j7XIvZQJ1
kHJ7jDGmdlZkVc37XeQJDOQRnhETf2MPojFC9BSdmZvLxHJydK1Gh3Eqq3w0t+X6
vue/VoUTUY0aQme5nk4p7fDgzPOoSllkA1bTaeQS7V7YVsRrrTGDxTT+NKajPcgw
e0R0imEPghgiR6AaceUr/VwaCUtMGFXoiMyacRCA6ye6/ApcZq+3GNbhWVC62umF
ZSg5bhKjhlUSWYTEG8ZiVBAWLQhVRUkNrfJNkk0NFbfN60fcH9z2JXmhMAoFUrAL
5uHNm/j5LWWYyagjtPP4LjMzguFYGausD70xUMv0unijZAxbZQKCAQEA4GRwLkgY
BhoGoHi6MpsJGm+7ZSexpNFRvL7rVS7PjWfaEvb6bcmZKc/Np17yW9O/+MP9TD5E
ZvlsA23epczdTgil93c2wA3FGrHLfJcb1Cm98DfH4GuWuOMRIsUbglMXiQdMionN
muJmeBbzUbKnynVlQ2cQ8ieruiAOtf2JLrp/VA8Fy3z10osMS2fbrLvlilXk7Oeh
2uX+GJAn8wEfVc4DEPBlFmxPaHDVV17S89cgW2BELXKKi+D5bh3TElKPQINqxMMy
R4XajSPvAfWrUdraK248gOSwhWWUwgc55qk3ouWG5+66r3APr5DiyuPVIWbodtFR
WbudAZ66mN0yWwKCAQEAy+W9e6oErJT+WM0IK4Q3rtWPq6odBLOWpkAlZNFOqQrF
oFOB/UYD28Jo8eByZUGT6dfASLys/t/b6fubW9L3VVOHjYoWMjpuOzD5eQqqC2Do
FRc7//TpMEiucwqpo5Geg6kz8Xd4yThttRBcUWZZOfF6+5w2hoVQFUHySHt2bYW6
Etxv4HoJP5Si18hpRlExI5hUsiqNLZJUc82WveKhAKhxQz9+N9UPIcU3SHHZA9rV
1WS789yei4+2G4l/7sV/pN6dNB1y+m+OymUBgT/pIM8ZsUyDbVvNLqak5yLF1rkq
tTHi7sjm2zTKgL9ujHCOAxzkR5LrInD5fHI+yDObqwKCAQAQv+6PxdnlOrQb7mcX
t/LqAkUeu+PKxcUwLaArb58WkxX4LG9tGotQT66uO6wjpp7B1v9lXZCFj0C80Mq5
3Zz3YjXhTwjLNTvtTjA+UJlHMQZ4k/FbkyxhCLISwlvXUpSDKgRQoxCx6wxzEmN/
SbK0rias7ViaetOaxSHT++TtlAZFJAnuspk26iI0zr/bGPt9ZAP+mxu910UsWDZj
MUZBhWs9LHSfQ955kVKbfvyILxq/N2RfcBlfb5+u1JMY3wjKbATgQyHRrPAvcOs6
Xjg3l/qJG2VwyaFTk82S3i84RJFem3uCfJN63QaOuGX1dyRFpgTcOau04whJQqvz
FS/ZAoIBACW1Hw23q5WfGEqj4Fbr5MBYjnFYzHUKIVotkrzKWJm6fEHiL6//Osdx
dKcGsyE2t3wPgmfxPK0GCU2EUn7wzPzirbbnSqOPWsNnkztnVXs1OLZtzZLQYC+l
wV1VI3o0e/QtXoXRam0MGn3ZwLjtDKNCnVPa19qY3O6Re5UcABLRIFBMpAQHjZ5M
hrERRybR7RyAF2m7pNBAt42S/8WuG8hbXlbpbxY0ONAuAMNNCq1D0zfboUbADjIR
KcR98ULfE3QSF1CEAbd+2hFkDJsyjpKauk2nsMaQq6ls+CImtB+5rsRDrYwra0Yu
xSR6Qp+GUT0SIsV5zjJYYks5dZLguHECggEBAJarI+exfGrx5CH4OnBUH6Hsw/oF
ZIM6hf2ecRm8Pcbis9VobEznuHUNgkVSefawUar+q5glSDmoI+Eha2rU98rmkWq6
llndZ2RXq/lfpEykOTsuH3Lg2O3gtoCchQJQ1Rna/kv4IEC92IO99DBVGmID57ni
hvejTBVO75npsegQiMk4KL4Juc2huYAfIjODac0RjQlnjmJ29jzN9GQvjY7En8iv
aT4jcOUCZSZ4oXmEF4kiZB9CQ4v9QZjmIoLLcLlnuPCieyycUZNkAdfBLYOCKUuW
RQOVOqaSn0GWz7JoIAOHsw2I8o/yxNu55kH52+lxbIq1kYcNkbxzkJx0l1c=
-----END RSA PRIVATE KEY-----
');

-- ----------------------------
-- Table structure for enforcer
-- ----------------------------
DROP TABLE IF EXISTS "public"."enforcer";
CREATE TABLE "public"."enforcer" (
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
-- Records of enforcer
-- ----------------------------
INSERT INTO "public"."enforcer" VALUES ('built-in', 'api-enforcer-built-in', '2026-03-17T01:59:49Z', '2026-03-17 01:59:49', 'API Enforcer', '', 'built-in/api-model-built-in', 'built-in/api-adapter-built-in', NULL);
INSERT INTO "public"."enforcer" VALUES ('built-in', 'user-enforcer-built-in', '2026-03-17T01:59:49Z', '2026-03-17 01:59:49', 'User Enforcer', '', 'built-in/user-model-built-in', 'built-in/user-adapter-built-in', NULL);

-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS "public"."group";
CREATE TABLE "public"."group" (
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
-- Records of group
-- ----------------------------

-- ----------------------------
-- Table structure for invitation
-- ----------------------------
DROP TABLE IF EXISTS "public"."invitation";
CREATE TABLE "public"."invitation" (
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
-- Records of invitation
-- ----------------------------

-- ----------------------------
-- Table structure for ldap
-- ----------------------------
DROP TABLE IF EXISTS "public"."ldap";
CREATE TABLE "public"."ldap" (
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
-- Records of ldap
-- ----------------------------
INSERT INTO "public"."ldap" VALUES ('ldap-built-in', 'built-in', '2026-03-17T01:59:49Z', 'BuildIn LDAP Server', 'example.com', 389, 'f', 'f', 'cn=buildin,dc=example,dc=com', '123', 'ou=BuildIn,dc=example,dc=com', '', 'null', '', '', 0, '');

-- ----------------------------
-- Table structure for model
-- ----------------------------
DROP TABLE IF EXISTS "public"."model";
CREATE TABLE "public"."model" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "description" varchar(100) COLLATE "pg_catalog"."default",
  "model_text" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of model
-- ----------------------------
INSERT INTO "public"."model" VALUES ('built-in', 'api-model-built-in', '2026-03-17T01:59:49Z', 'API Model', '', '[request_definition]
r = subOwner, subName, method, urlPath, objOwner, objName

[policy_definition]
p = subOwner, subName, method, urlPath, objOwner, objName

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (r.subOwner == p.subOwner || p.subOwner == "*") && \
    (r.subName == p.subName || p.subName == "*" || r.subName != "anonymous" && p.subName == "!anonymous") && \
    (r.method == p.method || p.method == "*") && \
    (r.urlPath == p.urlPath || p.urlPath == "*") && \
    (r.objOwner == p.objOwner || p.objOwner == "*") && \
    (r.objName == p.objName || p.objName == "*") || \
    (r.subOwner == r.objOwner && r.subName == r.objName)');
INSERT INTO "public"."model" VALUES ('built-in', 'user-model-built-in', '2026-03-17T01:59:49Z', 'Built-in Model', '', '[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act');

-- ----------------------------
-- Table structure for organization
-- ----------------------------
DROP TABLE IF EXISTS "public"."organization";
CREATE TABLE "public"."organization" (
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
-- Records of organization
-- ----------------------------
INSERT INTO "public"."organization" VALUES ('admin', 'built-in', '2026-03-17T01:59:46Z', 'Built-in Organization', 'https://example.com', '', '', 'https://cdn.casbin.org/img/casbin/favicon.ico', 'f', 'plain', '', '["AtLeast6"]', '', '', 0, '["US","ES","FR","DE","GB","CN","JP","KR","VN","ID","SG","IN"]', 'https://cdn.casbin.org/img/casbin.svg', '', '[]', '[]', '["en","zh","es","fr","de","id","ja","ko","ru","vi","pt"]', NULL, '', '', '', '', 2000, 'f', 'f', 'f', 'f', '', 'null', 'null', 'null', '[{"name":"Organization","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"ID","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"Name","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Display name","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Avatar","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"User type","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Password","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Email","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Phone","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Country code","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Country/Region","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Location","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Affiliation","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Title","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Homepage","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Bio","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Tag","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Signup application","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Roles","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"Permissions","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"Groups","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"3rd-party logins","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Properties","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is admin","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is forbidden","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is deleted","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Multi-factor authentication","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"WebAuthn credentials","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Managed accounts","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"MFA accounts","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""}]');
INSERT INTO "public"."organization" VALUES ('admin', 'costrict_group', '2026-03-17T10:00:38+08:00', 'costrict', 'https://door.casdoor.com', '', '', 'https://cdn.casbin.org/img/favicon.png', 'f', 'plain', '', '["AtLeast6"]', 'Plain', '', 0, '["US"]', 'https://cdn.casbin.org/img/casbin.svg', '', 'null', '[]', '["en","es","fr","de","zh","id","ja","ko","ru","vi","pt","it","ms","tr","ar","he","nl","pl","fi","sv","uk","kk","fa","cs","sk"]', NULL, '', '', '', '', 0, 'f', 'f', 'f', 'f', '', 'null', 'null', 'null', '[{"name":"Organization","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"ID","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"Name","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Display name","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Avatar","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"User type","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Password","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Email","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Phone","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Country code","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Country/Region","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Location","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Address","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Affiliation","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Title","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"ID card type","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"ID card","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"ID card info","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Homepage","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Bio","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Tag","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Language","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Gender","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Birthday","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Education","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Score","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Karma","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Ranking","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Signup application","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"API key","visible":false,"viewRule":"","modifyRule":"Self","regex":""},{"name":"Groups","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Roles","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"Permissions","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"3rd-party logins","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Properties","visible":false,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is online","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is admin","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is forbidden","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is deleted","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Multi-factor authentication","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"WebAuthn credentials","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Managed accounts","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"MFA accounts","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""}]');

-- ----------------------------
-- Table structure for payment
-- ----------------------------
DROP TABLE IF EXISTS "public"."payment";
CREATE TABLE "public"."payment" (
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
-- Records of payment
-- ----------------------------

-- ----------------------------
-- Table structure for permission
-- ----------------------------
DROP TABLE IF EXISTS "public"."permission";
CREATE TABLE "public"."permission" (
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
-- Records of permission
-- ----------------------------
INSERT INTO "public"."permission" VALUES ('built-in', 'permission-built-in', '2026-03-17T01:59:46Z', 'Built-in Permission', 'Built-in Permission', '["built-in/*"]', '[]', '[]', '[]', 'user-model-built-in', '', 'Application', '["app-built-in"]', '["Read","Write","Admin"]', 'Allow', 'f', 'admin', 'admin', '2026-03-17T01:59:46Z', 'Approved');

-- ----------------------------
-- Table structure for permission_rule
-- ----------------------------
DROP TABLE IF EXISTS "public"."permission_rule";
CREATE TABLE "public"."permission_rule" (
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
-- Records of permission_rule
-- ----------------------------
INSERT INTO "public"."permission_rule" VALUES (1, 'p', 'built-in/*', 'app-built-in', 'read', 'allow', '', 'built-in/permission-built-in');
INSERT INTO "public"."permission_rule" VALUES (2, 'p', 'built-in/*', 'app-built-in', 'write', 'allow', '', 'built-in/permission-built-in');
INSERT INTO "public"."permission_rule" VALUES (3, 'p', 'built-in/*', 'app-built-in', 'admin', 'allow', '', 'built-in/permission-built-in');

-- ----------------------------
-- Table structure for plan
-- ----------------------------
DROP TABLE IF EXISTS "public"."plan";
CREATE TABLE "public"."plan" (
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
-- Records of plan
-- ----------------------------

-- ----------------------------
-- Table structure for pricing
-- ----------------------------
DROP TABLE IF EXISTS "public"."pricing";
CREATE TABLE "public"."pricing" (
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
-- Records of pricing
-- ----------------------------

-- ----------------------------
-- Table structure for product
-- ----------------------------
DROP TABLE IF EXISTS "public"."product";
CREATE TABLE "public"."product" (
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

-- ----------------------------
-- Records of product
-- ----------------------------

-- ----------------------------
-- Table structure for provider
-- ----------------------------
DROP TABLE IF EXISTS "public"."provider";
CREATE TABLE "public"."provider" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "category" varchar(100) COLLATE "pg_catalog"."default",
  "type" varchar(100) COLLATE "pg_catalog"."default",
  "sub_type" varchar(100) COLLATE "pg_catalog"."default",
  "method" varchar(100) COLLATE "pg_catalog"."default",
  "client_id" varchar(200) COLLATE "pg_catalog"."default",
  "client_secret" varchar(3000) COLLATE "pg_catalog"."default",
  "client_id2" varchar(100) COLLATE "pg_catalog"."default",
  "client_secret2" varchar(500) COLLATE "pg_catalog"."default",
  "cert" varchar(100) COLLATE "pg_catalog"."default",
  "custom_auth_url" varchar(200) COLLATE "pg_catalog"."default",
  "custom_token_url" varchar(200) COLLATE "pg_catalog"."default",
  "custom_user_info_url" varchar(200) COLLATE "pg_catalog"."default",
  "custom_logo" varchar(200) COLLATE "pg_catalog"."default",
  "scopes" varchar(100) COLLATE "pg_catalog"."default",
  "user_mapping" varchar(500) COLLATE "pg_catalog"."default",
  "http_headers" varchar(500) COLLATE "pg_catalog"."default",
  "host" varchar(100) COLLATE "pg_catalog"."default",
  "port" int4,
  "disable_ssl" bool,
  "title" varchar(100) COLLATE "pg_catalog"."default",
  "content" varchar(2000) COLLATE "pg_catalog"."default",
  "receiver" varchar(100) COLLATE "pg_catalog"."default",
  "region_id" varchar(100) COLLATE "pg_catalog"."default",
  "sign_name" varchar(100) COLLATE "pg_catalog"."default",
  "template_code" varchar(100) COLLATE "pg_catalog"."default",
  "app_id" varchar(100) COLLATE "pg_catalog"."default",
  "endpoint" varchar(1000) COLLATE "pg_catalog"."default",
  "intranet_endpoint" varchar(100) COLLATE "pg_catalog"."default",
  "domain" varchar(100) COLLATE "pg_catalog"."default",
  "bucket" varchar(100) COLLATE "pg_catalog"."default",
  "path_prefix" varchar(100) COLLATE "pg_catalog"."default",
  "metadata" text COLLATE "pg_catalog"."default",
  "id_p" text COLLATE "pg_catalog"."default",
  "issuer_url" varchar(100) COLLATE "pg_catalog"."default",
  "enable_sign_authn_request" bool,
  "email_regex" varchar(200) COLLATE "pg_catalog"."default",
  "provider_url" varchar(200) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of provider
-- ----------------------------
INSERT INTO "public"."provider" VALUES ('admin', 'provider_captcha_default', '2026-03-17T01:59:46Z', 'Captcha Default', 'Captcha', 'Default', '', '', '', '', '', '', '', '', '', '', '', '', 'null', 'null', '', 0, 'f', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', 'f', '', '');

-- ----------------------------
-- Table structure for radius_accounting
-- ----------------------------
DROP TABLE IF EXISTS "public"."radius_accounting";
CREATE TABLE "public"."radius_accounting" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" timestamp(6),
  "username" varchar(255) COLLATE "pg_catalog"."default",
  "service_type" int8,
  "nas_id" varchar(255) COLLATE "pg_catalog"."default",
  "nas_ip_addr" varchar(255) COLLATE "pg_catalog"."default",
  "nas_port_id" varchar(255) COLLATE "pg_catalog"."default",
  "nas_port_type" int8,
  "nas_port" int8,
  "framed_ip_addr" varchar(255) COLLATE "pg_catalog"."default",
  "framed_ip_netmask" varchar(255) COLLATE "pg_catalog"."default",
  "acct_session_id" varchar(255) COLLATE "pg_catalog"."default",
  "acct_session_time" int8,
  "acct_input_total" int8,
  "acct_output_total" int8,
  "acct_input_packets" int8,
  "acct_output_packets" int8,
  "acct_terminate_cause" int8,
  "last_update" timestamp(6),
  "acct_start_time" timestamp(6),
  "acct_stop_time" timestamp(6)
)
;

-- ----------------------------
-- Records of radius_accounting
-- ----------------------------

-- ----------------------------
-- Table structure for record
-- ----------------------------
DROP TABLE IF EXISTS "public"."record";
CREATE TABLE "public"."record" (
  "id" int4 NOT NULL DEFAULT nextval('record_id_seq'::regclass),
  "owner" varchar(100) COLLATE "pg_catalog"."default",
  "name" varchar(100) COLLATE "pg_catalog"."default",
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "organization" varchar(100) COLLATE "pg_catalog"."default",
  "client_ip" varchar(100) COLLATE "pg_catalog"."default",
  "user" varchar(100) COLLATE "pg_catalog"."default",
  "method" varchar(100) COLLATE "pg_catalog"."default",
  "request_uri" varchar(1000) COLLATE "pg_catalog"."default",
  "action" varchar(1000) COLLATE "pg_catalog"."default",
  "language" varchar(100) COLLATE "pg_catalog"."default",
  "object" text COLLATE "pg_catalog"."default",
  "response" text COLLATE "pg_catalog"."default",
  "status_code" int4,
  "is_triggered" bool
)
;

-- ----------------------------
-- Records of record
-- ----------------------------

-- ----------------------------
-- Table structure for resource
-- ----------------------------
DROP TABLE IF EXISTS "public"."resource";
CREATE TABLE "public"."resource" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(180) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "user" varchar(100) COLLATE "pg_catalog"."default",
  "provider" varchar(100) COLLATE "pg_catalog"."default",
  "application" varchar(100) COLLATE "pg_catalog"."default",
  "tag" varchar(100) COLLATE "pg_catalog"."default",
  "parent" varchar(100) COLLATE "pg_catalog"."default",
  "file_name" varchar(255) COLLATE "pg_catalog"."default",
  "file_type" varchar(100) COLLATE "pg_catalog"."default",
  "file_format" varchar(100) COLLATE "pg_catalog"."default",
  "file_size" int4,
  "url" varchar(500) COLLATE "pg_catalog"."default",
  "description" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of resource
-- ----------------------------

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS "public"."role";
CREATE TABLE "public"."role" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "description" varchar(100) COLLATE "pg_catalog"."default",
  "users" text COLLATE "pg_catalog"."default",
  "groups" text COLLATE "pg_catalog"."default",
  "roles" text COLLATE "pg_catalog"."default",
  "domains" text COLLATE "pg_catalog"."default",
  "is_enabled" bool
)
;

-- ----------------------------
-- Records of role
-- ----------------------------

-- ----------------------------
-- Table structure for session
-- ----------------------------
DROP TABLE IF EXISTS "public"."session";
CREATE TABLE "public"."session" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "application" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "session_id" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of session
-- ----------------------------
INSERT INTO "public"."session" VALUES ('built-in', 'admin', 'app-built-in', '2026-03-17T01:59:55Z', '["5806fee5aced7c26f99975e848394209"]');
INSERT INTO "public"."session" VALUES ('costrict_group', 'costrict', 'costrict_login', '2026-03-17T02:26:18Z', '["5806fee5aced7c26f99975e848394209"]');

-- ----------------------------
-- Table structure for subscription
-- ----------------------------
DROP TABLE IF EXISTS "public"."subscription";
CREATE TABLE "public"."subscription" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "description" varchar(100) COLLATE "pg_catalog"."default",
  "user" varchar(100) COLLATE "pg_catalog"."default",
  "pricing" varchar(100) COLLATE "pg_catalog"."default",
  "plan" varchar(100) COLLATE "pg_catalog"."default",
  "payment" varchar(100) COLLATE "pg_catalog"."default",
  "start_time" timestamp(6),
  "end_time" timestamp(6),
  "period" varchar(100) COLLATE "pg_catalog"."default",
  "state" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of subscription
-- ----------------------------

-- ----------------------------
-- Table structure for syncer
-- ----------------------------
DROP TABLE IF EXISTS "public"."syncer";
CREATE TABLE "public"."syncer" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "organization" varchar(100) COLLATE "pg_catalog"."default",
  "type" varchar(100) COLLATE "pg_catalog"."default",
  "database_type" varchar(100) COLLATE "pg_catalog"."default",
  "ssl_mode" varchar(100) COLLATE "pg_catalog"."default",
  "ssh_type" varchar(100) COLLATE "pg_catalog"."default",
  "host" varchar(100) COLLATE "pg_catalog"."default",
  "port" int4,
  "user" varchar(100) COLLATE "pg_catalog"."default",
  "password" varchar(150) COLLATE "pg_catalog"."default",
  "ssh_host" varchar(100) COLLATE "pg_catalog"."default",
  "ssh_port" int4,
  "ssh_user" varchar(100) COLLATE "pg_catalog"."default",
  "ssh_password" varchar(150) COLLATE "pg_catalog"."default",
  "cert" varchar(100) COLLATE "pg_catalog"."default",
  "database" varchar(100) COLLATE "pg_catalog"."default",
  "table" varchar(100) COLLATE "pg_catalog"."default",
  "table_columns" text COLLATE "pg_catalog"."default",
  "affiliation_table" varchar(100) COLLATE "pg_catalog"."default",
  "avatar_base_url" varchar(100) COLLATE "pg_catalog"."default",
  "error_text" text COLLATE "pg_catalog"."default",
  "sync_interval" int4,
  "is_read_only" bool,
  "is_enabled" bool
)
;

-- ----------------------------
-- Records of syncer
-- ----------------------------

-- ----------------------------
-- Table structure for token
-- ----------------------------
DROP TABLE IF EXISTS "public"."token";
CREATE TABLE "public"."token" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "application" varchar(100) COLLATE "pg_catalog"."default",
  "organization" varchar(100) COLLATE "pg_catalog"."default",
  "user" varchar(100) COLLATE "pg_catalog"."default",
  "code" varchar(100) COLLATE "pg_catalog"."default",
  "access_token" text COLLATE "pg_catalog"."default",
  "refresh_token" text COLLATE "pg_catalog"."default",
  "access_token_hash" varchar(100) COLLATE "pg_catalog"."default",
  "refresh_token_hash" varchar(100) COLLATE "pg_catalog"."default",
  "expires_in" int4,
  "scope" varchar(100) COLLATE "pg_catalog"."default",
  "token_type" varchar(100) COLLATE "pg_catalog"."default",
  "code_challenge" varchar(100) COLLATE "pg_catalog"."default",
  "code_is_used" bool,
  "code_expire_in" int8
)
;

-- ----------------------------
-- Records of token
-- ----------------------------


-- ----------------------------
-- Table structure for transaction
-- ----------------------------
DROP TABLE IF EXISTS "public"."transaction";
CREATE TABLE "public"."transaction" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "provider" varchar(100) COLLATE "pg_catalog"."default",
  "category" varchar(100) COLLATE "pg_catalog"."default",
  "type" varchar(100) COLLATE "pg_catalog"."default",
  "product_name" varchar(100) COLLATE "pg_catalog"."default",
  "product_display_name" varchar(100) COLLATE "pg_catalog"."default",
  "detail" varchar(255) COLLATE "pg_catalog"."default",
  "tag" varchar(100) COLLATE "pg_catalog"."default",
  "currency" varchar(100) COLLATE "pg_catalog"."default",
  "amount" float8,
  "return_url" varchar(1000) COLLATE "pg_catalog"."default",
  "user" varchar(100) COLLATE "pg_catalog"."default",
  "application" varchar(100) COLLATE "pg_catalog"."default",
  "payment" varchar(100) COLLATE "pg_catalog"."default",
  "state" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of transaction
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS "public"."user";
CREATE TABLE "public"."user" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "updated_time" varchar(100) COLLATE "pg_catalog"."default",
  "deleted_time" varchar(100) COLLATE "pg_catalog"."default",
  "id" varchar(100) COLLATE "pg_catalog"."default",
  "external_id" varchar(100) COLLATE "pg_catalog"."default",
  "universal_id" varchar(100) COLLATE "pg_catalog"."default",
  "type" varchar(100) COLLATE "pg_catalog"."default",
  "password" varchar(150) COLLATE "pg_catalog"."default",
  "password_salt" varchar(100) COLLATE "pg_catalog"."default",
  "password_type" varchar(100) COLLATE "pg_catalog"."default",
  "display_name" varchar(100) COLLATE "pg_catalog"."default",
  "first_name" varchar(100) COLLATE "pg_catalog"."default",
  "last_name" varchar(100) COLLATE "pg_catalog"."default",
  "avatar" varchar(500) COLLATE "pg_catalog"."default",
  "avatar_type" varchar(100) COLLATE "pg_catalog"."default",
  "permanent_avatar" varchar(500) COLLATE "pg_catalog"."default",
  "email" varchar(100) COLLATE "pg_catalog"."default",
  "email_verified" bool,
  "phone" varchar(100) COLLATE "pg_catalog"."default",
  "country_code" varchar(6) COLLATE "pg_catalog"."default",
  "region" varchar(100) COLLATE "pg_catalog"."default",
  "location" varchar(100) COLLATE "pg_catalog"."default",
  "address" text COLLATE "pg_catalog"."default",
  "affiliation" varchar(100) COLLATE "pg_catalog"."default",
  "title" varchar(100) COLLATE "pg_catalog"."default",
  "id_card_type" varchar(100) COLLATE "pg_catalog"."default",
  "id_card" varchar(100) COLLATE "pg_catalog"."default",
  "homepage" varchar(100) COLLATE "pg_catalog"."default",
  "bio" varchar(100) COLLATE "pg_catalog"."default",
  "tag" varchar(100) COLLATE "pg_catalog"."default",
  "language" varchar(100) COLLATE "pg_catalog"."default",
  "gender" varchar(100) COLLATE "pg_catalog"."default",
  "birthday" varchar(100) COLLATE "pg_catalog"."default",
  "education" varchar(100) COLLATE "pg_catalog"."default",
  "score" int4,
  "karma" int4,
  "ranking" int4,
  "balance" float8,
  "currency" varchar(100) COLLATE "pg_catalog"."default",
  "vip" int4 DEFAULT 0,
  "vip_expire" varchar(100) COLLATE "pg_catalog"."default",
  "is_default_avatar" bool,
  "is_online" bool,
  "is_admin" bool,
  "is_forbidden" bool,
  "is_deleted" bool,
  "signup_application" varchar(100) COLLATE "pg_catalog"."default",
  "hash" varchar(100) COLLATE "pg_catalog"."default",
  "pre_hash" varchar(100) COLLATE "pg_catalog"."default",
  "access_key" varchar(100) COLLATE "pg_catalog"."default",
  "access_secret" varchar(100) COLLATE "pg_catalog"."default",
  "access_token" text COLLATE "pg_catalog"."default",
  "created_ip" varchar(100) COLLATE "pg_catalog"."default",
  "last_signin_time" varchar(100) COLLATE "pg_catalog"."default",
  "last_signin_ip" varchar(100) COLLATE "pg_catalog"."default",
  "github" varchar(100) COLLATE "pg_catalog"."default",
  "google" varchar(100) COLLATE "pg_catalog"."default",
  "qq" varchar(100) COLLATE "pg_catalog"."default",
  "wechat" varchar(100) COLLATE "pg_catalog"."default",
  "facebook" varchar(100) COLLATE "pg_catalog"."default",
  "dingtalk" varchar(100) COLLATE "pg_catalog"."default",
  "weibo" varchar(100) COLLATE "pg_catalog"."default",
  "gitee" varchar(100) COLLATE "pg_catalog"."default",
  "linkedin" varchar(100) COLLATE "pg_catalog"."default",
  "wecom" varchar(100) COLLATE "pg_catalog"."default",
  "lark" varchar(100) COLLATE "pg_catalog"."default",
  "gitlab" varchar(100) COLLATE "pg_catalog"."default",
  "adfs" varchar(100) COLLATE "pg_catalog"."default",
  "baidu" varchar(100) COLLATE "pg_catalog"."default",
  "alipay" varchar(100) COLLATE "pg_catalog"."default",
  "casdoor" varchar(100) COLLATE "pg_catalog"."default",
  "infoflow" varchar(100) COLLATE "pg_catalog"."default",
  "apple" varchar(100) COLLATE "pg_catalog"."default",
  "azuread" varchar(100) COLLATE "pg_catalog"."default",
  "azureadb2c" varchar(100) COLLATE "pg_catalog"."default",
  "slack" varchar(100) COLLATE "pg_catalog"."default",
  "steam" varchar(100) COLLATE "pg_catalog"."default",
  "bilibili" varchar(100) COLLATE "pg_catalog"."default",
  "okta" varchar(100) COLLATE "pg_catalog"."default",
  "douyin" varchar(100) COLLATE "pg_catalog"."default",
  "kwai" varchar(100) COLLATE "pg_catalog"."default",
  "line" varchar(100) COLLATE "pg_catalog"."default",
  "amazon" varchar(100) COLLATE "pg_catalog"."default",
  "auth0" varchar(100) COLLATE "pg_catalog"."default",
  "battlenet" varchar(100) COLLATE "pg_catalog"."default",
  "bitbucket" varchar(100) COLLATE "pg_catalog"."default",
  "box" varchar(100) COLLATE "pg_catalog"."default",
  "cloudfoundry" varchar(100) COLLATE "pg_catalog"."default",
  "dailymotion" varchar(100) COLLATE "pg_catalog"."default",
  "deezer" varchar(100) COLLATE "pg_catalog"."default",
  "digitalocean" varchar(100) COLLATE "pg_catalog"."default",
  "discord" varchar(100) COLLATE "pg_catalog"."default",
  "dropbox" varchar(100) COLLATE "pg_catalog"."default",
  "eveonline" varchar(100) COLLATE "pg_catalog"."default",
  "fitbit" varchar(100) COLLATE "pg_catalog"."default",
  "gitea" varchar(100) COLLATE "pg_catalog"."default",
  "heroku" varchar(100) COLLATE "pg_catalog"."default",
  "influxcloud" varchar(100) COLLATE "pg_catalog"."default",
  "instagram" varchar(100) COLLATE "pg_catalog"."default",
  "intercom" varchar(100) COLLATE "pg_catalog"."default",
  "kakao" varchar(100) COLLATE "pg_catalog"."default",
  "lastfm" varchar(100) COLLATE "pg_catalog"."default",
  "mailru" varchar(100) COLLATE "pg_catalog"."default",
  "meetup" varchar(100) COLLATE "pg_catalog"."default",
  "microsoftonline" varchar(100) COLLATE "pg_catalog"."default",
  "naver" varchar(100) COLLATE "pg_catalog"."default",
  "nextcloud" varchar(100) COLLATE "pg_catalog"."default",
  "onedrive" varchar(100) COLLATE "pg_catalog"."default",
  "oura" varchar(100) COLLATE "pg_catalog"."default",
  "patreon" varchar(100) COLLATE "pg_catalog"."default",
  "paypal" varchar(100) COLLATE "pg_catalog"."default",
  "salesforce" varchar(100) COLLATE "pg_catalog"."default",
  "shopify" varchar(100) COLLATE "pg_catalog"."default",
  "soundcloud" varchar(100) COLLATE "pg_catalog"."default",
  "spotify" varchar(100) COLLATE "pg_catalog"."default",
  "strava" varchar(100) COLLATE "pg_catalog"."default",
  "stripe" varchar(100) COLLATE "pg_catalog"."default",
  "tiktok" varchar(100) COLLATE "pg_catalog"."default",
  "tumblr" varchar(100) COLLATE "pg_catalog"."default",
  "twitch" varchar(100) COLLATE "pg_catalog"."default",
  "twitter" varchar(100) COLLATE "pg_catalog"."default",
  "typetalk" varchar(100) COLLATE "pg_catalog"."default",
  "uber" varchar(100) COLLATE "pg_catalog"."default",
  "vk" varchar(100) COLLATE "pg_catalog"."default",
  "wepay" varchar(100) COLLATE "pg_catalog"."default",
  "xero" varchar(100) COLLATE "pg_catalog"."default",
  "yahoo" varchar(100) COLLATE "pg_catalog"."default",
  "yammer" varchar(100) COLLATE "pg_catalog"."default",
  "yandex" varchar(100) COLLATE "pg_catalog"."default",
  "zoom" varchar(100) COLLATE "pg_catalog"."default",
  "metamask" varchar(100) COLLATE "pg_catalog"."default",
  "web3onboard" varchar(100) COLLATE "pg_catalog"."default",
  "custom" varchar(100) COLLATE "pg_catalog"."default",
  "webauthnCredentials" bytea,
  "preferred_mfa_type" varchar(100) COLLATE "pg_catalog"."default",
  "recovery_codes" varchar(1000) COLLATE "pg_catalog"."default",
  "totp_secret" varchar(100) COLLATE "pg_catalog"."default",
  "mfa_phone_enabled" bool,
  "mfa_email_enabled" bool,
  "invitation" varchar(100) COLLATE "pg_catalog"."default",
  "invitation_code" varchar(100) COLLATE "pg_catalog"."default",
  "face_ids" text COLLATE "pg_catalog"."default",
  "ldap" varchar(100) COLLATE "pg_catalog"."default",
  "properties" text COLLATE "pg_catalog"."default",
  "roles" text COLLATE "pg_catalog"."default",
  "permissions" text COLLATE "pg_catalog"."default",
  "groups" varchar(1000) COLLATE "pg_catalog"."default",
  "last_change_password_time" varchar(100) COLLATE "pg_catalog"."default",
  "last_signin_wrong_time" varchar(100) COLLATE "pg_catalog"."default",
  "signin_wrong_times" int4,
  "managedAccounts" bytea,
  "mfaAccounts" bytea,
  "need_update_password" bool,
  "ip_whitelist" varchar(200) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO "public"."user" VALUES ('built-in', 'admin', '2026-03-17T01:59:46Z', '2026-03-17T01:59:46Z', '', '6c8c7645-d4f4-4793-8811-249be3e87854', '', 'd768a1d5-a00e-49f1-adf8-22a08557eeda', 'normal-user', '123', '', 'plain', 'Admin', '', '', 'https://cdn.casbin.org/img/casbin.svg', '', '', 'admin@example.com', 'f', '12345678910', 'US', '', '', '[]', 'Example Inc.', '', '', '', '', '', 'staff', '', '', '', '', 2000, 0, 1, 0, '', 0, '', 'f', 'f', 'f', 'f', 'f', 'app-built-in', '', '', '', '', '', '127.0.0.1', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', E'null', '', 'null', '', 'f', 'f', '', '', 'null', '', '{}', 'null', 'null', 'null', '', '', 0, E'null', E'null', 'f', '');
INSERT INTO "public"."user" VALUES ('costrict_group', 'costrict', '2026-03-17T10:25:50+08:00', '2026-03-17T02:26:07Z', '', 'cfed9e44-b37c-4249-94ac-af28d345507a', '', '71a1332e-3641-462e-9b7b-ce0f45198f77', 'normal-user', '123', '', 'plain', 'costrict测试账户', '', '', 'https://cdn.casbin.org/img/casbin.svg', '', '', 'p19jk8@example.com', 'f', '56657247671', 'US', '', '', '[]', 'Example Inc.', '', '', '', '', '', 'staff', '', '', '', '', 0, 0, 1, 0, '', 0, '', 'f', 'f', 'f', 'f', 'f', 'costrict_login', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', E'null', '', 'null', '', 'f', 'f', '', '', 'null', '', '{}', 'null', 'null', '[]', '', '', 0, E'null', E'null', 'f', '');

-- ----------------------------
-- Table structure for user_identity_binding
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_identity_binding";
CREATE TABLE "public"."user_identity_binding" (
  "id" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "universal_id" varchar(100) COLLATE "pg_catalog"."default",
  "auth_type" varchar(50) COLLATE "pg_catalog"."default",
  "auth_value" varchar(255) COLLATE "pg_catalog"."default",
  "created_time" varchar(100) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Records of user_identity_binding
-- ----------------------------
INSERT INTO "public"."user_identity_binding" VALUES ('b66d1d4e-6937-4e6c-b90a-b0da2e994be3', 'd768a1d5-a00e-49f1-adf8-22a08557eeda', 'password', 'built-in/admin', '2026-03-17T01:59:46Z');
INSERT INTO "public"."user_identity_binding" VALUES ('c47b02f6-087d-461b-b636-e23765e8cb2d', '71a1332e-3641-462e-9b7b-ce0f45198f77', 'password', 'costrict_group/user_p19jk8', '2026-03-17T02:25:49Z');

-- ----------------------------
-- Table structure for verification_record
-- ----------------------------
DROP TABLE IF EXISTS "public"."verification_record";
CREATE TABLE "public"."verification_record" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "remote_addr" varchar(100) COLLATE "pg_catalog"."default",
  "type" varchar(10) COLLATE "pg_catalog"."default",
  "user" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "provider" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "receiver" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "code" varchar(10) COLLATE "pg_catalog"."default" NOT NULL,
  "time" int8 NOT NULL,
  "is_used" bool NOT NULL
)
;

-- ----------------------------
-- Records of verification_record
-- ----------------------------

-- ----------------------------
-- Table structure for webhook
-- ----------------------------
DROP TABLE IF EXISTS "public"."webhook";
CREATE TABLE "public"."webhook" (
  "owner" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
  "created_time" varchar(100) COLLATE "pg_catalog"."default",
  "organization" varchar(100) COLLATE "pg_catalog"."default",
  "url" varchar(200) COLLATE "pg_catalog"."default",
  "method" varchar(100) COLLATE "pg_catalog"."default",
  "content_type" varchar(100) COLLATE "pg_catalog"."default",
  "headers" text COLLATE "pg_catalog"."default",
  "events" varchar(1000) COLLATE "pg_catalog"."default",
  "token_fields" varchar(1000) COLLATE "pg_catalog"."default",
  "object_fields" varchar(1000) COLLATE "pg_catalog"."default",
  "is_user_extended" bool,
  "single_org_only" bool,
  "is_enabled" bool
)
;

-- ----------------------------
-- Records of webhook
-- ----------------------------

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
SELECT setval('"public"."record_id_seq"', 14, true);

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
