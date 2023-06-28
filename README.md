# dry

*本项目仅限于安全研究和教学，严禁用于非法用途！*

## Features

- 跨平台支持
- 多元化输入
- 定制化输出
- 自定义规则

## Install

### Binary

```url
https://github.com/DongHuangT1/dry/releases
```

### Golang

```bash
go install github.com/DongHuangT1/dry@latest
```

### Source

```bash
git clone https://github.com/DongHuangT1/dry.git
cd dry && make
```

## Usage

```
Usage: dry [-options] [- | url | file] [rules...]
  -0 string
        Set output prefix
  -9 string
        Set output suffix
  -d string
        Decompress (lzw, gzip, zlib, bzip2, flate)

Support: Base64 CIDR Domain Domain:Port Email Hash Hex IP IP:Port JWT MAC MD5 Phone SHA1 SHA224 SHA256 SHA384 SHA512 URL URL+ UUID Unicode
```

## Screenshot

![](https://i.postimg.cc/fystYXkN/dry.jpg)

## License

[GPL 3.0](https://github.com/DongHuangT1/dry/blob/master/LICENSE)
