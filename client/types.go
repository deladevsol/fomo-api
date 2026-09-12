package client

import "encoding/json"

type User struct {
	ID             string  `json:"id"`
	Handle         string  `json:"handle"`
	Name           *string `json:"name"`
	Bio            *string `json:"bio"`
	Banner         *string `json:"banner"`
	ProfilePicture *string `json:"profilePicture"`
	Twitter        *string `json:"twitter"`
	SolanaAddress  *string `json:"solanaAddress"`
	EVMAddress     *string `json:"evmAddress"`
}
type Thesis struct {
	ID                      string            `json:"id"`
	TokenAddress            *string           `json:"tokenAddress"`
	TokenNetwork            *string           `json:"tokenNetwork"`
	TokenSymbol             *string           `json:"tokenSymbol"`
	AuthorID                *string           `json:"authorId"`
	AuthorHandle            *string           `json:"authorHandle"`
	AuthorName              *string           `json:"authorName"`
	AuthorIsDev             *bool             `json:"authorIsDev"`
	Text                    *string           `json:"thesis"`
	Segments                []json.RawMessage `json:"segments"`
	LikeCount               *int64            `json:"likeCount"`
	HoldingsUSD             *float64          `json:"holdingsUsd"`
	AuthorTradeUSD          *float64          `json:"authorTradeUsd"`
	PnL                     *float64          `json:"pnl"`
	RealizedPnLUSD          *float64          `json:"realizedPnlUsd"`
	UnrealizedPnLUSD        *float64          `json:"unrealizedPnlUsd"`
	PercentageRealizedPnL   *float64          `json:"percentageRealizedPnl"`
	PercentageUnrealizedPnL *float64          `json:"percentageUnrealizedPnl"`
	TokenAmount             *float64          `json:"tokenAmount"`
	ClosedAt                *int64            `json:"closedAt"`
	FomoCreatedAt           *int64            `json:"fomoCreatedAt"`
	UpdatedAt               *int64            `json:"updatedAt"`
	AuthorAvatarURL         *string           `json:"authorAvatarUrl,omitempty"`
	TokenImageURL           *string           `json:"tokenImageUrl,omitempty"`
}
type ThesisPage struct {
	UserID       string   `json:"userId,omitempty"`
	TokenAddress *string  `json:"tokenAddress,omitempty"`
	TokenNetwork *string  `json:"tokenNetwork,omitempty"`
	Symbol       *string  `json:"symbol,omitempty"`
	Count        int      `json:"count"`
	HasMore      bool     `json:"hasMore"`
	NextBefore   *string  `json:"nextBefore"`
	UpdatedAt    *int64   `json:"updatedAt"`
	Items        []Thesis `json:"items"`
}
