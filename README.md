# kgs

key generation service for `tinyURL`

## Environment Variable Examples

### INIT_REDIS_FREE
NOTICE: `kgs` will initial all possible `tinyURL` (alphanumeric, `[a-zA-Z0-9]`) in `redis_free`, which will block.

Determine if init redis free, if not it will add all possible tinyURL in redis.

"false"

### KEY_LENGTH

The length of tinyURL

"4" 

### REDIS_FREE_PASSWORD

"supersecretpassword"
