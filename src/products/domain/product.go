package domain

type Product struct {
	Id       int32   `json:"id"`
	Nombre   string  `json:"nombre"`
	Precio   float32 `json:"precio"`
	Cantidad float32 `json:"cantidad"`
}
