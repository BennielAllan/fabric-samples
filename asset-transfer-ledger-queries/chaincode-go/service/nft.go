package service

// Milestones 里程碑表
type Milestones struct {
	DocType     string  `json:"docType"`
	ID          uint64  `json:"id"`
	Threshold   float64 `json:"threshold"` // 触发条件
	Description string  `json:"description"`
}

// NFTs NFT荣誉证书表
type NFTs struct {
	DocType         string `json:"docType"`
	ID              string `json:"id"`
	DonorID         string `json:"donor_id"`
	MetadataURL     string `json:"metadata_url"`
	ContractAddress string `json:"contract_address"`
	MintedAt        uint   `json:"minted_at"`
}
