// arm64 NEON 内核. 每个内核只处理整块元素(f64 每 4 个, f32 每 8 个),
// 不满一块的尾部由调用方在 Go 侧补齐, 因此内核内部没有标量尾循环.

#include "textflag.h"

// 常量
DATA ·neonOneF64+0(SB)/8, $0x3ff0000000000000
GLOBL ·neonOneF64(SB), RODATA, $8

DATA ·neonOneF32+0(SB)/4, $0x3f800000
GLOBL ·neonOneF32(SB), RODATA, $4

// ---------------------------------------------------------------- f64 元素级

// a[i] = a[i] op b[i]
#define ELEM2_F64(sym, OP) \
TEXT sym(SB), NOSPLIT, $0-24 \
	MOVD a+0(FP), R0 \
	MOVD b+8(FP), R1 \
	MOVD n+16(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1  (R0), [V0.D2, V1.D2] \
	VLD1  (R1), [V2.D2, V3.D2] \
	OP    V2.D2, V0.D2, V0.D2 \
	OP    V3.D2, V1.D2, V1.D2 \
	VST1  [V0.D2, V1.D2], (R0) \
	ADD   $32, R0 \
	ADD   $32, R1 \
	SUB   $4, R2 \
	CBNZ  R2, loop \
done: \
	RET

// a[i] = a[i] op v
#define ELEMNUM_F64(sym, OP) \
TEXT sym(SB), NOSPLIT, $0-24 \
	MOVD  a+0(FP), R0 \
	FMOVD v+8(FP), F0 \
	VDUP  V0.D[0], V4.D2 \
	MOVD  n+16(FP), R2 \
	CBZ   R2, done \
loop: \
	VLD1  (R0), [V0.D2, V1.D2] \
	OP    V4.D2, V0.D2, V0.D2 \
	OP    V4.D2, V1.D2, V1.D2 \
	VST1  [V0.D2, V1.D2], (R0) \
	ADD   $32, R0 \
	SUB   $4, R2 \
	CBNZ  R2, loop \
done: \
	RET

// a[i] = op(a[i])
#define UNARY_F64(sym, OP) \
TEXT sym(SB), NOSPLIT, $0-16 \
	MOVD a+0(FP), R0 \
	MOVD n+8(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1  (R0), [V0.D2, V1.D2] \
	OP    V0.D2, V0.D2 \
	OP    V1.D2, V1.D2 \
	VST1  [V0.D2, V1.D2], (R0) \
	ADD   $32, R0 \
	SUB   $4, R2 \
	CBNZ  R2, loop \
done: \
	RET

ELEM2_F64(·neonAddF64, VFADD)
ELEM2_F64(·neonSubF64, VFSUB)
ELEM2_F64(·neonMulF64, VFMUL)
ELEM2_F64(·neonDivF64, VFDIV)
ELEM2_F64(·neonMinimumF64, VFMIN)
ELEM2_F64(·neonMaximumF64, VFMAX)

ELEMNUM_F64(·neonAddNumF64, VFADD)
ELEMNUM_F64(·neonSubNumF64, VFSUB)
ELEMNUM_F64(·neonMulNumF64, VFMUL)
ELEMNUM_F64(·neonDivNumF64, VFDIV)
ELEMNUM_F64(·neonMinNumF64, VFMIN)
ELEMNUM_F64(·neonMaxNumF64, VFMAX)

UNARY_F64(·neonAbsF64, VFABS)
UNARY_F64(·neonNegF64, VFNEG)
UNARY_F64(·neonSqrtF64, VFSQRT)
UNARY_F64(·neonRoundF64, VFRINTN)
UNARY_F64(·neonFloorF64, VFRINTM)
UNARY_F64(·neonCeilF64, VFRINTP)

// a[i] = 1 / a[i]
TEXT ·neonInvF64(SB), NOSPLIT, $0-16
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	MOVD ·neonOneF64(SB), R3
	VDUP R3, V4.D2
	CBZ  R2, done
loop:
	VLD1  (R0), [V0.D2, V1.D2]
	VFDIV V0.D2, V4.D2, V0.D2
	VFDIV V1.D2, V4.D2, V1.D2
	VST1  [V0.D2, V1.D2], (R0)
	ADD   $32, R0
	SUB   $4, R2
	CBNZ  R2, loop
done:
	RET

// ---------------------------------------------------------------- f64 归约

// 双通道累加: 循环体为 OP V2.D2, V0.D2, V0.D2
#define REDUCE_F64(sym, OP) \
TEXT sym(SB), NOSPLIT, $0-24 \
	MOVD a+0(FP), R0 \
	MOVD n+8(FP), R2 \
	VEOR V0.B16, V0.B16, V0.B16 \
	VEOR V1.B16, V1.B16, V1.B16 \
	CBZ  R2, fold \
loop: \
	VLD1 (R0), [V2.D2, V3.D2] \
	OP   V2.D2, V0.D2, V0.D2 \
	OP   V3.D2, V1.D2, V1.D2 \
	ADD  $32, R0 \
	SUB  $4, R2 \
	CBNZ R2, loop \
fold: \
	VFADD   V1.D2, V0.D2, V0.D2 \
	VFADDP V0.D2, V0.D2, V0.D2 \
	VMOV   V0.D[0], R3 \
	MOVD   R3, ret+16(FP) \
	RET

REDUCE_F64(·neonSumF64, VFADD)

// 平方和
TEXT ·neonSumSqF64(SB), NOSPLIT, $0-24
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	VEOR V0.B16, V0.B16, V0.B16
	VEOR V1.B16, V1.B16, V1.B16
	CBZ  R2, fold
loop:
	VLD1  (R0), [V2.D2, V3.D2]
	VFMLA V2.D2, V2.D2, V0.D2
	VFMLA V3.D2, V3.D2, V1.D2
	ADD   $32, R0
	SUB   $4, R2
	CBNZ  R2, loop
fold:
	VFADD   V1.D2, V0.D2, V0.D2
	VFADDP V0.D2, V0.D2, V0.D2
	VMOV   V0.D[0], R3
	MOVD   R3, ret+16(FP)
	RET

// 点积
TEXT ·neonDotF64(SB), NOSPLIT, $0-32
	MOVD a+0(FP), R0
	MOVD b+8(FP), R1
	MOVD n+16(FP), R2
	VEOR V0.B16, V0.B16, V0.B16
	VEOR V1.B16, V1.B16, V1.B16
	CBZ  R2, fold
loop:
	VLD1  (R0), [V2.D2, V3.D2]
	VLD1  (R1), [V4.D2, V5.D2]
	VFMLA V2.D2, V4.D2, V0.D2
	VFMLA V3.D2, V5.D2, V1.D2
	ADD   $32, R0
	ADD   $32, R1
	SUB   $4, R2
	CBNZ  R2, loop
fold:
	VFADD   V1.D2, V0.D2, V0.D2
	VFADDP V0.D2, V0.D2, V0.D2
	VMOV   V0.D[0], R3
	MOVD   R3, ret+24(FP)
	RET

// 绝对值之和
TEXT ·neonAbsSumF64(SB), NOSPLIT, $0-24
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	VEOR V0.B16, V0.B16, V0.B16
	VEOR V1.B16, V1.B16, V1.B16
	CBZ  R2, fold
loop:
	VLD1  (R0), [V2.D2, V3.D2]
	VFABS V2.D2, V4.D2
	VFABS V3.D2, V5.D2
	VFADD  V4.D2, V0.D2, V0.D2
	VFADD  V5.D2, V1.D2, V1.D2
	ADD   $32, R0
	SUB   $4, R2
	CBNZ  R2, loop
fold:
	VFADD   V1.D2, V0.D2, V0.D2
	VFADDP V0.D2, V0.D2, V0.D2
	VMOV   V0.D[0], R3
	MOVD   R3, ret+16(FP)
	RET

// 差值平方和
TEXT ·neonDiffSqSumF64(SB), NOSPLIT, $0-32
	MOVD a+0(FP), R0
	MOVD b+8(FP), R1
	MOVD n+16(FP), R2
	VEOR V0.B16, V0.B16, V0.B16
	VEOR V1.B16, V1.B16, V1.B16
	CBZ  R2, fold
loop:
	VLD1  (R0), [V2.D2, V3.D2]
	VLD1  (R1), [V4.D2, V5.D2]
	VFSUB V4.D2, V2.D2, V6.D2
	VFSUB V5.D2, V3.D2, V7.D2
	VFMLA V6.D2, V6.D2, V0.D2
	VFMLA V7.D2, V7.D2, V1.D2
	ADD   $32, R0
	ADD   $32, R1
	SUB   $4, R2
	CBNZ  R2, loop
fold:
	VFADD   V1.D2, V0.D2, V0.D2
	VFADDP V0.D2, V0.D2, V0.D2
	VMOV   V0.D[0], R3
	MOVD   R3, ret+24(FP)
	RET

// 差值绝对值之和
TEXT ·neonAbsDiffSumF64(SB), NOSPLIT, $0-32
	MOVD a+0(FP), R0
	MOVD b+8(FP), R1
	MOVD n+16(FP), R2
	VEOR V0.B16, V0.B16, V0.B16
	VEOR V1.B16, V1.B16, V1.B16
	CBZ  R2, fold
loop:
	VLD1  (R0), [V2.D2, V3.D2]
	VLD1  (R1), [V4.D2, V5.D2]
	VFSUB V4.D2, V2.D2, V6.D2
	VFSUB V5.D2, V3.D2, V7.D2
	VFABS V6.D2, V6.D2
	VFABS V7.D2, V7.D2
	VFADD  V6.D2, V0.D2, V0.D2
	VFADD  V7.D2, V1.D2, V1.D2
	ADD   $32, R0
	ADD   $32, R1
	SUB   $4, R2
	CBNZ  R2, loop
fold:
	VFADD   V1.D2, V0.D2, V0.D2
	VFADDP V0.D2, V0.D2, V0.D2
	VMOV   V0.D[0], R3
	MOVD   R3, ret+24(FP)
	RET

// 乘积 (要求 n >= 4)
TEXT ·neonProdF64(SB), NOSPLIT, $0-24
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	VLD1  (R0), [V0.D2, V1.D2]
	ADD   $32, R0
	SUB   $4, R2
	CBZ   R2, fold
loop:
	VLD1  (R0), [V2.D2, V3.D2]
	VFMUL V2.D2, V0.D2, V0.D2
	VFMUL V3.D2, V1.D2, V1.D2
	ADD   $32, R0
	SUB   $4, R2
	CBNZ  R2, loop
fold:
	VFMUL V1.D2, V0.D2, V0.D2
	VMOV  V0.D[0], R3
	VMOV  V0.D[1], R4
	FMOVD R3, F0
	FMOVD R4, F1
	FMULD F1, F0, F0
	FMOVD F0, ret+16(FP)
	RET

// 最小值 (要求 n >= 4)
TEXT ·neonMinF64(SB), NOSPLIT, $0-24
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	VLD1  (R0), [V0.D2, V1.D2]
	ADD   $32, R0
	SUB   $4, R2
	CBZ   R2, fold
loop:
	VLD1  (R0), [V2.D2, V3.D2]
	VFMIN V2.D2, V0.D2, V0.D2
	VFMIN V3.D2, V1.D2, V1.D2
	ADD   $32, R0
	SUB   $4, R2
	CBNZ  R2, loop
fold:
	VFMIN  V1.D2, V0.D2, V0.D2
	VFMINP V0.D2, V0.D2, V0.D2
	VMOV   V0.D[0], R3
	MOVD   R3, ret+16(FP)
	RET

// 最大值 (要求 n >= 4)
TEXT ·neonMaxF64(SB), NOSPLIT, $0-24
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	VLD1  (R0), [V0.D2, V1.D2]
	ADD   $32, R0
	SUB   $4, R2
	CBZ   R2, fold
loop:
	VLD1  (R0), [V2.D2, V3.D2]
	VFMAX V2.D2, V0.D2, V0.D2
	VFMAX V3.D2, V1.D2, V1.D2
	ADD   $32, R0
	SUB   $4, R2
	CBNZ  R2, loop
fold:
	VFMAX  V1.D2, V0.D2, V0.D2
	VFMAXP V0.D2, V0.D2, V0.D2
	VMOV   V0.D[0], R3
	MOVD   R3, ret+16(FP)
	RET

// 重复填充
TEXT ·neonRepeatF64(SB), NOSPLIT, $0-24
	MOVD  a+0(FP), R0
	FMOVD v+8(FP), F0
	VDUP  V0.D[0], V1.D2
	VMOV  V1.B16, V2.B16
	MOVD  n+16(FP), R2
	CBZ   R2, done
loop:
	VST1 [V1.D2, V2.D2], (R0)
	ADD  $32, R0
	SUB  $4, R2
	CBNZ R2, loop
done:
	RET

// 查找 (返回首个下标, 未找到返回 -1)
TEXT ·neonFindF64(SB), NOSPLIT, $0-32
	MOVD  a+0(FP), R0
	MOVD  n+8(FP), R1
	FMOVD v+16(FP), F0
	VDUP  V0.D[0], V1.D2
	MOVD  $-1, R8
	MOVD  $0, R2
	CBZ   R1, done
loop:
	VLD1   (R0), [V2.D2, V3.D2]
	VFCMEQ V1.D2, V2.D2, V4.D2
	VFCMEQ V1.D2, V3.D2, V5.D2
	VMOV   V4.D[0], R4
	VMOV   V4.D[1], R5
	VMOV   V5.D[0], R6
	VMOV   V5.D[1], R7
	CBNZ   R4, hit0
	CBNZ   R5, hit1
	CBNZ   R6, hit2
	CBNZ   R7, hit3
	ADD    $4, R2
	ADD    $32, R0
	SUB    $4, R1
	CBNZ   R1, loop
done:
	MOVD R8, ret+24(FP)
	RET
hit0:
	MOVD R2, ret+24(FP)
	RET
hit1:
	ADD  $1, R2
	MOVD R2, ret+24(FP)
	RET
hit2:
	ADD  $2, R2
	MOVD R2, ret+24(FP)
	RET
hit3:
	ADD  $3, R2
	MOVD R2, ret+24(FP)
	RET

// ---------------------------------------------------------------- f64 比较

// dst[i] = x[i] OP y[i]
#define CMP_XOPY_F64(sym, MN) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD a+0(FP), R0 \
	MOVD b+8(FP), R1 \
	MOVD d+16(FP), R3 \
	MOVD n+24(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1 (R0), [V0.D2, V1.D2] \
	VLD1 (R1), [V2.D2, V3.D2] \
	MN   V2.D2, V0.D2, V4.D2 \
	MN   V3.D2, V1.D2, V5.D2 \
	VMOV V4.D[0], R4 \
	VMOV V4.D[1], R5 \
	VMOV V5.D[0], R6 \
	VMOV V5.D[1], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	ADD  $32, R0 \
	ADD  $32, R1 \
	ADD  $4, R3 \
	SUB  $4, R2 \
	CBNZ R2, loop \
done: \
	RET

// dst[i] = y[i] OP x[i]
#define CMP_YOPX_F64(sym, MN) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD a+0(FP), R0 \
	MOVD b+8(FP), R1 \
	MOVD d+16(FP), R3 \
	MOVD n+24(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1 (R0), [V0.D2, V1.D2] \
	VLD1 (R1), [V2.D2, V3.D2] \
	MN   V0.D2, V2.D2, V4.D2 \
	MN   V1.D2, V3.D2, V5.D2 \
	VMOV V4.D[0], R4 \
	VMOV V4.D[1], R5 \
	VMOV V5.D[0], R6 \
	VMOV V5.D[1], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	ADD  $32, R0 \
	ADD  $32, R1 \
	ADD  $4, R3 \
	SUB  $4, R2 \
	CBNZ R2, loop \
done: \
	RET

// dst[i] = !(y[i] OP x[i])
#define CMPN_YOPX_F64(sym, MN) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD a+0(FP), R0 \
	MOVD b+8(FP), R1 \
	MOVD d+16(FP), R3 \
	MOVD n+24(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1 (R0), [V0.D2, V1.D2] \
	VLD1 (R1), [V2.D2, V3.D2] \
	MN   V0.D2, V2.D2, V4.D2 \
	MN   V1.D2, V3.D2, V5.D2 \
	VMOV V4.D[0], R4 \
	VMOV V4.D[1], R5 \
	VMOV V5.D[0], R6 \
	VMOV V5.D[1], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	EOR  $1, R4, R4 \
	EOR  $1, R5, R5 \
	EOR  $1, R6, R6 \
	EOR  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	ADD  $32, R0 \
	ADD  $32, R1 \
	ADD  $4, R3 \
	SUB  $4, R2 \
	CBNZ R2, loop \
done: \
	RET

// 与标量比较: dst[i] = v OP x[i]
#define CMPNUM_XOPV_F64(sym, MN) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD  a+0(FP), R0 \
	MOVD  d+8(FP), R3 \
	MOVD  n+16(FP), R2 \
	FMOVD v+24(FP), F0 \
	VDUP  V0.D[0], V4.D2 \
	CBZ   R2, done \
loop: \
	VLD1 (R0), [V0.D2, V1.D2] \
	MN   V4.D2, V0.D2, V2.D2 \
	MN   V4.D2, V1.D2, V3.D2 \
	VMOV V2.D[0], R4 \
	VMOV V2.D[1], R5 \
	VMOV V3.D[0], R6 \
	VMOV V3.D[1], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	ADD  $32, R0 \
	ADD  $4, R3 \
	SUB  $4, R2 \
	CBNZ R2, loop \
done: \
	RET

// 与标量比较: dst[i] = x[i] OP v
#define CMPNUM_VOPX_F64(sym, MN) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD  a+0(FP), R0 \
	MOVD  d+8(FP), R3 \
	MOVD  n+16(FP), R2 \
	FMOVD v+24(FP), F0 \
	VDUP  V0.D[0], V4.D2 \
	CBZ   R2, done \
loop: \
	VLD1 (R0), [V0.D2, V1.D2] \
	MN   V0.D2, V4.D2, V2.D2 \
	MN   V1.D2, V4.D2, V3.D2 \
	VMOV V2.D[0], R4 \
	VMOV V2.D[1], R5 \
	VMOV V3.D[0], R6 \
	VMOV V3.D[1], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	ADD  $32, R0 \
	ADD  $4, R3 \
	SUB  $4, R2 \
	CBNZ R2, loop \
done: \
	RET

// 与标量比较: dst[i] = !(x[i] OP v)
#define CMPNUMN_VOPX_F64(sym, MN) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD  a+0(FP), R0 \
	MOVD  d+8(FP), R3 \
	MOVD  n+16(FP), R2 \
	FMOVD v+24(FP), F0 \
	VDUP  V0.D[0], V4.D2 \
	CBZ   R2, done \
loop: \
	VLD1 (R0), [V0.D2, V1.D2] \
	MN   V0.D2, V4.D2, V2.D2 \
	MN   V1.D2, V4.D2, V3.D2 \
	VMOV V2.D[0], R4 \
	VMOV V2.D[1], R5 \
	VMOV V3.D[0], R6 \
	VMOV V3.D[1], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	EOR  $1, R4, R4 \
	EOR  $1, R5, R5 \
	EOR  $1, R6, R6 \
	EOR  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	ADD  $32, R0 \
	ADD  $4, R3 \
	SUB  $4, R2 \
	CBNZ R2, loop \
done: \
	RET

CMP_YOPX_F64(·neonLtF64, VFCMGT)
CMP_YOPX_F64(·neonLteF64, VFCMGE)
CMP_XOPY_F64(·neonGtF64, VFCMGT)
CMP_XOPY_F64(·neonGteF64, VFCMGE)
CMP_YOPX_F64(·neonEqF64, VFCMEQ)
CMPN_YOPX_F64(·neonNeqF64, VFCMEQ)

CMPNUM_VOPX_F64(·neonLtNumF64, VFCMGT)
CMPNUM_VOPX_F64(·neonLteNumF64, VFCMGE)
CMPNUM_XOPV_F64(·neonGtNumF64, VFCMGT)
CMPNUM_XOPV_F64(·neonGteNumF64, VFCMGE)
CMPNUM_XOPV_F64(·neonEqNumF64, VFCMEQ)
CMPNUMN_VOPX_F64(·neonNeqNumF64, VFCMEQ)

// ---------------------------------------------------------------- f32 元素级

#define ELEM2_F32(sym, OP) \
TEXT sym(SB), NOSPLIT, $0-24 \
	MOVD a+0(FP), R0 \
	MOVD b+8(FP), R1 \
	MOVD n+16(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1  (R0), [V0.S4, V1.S4] \
	VLD1  (R1), [V2.S4, V3.S4] \
	OP    V2.S4, V0.S4, V0.S4 \
	OP    V3.S4, V1.S4, V1.S4 \
	VST1  [V0.S4, V1.S4], (R0) \
	ADD   $32, R0 \
	ADD   $32, R1 \
	SUB   $8, R2 \
	CBNZ  R2, loop \
done: \
	RET

#define ELEMNUM_F32(sym, OP) \
TEXT sym(SB), NOSPLIT, $0-24 \
	MOVD  a+0(FP), R0 \
	FMOVS v+8(FP), F0 \
	VDUP  V0.S[0], V4.S4 \
	MOVD  n+16(FP), R2 \
	CBZ   R2, done \
loop: \
	VLD1  (R0), [V0.S4, V1.S4] \
	OP    V4.S4, V0.S4, V0.S4 \
	OP    V4.S4, V1.S4, V1.S4 \
	VST1  [V0.S4, V1.S4], (R0) \
	ADD   $32, R0 \
	SUB   $8, R2 \
	CBNZ  R2, loop \
done: \
	RET

#define UNARY_F32(sym, OP) \
TEXT sym(SB), NOSPLIT, $0-16 \
	MOVD a+0(FP), R0 \
	MOVD n+8(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1  (R0), [V0.S4, V1.S4] \
	OP    V0.S4, V0.S4 \
	OP    V1.S4, V1.S4 \
	VST1  [V0.S4, V1.S4], (R0) \
	ADD   $32, R0 \
	SUB   $8, R2 \
	CBNZ  R2, loop \
done: \
	RET

ELEM2_F32(·neonAddF32, VFADD)
ELEM2_F32(·neonSubF32, VFSUB)
ELEM2_F32(·neonMulF32, VFMUL)
ELEM2_F32(·neonDivF32, VFDIV)
ELEM2_F32(·neonMinimumF32, VFMIN)
ELEM2_F32(·neonMaximumF32, VFMAX)

ELEMNUM_F32(·neonAddNumF32, VFADD)
ELEMNUM_F32(·neonSubNumF32, VFSUB)
ELEMNUM_F32(·neonMulNumF32, VFMUL)
ELEMNUM_F32(·neonDivNumF32, VFDIV)
ELEMNUM_F32(·neonMinNumF32, VFMIN)
ELEMNUM_F32(·neonMaxNumF32, VFMAX)

UNARY_F32(·neonAbsF32, VFABS)
UNARY_F32(·neonNegF32, VFNEG)
UNARY_F32(·neonSqrtF32, VFSQRT)
UNARY_F32(·neonRoundF32, VFRINTN)
UNARY_F32(·neonFloorF32, VFRINTM)
UNARY_F32(·neonCeilF32, VFRINTP)

// a[i] = 1 / a[i]
TEXT ·neonInvF32(SB), NOSPLIT, $0-16
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	MOVW ·neonOneF32(SB), R3
	VDUP R3, V4.S4
	CBZ  R2, done
loop:
	VLD1  (R0), [V0.S4, V1.S4]
	VFDIV V0.S4, V4.S4, V0.S4
	VFDIV V1.S4, V4.S4, V1.S4
	VST1  [V0.S4, V1.S4], (R0)
	ADD   $32, R0
	SUB   $8, R2
	CBNZ  R2, loop
done:
	RET

// ---------------------------------------------------------------- f32 归约

#define REDUCE_F32(sym, OP) \
TEXT sym(SB), NOSPLIT, $0-20 \
	MOVD a+0(FP), R0 \
	MOVD n+8(FP), R2 \
	VEOR V0.B16, V0.B16, V0.B16 \
	VEOR V1.B16, V1.B16, V1.B16 \
	CBZ  R2, fold \
loop: \
	VLD1 (R0), [V2.S4, V3.S4] \
	OP   V2.S4, V0.S4, V0.S4 \
	OP   V3.S4, V1.S4, V1.S4 \
	ADD  $32, R0 \
	SUB  $8, R2 \
	CBNZ R2, loop \
fold: \
	VFADD   V1.S4, V0.S4, V0.S4 \
	VFADDP V0.S4, V0.S4, V0.S4 \
	VFADDP V0.S4, V0.S4, V0.S4 \
	VMOV   V0.S[0], R3 \
	MOVW   R3, ret+16(FP) \
	RET

REDUCE_F32(·neonSumF32, VFADD)

// 平方和
TEXT ·neonSumSqF32(SB), NOSPLIT, $0-20
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	VEOR V0.B16, V0.B16, V0.B16
	VEOR V1.B16, V1.B16, V1.B16
	CBZ  R2, fold
loop:
	VLD1  (R0), [V2.S4, V3.S4]
	VFMLA V2.S4, V2.S4, V0.S4
	VFMLA V3.S4, V3.S4, V1.S4
	ADD   $32, R0
	SUB   $8, R2
	CBNZ  R2, loop
fold:
	VFADD   V1.S4, V0.S4, V0.S4
	VFADDP V0.S4, V0.S4, V0.S4
	VFADDP V0.S4, V0.S4, V0.S4
	VMOV   V0.S[0], R3
	MOVW   R3, ret+16(FP)
	RET

// 点积
TEXT ·neonDotF32(SB), NOSPLIT, $0-28
	MOVD a+0(FP), R0
	MOVD b+8(FP), R1
	MOVD n+16(FP), R2
	VEOR V0.B16, V0.B16, V0.B16
	VEOR V1.B16, V1.B16, V1.B16
	CBZ  R2, fold
loop:
	VLD1  (R0), [V2.S4, V3.S4]
	VLD1  (R1), [V4.S4, V5.S4]
	VFMLA V2.S4, V4.S4, V0.S4
	VFMLA V3.S4, V5.S4, V1.S4
	ADD   $32, R0
	ADD   $32, R1
	SUB   $8, R2
	CBNZ  R2, loop
fold:
	VFADD   V1.S4, V0.S4, V0.S4
	VFADDP V0.S4, V0.S4, V0.S4
	VFADDP V0.S4, V0.S4, V0.S4
	VMOV   V0.S[0], R3
	MOVW   R3, ret+24(FP)
	RET

// 绝对值之和
TEXT ·neonAbsSumF32(SB), NOSPLIT, $0-20
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	VEOR V0.B16, V0.B16, V0.B16
	VEOR V1.B16, V1.B16, V1.B16
	CBZ  R2, fold
loop:
	VLD1  (R0), [V2.S4, V3.S4]
	VFABS V2.S4, V4.S4
	VFABS V3.S4, V5.S4
	VFADD  V4.S4, V0.S4, V0.S4
	VFADD  V5.S4, V1.S4, V1.S4
	ADD   $32, R0
	SUB   $8, R2
	CBNZ  R2, loop
fold:
	VFADD   V1.S4, V0.S4, V0.S4
	VFADDP V0.S4, V0.S4, V0.S4
	VFADDP V0.S4, V0.S4, V0.S4
	VMOV   V0.S[0], R3
	MOVW   R3, ret+16(FP)
	RET

// 差值平方和
TEXT ·neonDiffSqSumF32(SB), NOSPLIT, $0-28
	MOVD a+0(FP), R0
	MOVD b+8(FP), R1
	MOVD n+16(FP), R2
	VEOR V0.B16, V0.B16, V0.B16
	VEOR V1.B16, V1.B16, V1.B16
	CBZ  R2, fold
loop:
	VLD1  (R0), [V2.S4, V3.S4]
	VLD1  (R1), [V4.S4, V5.S4]
	VFSUB V4.S4, V2.S4, V6.S4
	VFSUB V5.S4, V3.S4, V7.S4
	VFMLA V6.S4, V6.S4, V0.S4
	VFMLA V7.S4, V7.S4, V1.S4
	ADD   $32, R0
	ADD   $32, R1
	SUB   $8, R2
	CBNZ  R2, loop
fold:
	VFADD   V1.S4, V0.S4, V0.S4
	VFADDP V0.S4, V0.S4, V0.S4
	VFADDP V0.S4, V0.S4, V0.S4
	VMOV   V0.S[0], R3
	MOVW   R3, ret+24(FP)
	RET

// 差值绝对值之和
TEXT ·neonAbsDiffSumF32(SB), NOSPLIT, $0-28
	MOVD a+0(FP), R0
	MOVD b+8(FP), R1
	MOVD n+16(FP), R2
	VEOR V0.B16, V0.B16, V0.B16
	VEOR V1.B16, V1.B16, V1.B16
	CBZ  R2, fold
loop:
	VLD1  (R0), [V2.S4, V3.S4]
	VLD1  (R1), [V4.S4, V5.S4]
	VFSUB V4.S4, V2.S4, V6.S4
	VFSUB V5.S4, V3.S4, V7.S4
	VFABS V6.S4, V6.S4
	VFABS V7.S4, V7.S4
	VFADD  V6.S4, V0.S4, V0.S4
	VFADD  V7.S4, V1.S4, V1.S4
	ADD   $32, R0
	ADD   $32, R1
	SUB   $8, R2
	CBNZ  R2, loop
fold:
	VFADD   V1.S4, V0.S4, V0.S4
	VFADDP V0.S4, V0.S4, V0.S4
	VFADDP V0.S4, V0.S4, V0.S4
	VMOV   V0.S[0], R3
	MOVW   R3, ret+24(FP)
	RET

// 乘积 (要求 n >= 8)
TEXT ·neonProdF32(SB), NOSPLIT, $0-20
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	VLD1  (R0), [V0.S4, V1.S4]
	ADD   $32, R0
	SUB   $8, R2
	CBZ   R2, fold
loop:
	VLD1  (R0), [V2.S4, V3.S4]
	VFMUL V2.S4, V0.S4, V0.S4
	VFMUL V3.S4, V1.S4, V1.S4
	ADD   $32, R0
	SUB   $8, R2
	CBNZ  R2, loop
fold:
	VFMUL V1.S4, V0.S4, V0.S4
	VMOV  V0.S[0], R3
	VMOV  V0.S[1], R4
	VMOV  V0.S[2], R5
	VMOV  V0.S[3], R6
	FMOVS R3, F0
	FMOVS R4, F1
	FMOVS R5, F2
	FMOVS R6, F3
	FMULS F1, F0, F0
	FMULS F2, F0, F0
	FMULS F3, F0, F0
	FMOVS F0, ret+16(FP)
	RET

// 最小值 (要求 n >= 8)
TEXT ·neonMinF32(SB), NOSPLIT, $0-20
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	VLD1  (R0), [V0.S4, V1.S4]
	ADD   $32, R0
	SUB   $8, R2
	CBZ   R2, fold
loop:
	VLD1  (R0), [V2.S4, V3.S4]
	VFMIN V2.S4, V0.S4, V0.S4
	VFMIN V3.S4, V1.S4, V1.S4
	ADD   $32, R0
	SUB   $8, R2
	CBNZ  R2, loop
fold:
	VFMIN  V1.S4, V0.S4, V0.S4
	VFMINP V0.S4, V0.S4, V0.S4
	VFMINP V0.S4, V0.S4, V0.S4
	VMOV   V0.S[0], R3
	MOVW   R3, ret+16(FP)
	RET

// 最大值 (要求 n >= 8)
TEXT ·neonMaxF32(SB), NOSPLIT, $0-20
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	VLD1  (R0), [V0.S4, V1.S4]
	ADD   $32, R0
	SUB   $8, R2
	CBZ   R2, fold
loop:
	VLD1  (R0), [V2.S4, V3.S4]
	VFMAX V2.S4, V0.S4, V0.S4
	VFMAX V3.S4, V1.S4, V1.S4
	ADD   $32, R0
	SUB   $8, R2
	CBNZ  R2, loop
fold:
	VFMAX  V1.S4, V0.S4, V0.S4
	VFMAXP V0.S4, V0.S4, V0.S4
	VFMAXP V0.S4, V0.S4, V0.S4
	VMOV   V0.S[0], R3
	MOVW   R3, ret+16(FP)
	RET

// 重复填充
TEXT ·neonRepeatF32(SB), NOSPLIT, $0-24
	MOVD  a+0(FP), R0
	FMOVS v+8(FP), F0
	VDUP  V0.S[0], V1.S4
	VMOV  V1.B16, V2.B16
	MOVD  n+16(FP), R2
	CBZ   R2, done
loop:
	VST1 [V1.S4, V2.S4], (R0)
	ADD  $32, R0
	SUB  $8, R2
	CBNZ R2, loop
done:
	RET

// 查找 (返回首个下标, 未找到返回 -1)
TEXT ·neonFindF32(SB), NOSPLIT, $0-32
	MOVD  a+0(FP), R0
	MOVD  n+8(FP), R1
	FMOVS v+16(FP), F0
	VDUP  V0.S[0], V1.S4
	MOVD  $-1, R8
	MOVD  $0, R2
	CBZ   R1, done
loop:
	VLD1   (R0), [V2.S4, V3.S4]
	VFCMEQ V1.S4, V2.S4, V4.S4
	VFCMEQ V1.S4, V3.S4, V5.S4
	VMOV   V4.S[0], R4
	VMOV   V4.S[1], R5
	VMOV   V4.S[2], R6
	VMOV   V4.S[3], R7
	CBNZ   R4, hit0
	CBNZ   R5, hit1
	CBNZ   R6, hit2
	CBNZ   R7, hit3
	VMOV   V5.S[0], R4
	VMOV   V5.S[1], R5
	VMOV   V5.S[2], R6
	VMOV   V5.S[3], R7
	CBNZ   R4, hit4
	CBNZ   R5, hit5
	CBNZ   R6, hit6
	CBNZ   R7, hit7
	ADD    $8, R2
	ADD    $32, R0
	SUB    $8, R1
	CBNZ   R1, loop
done:
	MOVD R8, ret+24(FP)
	RET
hit0:
	MOVD R2, ret+24(FP)
	RET
hit1:
	ADD  $1, R2
	MOVD R2, ret+24(FP)
	RET
hit2:
	ADD  $2, R2
	MOVD R2, ret+24(FP)
	RET
hit3:
	ADD  $3, R2
	MOVD R2, ret+24(FP)
	RET
hit4:
	ADD  $4, R2
	MOVD R2, ret+24(FP)
	RET
hit5:
	ADD  $5, R2
	MOVD R2, ret+24(FP)
	RET
hit6:
	ADD  $6, R2
	MOVD R2, ret+24(FP)
	RET
hit7:
	ADD  $7, R2
	MOVD R2, ret+24(FP)
	RET

// ---------------------------------------------------------------- f32 比较

#define CMP_XOPY_F32(sym, MN) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD a+0(FP), R0 \
	MOVD b+8(FP), R1 \
	MOVD d+16(FP), R3 \
	MOVD n+24(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1 (R0), [V0.S4, V1.S4] \
	VLD1 (R1), [V2.S4, V3.S4] \
	MN   V2.S4, V0.S4, V4.S4 \
	MN   V3.S4, V1.S4, V5.S4 \
	VMOV V4.S[0], R4 \
	VMOV V4.S[1], R5 \
	VMOV V4.S[2], R6 \
	VMOV V4.S[3], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	VMOV V5.S[0], R4 \
	VMOV V5.S[1], R5 \
	VMOV V5.S[2], R6 \
	VMOV V5.S[3], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	MOVB R4, 4(R3) \
	MOVB R5, 5(R3) \
	MOVB R6, 6(R3) \
	MOVB R7, 7(R3) \
	ADD  $32, R0 \
	ADD  $32, R1 \
	ADD  $8, R3 \
	SUB  $8, R2 \
	CBNZ R2, loop \
done: \
	RET

#define CMP_YOPX_F32(sym, MN) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD a+0(FP), R0 \
	MOVD b+8(FP), R1 \
	MOVD d+16(FP), R3 \
	MOVD n+24(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1 (R0), [V0.S4, V1.S4] \
	VLD1 (R1), [V2.S4, V3.S4] \
	MN   V0.S4, V2.S4, V4.S4 \
	MN   V1.S4, V3.S4, V5.S4 \
	VMOV V4.S[0], R4 \
	VMOV V4.S[1], R5 \
	VMOV V4.S[2], R6 \
	VMOV V4.S[3], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	VMOV V5.S[0], R4 \
	VMOV V5.S[1], R5 \
	VMOV V5.S[2], R6 \
	VMOV V5.S[3], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	MOVB R4, 4(R3) \
	MOVB R5, 5(R3) \
	MOVB R6, 6(R3) \
	MOVB R7, 7(R3) \
	ADD  $32, R0 \
	ADD  $32, R1 \
	ADD  $8, R3 \
	SUB  $8, R2 \
	CBNZ R2, loop \
done: \
	RET

#define CMPN_YOPX_F32(sym, MN) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD a+0(FP), R0 \
	MOVD b+8(FP), R1 \
	MOVD d+16(FP), R3 \
	MOVD n+24(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1 (R0), [V0.S4, V1.S4] \
	VLD1 (R1), [V2.S4, V3.S4] \
	MN   V0.S4, V2.S4, V4.S4 \
	MN   V1.S4, V3.S4, V5.S4 \
	VMOV V4.S[0], R4 \
	VMOV V4.S[1], R5 \
	VMOV V4.S[2], R6 \
	VMOV V4.S[3], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	EOR  $1, R4, R4 \
	EOR  $1, R5, R5 \
	EOR  $1, R6, R6 \
	EOR  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	VMOV V5.S[0], R4 \
	VMOV V5.S[1], R5 \
	VMOV V5.S[2], R6 \
	VMOV V5.S[3], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	EOR  $1, R4, R4 \
	EOR  $1, R5, R5 \
	EOR  $1, R6, R6 \
	EOR  $1, R7, R7 \
	MOVB R4, 4(R3) \
	MOVB R5, 5(R3) \
	MOVB R6, 6(R3) \
	MOVB R7, 7(R3) \
	ADD  $32, R0 \
	ADD  $32, R1 \
	ADD  $8, R3 \
	SUB  $8, R2 \
	CBNZ R2, loop \
done: \
	RET

CMP_YOPX_F32(·neonLtF32, VFCMGT)
CMP_YOPX_F32(·neonLteF32, VFCMGE)
CMP_XOPY_F32(·neonGtF32, VFCMGT)
CMP_XOPY_F32(·neonGteF32, VFCMGE)
CMP_YOPX_F32(·neonEqF32, VFCMEQ)
CMPN_YOPX_F32(·neonNeqF32, VFCMEQ)

// 与标量比较: dst[i] = MN(A, B) != 0
#define CMPNUM_F32(sym, MN, A1, B1, A2, B2) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD  a+0(FP), R0 \
	MOVD  d+8(FP), R3 \
	MOVD  n+16(FP), R2 \
	FMOVS v+24(FP), F0 \
	VDUP  V0.S[0], V4.S4 \
	CBZ   R2, done \
loop: \
	VLD1 (R0), [V0.S4, V1.S4] \
	MN   A1, B1, V2.S4 \
	MN   A2, B2, V3.S4 \
	VMOV V2.S[0], R4 \
	VMOV V2.S[1], R5 \
	VMOV V2.S[2], R6 \
	VMOV V2.S[3], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	VMOV V3.S[0], R4 \
	VMOV V3.S[1], R5 \
	VMOV V3.S[2], R6 \
	VMOV V3.S[3], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	MOVB R4, 4(R3) \
	MOVB R5, 5(R3) \
	MOVB R6, 6(R3) \
	MOVB R7, 7(R3) \
	ADD  $32, R0 \
	ADD  $8, R3 \
	SUB  $8, R2 \
	CBNZ R2, loop \
done: \
	RET

// 与标量比较: dst[i] = MN(A, B) == 0
#define CMPNUMN_F32(sym, MN, A1, B1, A2, B2) \
TEXT sym(SB), NOSPLIT, $0-32 \
	MOVD  a+0(FP), R0 \
	MOVD  d+8(FP), R3 \
	MOVD  n+16(FP), R2 \
	FMOVS v+24(FP), F0 \
	VDUP  V0.S[0], V4.S4 \
	CBZ   R2, done \
loop: \
	VLD1 (R0), [V0.S4, V1.S4] \
	MN   A1, B1, V2.S4 \
	MN   A2, B2, V3.S4 \
	VMOV V2.S[0], R4 \
	VMOV V2.S[1], R5 \
	VMOV V2.S[2], R6 \
	VMOV V2.S[3], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	EOR  $1, R4, R4 \
	EOR  $1, R5, R5 \
	EOR  $1, R6, R6 \
	EOR  $1, R7, R7 \
	MOVB R4, 0(R3) \
	MOVB R5, 1(R3) \
	MOVB R6, 2(R3) \
	MOVB R7, 3(R3) \
	VMOV V3.S[0], R4 \
	VMOV V3.S[1], R5 \
	VMOV V3.S[2], R6 \
	VMOV V3.S[3], R7 \
	AND  $1, R4, R4 \
	AND  $1, R5, R5 \
	AND  $1, R6, R6 \
	AND  $1, R7, R7 \
	EOR  $1, R4, R4 \
	EOR  $1, R5, R5 \
	EOR  $1, R6, R6 \
	EOR  $1, R7, R7 \
	MOVB R4, 4(R3) \
	MOVB R5, 5(R3) \
	MOVB R6, 6(R3) \
	MOVB R7, 7(R3) \
	ADD  $32, R0 \
	ADD  $8, R3 \
	SUB  $8, R2 \
	CBNZ R2, loop \
done: \
	RET

CMPNUM_F32(·neonLtNumF32, VFCMGT, V0.S4, V4.S4, V1.S4, V4.S4)
CMPNUM_F32(·neonLteNumF32, VFCMGE, V0.S4, V4.S4, V1.S4, V4.S4)
CMPNUM_F32(·neonGtNumF32, VFCMGT, V4.S4, V0.S4, V4.S4, V1.S4)
CMPNUM_F32(·neonGteNumF32, VFCMGE, V4.S4, V0.S4, V4.S4, V1.S4)
CMPNUM_F32(·neonEqNumF32, VFCMEQ, V0.S4, V4.S4, V1.S4, V4.S4)
CMPNUMN_F32(·neonNeqNumF32, VFCMEQ, V0.S4, V4.S4, V1.S4, V4.S4)

// ---------------------------------------------------------------- bool

// x[i] = !x[i]
TEXT ·neonNotB(SB), NOSPLIT, $0-16
	MOVD  a+0(FP), R0
	MOVD  n+8(FP), R2
	VMOVI $1, V4.B16
	CBZ   R2, done
loop:
	VLD1  (R0), [V0.B16]
	VEOR  V4.B16, V0.B16, V0.B16
	VST1  [V0.B16], (R0)
	ADD   $16, R0
	SUB   $16, R2
	CBNZ  R2, loop
done:
	RET

#define BOOL2(sym, OP) \
TEXT sym(SB), NOSPLIT, $0-24 \
	MOVD a+0(FP), R0 \
	MOVD b+8(FP), R1 \
	MOVD n+16(FP), R2 \
	CBZ  R2, done \
loop: \
	VLD1 (R0), [V0.B16] \
	VLD1 (R1), [V1.B16] \
	OP   V1.B16, V0.B16, V0.B16 \
	VST1 [V0.B16], (R0) \
	ADD  $16, R0 \
	ADD  $16, R1 \
	SUB  $16, R2 \
	CBNZ R2, loop \
done: \
	RET

BOOL2(·neonAndB, VAND)
BOOL2(·neonOrB, VORR)
BOOL2(·neonXorB, VEOR)

// ---------------------------------------------------------------- 矩阵

// dst(m*p) += x(m*n) * y(n*psimd), 行跨度 p
TEXT ·neonMatMulF64(SB), NOSPLIT, $0-56
	MOVD dst+0(FP), R0
	MOVD x+8(FP), R1
	MOVD y+16(FP), R2
	MOVD m+24(FP), R3
	MOVD n+32(FP), R4
	MOVD psimd+40(FP), R5
	MOVD p+48(FP), R6
	LSL  $3, R6, R13
	CBZ  R3, done
rowloop:
	MOVD R4, R8
	MOVD R2, R9
	CBZ  R8, nextrow
kloop:
	FMOVD (R1), F0
	VDUP  V0.D[0], V4.D2
	MOVD  R0, R10
	MOVD  R9, R11
	MOVD  R5, R12
jloop:
	VLD1  (R10), [V0.D2, V1.D2]
	VLD1  (R11), [V2.D2, V3.D2]
	VFMLA V4.D2, V2.D2, V0.D2
	VFMLA V4.D2, V3.D2, V1.D2
	VST1  [V0.D2, V1.D2], (R10)
	ADD   $32, R10
	ADD   $32, R11
	SUB   $4, R12
	CBNZ  R12, jloop
	ADD   $8, R1
	ADD   R13, R9
	SUB   $1, R8
	CBNZ  R8, kloop
nextrow:
	ADD  R13, R0
	SUB  $1, R3
	CBNZ R3, rowloop
done:
	RET

// dst(4*4) += x(4*4) * y(4*4)
TEXT ·neonMat4MulF64(SB), NOSPLIT, $0-24
	MOVD dst+0(FP), R0
	MOVD x+8(FP), R1
	MOVD y+16(FP), R2
	MOVD $4, R3
rowloop:
	VLD1 (R2), [V4.D2, V5.D2]
	ADD  $32, R2
	VLD1 (R2), [V6.D2, V7.D2]
	ADD  $32, R2
	VLD1 (R2), [V8.D2, V9.D2]
	ADD  $32, R2
	VLD1 (R2), [V10.D2, V11.D2]
	SUB  $96, R2
	FMOVD (R1), F0
	VDUP  V0.D[0], V12.D2
	VFMUL V12.D2, V4.D2, V20.D2
	VFMUL V12.D2, V5.D2, V21.D2
	FMOVD 8(R1), F0
	VDUP  V0.D[0], V12.D2
	VFMLA V12.D2, V6.D2, V20.D2
	VFMLA V12.D2, V7.D2, V21.D2
	FMOVD 16(R1), F0
	VDUP  V0.D[0], V12.D2
	VFMLA V12.D2, V8.D2, V20.D2
	VFMLA V12.D2, V9.D2, V21.D2
	FMOVD 24(R1), F0
	VDUP  V0.D[0], V12.D2
	VFMLA V12.D2, V10.D2, V20.D2
	VFMLA V12.D2, V11.D2, V21.D2
	VST1 [V20.D2, V21.D2], (R0)
	ADD  $32, R0
	ADD  $32, R1
	SUB  $1, R3
	CBNZ R3, rowloop
	RET

// dst(m*p) += x(m*n) * y(n*psimd), 行跨度 p
TEXT ·neonMatMulF32(SB), NOSPLIT, $0-56
	MOVD dst+0(FP), R0
	MOVD x+8(FP), R1
	MOVD y+16(FP), R2
	MOVD m+24(FP), R3
	MOVD n+32(FP), R4
	MOVD psimd+40(FP), R5
	MOVD p+48(FP), R6
	LSL  $2, R6, R13
	CBZ  R3, done
rowloop:
	MOVD R4, R8
	MOVD R2, R9
	CBZ  R8, nextrow
kloop:
	FMOVS (R1), F0
	VDUP  V0.S[0], V4.S4
	MOVD  R0, R10
	MOVD  R9, R11
	MOVD  R5, R12
jloop:
	VLD1  (R10), [V0.S4, V1.S4]
	VLD1  (R11), [V2.S4, V3.S4]
	VFMLA V4.S4, V2.S4, V0.S4
	VFMLA V4.S4, V3.S4, V1.S4
	VST1  [V0.S4, V1.S4], (R10)
	ADD   $32, R10
	ADD   $32, R11
	SUB   $8, R12
	CBNZ  R12, jloop
	ADD   $4, R1
	ADD   R13, R9
	SUB   $1, R8
	CBNZ  R8, kloop
nextrow:
	ADD  R13, R0
	SUB  $1, R3
	CBNZ R3, rowloop
done:
	RET

// dst(4*4) += x(4*4) * y(4*4)
TEXT ·neonMat4MulF32(SB), NOSPLIT, $0-24
	MOVD dst+0(FP), R0
	MOVD x+8(FP), R1
	MOVD y+16(FP), R2
	MOVD $4, R3
rowloop:
	VLD1  (R2), [V4.S4, V5.S4]
	ADD   $32, R2
	VLD1  (R2), [V6.S4, V7.S4]
	SUB   $32, R2
	FMOVS (R1), F0
	VDUP  V0.S[0], V12.S4
	VFMUL V12.S4, V4.S4, V20.S4
	FMOVS 4(R1), F0
	VDUP  V0.S[0], V12.S4
	VFMLA V12.S4, V5.S4, V20.S4
	FMOVS 8(R1), F0
	VDUP  V0.S[0], V12.S4
	VFMLA V12.S4, V6.S4, V20.S4
	FMOVS 12(R1), F0
	VDUP  V0.S[0], V12.S4
	VFMLA V12.S4, V7.S4, V20.S4
	VST1  [V20.S4], (R0)
	ADD   $16, R0
	ADD   $16, R1
	SUB   $1, R3
	CBNZ  R3, rowloop
	RET

// ---------------------------------------------------------------- f32 超越函数
//
// exp/log 采用 Cephes 单精度算法(与 amd64 内核同源), 4 路 NEON, 每次处理 8 个元素;
// 要求 n 为 8 的整数倍. 特殊值(NaN / ±Inf / 非正数)按掩码在向量内选择, 与 Go 语义一致.

// exp 常量
DATA ·neonExpConstF32+0(SB)/4, $0x42b0c0a5   // 88.3762626647949
DATA ·neonExpConstF32+4(SB)/4, $0xc2b0c0a5   // -88.3762626647949
DATA ·neonExpConstF32+8(SB)/4, $0x3fb8aa3b   // log2(e)
DATA ·neonExpConstF32+12(SB)/4, $0x3f000000  // 0.5
DATA ·neonExpConstF32+16(SB)/4, $0xbf318000  // -ln2 高位
DATA ·neonExpConstF32+20(SB)/4, $0x395e8083  // -ln2 低位
DATA ·neonExpConstF32+24(SB)/4, $0x39506967  // p0
DATA ·neonExpConstF32+28(SB)/4, $0x3ab743ce  // p1
DATA ·neonExpConstF32+32(SB)/4, $0x3c088908  // p2
DATA ·neonExpConstF32+36(SB)/4, $0x3d2aa9c1  // p3
DATA ·neonExpConstF32+40(SB)/4, $0x3e2aaaaa  // p4
DATA ·neonExpConstF32+44(SB)/4, $0x3f000000  // p5
DATA ·neonExpConstF32+48(SB)/4, $0x3f800000  // 1.0
DATA ·neonExpConstF32+52(SB)/4, $0x0000007f  // 127
GLOBL ·neonExpConstF32(SB), RODATA, $56

// log 常量
DATA ·neonLogConstF32+0(SB)/4, $0x007fffff   // 尾数掩码
DATA ·neonLogConstF32+4(SB)/4, $0x3f000000   // 0.5 的指数位
DATA ·neonLogConstF32+8(SB)/4, $0x0000007e   // 126
DATA ·neonLogConstF32+12(SB)/4, $0xbf800000  // -1.0
DATA ·neonLogConstF32+16(SB)/4, $0x3f3504f3  // sqrt(0.5)
DATA ·neonLogConstF32+20(SB)/4, $0x3d9021bb  // p0
DATA ·neonLogConstF32+24(SB)/4, $0xbdebd1b8  // p1
DATA ·neonLogConstF32+28(SB)/4, $0x3def251a  // p2
DATA ·neonLogConstF32+32(SB)/4, $0xbdfe5d4f  // p3
DATA ·neonLogConstF32+36(SB)/4, $0x3e11e9bf  // p4
DATA ·neonLogConstF32+40(SB)/4, $0xbe2aae50  // p5
DATA ·neonLogConstF32+44(SB)/4, $0x3e4cceac  // p6
DATA ·neonLogConstF32+48(SB)/4, $0xbe7ffffc  // p7
DATA ·neonLogConstF32+52(SB)/4, $0x3eaaaaaa  // p8
DATA ·neonLogConstF32+56(SB)/4, $0xb95e8083  // ln2 低位
DATA ·neonLogConstF32+60(SB)/4, $0x3f318000  // ln2 高位
DATA ·neonLogConstF32+64(SB)/4, $0x3f000000  // 0.5
DATA ·neonLogConstF32+68(SB)/4, $0x3f800000  // 结果比例: 1.0
DATA ·neonLogConstF32+72(SB)/4, $0x3fb8aa3b  // 结果比例: log2(e)
DATA ·neonLogConstF32+76(SB)/4, $0x3ede5bd9  // 结果比例: log10(e)
GLOBL ·neonLogConstF32(SB), RODATA, $80

// 特殊值
DATA ·neonSpecialF32+0(SB)/4, $0x7f800000    // +Inf
DATA ·neonSpecialF32+4(SB)/4, $0x7fc00000    // NaN
DATA ·neonSpecialF32+8(SB)/4, $0xff800000    // -Inf
GLOBL ·neonSpecialF32(SB), RODATA, $12

// 单个 4 路向量的 exp: 结果写回 X.
// x 先夹到 [-88.376, 88.376], 再用 k = floor(x*log2(e)+0.5) 做 2^k * exp(r) 分解.
#define NEON_EXP_F32_STEP(X, XS) \
	VFMIN  V16.S4, XS, V6.S4 \
	VFMAX  V17.S4, V6.S4, V6.S4 \
	VFMUL  V18.S4, V6.S4, V7.S4 \
	VFADD  V19.S4, V7.S4, V7.S4 \
	VFRINTM V7.S4, V7.S4 \
	VFMUL  V20.S4, V7.S4, V8.S4 \
	VFMLA  V21.S4, V7.S4, V8.S4 \
	VFADD  V6.S4, V8.S4, V8.S4 \
	VFMUL  V8.S4, V8.S4, V9.S4 \
	VMOV   V22.B16, V10.B16 \
	VFMUL  V8.S4, V10.S4, V10.S4 \
	VFADD  V23.S4, V10.S4, V10.S4 \
	VFMUL  V8.S4, V10.S4, V10.S4 \
	VFADD  V24.S4, V10.S4, V10.S4 \
	VFMUL  V8.S4, V10.S4, V10.S4 \
	VFADD  V25.S4, V10.S4, V10.S4 \
	VFMUL  V8.S4, V10.S4, V10.S4 \
	VFADD  V26.S4, V10.S4, V10.S4 \
	VFMUL  V8.S4, V10.S4, V10.S4 \
	VFADD  V27.S4, V10.S4, V10.S4 \
	VFMUL  V9.S4, V10.S4, V10.S4 \
	VFADD  V8.S4, V10.S4, V10.S4 \
	VFADD  V28.S4, V10.S4, V10.S4 \
	VFCVTZS V7.S4, V11.S4 \
	VADD   V29.S4, V11.S4, V11.S4 \
	VSHL   $23, V11.S4, V11.S4 \
	VFMUL  V11.S4, V10.S4, V10.S4 \
	VFCMGT V16.S4, XS, V12.S4 \
	VBSL   V10.B16, V14.B16, V12.B16 \
	VFCMEQ XS, XS, V13.S4 \
	VBSL   V3.B16, V12.B16, V13.B16 \
	VMOV   V13.B16, X.B16

TEXT ·neonExpF32(SB), NOSPLIT, $0-16
	MOVD a+0(FP), R0
	MOVD n+8(FP), R2
	CBZ  R2, done
	FMOVS ·neonExpConstF32+0(SB), F0
	VDUP  V0.S[0], V16.S4
	FMOVS ·neonExpConstF32+4(SB), F0
	VDUP  V0.S[0], V17.S4
	FMOVS ·neonExpConstF32+8(SB), F0
	VDUP  V0.S[0], V18.S4
	FMOVS ·neonExpConstF32+12(SB), F0
	VDUP  V0.S[0], V19.S4
	FMOVS ·neonExpConstF32+16(SB), F0
	VDUP  V0.S[0], V20.S4
	FMOVS ·neonExpConstF32+20(SB), F0
	VDUP  V0.S[0], V21.S4
	FMOVS ·neonExpConstF32+24(SB), F0
	VDUP  V0.S[0], V22.S4
	FMOVS ·neonExpConstF32+28(SB), F0
	VDUP  V0.S[0], V23.S4
	FMOVS ·neonExpConstF32+32(SB), F0
	VDUP  V0.S[0], V24.S4
	FMOVS ·neonExpConstF32+36(SB), F0
	VDUP  V0.S[0], V25.S4
	FMOVS ·neonExpConstF32+40(SB), F0
	VDUP  V0.S[0], V26.S4
	FMOVS ·neonExpConstF32+44(SB), F0
	VDUP  V0.S[0], V27.S4
	FMOVS ·neonExpConstF32+48(SB), F0
	VDUP  V0.S[0], V28.S4
	FMOVS ·neonExpConstF32+52(SB), F0
	VDUP  V0.S[0], V29.S4
	FMOVS ·neonSpecialF32+0(SB), F0
	VDUP  V0.S[0], V14.S4
	FMOVS ·neonSpecialF32+4(SB), F0
	VDUP  V0.S[0], V3.S4
loop:
	VLD1 (R0), [V4.S4, V5.S4]
	NEON_EXP_F32_STEP(V4, V4.S4)
	NEON_EXP_F32_STEP(V5, V5.S4)
	VST1 [V4.S4, V5.S4], (R0)
	ADD  $32, R0
	SUB  $8, R2
	CBNZ R2, loop
done:
	RET

// 单个 4 路向量的 log: 结果写回 X.
// frexp 拆出尾数 m ∈ [0.5,1) 与指数 e, 按 m 与 sqrt(0.5) 的关系把 e 归一到 atanh 多项式区间.
#define NEON_LOG_F32_STEP(X, XS, XB) \
	VUSHR  $23, XS, V6.S4 \
	VSUB   V18.S4, V6.S4, V6.S4 \
	VAND   V16.B16, XB, V7.B16 \
	VORR   V17.B16, V7.B16, V7.B16 \
	VFCMGT V7.S4, V20.S4, V8.S4 \
	VADD   V8.S4, V6.S4, V6.S4 \
	VFADD  V7.S4, V7.S4, V9.S4 \
	VFADD  V19.S4, V9.S4, V9.S4 \
	VFADD  V19.S4, V7.S4, V13.S4 \
	VMOV   V8.B16, V10.B16 \
	VBSL   V13.B16, V9.B16, V10.B16 \
	VFMUL  V10.S4, V10.S4, V11.S4 \
	VMOV   V21.B16, V12.B16 \
	VFMUL  V10.S4, V12.S4, V12.S4 \
	VFADD  V22.S4, V12.S4, V12.S4 \
	VFMUL  V10.S4, V12.S4, V12.S4 \
	VFADD  V23.S4, V12.S4, V12.S4 \
	VFMUL  V10.S4, V12.S4, V12.S4 \
	VFADD  V24.S4, V12.S4, V12.S4 \
	VFMUL  V10.S4, V12.S4, V12.S4 \
	VFADD  V25.S4, V12.S4, V12.S4 \
	VFMUL  V10.S4, V12.S4, V12.S4 \
	VFADD  V26.S4, V12.S4, V12.S4 \
	VFMUL  V10.S4, V12.S4, V12.S4 \
	VFADD  V27.S4, V12.S4, V12.S4 \
	VFMUL  V10.S4, V12.S4, V12.S4 \
	VFADD  V28.S4, V12.S4, V12.S4 \
	VFMUL  V10.S4, V12.S4, V12.S4 \
	VFADD  V29.S4, V12.S4, V12.S4 \
	VFMUL  V10.S4, V12.S4, V12.S4 \
	VFMUL  V11.S4, V12.S4, V12.S4 \
	VSCVTF V6.S4, V6.S4 \
	VFMLA  V30.S4, V6.S4, V12.S4 \
	VFMLS  V0.S4, V11.S4, V12.S4 \
	VFADD  V10.S4, V12.S4, V12.S4 \
	VFMLA  V31.S4, V6.S4, V12.S4 \
	VFMUL  V1.S4, V12.S4, V12.S4 \
	VFCMGT V2.S4, XS, V6.S4 \
	VFCMEQ V2.S4, XS, V7.S4 \
	VFCMEQ V14.S4, XS, V8.S4 \
	VBSL   V3.B16, V15.B16, V7.B16 \
	VBSL   V7.B16, V12.B16, V6.B16 \
	VBSL   V6.B16, V14.B16, V8.B16 \
	VMOV   V8.B16, X.B16

// log 族内核: SCALE 是结果比例常量的偏移(log 1.0 / log2 log2(e) / log10 log10(e))
#define LOG_F32_KERNEL(sym, SCALE) \
TEXT sym(SB), NOSPLIT, $0-16 \
	MOVD a+0(FP), R0 \
	MOVD n+8(FP), R2 \
	CBZ  R2, done \
	FMOVS ·neonLogConstF32+64(SB), F13 \
	VDUP  V13.S[0], V0.S4 \
	FMOVS ·neonLogConstF32+SCALE(SB), F13 \
	VDUP  V13.S[0], V1.S4 \
	VMOVI $0, V2.B16 \
	FMOVS ·neonSpecialF32+4(SB), F13 \
	VDUP  V13.S[0], V3.S4 \
	FMOVS ·neonSpecialF32+0(SB), F13 \
	VDUP  V13.S[0], V14.S4 \
	FMOVS ·neonSpecialF32+8(SB), F13 \
	VDUP  V13.S[0], V15.S4 \
	FMOVS ·neonLogConstF32+0(SB), F13 \
	VDUP  V13.S[0], V16.S4 \
	FMOVS ·neonLogConstF32+4(SB), F13 \
	VDUP  V13.S[0], V17.S4 \
	FMOVS ·neonLogConstF32+8(SB), F13 \
	VDUP  V13.S[0], V18.S4 \
	FMOVS ·neonLogConstF32+12(SB), F13 \
	VDUP  V13.S[0], V19.S4 \
	FMOVS ·neonLogConstF32+16(SB), F13 \
	VDUP  V13.S[0], V20.S4 \
	FMOVS ·neonLogConstF32+20(SB), F13 \
	VDUP  V13.S[0], V21.S4 \
	FMOVS ·neonLogConstF32+24(SB), F13 \
	VDUP  V13.S[0], V22.S4 \
	FMOVS ·neonLogConstF32+28(SB), F13 \
	VDUP  V13.S[0], V23.S4 \
	FMOVS ·neonLogConstF32+32(SB), F13 \
	VDUP  V13.S[0], V24.S4 \
	FMOVS ·neonLogConstF32+36(SB), F13 \
	VDUP  V13.S[0], V25.S4 \
	FMOVS ·neonLogConstF32+40(SB), F13 \
	VDUP  V13.S[0], V26.S4 \
	FMOVS ·neonLogConstF32+44(SB), F13 \
	VDUP  V13.S[0], V27.S4 \
	FMOVS ·neonLogConstF32+48(SB), F13 \
	VDUP  V13.S[0], V28.S4 \
	FMOVS ·neonLogConstF32+52(SB), F13 \
	VDUP  V13.S[0], V29.S4 \
	FMOVS ·neonLogConstF32+56(SB), F13 \
	VDUP  V13.S[0], V30.S4 \
	FMOVS ·neonLogConstF32+60(SB), F13 \
	VDUP  V13.S[0], V31.S4 \
loop: \
	VLD1 (R0), [V4.S4, V5.S4] \
	NEON_LOG_F32_STEP(V4, V4.S4, V4.B16) \
	NEON_LOG_F32_STEP(V5, V5.S4, V5.B16) \
	VST1 [V4.S4, V5.S4], (R0) \
	ADD  $32, R0 \
	SUB  $8, R2 \
	CBNZ R2, loop \
done: \
	RET

LOG_F32_KERNEL(·neonLogF32, 68)
LOG_F32_KERNEL(·neonLog2F32, 72)
LOG_F32_KERNEL(·neonLog10F32, 76)
