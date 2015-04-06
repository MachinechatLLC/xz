package lzbase

// errDist indicates that the distance is out of range.
var errDist = newError("distance out of range")

// ReaderDict represents the dictionary for reading. It is a ring buffer using
// the end field as head for the dictionary.
type ReaderDict struct {
	buffer
	bufferSize int64
}

// NewReaderDict creates a new reader dictionary. The capacity of the ring
// buffer will be the maximum of historySize and bufferSize.
func NewReaderDict(historySize, bufferSize int64) (rd *ReaderDict, err error) {
	if !(1 <= historySize && historySize < MaxDictSize) {
		return nil, newError("historySize out of range")
	}
	if bufferSize <= 0 {
		return nil, newError("bufferSize must be greater than zero")
	}
	capacity := historySize
	if bufferSize > capacity {
		capacity = bufferSize
	}
	rd = &ReaderDict{bufferSize: bufferSize}
	err = initBuffer(&rd.buffer, capacity)
	return
}

// Offset returns the offset of the dictionary head.
func (rd *ReaderDict) Offset() int64 {
	return rd.end
}

// WriteRep writes a repetition with the given distance. While distance is
// given here as int64 the actual limit is the maximum of the int type.
func (rd *ReaderDict) WriteRep(dist int64, n int) (written int, err error) {
	if !(1 <= dist && dist <= int64(rd.Len())) {
		return 0, errDist
	}
	return rd.WriteRepOff(n, rd.end-dist)
}

// Byte returns a byte at the given distance.
func (rd *ReaderDict) Byte(dist int64) byte {
	c, _ := rd.ReadByteAt(rd.end - dist)
	return c
}

// WriterDict is the dictionary used for writing. It includes also a hashtable.
// It is a ring buffer using the cursor offset for the dictionary head. The
// capacity for the buffer is the sum of historySize and bufferSize.
type WriterDict struct {
	buffer
	bufferSize int64
	t4         *hashTable
}

func NewWriterDict(historySize, bufferSize int64) (wd *WriterDict, err error) {
	if !(1 <= historySize && historySize < MaxDictSize) {
		return nil, newError("historySize out of range")
	}
	if bufferSize <= 0 {
		return nil, newError("bufferSize must be greater than zero")
	}
	capacity := historySize + bufferSize
	wd = &WriterDict{bufferSize: bufferSize}
	if err = initBuffer(&wd.buffer, capacity); err != nil {
		return nil, err
	}
	wd.writeLimit = bufferSize
	if wd.t4, err = newHashTable(historySize, 4); err != nil {
		return nil, err
	}
	return wd, nil
}

// HistorySize returns the history length.
func (wd *WriterDict) HistorySize() int64 {
	return wd.Cap() - wd.bufferSize
}

// Returns the byte at the given distance to the dictionary head.
func (wd *WriterDict) Byte(dist int) byte {
	c, _ := wd.ReadByteAt(wd.cursor - int64(dist))
	return c
}

// Offset returns the offset of the head.
func (wd *WriterDict) Offset() int64 {
	return wd.cursor
}

// PeekHead reads bytes from the Head without moving it.
func (wd *WriterDict) PeekHead(p []byte) (n int, err error) {
	return wd.ReadAt(p, wd.cursor)
}

// AdvanceHead moves the head n bytes forward.
func (wd *WriterDict) AdvanceHead(n int) (advanced int, err error) {
	return wd.Copy(wd.t4, n)
}

// Offsets returns all potential offsets for the byte slice. The function
// panics if len(p) is not 4.
func (wd *WriterDict) Offsets(p []byte) []int64 {
	return wd.t4.Offsets(p)
}
