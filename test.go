package main

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
	"strconv"
)

func main() {
    str := "c0e0e346-bef0-416c-8863-23e85ff53d52"
    hash := md5.Sum([]byte(str))
    md5Str := hex.EncodeToString(hash[:])
	decimal, _ := strconv.ParseInt(string(md5Str[0]), 16, 64)
    fmt.Println(decimal)
}
