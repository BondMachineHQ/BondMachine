package bcof

func NewBCOF(rsize uint32) *BCOFEntry {
	return &BCOFEntry{
		Rsize:     rsize,
		Id:        0,
		Signature: "",
		Data:      make([]*BCOFEntrySubentry, 0),
	}
}

func (b *BCOFEntry) SetId(id uint32) {
	b.Id = id
}

func (b *BCOFEntry) SetSignature(signature string) {
	b.Signature = signature
}
