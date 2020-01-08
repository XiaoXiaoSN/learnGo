/*
* @Author: suifengtec
* @Date:   2018-09-29 15:17:34
* @Last Modified by:   suifengtec
* @Last Modified time: 2018-09-29 18:41:03
**/
/*

go run main.go
go build -o a.exe main.go


支持输入一张图片的URL,返回一个标记了人脸的图片的URL。

仅用于测试,未作充分的容错处理。

使用示例

POST
http://127.0.0.1:6688/img-face-detect

RAW
{

    "ori_img":"http://d.ifengimg.com/mw978_mh598/p1.ifengimg.com/cmpp/2017/11/11/18/ca2af88e-b2a7-4c2e-9ad6-c6e4f1ee407a_size155_w1024_h707.jpg"


}

返回示例

{
    "app_id": 0,
    "face_count": 0,
    "ori_img": "http://d.ifengimg.com/mw978_mh598/p1.ifengimg.com/cmpp/2017/11/11/18/ca2af88e-b2a7-4c2e-9ad6-c6e4f1ee407a_size155_w1024_h707.jpg",
    "size": 84,
    "size_unit": "Kb",
    "marked_img": "http://127.0.0.1:6688/marked/ca2af88e-b2a7-4c2e-9ad6-c6e4f1ee407a_size155_w1024_h707.jpg_c.png"
}

*/
package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

const (
	BASE_URL                    = "http://127.0.0.1"
	PORT                        = "6688"
	IMG_FACE_DETECTION_ENDPOINT = "img-face-detect"
	TEMP_DIR                    = "temp"
	MARKED_DIR                  = "marked"
)

type Item struct {
	AppId     int    `json:"app_id"`
	FaceCount int    `json:"face_count"`
	OriImg    string `json:"ori_img"`
	Size      int    `json:"size"`
	SizeUnit  string `json:"size_unit"`
	MarkedImg string `json:"marked_img"`
}

func isPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func ImgFaceDetec(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalln("Error Read Data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalln("Error Close Request Body", err)
		return
	}

	var item Item

	if err := json.Unmarshal(body, &item); err != nil {
		w.WriteHeader(422)
		log.Println("L89")
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Fatalln("Error Item unmarshalling data", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	//log.Println(item.OriImg)

	err, imgSize := fetchImgByUrl(item.OriImg)
	if err != nil {
		log.Fatalln("Error Item unmarshalling data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	imgSize2 := strconv.Itoa(int(imgSize))
	log.Printf("下载到图片,图片尺寸: %skb", imgSize2)

	unit := "Kb"
	if imgSize > 1024 {
		imgSize /= 1024
		unit = "Mb"
	}

	//imgSizeStr := strconv.Itoa(int(imgSize))
	item.Size = imgSize
	item.SizeUnit = unit
	cwd, _ := os.Getwd()
	imgName0 := getFileName(item.OriImg)
	imgInputPath := path.Join(cwd, TEMP_DIR, imgName0)
	imgOutputPath := path.Join(cwd, MARKED_DIR, imgName0+"_c.png")
	cfPath := path.Join(cwd, "bin", "facefinder")
	if _, err1 := isPathExists(imgInputPath); err1 != nil {
		log.Fatal("图片不存在?")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err2 := isPathExists(cfPath); err2 != nil {
		log.Fatal("识别数据不存在?")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//使用pigo后台运行
	cmd := exec.Command("pigo", "-in", imgInputPath, "-out", imgOutputPath, "-cf", cfPath)
	stdoutStderr, err3 := cmd.CombinedOutput()
	if err3 != nil {
		log.Fatal(err3)
	} else {
		log.Printf("%s\n", stdoutStderr)
	}
	log.Println(imgOutputPath)

	//标记后的图片,以base64形式返回,此时,如果输入的图片很大,输出会很长,所以改为使用URL输出
	// item.MarkedImg = convertImgToBase64(imgOutputPath)
	//标记人脸后的图片,以URL形式返回
	item.MarkedImg = getMarkedImgUrl(imgOutputPath)
	//删除输入的图片
	deleteImgFile(imgInputPath)
	data, er := json.Marshal(item)
	if er != nil {
		w.WriteHeader(422)
		log.Println("LINE97")
		log.Println(er)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	//w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return
}

//删除输入的图片
func deleteImgFile(fileFullPath string) bool {

	if _, err := isPathExists(fileFullPath); err != nil {
		return true
	}
	var err = os.Remove(fileFullPath)
	if err != nil {
		return false
	}
	return true
}

//标记后的图片链接
func getMarkedImgUrl(imgFullPath string) string {

	segments := strings.Split(imgFullPath, "/")
	fileName := segments[len(segments)-1]
	//log.Println("getMarkedImgUrl")
	//log.Println(fileName)
	imgUrl := ""
	if "80" == PORT {
		imgUrl = BASE_URL + "/" + MARKED_DIR + "" + fileName
	} else {
		imgUrl = BASE_URL + ":" + PORT + "/" + MARKED_DIR + "/" + fileName
	}
	return imgUrl
}

//标记后的图片转换为base64,缺点是：图片如果很大,输出就很长,没必要
func convertImgToBase64(imgFullPath string) string {
	imgFile, err := os.Open(imgFullPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer imgFile.Close()
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	return base64.StdEncoding.EncodeToString(buf)

}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getFileName(URL string) string {
	fileUrl, err := url.Parse(URL)
	checkError(err)
	path := fileUrl.Path
	segments := strings.Split(path, "/")

	fileName := segments[len(segments)-1]
	//log.Println(fileName)
	return fileName
}

func fetchImgByUrl(URL string) (error, int) {

	response, e := http.Get(URL)
	if e != nil {
		return e, 0
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.New("L178"), 0
	}
	cwd, _ := os.Getwd()
	imgInputPath := path.Join(cwd, TEMP_DIR, getFileName(URL))
	file, err := os.Create(imgInputPath)
	if err != nil {
		return err, 0
	}

	b, er := io.Copy(file, response.Body)
	if er != nil {
		return er, 0
	}
	file.Close()

	imgSize := b / 1024
	imgSizeKb := int(imgSize)

	return nil, imgSizeKb

}

//静态文件所在目录
func myRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/" + TEMP_DIR + "/").Handler(http.StripPrefix("/"+TEMP_DIR+"/", http.FileServer(http.Dir("./"+TEMP_DIR+"/"))))
	router.PathPrefix("/" + MARKED_DIR + "/").Handler(http.StripPrefix("/"+MARKED_DIR+"/", http.FileServer(http.Dir("./"+MARKED_DIR+"/"))))

	return router
}

func createDirs() {
	//如果存放处理后的图片的目录不存在,先创建它
	cwd, _ := os.Getwd()
	tempDir := path.Join(cwd, TEMP_DIR)
	if b1, _ := isPathExists(tempDir); b1 == false {
		os.Mkdir(tempDir, 0777)
	}
	markedImgDir := path.Join(cwd, MARKED_DIR)
	if b2, _ := isPathExists(markedImgDir); b2 == false {
		os.Mkdir(markedImgDir, 0777)
	}
}

func main() {
	createDirs()
	router := myRouter()
	router.HandleFunc("/"+IMG_FACE_DETECTION_ENDPOINT, ImgFaceDetec).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+PORT, router))

}
