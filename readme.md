# 4H Recordbook Backend

## Deployment Secret Files

### `/internal/config/config.json`

```json
{
    "max_page_size": 512,
    "cosmos": {
        "production": {
            "endpoint": [...],
            "key": [...]
        },
        "development": {
            "endpoint": [...],
            "key": [...]
        }
    },
    "upc": {
        "endpoint": [...],
        "production": {
            "key": [...]
        },
        "development": {
            "key": [...]
        }
    },
    "auth0": {
        "domain": [...],
        "audience": [...]
    }
}
```
