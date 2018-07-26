//
// Test for AVX2 support in x86 processor. This includes
// vector instructions like vpgatherdd over 256-bit ymm registers
//

// func setAVX2()
TEXT ·setAVX2(SB),4,$0-0
	MOVL	$7,AX // eax=7,ecx=0 -> CPU extended features
	XORL	CX,CX
	CPUID
	SHRL	$5,BX // ebx bit 5 -> AVX2 present
	ANDL	$1,BX
	MOVB	BL,·hasAVX2+0(SB)
	RET
