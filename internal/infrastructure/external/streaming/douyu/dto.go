package douyu

type BetardResponse struct {
	Room struct {
		Nickname    string `json:"nickname"`
		OwnerAvatar string `json:"owner_avatar"`
		Status      string `json:"status"`
		ShowStatus  int    `json:"show_status"`
		ShowDetails string `json:"show_details"`
		RoomName    string `json:"room_name"`
		RoomPic     string `json:"room_pic"`
		CoverSrc    string `json:"coverSrc"`
		ShowTime    int64  `json:"show_time"` // Unix seconds 1762482715
		Avatar      struct {
			Big    string `json:"big"`
			Middle string `json:"middle"`
			Small  string `json:"small"`
		} `json:"avatar"`
		CateName      string `json:"cate_name"`
		SecondLvlName string `json:"second_lvl_name"`
		RoomBizAll    struct {
			Hot string `json:"hot"`
		} `json:"room_biz_all"`
	} `json:"room"`
	Column struct {
		CateId   string `json:"cate_id"`
		CateName string `json:"cate_name"`
	} `json:"column"`
}
