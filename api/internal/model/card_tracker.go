package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/qxsugar/pkg/kit"
)

// CardCounts 各牌面剩余数量映射，键为牌面（BJ/SJ/A/K/.../2），值为剩余张数。
// 实现 driver.Valuer / sql.Scanner 以 JSON 形式落库。
type CardCounts map[string]int

func (c CardCounts) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *CardCounts) Scan(src any) error {
	if src == nil {
		*c = CardCounts{}
		return nil
	}
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, c)
	case string:
		return json.Unmarshal([]byte(v), c)
	default:
		return fmt.Errorf("cannot scan %T into CardCounts", src)
	}
}

// CardTracker 扑克记牌器状态（每个用户一份）。
type CardTracker struct {
	Id        int64         `json:"id" gorm:"column:id"`
	UserId    int64         `json:"user_id" gorm:"column:user_id"`
	DeckCount int           `json:"deck_count" gorm:"column:deck_count"`
	Counts    CardCounts    `json:"counts" gorm:"column:counts"`
	CreatedAt kit.TimeStamp `json:"created_at" gorm:"column:created_at"`
	UpdatedAt kit.TimeStamp `json:"updated_at" gorm:"column:updated_at"`
}

func (CardTracker) TableName() string { return "bill_card_trackers" }

// CardRanks 牌面顺序（展示用），大小王在前，A→2。
var CardRanks = []string{"BJ", "SJ", "A", "K", "Q", "J", "10", "9", "8", "7", "6", "5", "4", "3", "2"}

// DefaultCounts 按牌副数生成默认剩余数量：大小王每副各 1 张，其余每副各 4 张。
func DefaultCounts(deckCount int) CardCounts {
	if deckCount < 1 {
		deckCount = 1
	}
	c := CardCounts{}
	for _, r := range CardRanks {
		if r == "BJ" || r == "SJ" {
			c[r] = 1 * deckCount
		} else {
			c[r] = 4 * deckCount
		}
	}
	return c
}
