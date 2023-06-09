package api

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	beegoAPI "StockMe/controllers"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"

	"github.com/beego/beego/v2/client/orm"
	"github.com/golang-jwt/jwt"
)

type API struct {
	beegoAPI.API
}

var (
	ImagePath, _    = beego.AppConfig.String("imagePath")
	PathCallData, _ = beego.AppConfig.String("pathCall")

	BackEndDataPath, _    = beego.AppConfig.String("backEndDataPath")
	GomoAccessKey, _      = beego.AppConfig.String("gomoAccessKey")
	FirebaseKey, _        = beego.AppConfig.String("firebaseKey")
	Ip2location, _        = beego.AppConfig.String("ip2location")
	DropDatePlusURL, _    = beego.AppConfig.String("dropDatePlusURL")
	UserDropDate, _       = beego.AppConfig.String("userDropDate")
	PasswordDropDate, _   = beego.AppConfig.String("passwordDropDate")
	AwsAccessKeyId, _     = beego.AppConfig.String("awsAccessKeyId")
	AwsSecretAccessKey, _ = beego.AppConfig.String("awsSecretAccessKey")
	AwsS3Bucket, _        = beego.AppConfig.String("awsS3Bucket")
	AwsS3Region, _        = beego.AppConfig.String("awsS3Region")
)

const (
	AccessKey                 = "VDeoa0934lkfaZ30ds"
	Success                   = "Success"
	BadRequest                = "Bad Request!"
	SomethingWentWrong        = "Something went wrong!"
	DuplicatedRequest         = "Duplicated Request!"
	RateLimits                = "Too many Request!"
	RequestTimedOut           = "Request timed out"
	InvalidArgument           = "Invalid argument"
	InvalidEmail              = "Invalid E-Mail Address"
	NotFound                  = "NOT FOUND"
	AccountOrPasswordNotFound = "Invalid email account or password."
	UserIsNotAdminLevel       = "User is not admin level."
	MaxFileSize               = 500000000

	// header
	HeaderAuthToken = "X-Auth-Token"
	HeaderMobileId  = "Mobile-Id"
	HeaderToken     = "Token"

	// user api error messages
	Unauthorized       = "Unauthorized"
	PermissionDeniedEn = "You don’t have permission to access"
	PermissionDenied   = "ขออภัย คุณไม่มีสิทธิ์เข้าถึงข้อมูลนี้"
	FileSizeLarge      = "ขออภัยขนาดไฟล์ใหญ่เกิน 500 MB"

	MissionNotFound           = "Mission Not Found"
	MissionTaskNotFound       = "Mission Task Not Found"
	MissionTaskResultNotFound = "Mission Task Result Not Found"

	//Role
	Admin = "admin"
	User  = "user"
)

func (api *API) GetAccessCredentials() string {
	return api.Ctx.Input.Header(HeaderAuthToken)
}

func (api *API) getHeaderAuthToken() string {
	return api.Ctx.Input.Header(HeaderAuthToken)
}

func (api *API) GetHeaderMobileId() string {
	return strings.TrimSpace(api.Ctx.Input.Header(HeaderMobileId))
}
func (api *API) GetHeaderToken() string {
	return api.Ctx.Input.Header(HeaderToken)
}

func NewAccessToken() *string {
	//token := jwt.New(jwt.SigningMethodHS512)
	//Set some claims
	//Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).UnixNano(),
		Issuer:    "ind-platform.com",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	//token.Claims["time"] = time.Now().Unix()
	//token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(AccessKey))
	if err != nil {
		logs.Error("Error cannot gen new token | ", err.Error())
		return nil
	}
	return &tokenString
}

type CustomClaim struct {
	Channel       string `json:"channel"`
	TransactionId string `json:"transactionId"`
	MobileId      string `json:"mobileId"`
	//Action        string `json:"action"`
	//Schema        Schema `json:"schema"`
	//Iat           int64  `json:"iat"`
	jwt.StandardClaims
}
type Schema struct {
	Home     string `json:"home"`
	Toggle   string `json:"toggle"`
	Register string `json:"register"`
}

func NewGomoToken(channel, mobileId string) (res, tranId string) {
	//token := jwt.New(jwt.SigningMethodHS512)
	//Set some claims
	//Create the Claims
	t := time.Now()
	y := strconv.Itoa(t.Year())
	m := strconv.Itoa(int(t.Month()))
	//m := t.Month()
	d := strconv.Itoa(t.Day())
	h := strconv.Itoa(t.Hour())
	min := strconv.Itoa(t.Minute())
	sec := strconv.Itoa(t.Second())
	//NU 2021 11 20 11 49 40 54321
	// NU 2018 08 20 16 02 01 54321
	tranId = channel + y + m + d + h + min + sec + "54321"
	now := (t).Unix()
	expire := (t.Add(time.Minute * 1)).Unix()
	//{
	//	"mobileId": "611e2492d3f89f56750de208",
	//	"action": "mission",
	//	"schema": {
	//		"home": "gomogame://home",
	//		"toggle": "gomogame://toggle-speed",
	//		"register": "gomogame://register-point"
	//},
	//	"iat": 1637138695,
	//	"exp": 1637142295
	//}

	//{
	//	"alg": "HS256",
	//	"typ": "JWT"
	//}
	claims := &CustomClaim{
		channel,
		tranId,
		mobileId,
		//"mission",
		//Schema {
		//	Home:     "gomogame://home",
		//	Toggle:   "gomogame://toggle-speed",
		//	Register: "gomogame://register-point",
		//},
		//expire,
		jwt.StandardClaims{
			IssuedAt:  now,
			ExpiresAt: expire,
			//Issuer:    "gomo-mission-app",
		},
	}
	//var buf bytes.Buffer
	//jsonBytes , err := json.Marshal(claims)
	//if err != nil {
	//	logs.Error("err marshall json:",err)
	//}
	//err = json.Compact(&buf,jsonBytes)
	//if err != nil {
	//	logs.Debug("err:",err)
	//}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	logs.Debug("expire:", expire)
	logs.Debug("token:", token)
	logs.Debug("token.Raw:", token.Raw)
	//encodedKey := base64.RawURLEncoding.EncodeToString([]byte(GomoAccessKey))
	//tokenString, err := token.SignedString([]byte(encodedKey))
	tokenString, err := token.SignedString([]byte(GomoAccessKey))
	logs.Debug("tokenString:", tokenString)
	if err != nil {
		logs.Error("Error cannot gen new token | ", err.Error())
		return "", ""
	}
	return tokenString, tranId
}

func (api *API) ResponseJSONWithCode(results interface{}, statusCode int, code int64, msg string) {
	if results == nil {
		results = struct{}{}
	}

	response := &beegoAPI.ResponseObjectWithCode{
		Code:           code,
		Message:        msg,
		ResponseObject: results,
	}

	api.Data["json"] = response
	api.Ctx.ResponseWriter.Header().Set("access-control-allow-headers", "Origin,Accept,Content-Length,Content-Type,X-Atmosphere-tracking-id,X-Atmosphere-Framework,X-Cache-Dat,Cache-Control,X-Requested-With,X-Auth-Token,Authorization,Access-Control-Allow-Origin")
	api.Ctx.ResponseWriter.Header().Set("access-control-expose-headers", "access-control-allow-origin;Content-Type")
	api.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	api.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "PUT,PATCH,GET,POST,DELETE,OPTIONS")
	api.Ctx.ResponseWriter.WriteHeader(statusCode)

	api.Ctx.Output.SetStatus(statusCode)
	api.ServeJSON()
	return
}

func (api *API) ResponseJSON(results interface{}, code int, msg string) {
	if results == nil {
		results = struct{}{}
	}
	response := &beegoAPI.ResponseObject{
		Code:           code,
		Message:        msg,
		ResponseObject: results,
	}

	api.Data["json"] = response
	api.Ctx.ResponseWriter.Header().Set("access-control-allow-headers", "Origin,Accept,Content-Length,Content-Type,X-Atmosphere-tracking-id,X-Atmosphere-Framework,X-Cache-Dat,Cache-Control,X-Requested-With,X-Auth-Token,Authorization,Access-Control-Allow-Origin")
	api.Ctx.ResponseWriter.Header().Set("access-control-expose-headers", "access-control-allow-origin;Content-Type")
	api.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	api.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "PUT,PATCH,GET,POST,DELETE,OPTIONS")
	api.Ctx.ResponseWriter.WriteHeader(code)

	api.Ctx.Output.SetStatus(code)
	api.ServeJSON()
	return
}

//	func (api *API) GetUser() *models.User {
//		token := api.getHeaderAuthToken()
//		user := models.GetUserByToken(token)
//		api.Data["user"] = user
//		return user
//	}
func (api *API) ValidatePassword(s string) error {
	if len(s) < 8 || len(s) > 64 {
		return errors.New("Password's length have to be between 8 - 64 characters.")
	}
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if regex.MatchString(s) == false {
		return errors.New("Password can contains only english characters and numbers.")
	}
	return nil
}

func (api *API) CheckTransaction(err error, to orm.TxOrmer) error {
	if err != nil {
		logs.Error("execute transaction's sql fail, rollback.", err)
		err = to.Rollback()
		if err != nil {
			logs.Error("roll back transaction failed", err)
		}
	} else {
		err = to.Commit()
		if err != nil {
			logs.Error("commit transaction failed.", err)
		}
	}
	return err
}

func (api *API) CheckAndCreatesDirectory(filePath string, creates bool) (err error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logs.Debug("no dir")
		if creates {
			err = CreatesDirectory(filePath)
		}
	}
	return
}

func CreatesDirectory(filePath string) (err error) {
	if err = os.Mkdir(filePath, os.ModePerm); err != nil {
		logs.Debug(err)
	}
	return
}
func (api *API) TrimString(message string) string {
	return strings.TrimSpace(message)
}

// ToDateTime
// support dd/mm/yyyy , dd/mm/yyyy hh:mm dd/mm/yyyy hh:mm:ss
// support yyyy-mm-dd hh:mm:ss
func (api *API) ToDateTime(s string) (t time.Time) {
	if s == "" {
		return time.Time{}
	}
	ss := strings.Split(s, " ")
	sections := strings.Split(s, ":")
	// hh:mm
	if len(ss) == 2 && len(ss[1]) == 5 {
		logs.Debug("case 0")
		s += ":00"
	} else if sections == nil || len(sections) == 1 {
		logs.Debug("case 1")
		s += " 00:00:00"
	}
	isDash := (strings.Replace(s, "-", "x", -1)) != s
	var err error
	if isDash {
		t, err = time.Parse("2006-01-02 15:04:05", s)
	} else {
		t, err = time.Parse("02/01/2006 15:04:05", s)
	}
	if err != nil {
		logs.Error("err parse date", err)
		return time.Time{}
	}
	return t
}

func (api *API) FormatDateUnix(t time.Time) string {
	if (t == time.Time{}) {
		return ""
	}
	return strconv.Itoa(int(t.UnixNano()))
}

func (api *API) FormatDateNoTime(t time.Time) string {
	if (t == time.Time{}) {
		return ""
	}
	return t.Format("02/01/2006")
}

func (api *API) DateToBuddhistEra(time time.Time) time.Time {
	date := time.AddDate(543, 0, 0)
	return date
}

func (api *API) ZipFolder(zippath string, path []string) (string, string, error) {
	logs.Debug("ZipPediaFolder :")
	rand := randString(16, "")
	zipFileName := rand + "-" + time.Now().Format("02012006150405")
	tmpDir := ImagePath + "/tmp/" + zipFileName + "/"
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		err = os.Mkdir(tmpDir, os.ModePerm)
		if err != nil {
			return "", "", err
		}
	}
	logs.Debug("zipFileName :", zipFileName)
	contentDir := moveFileToParentFolder(tmpDir, zipFileName, ImagePath, path)
	logs.Debug("> : contentDir :", contentDir)
	src, err := os.Open(contentDir)
	if err != nil {
		logs.Error("ERROR 'ZipPediaFolder' 1 :", err)
	}
	dst, err := os.Open(contentDir)
	if err != nil {
		logs.Error("ERROR 'ZipPediaFolder' 2 :", err)
	}
	io.Copy(dst, src)
	err = zipFilePath(contentDir, tmpDir, zipFileName+".zip")
	if err != nil {
		logs.Error("ERROR 'ZipPediaFolder' 3 :", err)
	}
	logs.Debug("return dir :", tmpDir+zipFileName+".zip")
	logs.Debug("tmpDir :", tmpDir)
	logs.Debug("zipFileName :", zipFileName)
	return tmpDir + zipFileName + ".zip", contentDir, nil
}

func randString(n int64, prefix string) (name string) {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	name = prefix + string(b)
	return
}

func moveFileToParentFolder(tmpDir, name string, base string, paths []string) string {
	rand := randString(16, "")
	dir := base + rand + "-" + time.Now().Format("02012006150405")
	if name != "" {
		dir = base + "/tmp/" + name + "/"
	}
	logs.Debug("dir:", dir)
	dir = tmpDir
	for _, path := range paths {
		ss := strings.Split(path, "/")
		fname := ss[len(ss)-1]
		src, err := os.Open(path)
		if err != nil {
			logs.Error("ERROR 'moveFileToParentFolder' 1 :", err)
			return ""
		}
		logs.Debug("dst path:", dir+fname)
		dst, err := os.Create(dir + fname)
		if err != nil {
			logs.Error("ERROR 'moveFileToParentFolder' 2 :", err)
			return ""
		}
		wb, err := io.Copy(dst, src)
		logs.Error("wb 'copy'  :", wb)
		if err != nil {
			logs.Error("ERROR 'copy'  :", err)

		}
	}
	return dir
}

func zipFilePath(path string, zipToPath string, zipName string) error {
	logs.Debug("ZipFilePath :", path, " -- ", zipName)
	zipFile, err := os.Create(zipToPath + zipName)
	if err != nil {
		logs.Debug("Create err:", err)
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileToZip, err := os.Open(path)

	if err != nil {
		logs.Debug("Open err:", err)
		return err
	}

	defer fileToZip.Close()
	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		logs.Debug("fileToZip err:", err)
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		logs.Debug("FileInfoHeader err:", err)
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	//TODO : Change file name to original file name here
	header.Name = zipName

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate
	//writer, err := zipWriter.CreateHeader(header)
	//if err != nil {
	//	logs.Debug("CreateHeader err:", err)
	//	return err
	//}

	addFileToZip(zipWriter, path, "", zipName)

	//_, err = io.Copy(writer, fileToZip)
	//logs.Debug("Copy err:", err)
	return err
}

func addFileToZip(w *zip.Writer, path string, baseInZip string, zipName string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		if file.Name() != zipName {
			fmt.Println(path + file.Name())
			if !file.IsDir() {
				dat, err := ioutil.ReadFile(path + file.Name())
				if err != nil {
					fmt.Println(err)
				}
				// Add some files to the archive.
				f, err := w.Create(baseInZip + file.Name())
				if err != nil {
					fmt.Println(err)
				}
				_, err = f.Write(dat)
				if err != nil {
					fmt.Println(err)
				}
			} else if file.IsDir() {
				// Recurse
				newBase := path + file.Name() + "/"
				fmt.Println("Recursing and Adding SubDir: " + file.Name())
				fmt.Println("Recursing and Adding SubDir: " + newBase)
				addFileToZip(w, newBase, baseInZip+file.Name()+"/", zipName)
			}
		}
	}
}

func (api *API) RandomString(digit int64) (res string) {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, digit)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	res = string(b)
	return
}

func (api *API) CreatesTrackDirectory(filePath string) (suc bool, err error) {

	err = CheckAndCreatesDirectory(filePath, true)
	if err != nil {
		suc = true
	}
	return
}

func CheckAndCreatesDirectory(filePath string, creates bool) (err error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logs.Debug("no dir")
		if creates {
			err = CreatesDirectory(filePath)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (api *API) Int64ToString(a int64) string {
	return strconv.FormatInt(a, 10)
}
