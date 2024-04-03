package main

func toMp4WithoutSizeCompressionArgs(inputFilepath string, destFilepath string) (args []string) {
	return []string{
		"-i", inputFilepath,
		"-vcodec", "h264",
		"-preset", "medium",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "64k",
		"-ac", "2",
		destFilepath,
	}
}
