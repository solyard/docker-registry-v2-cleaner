# Docker Registry V2 Helper
---
This simple tool helps you to mark images filtered by glob as need to be removed in next cycle of garbage collecting

To use this tool just run `docker-registry-v2-helper` with specified params:
```bash
root@machine:$ docker-registry-v2-helper -h

Usage of docker-registry-v2-helper:
  -glob string
        Glob for image filtering by tags (If empty will be match nothing)
  -password string
        Please set Password if basic auth enabled for Docker Registry V2
  -url string
        Docker Registry V2 address (HTTPS / HTTP)  (default "https://registry.docker.io")
  -username string
        Please set Username if basic auth enabled for Docker Registry V2
```

Example params for remove all images tags that match glob `develop-*`:
```bash
root@machine:$ docker-registry-v2-helper -glob "develop-*" -username "user" -password "pass" -url "https://mycool.registry.com"`
```

> ⚠️ After script is complete their work please notice that you need to run garbage collection task on your registry to clear disk space:
`registry garbage-collection /etc/docker/registry/config.yml`

---
## TODO:
- [x] Remove images tags by glob
- [ ] Remove images tags between two timestamps (period)
- [ ] Remove images tags by regexp
