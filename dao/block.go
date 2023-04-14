package dao

type Block struct {
	Id          int64  `json:"id"` //primary key
	BlockNumber string `json:"block_number"`
	BlockHash   string `json:"block_hash"`
	ParentHash  string `json:"parent_hash"`
	CreateTime  int64  `json:"create_time"` //the block time for create
	Fork        bool   `json:"fork"`        // is fork branch or not
}
