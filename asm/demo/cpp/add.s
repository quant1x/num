	.text
	.intel_syntax noprefix
	.file	"add.c"
	.globl	Add                             # -- Begin function Add
	.p2align	4, 0x90
	.type	Add,@function
Add:                                    # @Add
# %bb.0:
	push	rbp
	mov	rbp, rsp
	and	rsp, -8
	mov	eax, dword ptr [rdi]
	mov	dword ptr [rsi], eax
	mov	rsp, rbp
	pop	rbp
	ret
.Lfunc_end0:
	.size	Add, .Lfunc_end0-Add
                                        # -- End function
	.ident	"Apple clang version 15.0.0 (clang-1500.1.0.2.5)"
	.section	".note.GNU-stack","",@progbits
	.addrsig
