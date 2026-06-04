---
title: completion
description: Generate the autocompletion script for the specified shell
---
## korrel8rcli completion

Generate the autocompletion script for the specified shell

### Synopsis

Generate the autocompletion script for korrel8rcli for the specified shell.
See each sub-command's help for details on how to use the generated script.


### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
  -t, --bearer-token string                         Authhorization token, default from env KORREL8RCLI_BEARER_TOKEN or kube config.
      --debug                                       Enable debug output.
  -k, --insecure                                    Insecure connection, skip TLS verification.
  -o, --output enum(yaml,json-pretty,json,ndjson)   Output format (default yaml)
  -u, --url string                                  URL of remote korrel8r, default from env KORREL8RCLI_URL (default "http://localhost:8080")
```

### SEE ALSO
