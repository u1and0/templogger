SDカードにためたバイナリデータをテキスト(JSON形式)にして標準出力にdumpします。

## 使い方(Windows)
変換する目的のDATファイルを`templogger.exe`へドラッグ & ドロップしてください。 (複数選択可)

## Usage
単一のファイルをCSV化(data/12161037.csvが作成されます)

```
$ templogger data/12161037.DAT
```

複数のファイルをCSV化(data/12161037.csvが作成されます)

```
$ templogger data/12161037.DAT data/12161137.DAT
```

単一のファイルをJSON化

```
$ templogger -f json data/12161037.DAT
```

複数のファイルをJSON化

```
$ templogger --format json data/12161037.DAT data/12161237.DAT
```

すべてのDATファイルをJSON化

```
$ templogger --format json data/*.DAT
```

-tオプションで読みやすいようにインデントを入れます

```
$ templogger --format json -t data/*.DAT
```

## Options
* -f, -format: dump format "csv" or "json" (default "csv")
* -h, -help: show help message
* -t, -indent: indent to format output (must use with "--format json")
* -v, -version: show version
