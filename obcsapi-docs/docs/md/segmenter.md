


```
# full 580k lines mem 200-300Mb
# 200k lines mem 120Mb
# 100k lines mem 70Mb
# 20k lines mem 30Mb
# 10k lines mem 27Mb
# None mem < 20Mb
```


```go
func ChineseSegmenterTest() {
    // 初始化timefinder 对自然语言（中文）提取时间
	var segmenter = timefinder.New("./static/jieba_dict.txt,./static/dictionary-200k.txt")
	var msg string
	var extract []time.Time

	msg = " 6月9日有一场show要去观看"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "后天早上10:35的会议，需要及时参与"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明天下午三点的飞机，提醒我坐车"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "一个小时后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明天早上8:00喊我起床"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明天早上8点喊我起床"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明早十点喊我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明天早上十点喊我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "明天下午三点提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "一天后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "一年后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "一个月后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "一月后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "我要住到大后天"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "我要住到明天"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "下个月到上个月再到这个月"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "我要住到明天下午三点十分"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "帮我预定明天凌晨3点的飞机"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "今天13:00的飞机"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "3月15号的飞机"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "昨天凌晨2点"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)

	msg = "十分钟后提醒我喝水"
	extract = segmenter.TimeExtract(msg)
	fmt.Println(msg)
	fmt.Println(extract)
}
```