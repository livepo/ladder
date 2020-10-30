package cipher


var (
    ENC = []byte{
        130,4,133,49,108,178,125,95,35,126,41,
        129,229,48,6,94,69,20,194,236,79,156,
        67,100,239,152,149,93,91,56,8,183,42,
        148,114,59,57,5,112,151,54,97,109,145,
        228,196,250,104,169,107,86,64,98,181,
        200,58,199,70,138,179,60,249,34,123,30,
        22,124,240,201,132,218,21,74,83,39,223,
        73,88,136,27,0,10,89,51,215,251,255,3,
        235,241,19,102,71,38,166,220,110,23,232,
        25,172,210,142,211,121,242,75,208,195,
        203,226,253,176,17,66,158,231,237,99,254,
        173,221,117,139,213,90,85,45,187,84,92,44,
        164,247,122,32,127,177,170,155,111,185,171,
        61,76,184,234,192,16,106,160,204,153,161,
        186,131,28,137,37,216,248,55,72,50,26,46,53,
        224,7,217,189,120,219,167,119,11,252,65,135,
        96,222,68,144,214,227,101,207,103,212,175,
        157,141,168,82,163,47,52,15,113,230,245,116,
        43,80,246,33,198,197,146,193,13,31,24,143,12,
        18,118,14,62,154,78,81,134,162,105,63,244,77,
        190,209,150,233,159,202,191,40,87,180,188,36,
        238,9,140,128,147,174,1,2,182,243,29,115,205,
        225,165,206,
    }
    DEC = make([]byte, 256)
)


func init() {
    for i := 0; i < 256; i ++ {
        DEC[ENC[i]] = byte(i)
    }
}



func Decode(data []byte) {
    for i:= 0; i < len(data); i ++ {
        data[i] = DEC[data[i]]
    }
}


func Encode(data []byte) {
    for i := 0; i < len(data); i ++ {
        data[i] = ENC[data[i]]
    }
}