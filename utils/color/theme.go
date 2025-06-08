package color

var Theme = struct {
	Error   Color
	Warning Color
	Info    Color
	Debug   Color
	Pub     Color
	Sub     Color
}{
	Error:   NewColor(255, 0, 0),     // 红色
	Warning: NewColor(255, 165, 0),   // 橙色
	Info:    NewColor(0, 191, 255),   // 深天蓝
	Debug:   NewColor(128, 128, 128), // 灰色
	Pub:     NewColor(255, 255, 0),   // 亮黄色
	Sub:     NewColor(0, 255, 255),   // 亮蓝色
}
