## video-asr
Video Auto Speech Recognition, Auto translate. Auto remove hard subtitles, Add new subtitles.
视频自动语音识别，自动翻译成目标语言，自动模糊/删硬字幕或者软字幕，自动添加新字幕

#### 安装
1. bash ./install && bash ./subtitle/install ## todo

#### 功能
1. 支持腾讯语音识别
2. 支持google翻译
3. 调用ffmpeg处理音视频
4. 子项目subtitle进行字幕定位
5. 各个模块大多都设计成interface，易接入其他语音识别/翻译工具


#### 背景
调研视频处理工具可自动翻译语音对白并覆盖字幕

###### 功能拆解成以下几个部分
1. 需要自动识别语音， 两个方向， 一是语音识别，二是OCR处理字幕文件。 如果视频没有字幕第二点就走不通
2. 字幕文件翻译成目标语言
3. 覆盖以前的旧字幕
4. 添加翻译后的新字幕

#### 方案

###### 配置
1. CPU: AMD A8-5600K APU with Radeon(tm) HD Graphics
2. 内存: 8G
3. GPU: Radeon HD 7560D

###### 方案一
1. 腾讯asr识别拿到srt file
2. opencv识别字幕区域并直接抹掉字幕 -- 存在误差较大的问题
3. ffmpeg添加字幕

###### 方案二
1. 腾讯asr识别拿到srt file
2. opencv采样识别字幕区域，融合判断最大概率的区域
3. ffmpeg抹掉旧字母，同时添加字幕
4. 需要测试性能问题 todo

###### 方案三
1. 在方案二的基础上去逐个处理每帧，每帧拿到的字幕区域与预处理的结果搞个融合对比就行了
2. 这种方案是基于字幕的基本信息，因为人脑在看视频的时候，处理第N帧信息的时候，N-1帧的信息会作为输入，对第N帧进行综合的判断。简单讲就是下一帧字幕的所在区域的中心点与上一帧的字幕所在区域中心点差别不会太大，颜色字体亦是。

###### 语音识别方案对比
1. google
    - https://cloud.google.com/speech-to-text
    - pricing: https://cloud.google.com/speech-to-text/pricing ，第一次300$免费
    - 支持中文， 未找到翻译功能, 需另外开发 性能未知, 质量未知
    - 开发周期未知
2. software
    - https://learn.microsoft.com/en-us/azure/cognitive-services/speech-service/captioning-concepts?pivots=programming-language-python#caption-and-speech-synchronization
    - https://azure.microsoft.com/en-us/pricing/details/cognitive-services/speech-services/
    - 支持中->印尼语
3. alipyun
   - https://ai.aliyun.com/nls
   - 40小时 100.00/年起,
   - 支持中文， 未找到翻译功能, 需另外开发
4. 腾讯
   - https://cloud.tencent.com/product/asr
   - 60小时， 72/年
   - 支持中文， 未找到翻译功能, 需另外开发
5. baidu
   - http://www.baiemai.com/product/asr.htm

###### 翻译识别方案对比
1. https://cloud.tencent.com/product/tmt 
2. https://cloud.google.com/translate


###### 处理工具
1. ffmpeg
2. opencv


#### 相关知识/坑点
1. srt格式 https://docs.fileformat.com/zh/video/srt/
2. 硬字幕与软字幕的区别。 硬字幕是直接嵌入在视频frame上的，必须要进行自动化的图像处理。question/23231910
3. vlc打开视频的时候会自动加载同目录下，同名srt file
4. google translate 请求时候需要分页
5. 语音识别并加字幕的过程 https://learn.microsoft.com/zh-cn/azure/cognitive-services/speech-service/captioning-concepts?pivots=programming-language-python#caption-and-speech-synchronization


#### 存在的问题
1. 如果视频的字幕是直接硬编码进视频里的，清除/模糊原字幕会很麻烦
   - 解决方案（手动）： https://www.zhihu.com/question/23231910

2. 字幕校准问题/字幕。（时间序列准确度问题）
3. 字幕翻译精准问题
4. 原字幕清除前定位不精准的问题


#### 参考
1. 微软，语音识别同步的过程 https://learn.microsoft.com/zh-cn/azure/cognitive-services/speech-service/captioning-concepts?pivots=programming-language-cli#caption-and-speech-synchronization
2. 《数字图像处理》冈萨雷斯
3. subtitle location https://ieeexplore.ieee.org/abstract/document/4126288
4. An intelligent subtitle detection model for locating television commercials https://ieeexplore.ieee.org/abstract/document/4126288/
