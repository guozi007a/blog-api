package global

/* 项目中需要预先创建的目录和文件 */
import (
	"fmt"
	"os"
)

func PreDirs() {
	os.MkdirAll(fmt.Sprintf("%s/%s", StaticPath, SourceCateImg), 0666)   // 存储上传图片
	os.MkdirAll(fmt.Sprintf("%s/%s", StaticPath, SourceCateAV), 0666)    // 存储上传音频
	os.MkdirAll(fmt.Sprintf("%s/%s", StaticPath, SourceCateOther), 0666) // 存储上传其他文件
	// os.MkdirAll(TempPath, 0666)                                          // 存储临时切片
}
