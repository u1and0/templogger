SDカードにためたバイナリデータをテキスト(json形式)にして標準出力にdumpします。

## 使用方法:
単一のファイルをJSON化
```
$ templogger data/12161037.DAT
```

複数のファイルをJSON化
```
$ templogger data/12161037.DAT data/12161237.DAT
```

すべてのDATファイルをJSON化
```
$ templogger data/*.DAT
```

-tオプションで読みやすいようにインデントを入れます
```
$ templogger -t data/*.DAT
```

