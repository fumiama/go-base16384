//go:build amd64
// +build amd64

#include "textflag.h"

// func encode(offset, outlen int, b, encd []byte)
TEXT ·encode(SB), NOSPLIT, $0-64
    MOVQ ·offset+0(FP), R10
    MOVQ ·outlen+8(FP), AX
    MOVQ ·data+16(FP), DI
    MOVQ ·dlen+24(FP), R8
    MOVQ ·encd+40(FP), R9
    XORQ CX, CX
    XORQ SI, SI
    SUBQ $6, R8
    JLE  encrem
    MOVQ $4611404543450677248, BP
    MOVQ $70364449210368, BX
    MOVQ $5620578098173988352, R11

enclop:
    MOVBEQ (DI)(CX*1), DX
    INCQ SI
    ADDQ $7, CX
    MOVQ DX, R13
    MOVQ DX, R12
    SHRQ $2, R13
    SHRQ $4, R12
    ANDQ BX, R12
    ANDQ BP, R13
    ORQ  R12, R13
    MOVQ DX, R12
    SHRQ $8, DX
    SHRQ $6, R12
    ANDL $16383, DX
    ANDL $1073676288, R12
    ORQ  R13, R12
    ORQ  R12, DX
    ADDQ R11, DX
    MOVBEQ DX, -8(R9)(SI*8)
    CMPQ CX, R8
    JL   enclop

encrem:
    TESTQ R10, R10
    JE   encend

    MOVBLZX (DI)(CX*1), DX
    MOVL DX, R8
    SALQ $14, DX
    SHRB $2, R8
    ANDL $49152, DX
    MOVBLZX R8, R8
    ORQ  R8, DX
    CMPL R10, $1
    JE   encsav

    MOVBQSX 1(DI)(CX*1), R8
    MOVQ R8, R11
    SALQ $20, R8
    SALQ $6, R11
    ANDL $3145728, R8
    ANDL $16128, R11
    ORQ  R11, DX
    ORQ  R8, DX
    CMPL R10, $2
    JE   encsav

    MOVBQSX 2(DI)(CX*1), R8
    MOVQ R8, R11
    SALQ $28, R8
    SALQ $12, R11
    MOVL R8, R8
    ANDL $983040, R11
    ORQ  R11, R8
    ORQ  R8, DX
    CMPL R10, $3
    JE   encsav

    MOVQ $257698037760, BX
    MOVBQSX 3(DI)(CX*1), R8
    MOVQ R8, R11
    SALQ $34, R8
    SALQ $20, R11
    ANDQ BX, R8
    ANDL $251658240, R11
    ORQ  R11, R8
    ORQ  R8, DX
    CMPL R10, $4
    JE   encsav

    MOVQ $12884901888, BX
    MOVBQSX 4(DI)(CX*1), R8
    MOVQ R8, R11
    SALQ $42, R8
    SALQ $26, R11
    ANDQ BX, R11
    MOVQ $277076930199552, BX
    ANDQ BX, R8
    ORQ  R11, R8
    ORQ  R8, DX
    CMPL R10, $5
    JE   encsav

    MOVQ $3298534883328, R8
    MOVBQSX 5(DI)(CX*1), CX
    MOVQ CX, DI
    SALQ $48, CX
    SALQ $34, DI
    ANDQ R8, DI
    MOVQ $17732923532771328, R8
    ANDQ R8, CX
    ORQ  DI, CX
    ORQ  CX, DX

encsav:
    MOVQ $21955383195992142, CX
    ADDQ CX, DX
    MOVQ DX, 0(R9)(SI*8)
    MOVB $61, -2(R9)(AX*1)
    MOVB R10, -1(R9)(AX*1)

encend:
    RET


// func decode(offset, outlen int, b, decd []byte)
TEXT ·decode(SB), NOSPLIT, $0-64
    MOVQ ·offset+0(FP), BX
    MOVQ ·outlen+8(FP), R8
    MOVQ ·data+16(FP), DI
    MOVQ ·decd+40(FP), R9
    XORQ CX, CX
    XORQ SI, SI
    SUBQ $6, R8
    JLE  decrem
    MOVQ $-5620578098173988352, R12
    MOVQ $-1125899906842624, BP
    MOVQ $1125831187365888, R11
    MOVQ $68715282432, R10

declop:
    MOVBEQ (DI)(SI*8), DX
    INCQ SI
    ADDQ R12, DX
    MOVQ DX, R13
    LEAQ 0(DX*4), R14
    SALQ $4, R13
    ANDQ BP, R14
    ANDQ R11, R13
    ORQ  R13, R14
    MOVQ DX, R13
    SALQ $8, DX
    SALQ $6, R13
    ANDL $4194048, DX
    ANDQ R10, R13
    ORQ  R14, R13
    ORQ  R13, DX
    MOVBEQ DX, (R9)(CX*1)
    ADDQ $7, CX
    CMPQ CX, R8
    JL   declop

decrem:
    TESTQ BX, BX
    JE   decend

    MOVQ (DI)(SI*8), DI
    LEAQ -78(DI), SI
    MOVQ SI, DX
    SALL $2, SI
    SHRQ $14, DX
    ANDL $3, DX
    ORL  SI, DX
    MOVB DX, 0(R9)(CX*1)
    CMPL BX, $1
    JE   decend

    LEAQ -5111886(DI), DX
    MOVQ DX, SI
    MOVQ DX, R8
    SHRQ $6, SI
    SHRQ $20, R8
    ANDL $-4, SI
    ANDL $3, R8
    ORL  R8, SI
    MOVB SI, 1(R9)(CX*1)
    CMPL BX, $2
    JE   decend

    MOVQ DX, SI
    SHRQ $28, DX
    SHRQ $12, SI
    ANDL $15, DX
    ANDL $-16, SI
    ORL  SI, DX
    MOVB DX, 2(R9)(CX*1)
    CMPL BX, $3
    JE   decend

    MOVQ $-335012560974, DX
    ADDQ DI, DX
    MOVQ DX, SI
    MOVQ DX, R8
    SHRQ $20, SI
    SHRQ $34, R8
    ANDL $-16, SI
    ANDL $15, R8
    ORL  R8, SI
    MOVB SI, 3(R9)(CX*1)
    CMPL BX, $4
    JE   decend

    MOVQ DX, SI
    SHRQ $42, DX
    SHRQ $26, SI
    ANDL $63, DX
    ANDL $-64, SI
    ORL  SI, DX
    MOVB DX, 4(R9)(CX*1)
    CMPL BX, $5
    JE   decend

    MOVQ $-21955383195992142, DX
    ADDQ DX, DI
    MOVQ DI, DX
    SHRQ $48, DI
    SHRQ $34, DX
    ANDL $63, DI
    ANDL $-64, DX
    ORL  DI, DX
    MOVB DX, 5(R9)(CX*1)

decend:
    RET
