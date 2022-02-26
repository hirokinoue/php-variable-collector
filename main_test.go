package main

import (
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
