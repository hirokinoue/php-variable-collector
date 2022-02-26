# php-variable-collector
PHPの変数を収集するツール

## インストール方法
本ツールはGolangで書かれている。
`main.go`をコンパイルして、実行ファイルを生成する。

## 使い方
`./main -in testdata -out out -exclude vendor`  
`-in`で指定したディレクトリ配下にある（`-exclude`で指定したディレクトリは無視する）phpのコードから変数名を抜き出して`-out`で指定したディレクトリに結果を書き出す。
