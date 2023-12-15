package tables

/*
GiftType GiftTypeID 一个礼物只能有一个对应的type
普通礼物 1 表示价值[0, 1000)的礼物
高级礼物 2 表示价值[1000, 50000)的礼物
豪华礼物 3 表示价值[50000, ∞）的礼物
促销礼物 4 表示带(促)的礼物
幸运礼物 5 表示价值[0, 1000)且送出后可以获得若干个该礼物的礼物，俗称"卡奖礼物"
图片礼物 6 不是礼物，是一些图片，比如VIP、SVIP等常用的奖励图
*/

/*
ExtendsTypes 拓展分类，有些类型的礼物，在后来会增加为别的类型，该结构是用于类型的额外拓展
经典礼物 1001 表示类型为经典的礼物
玩法礼物 1002 表示类型为玩法的礼物，注意这里是分类，不是tag
特权礼物 1003 表示类型为特权的礼物
*/

/*
GiftTag GiftTagID 一个礼物可以有多个不同的tag
盲盒礼物 1 表示用于开盲盒的礼物
挖宝礼物 2 凡是可以用于中奖的礼物，都计入挖宝礼物。比如铲子、老鼠、盲盒、池子、福袋、抽签、水晶、扭蛋等。
折扣礼物 3 表示带(x折)的礼物
活动礼物 4 表示用于指定活动的礼物。如圣诞帽、圣诞树、纪念册等。
公益礼物 5 表示用于做公益的礼物
等级礼物 6 表示达到指定等级才可以赠送的礼物
粉丝礼物 7 表示关注了本房间房主的用户才可以赠送的礼物
*/

var InsertGifts = []KKGifts{
	{
		GiftID:   80,
		GiftName: "VIP",
	},
}
