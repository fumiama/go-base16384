//go:build arm64
// +build arm64

#include "textflag.h"

// func _encode(offset, outlen int, b, encd []byte) (sum uint64, &vals[n] uintptr)
TEXT ·_encode(SB), NOSPLIT, $0-81
    MOVD    ·offset+0(FP), R0
    MOVD    ·data+16(FP), R9
    MOVD    ·dlen+24(FP), R3
    MOVD    ·encd+40(FP), R5

    SUBW    $6, R3, R3
    CMPW    $0, R3
    BLE     enctil
    MOVW    $0x4e00, R11
    SUB     $8, R5, R14
    SUB     $4, R5, R13
    MOVD    $2, R8
    MOVW    $0, R10 // int32_t i = 0
    MOVK    $(0x4e00<<16), R11
enclop:
    MOVW    (R9), R4
    ADDW    $7, R10, R10
    MOVW    R8, R12
    CMPW    R3, R10
    REVW    R4, R4
    ADD     $7, R9, R9
    LSRW    $2, R4, R6
    UBFX    $4, R4, $14, R15
    ANDW    $0x3fff0000, R6, R6
    UBFIZW  $26, R4, $4, R7
    ORRW    R15, R6, R6
    ADDW    R11, R6, R6
    REVW    R6, R6
    MOVW    R6, (R14)(R8<<2)
    MOVW    -3(R9), R4
    REVW    R4, R4
    LSRW    $6, R4, R4
    ANDW    $0x3fffffc, R4, R4
    ORRW    R7, R4, R4
    ANDW    $0x3fff0000, R4, R6
    UBFX    $2, R4, $14, R4
    ORRW    R6, R4, R4
    ADDW    R11, R4, R4
    REVW    R4, R4
    MOVW    R4, (R13)(R8<<2)
    ADDW    $2, R8, R8
    BLT     enclop
encrem:
    ANDSW   $0xff, R0, R0
    BEQ     encret

    MOVBU   (R2)(R10.SXTW), R3
    UXTW    R12, R8
    CMPW    $1, R0
    SXTW    R10, R10
    ADD     R8<<2, R5, R7
    UBFIZW  $14, R3, $2, R4
    ORRW    R3>>2, R4, R3
    BEQ     encsum

    ADD     R10, R2, R9
    CMPW    $2, R0
    MOVBU   1(R9), R6
    LSLW    $6, R6, R4
    UBFIZW  $20, R6, $2, R6
    ANDW    $0x3f00, R4, R4
    ORRW    R3, R4, R3
    ORRW    R3, R6, R3
    BEQ     encsum

    MOVBU   2(R9), R4
    CMPW    $3, R0
    LSLW    $12, R4, R6
    ANDW    $0xf0000, R6, R6
    ORRW    R4<<28, R6, R4
    ORRW    R4, R3, R3
    BEQ     encsum

    ADD     $3, R10, R10
    ADDW    $1, R12, R12
    CMPW    $4, R0
    ADD     R12<<2, R5, R7
    MOVBU   (R2)(R10), R4
    LSLW    $20, R4, R4
    ANDW    $0xf000000, R4, R4
    ORRW    R3, R4, R3
    ADDW    $0x4e0000, R3, R3
    ADDW    $0x4e, R3, R3
    MOVW    R3, (R5)(R8<<2)
    MOVBU   (R2)(R10), R3
    UBFIZW  $2, R3, $4, R3
    BEQ     encsum

    MOVBU   4(R9), R4
    CMPW    $5, R0
    UBFIZW  $10, R4, $6, R2
    ORRW    R3, R2, R3
    ORRW    R4>>6, R3, R3
    BEQ     encsum

    MOVBU   5(R9), R4
    LSLW    $2, R4, R2
    UBFIZW  $16, R4, $6, R4
    ANDW    $0x300, R2, R2
    ORRW    R4, R2, R2
    ORRW    R2, R3, R3
encsum:
    ADDW    $0x4e0000, R3, R3
    ADDW    $0x4e, R3, R3
    MOVD    R3, ·sum+64(FP)
    MOVD    R7, ·n+72(FP)
encret:
    RET
enctil:
    MOVW    $0, R10
    MOVW    $0, R12
    JMP     encrem

// func _decode(offset, outlen int, b, decd []byte)
TEXT ·_decode(SB), NOSPLIT, $0-64
    MOVD    ·offset+0(FP), R0
    MOVD    ·outlen+8(FP), R1
    MOVD    ·data+16(FP), R2
    MOVD    ·decd+40(FP), R5

    SUBW    $6, R1, R1
    CMPW    $0, R1
    BLE     dectil
    MOVW    $0xb200, R11
    MOVD    R5, R9
    SUB     $8, R2, R14
    SUB     $4, R2, R13
    MOVD    $2, R8
    MOVW    $0, R10
    MOVK    $(0xb1ff<<16), R11
declop:
    MOVW    (R14)(R8<<2), R4
    ADDW    $7, R10, R10
    MOVW    (R13)(R8<<2), R3
    MOVW    R8, R12
    REVW    R4, R4
    CMPW    R1, R10
    ADDW    R11, R4, R4
    REVW    R3, R3
    ADDW    R11, R3, R3
    ADD     $2, R8, R8
    LSLW    $2, R4, R7
    UBFIZW  $4, R4, $14, R4
    LSLW    $6, R3, R6
    ANDW    $-262144, R7, R7
    ORRW    R4, R7, R7
    ANDW    $-4194304, R6, R4
    UBFIZW  $8, R3, $14, R6
    ORRW    R3>>26, R7, R3
    ORRW    R6, R4, R4
    REVW    R3, R3
    REVW    R4, R4
    STPW    (R3, R4), (R9)
    ADD     $7, R9, R9
    BLT     declop
decrem:
    CBZW    R0, decret
    MOVW    (R2)(R12.UXTW<<2), R1
    CMPW    $1, R0
    SUBW    $0x4e, R1, R3
    UBFX    $14, R3, $2, R4
    ORRW    R3<<2, R4, R3
    MOVB    R3, (R5)(R10.SXTW)
    BEQ     decret

    MOVW    $0xffb2, R7
    ADDW    $1, R10, R4
    MOVK    $(0xffb1<<16), R7
    ADDW    R7, R1, R1
    CMPW    $2, R0
    UBFX    $20, R1, $8, R6
    LSRW    $6, R1, R3
    ANDW    $3, R6, R8
    ANDW    $-4, R3, R3
    ORRW    R8, R3, R3
    MOVB    R3, (R5)(R4.SXTW)
    BEQ     decret

    ADDW    $2, R10, R3
    LSRW    $12, R1, R4
    ANDW    $-16, R4, R4
    CMPW    $3, R0
    ORRW    R1>>28, R4, R1
    MOVB    R1, (R5)(R3.SXTW)
    BEQ     decret

    ADDW    $3, R10, R1
    ADDW    $1, R12, R12
    ANDW    $0xf0, R6, R6
    CMPW    $4, R0
    MOVW    (R2)(R12<<2), R3
    SUBW    $0x4e, R3, R2
    UBFX    $2, R2, $4, R4
    ORRW    R6, R4, R4
    MOVB    R4, (R5)(R1.SXTW)
    BEQ     decret

    ADDW    $4, R10, R1
    UBFX    $10, R2, $6, R4
    ORRW    R2<<6, R4, R2
    CMPW    $5, R0
    MOVB    R2, (R5)(R1.SXTW)
    BEQ     decret

    ADDW    R7, R3, R3
    ADDW    $5, R10, R10
    LSRW    $2, R3, R0
    UBFX    $16, R3, $6, R3
    ANDW    $-64, R0, R0
    ORRW    R3, R0, R3
    MOVB    R3, (R5)(R10.SXTW)
decret:
    RET
dectil:
    MOVW    $0, R10
    MOVW    $0, R12
    JMP     decrem
