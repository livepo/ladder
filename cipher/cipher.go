package cipher

import (
    "math/rand"
)


var (
    ENC [256]byte
    DEC [256]byte
)


func init() {
    arrays := rand.Perm(256)
    for i := 0; i < len(arrays); i ++ {
        ENC[i] = byte(arrays[i])
        DEC[arrays[i]] = byte(i)
    }
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
