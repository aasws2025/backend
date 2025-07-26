package model

type DataEvent struct {
	ID        string `json:"id"`
	AuthorID  string `json:"author_id"`
	Judul     string `json:"judul"`
	Tanggal   string `json:"tanggal"`
	Harga     string `json:"harga"`
	Lokasi    string `json:"lokasi"`
	Deskripsi string `json:"deskripsi"`
	Kategori  string `json:"kategori"`
}
