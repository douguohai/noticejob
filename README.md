### 打包并推送
```bash
goreleaser release --rm-dist
```

### 本地测试
```bash
 goreleaser --snapshot --skip-publish --snapshot --rm-dist
```

### 使用
```bash
./noticejob image  -e dev -v v20221024 -n 测试项目 -t 1e336794c4b90331ba4999d -u  "192.168.10.239:8888/charles0320/xyx-zudp-auth:dev-2024-04-03-152841"
```

### 帮助文档
```bash
message push to Ding talk

Usage:
  push [flags]
  push [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  deploy      deploy notice
  help        Help about any command
  image       just build image notice
  package     package notice
  public      public message notice

Flags:
  -h, --help   help for push

Use "push [command] --help" for more information about a command.
```
