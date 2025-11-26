/*
 Navicat Premium Dump SQL - Data Insertion Only

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
-- Records of adapter
-- ----------------------------
INSERT INTO "public"."adapter" VALUES ('built-in', 'api-adapter-built-in', '2025-07-31T18:18:34Z', 'casbin_api_rule', 'f', '', '', '', 0, '', '', '');
INSERT INTO "public"."adapter" VALUES ('built-in', 'user-adapter-built-in', '2025-07-31T18:18:34Z', 'casbin_user_rule', 'f', '', '', '', 0, '', '', '');

-- ----------------------------
-- Records of application
-- ----------------------------
INSERT INTO "public"."application" VALUES ('admin', 'loginApp', '2025-08-01T02:42:53+08:00', 'loginApp', 'https://costrict.ai/plugin/show/costrict.svg', '', '', 'user-group', 'cert-built-in', '', '', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', '', '', '[{"owner":"","name":"Oauth","canSignUp":true,"canSignIn":true,"canUnlink":true,"countryCodes":null,"prompted":false,"signupGroup":"","rule":"None","provider":null},{"owner":"","name":"SMS","canSignUp":true,"canSignIn":true,"canUnlink":true,"countryCodes":null,"prompted":false,"signupGroup":"","rule":"All","provider":null}]', '[{"name":"Verification code","displayName":"Verification code","rule":"Phone only"},{"name":"Password","displayName":"Password","rule":"All"}]', '[{"name":"ID","visible":false,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"Random"},{"name":"Username","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Display name","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Password","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Confirm password","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Phone","visible":true,"required":false,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"No verification"},{"name":"Agreement","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Signup button","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Providers","visible":true,"required":true,"prompted":false,"type":"","customCss":".provider-img {\n width: 30px;\n margin: 5px;\n }\n .provider-big-img {\n margin-bottom: 10px;\n }\n ","label":"","placeholder":"","options":null,"regex":"","rule":"small"}]', '[{"name":"Back button","visible":true,"label":"","customCss":".back-button {\n      top: 65px;\n      left: 15px;\n      position: absolute;\n}\n.back-inner-button{}","placeholder":"","rule":"None","isCustom":false},{"name":"Languages","visible":false,"label":"","customCss":".login-languages {\n    top: 55px;\n    right: 5px;\n    position: absolute;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Logo","visible":true,"label":"","customCss":".login-logo-box {}","placeholder":"","rule":"None","isCustom":false},{"name":"Signin methods","visible":true,"label":"","customCss":".signin-methods {}","placeholder":"","rule":"None","isCustom":false},{"name":"Username","visible":true,"label":"","customCss":".login-username {}\n.login-username-input{}","placeholder":"","rule":"None","isCustom":false},{"name":"Password","visible":true,"label":"","customCss":".login-password {}\n.login-password-input{}","placeholder":"","rule":"None","isCustom":false},{"name":"Agreement","visible":false,"label":"","customCss":".login-agreement {}","placeholder":"","rule":"None","isCustom":false},{"name":"Forgot password?","visible":false,"label":"","customCss":".login-forget-password {\n    display: inline-flex;\n    justify-content: space-between;\n    width: 320px;\n    margin-bottom: 25px;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Login button","visible":true,"label":"","customCss":".login-button-box {\n    margin-bottom: 5px;\n}\n.login-button {\n    width: 100%;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Signup link","visible":false,"label":"","customCss":".login-signup-link {\n    margin-bottom: 24px;\n    display: flex;\n    justify-content: end;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Providers","visible":true,"label":"","customCss":".provider-img {\n      width: 30px;\n      margin: 5px;\n}\n.provider-big-img {\n      margin-bottom: 10px;\n}","placeholder":"","rule":"small","isCustom":false}]', '["authorization_code","password","client_credentials","token","id_token","refresh_token"]', '[]', 'null', 'f', '', '{{OIDC_AUTH_CLIENT_ID}}', '{{OIDC_AUTH_CLIENT_SECRET}}', '["{{COSTRICT_BACKEND_BASEURL}}/oidc-auth/api/v1/plugin/login/callback"]', '', 'JWT', '', '[]', 1200, 1200, '', '', '', '', '', '', '', '', NULL, '<style>\n    #footer {\n        display: none;\n    }\n<style>\n  ', '', '', 2, '', '', '', 5, 15);
INSERT INTO "public"."application" VALUES ('admin', 'app-built-in', '2025-07-31T18:18:32Z', 'managerApp', 'https://costrict.ai/plugin/show/costrict.svg', '', '', 'built-in', 'cert-built-in', '', '', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', 'f', '', '', '[]', '[{"name":"Password","displayName":"Password","rule":"All"}]', '[{"name":"ID","visible":false,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"Random"},{"name":"Username","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Display name","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Password","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Confirm password","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Email","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"Normal"},{"name":"Phone","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Agreement","visible":true,"required":true,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":"None"},{"name":"Signup button","visible":true,"required":false,"prompted":false,"type":"","customCss":"","label":"","placeholder":"","options":null,"regex":"","rule":""}]', '[{"name":"Back button","visible":true,"label":"","customCss":".back-button {\n      top: 65px;\n      left: 15px;\n      position: absolute;\n}\n.back-inner-button{}","placeholder":"","rule":"None","isCustom":false},{"name":"Languages","visible":false,"label":"","customCss":".login-languages {\n    top: 55px;\n    right: 5px;\n    position: absolute;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Logo","visible":true,"label":"","customCss":".login-logo-box {}","placeholder":"","rule":"None","isCustom":false},{"name":"Signin methods","visible":true,"label":"","customCss":".signin-methods {}","placeholder":"","rule":"None","isCustom":false},{"name":"Username","visible":true,"label":"","customCss":".login-username {}\n.login-username-input{}","placeholder":"","rule":"None","isCustom":false},{"name":"Password","visible":true,"label":"","customCss":".login-password {}\n.login-password-input{}","placeholder":"","rule":"None","isCustom":false},{"name":"Agreement","visible":false,"label":"","customCss":".login-agreement {}","placeholder":"","rule":"None","isCustom":false},{"name":"Forgot password?","visible":false,"label":"","customCss":".login-forget-password {\n    display: inline-flex;\n    justify-content: space-between;\n    width: 320px;\n    margin-bottom: 25px;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Login button","visible":true,"label":"","customCss":".login-button-box {\n    margin-bottom: 5px;\n}\n.login-button {\n    width: 100%;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Signup link","visible":false,"label":"","customCss":".login-signup-link {\n    margin-bottom: 24px;\n    display: flex;\n    justify-content: end;\n}","placeholder":"","rule":"None","isCustom":false},{"name":"Providers","visible":true,"label":"","customCss":".provider-img {\n      width: 30px;\n      margin: 5px;\n}\n.provider-big-img {\n      margin-bottom: 10px;\n}","placeholder":"","rule":"small","isCustom":false}]', '["authorization_code"]', '[]', 'null', 'f', '', '3b215b8c312b0c47ccd7', 'e2f78c894cafc828a0b0161059d2c0b8d81deac5', '[]', '', 'JWT', '', '[]', 168, 0, '', '', '', '', '', '', '', '', NULL, '', '', '', 2, '', '', '', 5, 15);

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
-- Records of casbin_rule
-- ----------------------------

-- ----------------------------
-- Records of casbin_user_rule
-- ----------------------------

-- ----------------------------
-- Records of cert
-- ----------------------------
INSERT INTO "public"."cert" VALUES ('admin', 'cert-built-in', '2025-07-31T18:18:32Z', 'Built-in Cert', 'JWT', 'x509', 'RS256', 4096, 20, '-----BEGIN CERTIFICATE-----
MIIE3TCCAsWgAwIBAgIDAeJAMA0GCSqGSIb3DQEBCwUAMCgxDjAMBgNVBAoTBWFk
bWluMRYwFAYDVQQDEw1jZXJ0LWJ1aWx0LWluMB4XDTI1MDczMTE4MTgzNFoXDTQ1
MDczMTE4MTgzNFowKDEOMAwGA1UEChMFYWRtaW4xFjAUBgNVBAMTDWNlcnQtYnVp
bHQtaW4wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC2uJL1rSjrkOTv
B86kXOX8heyRC2oHAjih8V+7/OQIsI69/4azcDm1yGRSqxUAgb8VSYqw/hs3sPLv
8XdgAUFSuThh3hqLiWgIjKpju/+1h2W7j2+wog5g9SzRqhjVL5Nxk/QMeFzhhlzy
kzNhBEolXs5fd33KcLiaVhE10AqgJ+9zkonwLz/ZyyxtDMiug4Wi9Nd403Oc7S3T
urF4CaQUi+488SRhR/ltJ8v7OC0+tB7EWYU7ipBk+RgQssNRowb8HDPxDQBh+d/e
z1NNjsRAf5iDWKZSTFy1dPnCA5KrL7T+FWhM+zpS9oA5IiPdIiP2QE1jhN/Fk9TL
YqpQ09lcKXJSUGy9nrWu0Fk3qCRtF5yFskPSEaOHlqKJxh0YoNrV/KpQsOHcHlzE
VGu16TADnbxcNaO+BOt3/ZcZiRbW1g3/bYWsUJOK8sUGhVx2rmxhTQDXlEBcRG1V
+zPDw7H9VFMR2ZDEDTKjQvK/krSKd2RIUX/JSASO+oZSw0o2pyI5O79vtw5cB3TR
3YEsXM8+gIbAUAFplvcBDXuALT9tsnjpf2B2I7eq7YMQAvo2TTlncYdnqTPdmgT2
5hAnny2o1C9WLUGr7O20o+0gIeNj/tjdD3on08lIIKlNxkD2rBnkcP64VuKxAX/4
gYq3YYFokYdsMuFM03pZFUrnAEY7WQIDAQABoxAwDjAMBgNVHRMBAf8EAjAAMA0G
CSqGSIb3DQEBCwUAA4ICAQCSDOAHlyhrwZrcrP85T79AmhPq8Hhj0HmmR8cdiMb3
RuzRCCkWHolK1NVnnRLhgVBG2OJH/PyZ05jRmlsAvaG8RCgcJyRJ9gFdOo3UlerT
fb+VfZrGAwxVFuywGygctNVgN1arZ0nil2TAVQ93Z+ENZpRgKoca5xYXEhzWfJ0T
xrnZM6ZYC+fjRfXoqvQRDHpoUvlkeWvqIswyr3fBK19JJoPGqgOa7iziqJbBdAs1
vZ6gJoeEiocgb2qaETQi8GYymR3tHhnnCW7NC2ZRxGlQsFGGmhZwe5vSbA8jtrtV
6owiQRT6arcp+nTQKPJwqPrGYqA0XF3uAbipjctapHirs+Pp3dVeSAeK7DNOVtSy
771Pv2QhXjeo5kPwaMaT899jzcjzj6SpcZOio4cLZs7yYwuWJ27zgvN/Msr64QdK
gh2R3uJLoORVJQRW/hCIYDo41e1kvUzWsIqnaCIonwtdaDzHdp6Q+7Fn90xX4Qq6
eJ05byQxIduUNhII+8d8d7Cs0VxP8GMFEiGtwP4JPMU3d/Hp0AAgLWHe/92nkQZ1
aFBBx31qIrrnpc0NWcrS/K5LibsnpRReHU5wPTrOf9HweOQ0HviVFmFcl17UfKVH
knO89BoIM+PAxgRaCB0MPgPDPajdUFWup1craNB691D74ZoLgttIiKY1Kkf3M/hZ
WQ==
-----END CERTIFICATE-----
', '-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
');

-- ----------------------------
-- Records of enforcer
-- ----------------------------
INSERT INTO "public"."enforcer" VALUES ('built-in', 'api-enforcer-built-in', '2025-07-31T18:18:34Z', '2025-07-31 18:18:34', 'API Enforcer', '', 'built-in/api-model-built-in', 'built-in/api-adapter-built-in', NULL);
INSERT INTO "public"."enforcer" VALUES ('built-in', 'user-enforcer-built-in', '2025-07-31T18:18:34Z', '2025-07-31 18:18:34', 'User Enforcer', '', 'built-in/user-model-built-in', 'built-in/user-adapter-built-in', NULL);

-- ----------------------------
-- Records of group
-- ----------------------------

-- ----------------------------
-- Records of invitation
-- ----------------------------
INSERT INTO "public"."invitation" VALUES ('built-in', 'invitation_i3bahu', '2025-08-05T14:40:22+08:00', '2025-08-05T14:40:22+08:00', 'New Invitation - i3bahu', '62fgg95y8n', 'f', 1, 0, 'All', '', '', '', '', '62fgg95y8n', 'Active');

-- ----------------------------
-- Records of ldap
-- ----------------------------
INSERT INTO "public"."ldap" VALUES ('ldap-built-in', 'built-in', '2025-07-31T18:18:34Z', 'BuildIn LDAP Server', 'example.com', 389, 'f', 'f', 'cn=buildin,dc=example,dc=com', '123', 'ou=BuildIn,dc=example,dc=com', '', 'null', '', '', 0, '');

-- ----------------------------
-- Records of model
-- ----------------------------
INSERT INTO "public"."model" VALUES ('built-in', 'api-model-built-in', '2025-07-31T18:18:34Z', 'API Model', '', '[request_definition]
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
INSERT INTO "public"."model" VALUES ('built-in', 'user-model-built-in', '2025-07-31T18:18:34Z', 'Built-in Model', '', '[request_definition]
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
-- Records of organization
-- ----------------------------
INSERT INTO "public"."organization" VALUES ('admin', 'built-in', '2025-07-31T18:18:31Z', 'manager', '', 'https://costrict.ai/plugin/show/costrict.svg', '', 'https://costrict.ai/plugin/show/costrict.svg', 'f', 'plain', '', '["AtLeast6"]', '', '', 0, '["US","ES","FR","DE","GB","CN","JP","KR","VN","ID","SG","IN"]', '', '', '[]', '[]', '["en","zh","es","fr","de","id","ja","ko","ru","vi","pt"]', NULL, '', '', '', '', 2000, 'f', 'f', 'f', 'f', '', 'null', '["theme","ai-assistant","language"]', 'null', '[{"name":"Organization","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"ID","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"Name","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Display name","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Avatar","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"User type","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Password","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Email","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Phone","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Country code","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Country/Region","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Location","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Affiliation","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Title","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Signup application","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Roles","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"Permissions","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"Groups","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"3rd-party logins","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Properties","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is admin","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is forbidden","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is deleted","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Multi-factor authentication","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"WebAuthn credentials","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Managed accounts","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"MFA accounts","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""}]');
INSERT INTO "public"."organization" VALUES ('admin', 'user-group', '2025-08-01T02:39:12+08:00', 'costrict', '', 'https://costrict.ai/plugin/show/costrict.svg', '', 'https://costrict.ai/plugin/show/costrict.svg', 'f', 'plain', '', '["AtLeast6"]', 'Plain', '', 0, '["CN"]', '', '', 'null', '[]', '["en","es","fr","de","zh","id","ja","ko","ru","vi","pt","it","ms","tr","ar","he","nl","pl","fi","sv","uk","kk","fa","cs","sk"]', NULL, '', '', '', '', 0, 'f', 'f', 'f', 'f', '', '["/home-top","/orgs-top","/","/shortcuts","/apps","/organizations","/groups","/users","/invitations","/applications-top","/applications","/providers","/resources","/certs"]', 'null', 'null', '[{"name":"Organization","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"ID","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"Name","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Display name","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Avatar","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"User type","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Password","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Email","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Phone","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Country code","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Country/Region","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Location","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Address","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Affiliation","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Title","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"ID card type","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"ID card","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"ID card info","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Homepage","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Bio","visible":true,"viewRule":"Public","modifyRule":"Self","regex":""},{"name":"Tag","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Language","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Gender","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Birthday","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Education","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Score","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Karma","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Ranking","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Signup application","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"API key","visible":false,"viewRule":"","modifyRule":"Self","regex":""},{"name":"Groups","visible":true,"viewRule":"Public","modifyRule":"Admin","regex":""},{"name":"Roles","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"Permissions","visible":true,"viewRule":"Public","modifyRule":"Immutable","regex":""},{"name":"3rd-party logins","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Properties","visible":false,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is online","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is admin","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is forbidden","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Is deleted","visible":true,"viewRule":"Admin","modifyRule":"Admin","regex":""},{"name":"Multi-factor authentication","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"WebAuthn credentials","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"Managed accounts","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""},{"name":"MFA accounts","visible":true,"viewRule":"Self","modifyRule":"Self","regex":""}]');

-- ----------------------------
-- Records of payment
-- ----------------------------

-- ----------------------------
-- Records of permission
-- ----------------------------
INSERT INTO "public"."permission" VALUES ('built-in', 'permission-built-in', '2025-07-31T18:18:31Z', 'Built-in Permission', 'Built-in Permission', '["built-in/*"]', '[]', '[]', '[]', 'user-model-built-in', '', 'Application', '["app-built-in"]', '["Read","Write","Admin"]', 'Allow', 'f', 'admin', 'admin', '2025-07-31T18:18:31Z', 'Approved');

-- ----------------------------
-- Records of permission_rule
-- ----------------------------
INSERT INTO "public"."permission_rule" VALUES (1, 'p', 'built-in/*', 'app-built-in', 'read', 'allow', '', 'built-in/permission-built-in');
INSERT INTO "public"."permission_rule" VALUES (2, 'p', 'built-in/*', 'app-built-in', 'write', 'allow', '', 'built-in/permission-built-in');
INSERT INTO "public"."permission_rule" VALUES (3, 'p', 'built-in/*', 'app-built-in', 'admin', 'allow', '', 'built-in/permission-built-in');

-- ----------------------------
-- Records of plan
-- ----------------------------

-- ----------------------------
-- Records of pricing
-- ----------------------------

-- ----------------------------
-- Records of product
-- ----------------------------

-- ----------------------------
-- Records of provider
-- ----------------------------
INSERT INTO "public"."provider" VALUES ('user-group', 'Oauth', '2025-08-04T10:35:22+08:00', 'Oauth', 'OAuth', 'Custom', '', 'Normal', '1239280978', '49a2e85e8fbe81ce5bf768889c8e2a9b', '', '', '', '', '', '', '', 'openid profile email', '{"avatarUrl":"","displayName":"username","email":"phone_number","id":"employee_number","username":"username"}', 'null', '', 0, 'f', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', 'f', '', '');
INSERT INTO "public"."provider" VALUES ('user-group', 'SMS', '2025-08-01T02:40:43+08:00', 'SMS', 'SMS', 'Custom HTTP SMS', '', 'POST', '', '', '', '', '', '', '', '', '', '', '{"avatarUrl":"avatarUrl","displayName":"displayName","email":"email","id":"id","username":"username"}', 'null', '', 0, 'f', 'code', '', '', '', '', '', '', '', '', '', '', '', '', '', '', 'f', '', '');

-- ----------------------------
-- Records of radius_accounting
-- ----------------------------

-- ----------------------------
-- Records of record
-- ----------------------------

-- ----------------------------
-- Records of resource
-- ----------------------------

-- ----------------------------
-- Records of role
-- ----------------------------

-- ----------------------------
-- Records of session
-- ----------------------------
INSERT INTO "public"."session" VALUES ('built-in', 'admin', 'app-built-in', '2025-08-06T01:09:24Z', '["f187231a9ac757d4922702a9601dab24"]');
INSERT INTO "public"."session" VALUES ('user-group', 'demo', 'loginApp', '2025-08-06T01:44:58Z', '["f187231a9ac757d4922702a9601dab24"]');

-- ----------------------------
-- Records of subscription
-- ----------------------------

-- ----------------------------
-- Records of syncer
-- ----------------------------

-- ----------------------------
-- Records of token
-- ----------------------------

-- ----------------------------
-- Records of transaction
-- ----------------------------

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO "public"."user" VALUES ('built-in', 'admin', '2025-07-31T18:18:32Z', '2025-07-31T18:18:32Z', '', 'a46161cb-8d7a-47ef-94ad-02a3d828d856', '', 'd919d6d8-c0e8-49d9-8a90-3880f55bcb91', 'normal-user', '123', '', 'plain', 'Admin', '', '', 'https://cdn.casbin.org/img/casbin.svg', '', '', 'admin@example.com', 'f', '12345678910', 'US', '', '', '[]', 'Example Inc.', '', '', '', '', '', 'staff', '', '', '', '', 2000, 0, 1, 0, '', 'f', 'f', 'f', 'f', 'f', 'app-built-in', '', '', '', '', '', '127.0.0.1', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', E'null', '', 'null', '', 'f', 'f', '', '', 'null', '', '{}', 'null', 'null', 'null', '', '', 0, E'null', E'null', 'f', '');
INSERT INTO "public"."user" VALUES ('user-group', 'demo', '2025-08-05T02:25:13Z', '2025-08-05T15:43:41+08:00', '', '3f88bbf1-671b-496d-93dc-4e81c022d943', '', '39b67474-a3f8-4114-8ba0-9c57f6d09596', 'normal-user', 'test123', '', 'plain', 'demo', '', '', 'https://cdn.casbin.org/img/casbin.svg', '', '', '', 'f', '13410000000', 'CN', '', '', '[]', '', '', '', '', '', '', '', '', '', '', '', 0, 0, 1, 0, '', 'f', 'f', 'f', 'f', 'f', 'loginApp', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', '', E'null', '', 'null', '', 'f', 'f', '', '', 'null', '', '{}', 'null', 'null', '[]', '', '', 0, E'null', E'null', 'f', '');

-- ----------------------------
-- Records of user_identity_binding
-- ----------------------------
INSERT INTO "public"."user_identity_binding" VALUES ('364a1db3-8373-4af6-a5a9-781dbc8cbe53', 'd919d6d8-c0e8-49d9-8a90-3880f55bcb91', 'password', 'built-in/admin', '2025-07-31T18:18:32Z');
INSERT INTO "public"."user_identity_binding" VALUES ('d38b8690-f846-4251-81de-2835365369e6', '39b67474-a3f8-4114-8ba0-9c57f6d09596', 'password', 'user-group/test-user', '2025-08-05T02:25:13Z');

-- ----------------------------
-- Records of verification_record
-- ----------------------------

-- ----------------------------
-- Records of webhook
-- ----------------------------
