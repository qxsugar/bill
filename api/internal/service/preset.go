package service

import "math/rand"

// PresetAvatars 预设头像列表（前端「修改信息-头像选择」展示同一份）。
// 使用 emoji 占位，前端可直接渲染；后续可替换为 CDN 图片 URL。
var PresetAvatars = []string{
	"🐶", "🐱", "🐭", "🐹", "🐰", "🦊",
	"🐻", "🐼", "🐨", "🐯", "🦁", "🐮",
	"🐷", "🐸", "🐵", "🐔", "🐧", "🐦",
}

// nicknameWords 随机昵称词库，组合成「形容词+名词」。
var nicknameAdjectives = []string{"快乐", "悠闲", "机智", "勇敢", "神秘", "可爱", "淡定", "幸运"}
var nicknameNouns = []string{"玩家", "牌神", "高手", "选手", "大侠", "队友", "新星", "黑马"}

// randomNickname 生成随机昵称，例如「快乐牌神」。
func randomNickname() string {
	adj := nicknameAdjectives[rand.Intn(len(nicknameAdjectives))]
	noun := nicknameNouns[rand.Intn(len(nicknameNouns))]
	return adj + noun
}

// randomAvatar 从预设头像中随机取一个。
func randomAvatar() string {
	return PresetAvatars[rand.Intn(len(PresetAvatars))]
}
