package internal

//nolint
//export proxy_on_memory_allocate
func proxyOnMemoryAllocate(size uint) *byte {
	buffer := make([]byte, size)
	return &buffer[0]
}
