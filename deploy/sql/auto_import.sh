#!/bin/bash
mysql -h w8t-mysql -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < /sql/notice_template_examples.sql
mysql -h w8t-mysql -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < /sql/rule_template_groups.sql
mysql -h w8t-mysql -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < /sql/rule_templates.sql
mysql -h w8t-mysql -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < /sql/user_roles.sql
mysql -h w8t-mysql -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < /sql/tenants.sql
mysql -h w8t-mysql -u root -pw8t.123 --default-character-set=utf8mb4 -D watchalert < /sql/tenants_linked_users.sql
