package play_2399

type CardAwardInfo struct {
	PrizeId   int    `json:"prizeId"`
	PrizeName string `json:"prizeName"`
}

var (
	ACTIVITY_BRANCH       = "play_2399"
	DAY_CHARGE_LIMIT      = 10
	SIGN_AWARD_POOL_LIMIT = 200000 // 元
	SIGN_STABLE_AWARD     = "感恩节勋章*1天"
	SIGN_RANDOM_AWARDS    = []string{"", "狂欢盲盒*1", "中级盲盒*1", "恋爱盲盒*1"}
	ROUND_LIMIT           = 5000                        // 元
	CARDS_LIMIT           = []int{100, 500, 2000, 5000} // 元
	CARDS_AWARD1          = []CardAwardInfo{
		{PrizeId: 939, PrizeName: "感恩节勋章1天"},
		{PrizeId: 100004, PrizeName: "SVIP1天"},
		{PrizeId: 40002542, PrizeName: "小黄鸭"},
	}
	CARDS_AWARD2 = []CardAwardInfo{
		{PrizeId: 40002836, PrizeName: "恋爱三重奏"},
		{PrizeId: 40003911, PrizeName: "可爱鬼"},
		{PrizeId: 40001340, PrizeName: "直升机"},
	}
	CARDS_AWARD3 = []CardAwardInfo{
		{PrizeId: 10000049, PrizeName: "强吻"},
		{PrizeId: 40003461, PrizeName: "独角兽"},
		{PrizeId: 40000053, PrizeName: "流星雨"},
	}
	CARDS_AWARD4 = []CardAwardInfo{
		{PrizeId: 40001340, PrizeName: "直升机"},
		{PrizeId: 40003066, PrizeName: "旋转木马"},
		{PrizeId: 40003271, PrizeName: "爱满星河"},
	}
	CARDS_AWARDS = [][]CardAwardInfo{CARDS_AWARD1, CARDS_AWARD2, CARDS_AWARD3, CARDS_AWARD4} // 翻卡奖励
)
