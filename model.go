package qimage

type Imager interface {
	SetID(id uint32)

	SetName(name string)

	SetFileType(fileType string)

	GetRaw() []byte
	SetRaw(raw []byte)

	GetOID() uint32
	SetOID(oid uint32)

	GetSize() []int
	SetSize(size int)

	SetSortIndex(index int)
}
