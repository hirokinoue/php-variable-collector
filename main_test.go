package main

import (
	"reflect"
	"runtime"
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
		inDir   string
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
				inDir:   "testdata",
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
				inDir:   "testdata",
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
				inDir:   "testdata",
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
				inDir:   "heavensdoor",
				exclude: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filePaths(tt.args.inDir, tt.args.exclude)
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

func Test_phpFilePaths(t *testing.T) {
	type args struct {
		inDir   string
		exclude string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "指定したディレクトリを再帰的に探索してすべてのphpのファイルパスをスライスとして返せる",
			args: args{
				inDir:   "testdata",
				exclude: "",
			},
			want: []string{
				"testdata/README.php",
				"testdata/src/Test.php",
				"testdata/vendor/Exclude1.php",
				"testdata/vendor/Exclude2.php",
			},
			wantErr: false,
		},
		{
			name: "指定したディレクトリを再帰的に探索して除外ディレクトリ中のファイルを除きすべてのphpのファイルパスをスライスとして返せる",
			args: args{
				inDir:   "testdata",
				exclude: "vendor",
			},
			want: []string{
				"testdata/README.php",
				"testdata/src/Test.php",
			},
			wantErr: false,
		},
		{
			name: "指定したディレクトリを再帰的に探索して除外ファイルを除きすべてのphpのファイルパスをスライスとして返せる",
			args: args{
				inDir:   "testdata",
				exclude: "Exclude1.php",
			},
			want: []string{
				"testdata/README.php",
				"testdata/src/Test.php",
				"testdata/vendor/Exclude2.php",
			},
			wantErr: false,
		},
		{
			name: "存在しないディレクトリを指定した時エラーを返せる",
			args: args{
				inDir:   "cheaptrick",
				exclude: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := phpFilePaths(tt.args.inDir, tt.args.exclude)
			if (err != nil) != tt.wantErr {
				t.Errorf("phpFilePaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("phpFilePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dict_add(t *testing.T) {
	type args struct {
		variable string
	}
	tests := []struct {
		name string
		d    *dict
		args args
		want map[string]bool
	}{
		{
			name: "dict.valueに文字列を追加できる",
			d: func() *dict {
				return newDict()
			}(),
			args: args{
				variable: "poco",
			},
			want: map[string]bool{"poco": true},
		},
		{
			name: "dict.valueに既に存在する文字列は追加されない",
			d: func() *dict {
				d := newDict()
				d.add("poco")
				return d
			}(),
			args: args{
				variable: "poco",
			},
			want: map[string]bool{"poco": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.add(tt.args.variable)
			got := tt.d.value
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dict.add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dict_sortValue(t *testing.T) {
	tests := []struct {
		name string
		d    *dict
		want []string
	}{
		{
			name: "dict.valueをソートしてスライスとして返せる",
			d: func() *dict {
				d := newDict()
				d.add("wham!")
				d.add("acdc")
				d.add("cars")
				return d
			}(),
			want: []string{
				"acdc",
				"cars",
				"wham!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sorted := tt.d.sortValue()
			for i, got := range sorted {
				if got != tt.want[i] {
					t.Errorf("dict.sortValue() = %v, want %v", got, tt.want[i])
				}
			}
		})
	}
}

func Test_collectPhpVariable(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "指定したファイルからPHPの変数を取り出してスライスとして返せる",
			args: args{
				filePath: "testdata/README.php",
			},
			want:    []string{"$readmePhp", "$foo"},
			wantErr: false,
		},
		{
			name: "存在しないディレクトリを指定した時エラーを返せる",
			args: args{
				filePath: "highway/star.php",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ch := make(chan []string)
		e := make(chan error)
		semaphore := make(chan struct{}, runtime.NumCPU())
		defer func() {
			close(ch)
			close(e)
			close(semaphore)
		}()

		t.Run(tt.name, func(t *testing.T) {
			go collectPhpVariable(tt.args.filePath, ch, e, semaphore)
			select {
			case got := <-ch:
				if tt.wantErr != false {
					t.Errorf("collectPhpVariable() = %v, want error", got)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("collectPhpVariable() = %v, want %v", got, tt.want)
				}
			case err := <-e:
				if tt.wantErr == false || tt.want != nil {
					t.Errorf("collectPhpVariable() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
