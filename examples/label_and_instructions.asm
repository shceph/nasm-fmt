_start:
	extern page_directory
	mov eax, page_directory
	mov cr3, eax

	mov eax, cr0
	or eax, 0x80000001
	mov cr0, eax
