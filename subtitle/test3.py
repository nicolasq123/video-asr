#Imports
import cv2
import os
import numpy as np
# import pytesseract as pt
# from pytesseract import Output

def is_ok(contour):
    x,y,w,h = cv2.boundingRect(contour)
    ar = cv2.contourArea(contour)
    return ar/ (h *w)  > 0.8

def rectratio(contour):
    x,y,w,h = cv2.boundingRect(contour)
    ar = cv2.contourArea(contour)
    return round(ar / (h *w) , 2)

# 50%以下区域
def is_ok2(c, height):
    x,y,w,h = cv2.boundingRect(c)
    return (y+h) > height//2


## 必须是矩形，矩形里面积最大， 且在50%以下高度区域
def judge(contours, img_height):
    for i,c in enumerate(contours):
        x,y,w,h = cv2.boundingRect(c) # 外接最小垂直矩形
        print("---i: ",i, is_ok(c), is_ok2(c, img_height), y+h, img_height//2, rectratio(c), cv2.contourArea(c))

    cons = [c for c in contours if is_ok(c) and is_ok2(c, img_height)]
    if len(cons) == 0:
        cons = contours
    c = max(cons, key=cv2.contourArea)
    x,y,w,h = cv2.boundingRect(c)
    return x,y,w,h

def shuiyin(videofile):
    capture=cv2.VideoCapture(videofile)   #读取本机摄像头
    while True:
        ret,frame=capture.read()   #ret状态  frame：这一针的图像
        logal_image=cv2.imread('tsww.jpg')
        w1, h1, c1 = frame.shape
        w2,h2,c2=logal_image.shape

        roi=frame[w1-w2:w1,h1-h2:h1]
        #灰度化
        gray_logol=cv2.cvtColor(logal_image,cv2.COLOR_BGR2GRAY)
        #黑化
        _,black_logol=cv2.threshold(gray_logol,170,255,cv2.THRESH_BINARY)
        imag_tsw=cv2.bitwise_and(roi,roi,mask=black_logol)
        #白化
        _,white_logal=cv2.threshold(gray_logol,170,255,cv2.THRESH_BINARY_INV)
        imag_tsw1=cv2.bitwise_and(logal_image,logal_image,mask=white_logal)
        imag_tsw2=cv2.add(roi,imag_tsw)
        roii=cv2.add(imag_tsw,imag_tsw1)
        roi[:]=roii
        cv2.imshow('roi',frame)
    
        if cv2.waitKey(30) & 0xFF == 27:  # waitKey延迟作用，130有点卡，一般30或者60
            break
    capture.release()
    cv2.destroyAllWindows()


def clean_img(img, show=True):
    #Main Process
    #Detect subtitles, create mask, inpaint that area
    mask = np.zeros(img.shape, np.uint8)
    recogImg = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    recogImg = cv2.threshold(recogImg, 240, 255, cv2.THRESH_BINARY)[1]
    kernel = cv2.getStructuringElement(cv2.MORPH_RECT, (7,5))
    recogImg = cv2.morphologyEx(recogImg, cv2.MORPH_CLOSE, kernel)
    recogImg = cv2.dilate(recogImg, kernel, iterations=3)
    contours, hierarchy = cv2.findContours(recogImg, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_NONE)
    if len(contours) == 0:
        return img,img,img,img

    # print("--", type(contours[0]))
    # if len(contours) != 0:
    #     c = max(contours, key=cv2.contourArea)
    #     x,y,w,h = cv2.boundingRect(c)
    x,y,w,h = judge(contours, img.shape[0])

    print("area is--------", img.shape[0],img.shape[1], x,y,w,h)

    mask = cv2.cvtColor(mask, cv2.COLOR_BGR2GRAY)
    mask[y:y+h, x:x+w] = recogImg[y:y+h, x:x+w]
    mask = cv2.erode(mask, kernel, iterations=1)
    mask = cv2.GaussianBlur(mask, (3,3), 0)

    cleanedImg = cv2.inpaint(img, mask, 3, cv2.INPAINT_TELEA)

    #Show all stages of process

    # cv2.drawContours(img, contours, -1, (0, 255, 0), 3)
    # cv2.imshow('Contours', img)
    # cv2.imshow("b&w", recogImg)
    # cv2.waitKey(0)
    # cv2.destroyAllWindows()

    if show:
        cv2.imshow("original", img)
        cv2.imshow("b&w", recogImg)
        cv2.imshow("mask", mask)
        cv2.imshow("clean", cleanedImg)

        #Press key to close
        cv2.waitKey(0)
        cv2.destroyAllWindows()
    
    return img, recogImg, mask, cleanedImg


def show_img_area(img, x,y,w,h):
    print("xxx: ", (x, y), (w+x, h+y))
    draw_0 = cv2.rectangle(img, (x, y), (x+w, y+h), (0,0,255), 2)
    cv2.imshow("draw_0", draw_0)#显示画过矩形框的图片
    cv2.waitKey(2000)
    cv2.destroyWindow("draw_0")


def findzimuarea(img, show=False):
    #Main Process
    #Detect subtitles, create mask, inpaint that area
    mask = np.zeros(img.shape, np.uint8)
    recogImg = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY) # 灰度
    grey = recogImg.copy()
    recogImg = cv2.threshold(recogImg, 240, 255, cv2.THRESH_BINARY)[1] #二值
    binary = recogImg.copy()
    kernel = cv2.getStructuringElement(cv2.MORPH_RECT, (7,5)) # 卷积核， np.ones((3,3), np.uint8)
    recogImg = cv2.morphologyEx(recogImg, cv2.MORPH_CLOSE, kernel) # 闭运算，先膨胀后腐蚀
    close  = recogImg.copy()
    recogImg = cv2.dilate(recogImg, kernel, iterations=3) # 膨胀3次
    dilate3 = recogImg.copy()
    contours, hierarchy = cv2.findContours(recogImg, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_NONE)
    if len(contours) == 0:
        return None

    x,y,w,h = judge(contours, img.shape[0])
    if show:
        mask = np.zeros(img.shape, np.uint8)
        mask = cv2.cvtColor(mask, cv2.COLOR_BGR2GRAY)
        mask[y:y+h, x:x+w] = recogImg[y:y+h, x:x+w]
        mask = cv2.erode(mask, kernel, iterations=1)
        mask = cv2.GaussianBlur(mask, (3,3), 0)
        cleanedImg = cv2.inpaint(img, mask, 3, cv2.INPAINT_TELEA)
        cv2.namedWindow('win', cv2.WINDOW_AUTOSIZE)
        cv2.imshow('grey', grey)
        cv2.imshow('binary', binary)
        cv2.imshow('close', close)
        cv2.imshow('dilate3', dilate3)
        cv2.imshow('clean', cleanedImg)
        cv2.waitKey(2000)

    return [x,y,w,h]

def find_max_area(areas):
    if len(areas) == 0:
        return None,None,None,None,

    x0,y0,w0, h0 = areas[0][0],areas[0][1],areas[0][2],areas[0][3]
    x1, y1 = x0+w0, y0+h0

    for x,y,w,h in areas:
        x0 = min(x0, x)
        y0 = min(y0, y)
        x1 = max(x1, x+w)
        y1 = max(y1, y+h)
    
    w0 = x1-x0
    h0 = y1-y0
    return x0,y0,w0,h0


def makevideo(infile, outfile):
    videoinpath  = infile
    videooutpath = outfile
    capture     = cv2.VideoCapture(videoinpath)
    #fourcc      = cv2.VideoWriter_fourcc(*'mp4') # todo 文件格式
    writer      = cv2.VideoWriter(videooutpath, -1, 20.0, (1280,960), False)
    if capture.isOpened():
        while True:
            ret,img_src=capture.read()
            if not ret:break
            _, _, _, img_out = clean_img(img_src, False)    # 自己写函数op_one_img()逐帧处理
            writer.write(img_out)
    else:
        raise Exception('视频打开失败！')
    writer.release()

def video_sample(videof, sample_num=10, show=False):
    capture = cv2.VideoCapture(videof)
    if not capture.isOpened():
        raise Exception('视频打开失败！')

    imgs = []
    framenum = capture.get(cv2.CAP_PROP_FRAME_COUNT)
    start = framenum // (sample_num+1)
    gap = start

    i = 0
    while i < sample_num and start < framenum:
        capture.set(cv2.CAP_PROP_POS_FRAMES, start)  #设置要获取的帧号
        _, img=capture.read()  #read方法返回一个布尔值和一个视频帧。若帧读取成功，则返回True
        #cv2.imshow('b', b)
        #cv2.waitKey(1000)
        i += 1
        start += gap
        imgs.append(img)
        # yeild img
        if show:
            cv2.imshow('sample_img', img)
            cv2.waitKey(20)

    return imgs



from functools import wraps
import time
def time_consume(func):
    @wraps(func)
    def wrapper():
        start = time.time()
        print(f'Time start: {start}')
        func()
        end = time.time()
        print(f'Time end: {end}')
        print(f'Time consumed: {end - start}')
    return wrapper


# 可以用离群因子检测，来去除非字幕的错误区域？
# 去掉最大最小值
def remove_minmaxval(res, num):
    length = len(res)
    if 2*num >= length:
        num = length//2-1

    if num < 0 :
        return res

    def keyfunc(item):
        return item[2]*item[3]
    res.sort(key=keyfunc)
    return res[num:length-num]

# 选取最多的y
# 最小的x和最大的w
def find_areav2(res):
    ys = [r[1] for r in res]
    counts = np.bincount(ys)
    y = np.argmax(counts)
    x = min([r[0] for r in res])
    w = max([r[2]+r[0] for r in res])-x
    hs = [r[3] for r in res]
    counts = np.bincount(hs)
    h = np.argmax(counts)
    return x,y ,w,h


#@time_consume
def run_parse(input_video, show=False, showres=False):
    # msg = ""
    if not input_video:
        raise ValueError("input_video cant be empty")
    rate = 100
    imgs = video_sample(input_video, rate, False)

    height, width = imgs[0].shape[0], imgs[0].shape[1]
    res = [findzimuarea(img, show) for img in imgs]
    rres = [r for r in res if r is not None]
    
    res = remove_minmaxval(rres, rate//10) #
    print("height, width is: ", height, width)
    print("img len is: ", len(imgs))
    print("res len is: ", len(res), res)
    area = find_max_area(res)
    print("max area is :", area, "\n")
    area = find_areav2(res)
    print("max area v2 is :", area, "\n")

    print("({},{},{},{},{},{})".format(area[0], area[1], area[2], area[3], imgs[0].shape[0], imgs[0].shape[1]))
    if showres:
        for img in imgs:
            show_img_area(img, area[0], area[1], area[2], area[3])

    cv2.destroyAllWindows()

    # makevideo(input_video, output_video)
    # print("test-------")

import argparse


if __name__ == "__main__":
    # 方案二 opencv识别并拿到最大的字幕区域， 后面的模糊和加字幕交给ffmpeg处理


    parser = argparse.ArgumentParser(description='Process some params.')
    parser.add_argument('-i','--input', help='input video')
    parser.add_argument('-s','--show', help='input video')
    parser.add_argument('--showres', help='showres video')

    args = parser.parse_args()
    print(args)

    run_parse(input_video=args.input, show=args.show, showres=args.showres)


# from flask import Flask
# app = Flask(__name__)
# @app.route('/')
# def hello_world():
#     return 'Hello World'
    
# if __name__ == '__main__':
#     app.run()

"""
相关知识点
1. https://www.tutorialkart.com/opencv/python/opencv-python-get-image-size/#gsc.tab=0
   - img.shape: 高/宽/； number of channels at index 2.
2. opencv的坐标系 https://blog.csdn.net/oqqENvY12/article/details/71933651

"""