package utils

// TODO: Step 1 - Implement string utilities:
//   - GenerateUUID() string - Generate UUID v4
//   - GenerateRandomString(length int) string - Cryptographically secure random strings
//   - SlugifyString(input string) string - URL-safe slug generation
//   - TruncateString(s string, length int) string - Truncate with ellipsis
//   - SanitizeHTML(input string) string - Remove/escape HTML tags
//   - ValidateEmail(email string) bool - Email format validation
//   - MaskEmail(email string) string - Mask email for privacy (j***@example.com)
// TODO: Step 2 - Implement time utilities:
//   - ParseTimeWithTimezone(timeStr, timezone string) (time.Time, error)
//   - FormatDuration(duration time.Duration) string - Human-readable duration
//   - TimeAgo(t time.Time) string - "2 hours ago" format
//   - BusinessDaysInRange(start, end time.Time) int - Calculate business days
//   - AddBusinessDays(t time.Time, days int) time.Time - Add business days only
// TODO: Step 3 - Implement validation utilities:
//   - ValidateStruct(s interface{}) error - Struct validation using tags
//   - ValidatePassword(password string) error - Password strength validation
//   - ValidatePhoneNumber(phone string) bool - Phone number format validation
//   - ValidateURL(url string) bool - URL format validation
//   - SanitizeInput(input string) string - General input sanitization
// TODO: Step 4 - Implement conversion utilities:
//   - ToJSON(v interface{}) ([]byte, error) - Safe JSON marshaling
//   - FromJSON(data []byte, v interface{}) error - Safe JSON unmarshaling
//   - StringToInt(s string, defaultValue int) int - Safe string to int conversion
//   - StringToBool(s string) bool - String to boolean conversion
//   - MapToStruct(m map[string]interface{}, s interface{}) error - Map to struct conversion
// TODO: Step 5 - Implement crypto utilities:
//   - HashSHA256(data []byte) string - SHA256 hashing
//   - GenerateHMAC(message, secret string) string - HMAC generation
//   - EncryptAES(plaintext, key string) (string, error) - AES encryption
//   - DecryptAES(ciphertext, key string) (string, error) - AES decryption
//   - GenerateRSAKeyPair() (privateKey, publicKey string, error) - RSA key generation
// TODO: Step 6 - Implement file utilities:
//   - EnsureDir(path string) error - Create directory if not exists
//   - WriteJSONFile(path string, data interface{}) error - Write JSON to file
//   - ReadJSONFile(path string, v interface{}) error - Read JSON from file
//   - GetFileSize(path string) (int64, error) - Get file size safely
//   - IsValidImageType(filename string) bool - Validate image file extensions
// TODO: Step 7 - Implement slice/array utilities:
//   - Contains(slice []string, item string) bool - Check if slice contains item
//   - RemoveDuplicates(slice []string) []string - Remove duplicate strings
//   - ChunkSlice(slice []interface{}, size int) [][]interface{} - Split slice into chunks
//   - FilterSlice(slice []string, predicate func(string) bool) []string - Filter slice items

// Flow: All services -> shared utilities -> common operations and validations
// Dependencies: crypto/rand, encoding/json, time, regexp, crypto packages
// Used by: All microservices for common operations and data processing
