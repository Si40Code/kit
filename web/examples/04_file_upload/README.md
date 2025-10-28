# 文件上传示例

展示如何处理文件上传，包括单文件和多文件上传。

## 功能

- 单文件上传
- 多文件上传
- 文件元信息记录

## 准备上传目录

```bash
mkdir -p uploads
```

## 运行

```bash
go run main.go
```

## 测试

### 单文件上传

```bash
curl -X POST http://localhost:8080/upload/single \
  -F "file=@/path/to/your/file.txt"
```

### 多文件上传

```bash
curl -X POST http://localhost:8080/upload/multiple \
  -F "files=@/path/to/file1.txt" \
  -F "files=@/path/to/file2.txt"
```
