//go:build arm64
// +build arm64

#include "textflag.h"

// func _encode(offset, b, encd []byte) (sum uint64, n int)
TEXT ·_encode(SB), NOSPLIT, $0-72
    MOVD    ·offset+0(FP), R0
    MOVD    ·data+8(FP), R9
    MOVD    ·dlen+16(FP), R3
    MOVD    ·encd+32(FP), R5

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
    SUB     R5, R7, R7
    MOVD    R3, ·sum+56(FP)
    MOVD    R7, ·n+64(FP)
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

    SUBW    $6, R1, R1                  // sub     w1, w1, #6
    CMPW    $0, R1                      // cmp     w1, 0
    BLE     dectil                      // ble     .L7
    MOVW    $0xb200, R11                // mov     w11, 45568
    MOVD    R5, R9                      // mov     x9, x5
    SUB     $8, R2, R14                 // sub     x14, x2, #8
    SUB     $4, R2, R13                 // sub     x13, x2, #4
    MOVD    $2, R8                      // mov     x8, 2
    MOVW    $0, R10                     // mov     w10, 0
    MOVK    $(0xb1ff<<16), R11          // movk    w11, 0xb1ff, lsl 16
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
    CBZW    R0, decret                  // cbz     w0, .L1
    MOVW    (R2)(R12.UXTW<<2), R1       // ldr     w1, [x2, w12, uxtw 2]
    CMPW    $1, R0                      // cmp     w0, 1
    SUBW    $0x4e, R1, R3               // sub     w3, w1, #78
    UBFX    $14, R3, $2, R4             // ubfx    x4, x3, 14, 2
    ORRW    R3<<2, R4, R3               // orr     w3, w4, w3, lsl 2
    MOVB    R3, (R5)(R10.SXTW)          // strb    w3, [x5, w10, sxtw]
    BEQ     decret                      // beq     .L1

    MOVW    $0xffb2, R7                 // mov     w7, 65458
    ADDW    $1, R10, R6                 // add     w6, w10, 1
    MOVK    $(0xffb1<<16), R7           // movk    w7, 0xffb1, lsl 16
    ADDW    R7, R1, R1                  // add     w1, w1, w7
    CMPW    $2, R0                      // cmp     w0, 2
    UBFX    $20, R1, $8, R4             // ubfx    x4, x1, 20, 8
    LSRW    $6, R1, R3                  // lsr     w3, w1, 6
    ANDW    $3, R4, R8                  // and     w8, w4, 3
    ANDW    $-4, R3, R3                 // and     w3, w3, -4
    ORRW    R8, R3, R3                  // orr     w3, w3, w8
    MOVB    R3, (R5)(R6.SXTW)           // strb    w3, [x5, w6, sxtw]
    BEQ     decret                      // beq     .L1

    ADDW    $2, R10, R3                 // add     w3, w10, 2
    LSRW    $12, R1, R6                 // lsr     w6, w1, 12
    ANDW    $-16, R6, R6                // and     w6, w6, -16
    CMPW    $3, R0                      // cmp     w0, 3
    ORRW    R1>>28, R6, R1              // orr     w1, w6, w1, lsr 28
    MOVB    R1, (R5)(R3.SXTW)           // strb    w1, [x5, w3, sxtw]
    BEQ     decret                      // beq     .L1

    ADDW    $3, R10, R1                 // add     w1, w10, 3
    ADDW    $1, R12, R12                // add     w12, w12, 1
    ANDW    $0xf0, R4, R4               // and     w4, w4, 240
    CMPW    $4, R0                      // cmp     w0, 4
    MOVW    (R2)(R12<<2), R3            // ldr     w3, [x2, x12, lsl 2]
    SUBW    $0x4e, R3, R2               // sub     w2, w3, #78
    UBFX    $2, R2, $4, R6              // ubfx    x6, x2, 2, 4
    ORRW    R6, R4, R4                  // orr     w4, w4, w6
    MOVB    R4, (R5)(R1.SXTW)           // strb    w4, [x5, w1, sxtw]
    BEQ     decret                      // beq     .L1

    ADDW    $4, R10, R1                 // add     w1, w10, 4
    UBFX    $10, R2, $6, R4             // ubfx    x4, x2, 10, 6
    ORRW    R2<<6, R4, R2               // orr     w2, w4, w2, lsl 6
    CMPW    $5, R0                      // cmp     w0, 5
    MOVB    R2, (R5)(R1.SXTW)           // strb    w2, [x5, w1, sxtw]
    BEQ     decret                      // beq     .L1

    ADDW    R7, R3, R3                  // add     w3, w3, w7
    ADDW    $5, R10, R10                // add     w10, w10, 5
    LSRW    $2, R3, R0                  // lsr     w0, w3, 2
    UBFX    $16, R3, $6, R3             // ubfx    x3, x3, 16, 6
    ANDW    $-64, R0, R0                // and     w0, w0, -64
    ORRW    R3, R0, R3                  // orr     w3, w0, w3
    MOVB    R3, (R5)(R10.SXTW)          // strb    w3, [x5, w10, sxtw]
decret:
    RET
dectil:
    MOVW    $0, R10
    MOVW    $0, R12
    JMP     decrem
