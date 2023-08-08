package entity


type PieceWork struct {
	Index uint
	Hash [20]byte
	Length uint
}

type PieceResult struct {
	Index int 
	Buf []byte
}

