package structs

type UpdateDistributionKeysPayload struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}
