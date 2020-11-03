package cipher

import (
    "bufio"
    "encoding/base64"
    "fmt"
    "math/rand"
    "os"
    "strings"
    "time"
)


var (
    ENC [256]byte
    DEC [256]byte
)


func InitPasswd() string {
    var arrays []byte
    f, err := os.Open(".gosocks")
    fileExist := false
    if err == nil {
        fileExist = true
        reader := bufio.NewReader(f)
        b64str, _ := reader.ReadString('\n')
        arrays, _ = base64.StdEncoding.DecodeString(b64str)
    } else {
        rand.Seed(time.Now().UnixNano())
        arrays = make([]byte, 256)
        arrints := rand.Perm(256)
        for idx, v := range arrints {
            arrays[idx] = byte(v)
        }
    }
    for i := 0; i < len(arrays); i ++ {
        ENC[i] = arrays[i]
        DEC[arrays[i]] = byte(i)
    }
    b64str := base64.StdEncoding.EncodeToString(ENC[:])
    fmt.Printf("secret: %s\n", b64str)

    if !fileExist {
        f, err := os.Create(".gosocks")
        if err != nil {
            fmt.Println("create secret file error", err)
            return b64str
        }
        f.Write([]byte(b64str + "\n"))
        f.Close()
    }
    return b64str
}

func Encode(data []byte) {
    for i := 0; i < len(data); i ++ {
        data[i] = ENC[data[i]]
    }
}

func Decode(data []byte) {
    for i := 0; i < len(data); i ++ {
        data[i] = DEC[data[i]]
    }
}

func LoadSecret() string {
    f, err := os.Open(".gosocks")
    if err != nil {
        return ""
    }
    reader := bufio.NewReader(f)
    b64str, err := reader.ReadString('\n')
    if err != nil {
        return ""
    }
    b64str = strings.TrimSpace(b64str)
    return b64str
}

func UpdateSecret(secret string) {
    arrays, _ := base64.StdEncoding.DecodeString(secret)
    for i := 0; i < len(arrays); i ++ {
        ENC[i] = arrays[i]
        DEC[arrays[i]] = byte(i)
    }
}
