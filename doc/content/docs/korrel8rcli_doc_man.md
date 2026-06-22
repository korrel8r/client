---
title: doc man
description: Generate man pages in directory DIR
---
## korrel8rcli doc man

Generate man pages in directory DIR

```
korrel8rcli doc man DIR [flags]
```

### Options

```
  -h, --help   help for man
```

### Options inherited from parent commands

```
  -t, --bearer-token string                         Authhorization token, default from env KORREL8RCLI_BEARER_TOKEN or kube config.
      --debug                                       Enable debug output.
  -k, --insecure                                    Insecure connection, skip TLS verification.
  -o, --output enum(yaml,json-pretty,json,ndjson)   Output format (default yaml)
  -u, --url string                                  URL of remote korrel8r, default from env KORREL8RCLI_URL (default "http://localhost:8080")
```

