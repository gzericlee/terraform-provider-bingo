#!/usr/bin/env bash
set -e

run_with_cache() {
  file=$1
  cmd=$2

  file_sum=${file// /+}
  file_sum=${file_sum//\//_}.sum

  sum=$(eval $(echo find $file -type f) | xargs md5sum | md5sum | awk '{printf $1}')
  if [[ $(cat $file_sum) == "$sum" ]]; then
    echo $cmd is cached by $file_sum
    return 0
  fi

  eval $cmd
  echo $sum >$file_sum
}

go env -w GO111MODULE=on
go env -w GOPRIVATE=gitlab.bingosoft.net
go env -w GOPROXY=https://goproxy.cn,direct
run_with_cache go.mod "go mod tidy"

run_with_cache "api assets cmd global internal mock tool util pkg" "go build -o ./terraform_provider_bingo .main.go"
