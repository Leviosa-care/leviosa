go test -v -count=1 ./internal/adapters/...
FAIL	github.com/Leviosa-care/settings/internal/adapters/s3 [setup failed]
?   	github.com/Leviosa-care/settings/internal/adapters/http	[no test files]
2025/10/15 01:19:57 Creating postgres container.
2025/10/15 01:19:57 github.com/testcontainers/testcontainers-go - Connected to docker: 
  Server Version: 27.2.1
  API Version: 1.47
  Operating System: Ubuntu 24.04.2 LTS
  Total Memory: 6830 MB
  Testcontainers for Go Version: v0.38.0
  Resolved Docker Host: unix:///var/run/docker.sock
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: daaeaa27f90081d814095e6ce0ee8ad5b2a8ca7af3785d08c2c6ed8280218084
  Test ProcessID: bc8bbd2b-e1f5-47ab-a0b9-fcd3a7f761ee
2025/10/15 01:19:57 🐳 Creating container for image postgres:17.5-alpine3.21
2025/10/15 01:19:57 🐳 Creating container for image testcontainers/ryuk:0.12.0
2025/10/15 01:19:57 ✅ Container created: 3fd2f25c2895
2025/10/15 01:19:57 🐳 Starting container: 3fd2f25c2895
2025/10/15 01:19:57 ✅ Container started: 3fd2f25c2895
2025/10/15 01:19:57 ⏳ Waiting for container id 3fd2f25c2895 image: testcontainers/ryuk:0.12.0. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms skipInternalCheck:false skipExternalCheck:false}
2025/10/15 01:19:57 🔔 Container is ready: 3fd2f25c2895
2025/10/15 01:19:57 ✅ Container created: ea7914247a60
2025/10/15 01:19:57 🐳 Starting container: ea7914247a60
2025/10/15 01:19:58 ✅ Container started: ea7914247a60
2025/10/15 01:19:58 ⏳ Waiting for container id ea7914247a60 image: postgres:17.5-alpine3.21. Waiting for: &{Port:5432/tcp timeout:<nil> PollInterval:100ms skipInternalCheck:false skipExternalCheck:false}
2025/10/15 01:20:01 🔔 Container is ready: ea7914247a60
2025/10/15 01:20:01 Postgres container successfully created.
2025/10/15 01:20:01 Creating pgxpool...
2025/10/15 01:20:01 pgxpool created.
2025/10/15 01:20:01 Database pool ping successful.
2025/10/15 01:20:01 Applying database migrations...
2025/10/15 01:20:01 OK   20250602201515_catalog_init_schema.sql (96.37ms)
2025/10/15 01:20:01 OK   20250816092305_settings_init_schema.sql (22.31ms)
2025/10/15 01:20:01 OK   20250828120445_auth_init_schema.sql (59.29ms)
2025/10/15 01:20:01 OK   20250927213930_booking_init_schema.sql (63.15ms)
2025/10/15 01:20:01 goose: successfully migrated database to version: 20250927213930
2025/10/15 01:20:01 Migrations applied.
=== RUN   TestGetEncryptedSetting
=== RUN   TestGetEncryptedSetting/successful_retrieval_of_encrypted_setting
=== RUN   TestGetEncryptedSetting/successful_retrieval_with_empty_encrypted_value
=== RUN   TestGetEncryptedSetting/unsuccessful_retrieval_due_to_not_null_violation_thanks_to_nil_encrypted_value
=== RUN   TestGetEncryptedSetting/successful_retrieval_with_large_encrypted_data
=== RUN   TestGetEncryptedSetting/successful_retrieval_with_special_characters_in_key
=== RUN   TestGetEncryptedSetting/key_not_found_returns_repository_not_found_error
=== RUN   TestGetEncryptedSetting/context_cancellation
=== RUN   TestGetEncryptedSetting/successful_retrieval_with_zero_key_version
--- PASS: TestGetEncryptedSetting (0.03s)
    --- PASS: TestGetEncryptedSetting/successful_retrieval_of_encrypted_setting (0.01s)
    --- PASS: TestGetEncryptedSetting/successful_retrieval_with_empty_encrypted_value (0.00s)
    --- PASS: TestGetEncryptedSetting/unsuccessful_retrieval_due_to_not_null_violation_thanks_to_nil_encrypted_value (0.00s)
    --- PASS: TestGetEncryptedSetting/successful_retrieval_with_large_encrypted_data (0.00s)
    --- PASS: TestGetEncryptedSetting/successful_retrieval_with_special_characters_in_key (0.00s)
    --- PASS: TestGetEncryptedSetting/key_not_found_returns_repository_not_found_error (0.00s)
    --- PASS: TestGetEncryptedSetting/context_cancellation (0.00s)
    --- PASS: TestGetEncryptedSetting/successful_retrieval_with_zero_key_version (0.00s)
=== RUN   TestGetInt
=== RUN   TestGetInt/successful_retrieval_of_int_setting
=== RUN   TestGetInt/successful_retrieval_of_negative_int
=== RUN   TestGetInt/successful_retrieval_of_zero_value
=== RUN   TestGetInt/key_not_found_returns_repository_not_found_error
=== RUN   TestGetInt/invalid_int_value_returns_conversion_error
=== RUN   TestGetInt/float_value_stored_as_string_returns_conversion_error
=== RUN   TestGetInt/empty_string_value_returns_conversion_error
=== RUN   TestGetInt/context_cancellation
=== RUN   TestGetInt/very_large_int_value
=== RUN   TestGetInt/whitespace_in_value_returns_conversion_error
--- PASS: TestGetInt (0.03s)
    --- PASS: TestGetInt/successful_retrieval_of_int_setting (0.01s)
    --- PASS: TestGetInt/successful_retrieval_of_negative_int (0.00s)
    --- PASS: TestGetInt/successful_retrieval_of_zero_value (0.00s)
    --- PASS: TestGetInt/key_not_found_returns_repository_not_found_error (0.00s)
    --- PASS: TestGetInt/invalid_int_value_returns_conversion_error (0.00s)
    --- PASS: TestGetInt/float_value_stored_as_string_returns_conversion_error (0.00s)
    --- PASS: TestGetInt/empty_string_value_returns_conversion_error (0.00s)
    --- PASS: TestGetInt/context_cancellation (0.00s)
    --- PASS: TestGetInt/very_large_int_value (0.00s)
    --- PASS: TestGetInt/whitespace_in_value_returns_conversion_error (0.00s)
=== RUN   TestGetString
=== RUN   TestGetString/successful_retrieval_of_string_setting
=== RUN   TestGetString/successful_retrieval_of_empty_string
=== RUN   TestGetString/successful_retrieval_of_string_with_whitespace
=== RUN   TestGetString/successful_retrieval_of_string_with_special_characters
=== RUN   TestGetString/successful_retrieval_of_unicode_string
=== RUN   TestGetString/successful_retrieval_of_multiline_string
=== RUN   TestGetString/successful_retrieval_of_JSON_string
=== RUN   TestGetString/successful_retrieval_of_numeric_string
=== RUN   TestGetString/successful_retrieval_of_very_long_string
=== RUN   TestGetString/key_not_found_returns_repository_not_found_error
=== RUN   TestGetString/context_cancellation
=== RUN   TestGetString/successful_retrieval_with_tab_characters
=== RUN   TestGetString/successful_retrieval_of_single_character
=== RUN   TestGetString/successful_retrieval_of_string_with_SQL_injection_attempt
--- PASS: TestGetString (0.04s)
    --- PASS: TestGetString/successful_retrieval_of_string_setting (0.00s)
    --- PASS: TestGetString/successful_retrieval_of_empty_string (0.00s)
    --- PASS: TestGetString/successful_retrieval_of_string_with_whitespace (0.00s)
    --- PASS: TestGetString/successful_retrieval_of_string_with_special_characters (0.00s)
    --- PASS: TestGetString/successful_retrieval_of_unicode_string (0.00s)
    --- PASS: TestGetString/successful_retrieval_of_multiline_string (0.00s)
    --- PASS: TestGetString/successful_retrieval_of_JSON_string (0.00s)
    --- PASS: TestGetString/successful_retrieval_of_numeric_string (0.00s)
    --- PASS: TestGetString/successful_retrieval_of_very_long_string (0.00s)
    --- PASS: TestGetString/key_not_found_returns_repository_not_found_error (0.00s)
    --- PASS: TestGetString/context_cancellation (0.00s)
    --- PASS: TestGetString/successful_retrieval_with_tab_characters (0.00s)
    --- PASS: TestGetString/successful_retrieval_of_single_character (0.00s)
    --- PASS: TestGetString/successful_retrieval_of_string_with_SQL_injection_attempt (0.00s)
=== RUN   TestSetPhone
=== RUN   TestSetPhone/successful_phone_setting_creation
=== RUN   TestSetPhone/successful_phone_setting_with_empty_encrypted_values
=== RUN   TestSetPhone/successful_phone_setting_with_nil_encrypted_values
HERE GETTING THE RIGHT ERROR VALUE
=== RUN   TestSetPhone/context_cancellation_returns_error
=== RUN   TestSetPhone/setting_with_very_large_encrypted_values
=== RUN   TestSetPhone/setting_with_different_key_versions
=== RUN   TestSetPhone/setting_with_different_key_versions/key_version_0
=== RUN   TestSetPhone/setting_with_different_key_versions/key_version_1
=== RUN   TestSetPhone/setting_with_different_key_versions/key_version_10
=== RUN   TestSetPhone/setting_with_different_key_versions/key_version_100
=== RUN   TestSetPhone/setting_with_different_key_versions/key_version_999
--- PASS: TestSetPhone (0.03s)
    --- PASS: TestSetPhone/successful_phone_setting_creation (0.01s)
    --- PASS: TestSetPhone/successful_phone_setting_with_empty_encrypted_values (0.01s)
    --- PASS: TestSetPhone/successful_phone_setting_with_nil_encrypted_values (0.00s)
    --- PASS: TestSetPhone/context_cancellation_returns_error (0.00s)
    --- PASS: TestSetPhone/setting_with_very_large_encrypted_values (0.00s)
    --- PASS: TestSetPhone/setting_with_different_key_versions (0.02s)
        --- PASS: TestSetPhone/setting_with_different_key_versions/key_version_0 (0.01s)
        --- PASS: TestSetPhone/setting_with_different_key_versions/key_version_1 (0.00s)
        --- PASS: TestSetPhone/setting_with_different_key_versions/key_version_10 (0.00s)
        --- PASS: TestSetPhone/setting_with_different_key_versions/key_version_100 (0.00s)
        --- PASS: TestSetPhone/setting_with_different_key_versions/key_version_999 (0.00s)
=== RUN   TestSetInt
=== RUN   TestSetInt/successful_insertion
=== RUN   TestSetInt/successful_insertion_with_negative_value
=== RUN   TestSetInt/successful_insertion_with_zero_value
=== RUN   TestSetInt/duplicate_key_insertion_should_update_existing_value
=== RUN   TestSetInt/nil_setting_should_panic_or_cause_error
=== RUN   TestSetInt/context_cancellation_should_return_error
=== RUN   TestSetInt/context_timeout_should_return_error
=== RUN   TestSetInt/empty_key_should_work_if_allowed_by_schema
=== RUN   TestSetInt/very_large_integer_values
=== RUN   TestSetInt/database_connection_closed_should_return_error
    set_int_test.go:240: Database connection failure test requires specific setup
--- PASS: TestSetInt (0.02s)
    --- PASS: TestSetInt/successful_insertion (0.00s)
    --- PASS: TestSetInt/successful_insertion_with_negative_value (0.00s)
    --- PASS: TestSetInt/successful_insertion_with_zero_value (0.00s)
    --- PASS: TestSetInt/duplicate_key_insertion_should_update_existing_value (0.00s)
    --- PASS: TestSetInt/nil_setting_should_panic_or_cause_error (0.00s)
    --- PASS: TestSetInt/context_cancellation_should_return_error (0.00s)
    --- PASS: TestSetInt/context_timeout_should_return_error (0.00s)
    --- PASS: TestSetInt/empty_key_should_work_if_allowed_by_schema (0.00s)
    --- PASS: TestSetInt/very_large_integer_values (0.00s)
    --- SKIP: TestSetInt/database_connection_closed_should_return_error (0.00s)
=== RUN   TestSetString
=== RUN   TestSetString/successful_insertion
=== RUN   TestSetString/successful_insertion_with_empty_string
=== RUN   TestSetString/successful_insertion_with_special_characters
=== RUN   TestSetString/successful_insertion_with_unicode_characters
=== RUN   TestSetString/successful_insertion_with_multiline_string
=== RUN   TestSetString/successful_insertion_with_very_long_string
=== RUN   TestSetString/successful_insertion_with_JSON-like_string
=== RUN   TestSetString/duplicate_key_insertion_should_update_existing_value
=== RUN   TestSetString/setting_existing_migration_values_should_update_them
=== RUN   TestSetString/nil_setting_should_panic_or_cause_error
=== RUN   TestSetString/context_cancellation_should_return_error
=== RUN   TestSetString/context_timeout_should_return_error
=== RUN   TestSetString/empty_key_should_work_if_allowed_by_schema
=== RUN   TestSetString/whitespace-only_key_should_work_if_allowed
=== RUN   TestSetString/database_connection_closed_should_return_error
    set_string_test.go:377: Database connection failure test requires specific setup
--- PASS: TestSetString (0.10s)
    --- PASS: TestSetString/successful_insertion (0.00s)
    --- PASS: TestSetString/successful_insertion_with_empty_string (0.00s)
    --- PASS: TestSetString/successful_insertion_with_special_characters (0.00s)
    --- PASS: TestSetString/successful_insertion_with_unicode_characters (0.00s)
    --- PASS: TestSetString/successful_insertion_with_multiline_string (0.00s)
    --- PASS: TestSetString/successful_insertion_with_very_long_string (0.05s)
    --- PASS: TestSetString/successful_insertion_with_JSON-like_string (0.01s)
    --- PASS: TestSetString/duplicate_key_insertion_should_update_existing_value (0.01s)
    --- PASS: TestSetString/setting_existing_migration_values_should_update_them (0.01s)
    --- PASS: TestSetString/nil_setting_should_panic_or_cause_error (0.00s)
    --- PASS: TestSetString/context_cancellation_should_return_error (0.00s)
    --- PASS: TestSetString/context_timeout_should_return_error (0.00s)
    --- PASS: TestSetString/empty_key_should_work_if_allowed_by_schema (0.00s)
    --- PASS: TestSetString/whitespace-only_key_should_work_if_allowed (0.01s)
    --- SKIP: TestSetString/database_connection_closed_should_return_error (0.00s)
PASS
2025/10/15 01:20:01 Test(s) executed
ok  	github.com/Leviosa-care/settings/internal/adapters/postgres	4.804s
?   	github.com/Leviosa-care/settings/internal/adapters/rabbitmq	[no test files]
FAIL
