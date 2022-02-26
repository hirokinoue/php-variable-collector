package main

import (
	"reflect"
	"testing"
)

func Test_isPhpFile(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "PHPファイルであることが判定できる",
			args: args{
				s: "hirose.php",
			},
			want: true,
		},
		{
			name: "PHPファイルでないことが判定できる",
			args: args{
				s: "yamagishi.txt",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPhpFile(tt.args.s); got != tt.want {
				t.Errorf("isPhpFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPhpVariable(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "PHPの変数であることが判定できる",
			args: args{
				s: "$hazekura",
			},
			want: true,
		},
		{
			name: "$this->はPHPの変数とみなさない",
			args: args{
				s: "$this->hazekura",
			},
			want: false,
		},
		{
			name: "PHPの変数でないことが判定できる",
			args: args{
				s: "hazekura",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPhpVariable(tt.args.s); got != tt.want {
				t.Errorf("isPhpVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeSymbolFromVariable(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "変数名に続く記号（[）を取り除ける",
			args: args{
				s: "risarisa[",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（]）を取り除ける",
			args: args{
				s: "risarisa]",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（.）を取り除ける",
			args: args{
				s: "risarisa.",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（,）を取り除ける",
			args: args{
				s: "risarisa,",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（;）を取り除ける",
			args: args{
				s: "risarisa;",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（!）を取り除ける",
			args: args{
				s: "risarisa!",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（\"）を取り除ける",
			args: args{
				s: "risarisa\"",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（'）を取り除ける",
			args: args{
				s: "risarisa'",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（)）を取り除ける",
			args: args{
				s: "risarisa)",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（(）を取り除ける",
			args: args{
				s: "risarisa(",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（:）を取り除ける",
			args: args{
				s: "risarisa:",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（%）を取り除ける",
			args: args{
				s: "risarisa%",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（+）を取り除ける",
			args: args{
				s: "risarisa+",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号（-）を取り除ける",
			args: args{
				s: "risarisa-",
			},
			want: "risarisa",
		},
		{
			name: "変数名に続く記号がない時、入力値をそのまま出力する",
			args: args{
				s: "risarisa",
			},
			want: "risarisa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeSymbolFromVariable(tt.args.s); got != tt.want {
				t.Errorf("removeSymbolFromVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filePaths(t *testing.T) {
	type args struct {
		in      string
		exclude string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "指定したディレクトリを再帰的に探索してすべてのファイルパスをスライスとして返せる",
			args: args{
				in:      "testdata",
				exclude: "",
			},
			want: []string{
				"testdata/README.md",
				"testdata/README.php",
				"testdata/src/Test.php",
				"testdata/src/Test.txt",
				"testdata/vendor/Exclude1.php",
				"testdata/vendor/Exclude2.php",
			},
			wantErr: false,
		},
		{
			name: "指定したディレクトリを再帰的に探索して除外ディレクトリ中のファイルを除きすべてのファイルパスをスライスとして返せる",
			args: args{
				in:      "testdata",
				exclude: "vendor",
			},
			want: []string{
				"testdata/README.md",
				"testdata/README.php",
				"testdata/src/Test.php",
				"testdata/src/Test.txt",
			},
			wantErr: false,
		},
		{
			name: "指定したディレクトリを再帰的に探索して除外ファイルを除きすべてのファイルパスをスライスとして返せる",
			args: args{
				in:      "testdata",
				exclude: "Exclude1.php",
			},
			want: []string{
				"testdata/README.md",
				"testdata/README.php",
				"testdata/src/Test.php",
				"testdata/src/Test.txt",
				"testdata/vendor/Exclude2.php",
			},
			wantErr: false,
		},
		{
			name: "存在しないディレクトリを指定した時エラーを返せる",
			args: args{
				in:      "heavensdoor",
				exclude: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filePaths(tt.args.in, tt.args.exclude)
			if (err != nil) != tt.wantErr {
				t.Errorf("filePaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
