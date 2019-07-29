package adapter

// 经测试，并没有快多少，可用可不用
// 10W个json数组解析
// 循环解析：243ms
// 转换解析：205ms
// 将二维数组转换为一维数组（json格式，方便Json解析）
func TwoToOne(contents [][]byte) []byte {
	if len(contents) == 0 {
		return nil
	}
	var size = len(contents)*(len(contents[0])+1) + 10
	var one = make([]byte, 0, size)
	one = append(one, '[')
	for _, b := range contents {
		one = append(one, b...)
		one = append(one, ',')
	}
	one[len(one)-1] = ']'
	return one
}
