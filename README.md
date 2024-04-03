### 打包并推送
```bash
goreleaser release --rm-dist
```

### 本地测试
```bash
 goreleaser --snapshot --skip-publish --snapshot --rm-dist
```