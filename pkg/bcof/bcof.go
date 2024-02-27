package bcof

import "fmt"

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

func NewBCOFData(rsize uint32) *BCOFData {
	return &BCOFData{
		Rsize:     rsize,
		Id:        0,
		Signature: "",
		Payload:   make([]byte, 0),
	}
}

func (b *BCOFData) SetId(id uint32) {
	b.Id = id
}

func (b *BCOFData) SetSignature(signature string) {
	b.Signature = signature
}

func (b *BCOFData) Dump() string {
	regBytes := (b.Rsize + 7) / 8
	result := fmt.Sprintf("Rsize: %d\nId: %d\nSignature: %s\nPayload:", b.Rsize, b.Id, b.Signature)
	for i, ch := range b.Payload {
		if i%int(regBytes) == 0 {
			result += "\n\t0x"
		}
		result += fmt.Sprintf("%02x", ch)
	}
	return result
}

func (b *BCOFEntry) AddData(data *BCOFData) {
	newSubEntry := new(BCOFEntrySubentry)
	bin := new(BCOFEntrySubentry_Binary)
	bin.Binary = data
	newSubEntry.Pl = bin
	b.Data = append(b.Data, newSubEntry)
}

func (b *BCOFEntry) SearchData(s string) *BCOFData {
	for _, subEntry := range b.Data {
		//fmt.Println("subEntry: ", subEntry)
		if d := subEntry.GetLeaf(); d != nil {
			//TODO
		}
		if d := subEntry.GetBinary(); d != nil {
			if d.Signature == s {
				return d
			}
		}
	}

	return nil
}
