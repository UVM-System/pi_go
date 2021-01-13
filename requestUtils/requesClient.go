
/*
http 请求的封装，用于voguemannner系统内部微服务的访问
 */
package requestutils

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"
)


/*
发送一个图像字节流的用法:
bytesH,_:=gocv.IMEncode(".jpg",img)


	imageByte := requestutils.ImageBytes{
		FileName:"1.jpg",
		FieldName:"image",
		Content:bytesH,
		ContentType:"image/jpeg",
	}
	formData := requestutils.PostBytes{
		FileMap:[]requestutils.ImageBytes{imageByte},
		FieldMap:nil,
	}
	contentType,bodyBuffer,_ :=requestutils.CreateFormDataFromBytes(formData)
	response, err := http.Post("http://127.0.0.1:5000/test", contentType, bodyBuffer)
 */

type ImageBytes struct {
	FieldName string
	FileName string
	Content []byte
	ContentType string
}
type PostBytes struct {
	FileMap []ImageBytes
	FieldMap map[string]string
}
func CreateFormDataFromBytes(formdata PostBytes) (string,*bytes.Buffer,error)  {
	bodyBuf := &bytes.Buffer{}

	bodyWriter :=multipart.NewWriter(bodyBuf)

	//写入文件
	for _,imageBytes :=range formdata.FileMap{
		bufferWriter,_ :=bodyWriter.CreatePart(mimeHeader(imageBytes.FieldName,imageBytes.FileName,imageBytes.ContentType))

		bufferWriter.Write(imageBytes.Content)
	}
	for key, val := range formdata.FieldMap {
		_ = bodyWriter.WriteField(key, val)
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	return contentType,bodyBuf,nil
}

////TODU:将返回json 格式化为
//func Postformdata2json(contenttype,url string,bodybuffer *bytes.Buffer,result interface{}) error  {
//	response, err := http.Post(url, contenttype, bodybuffer)
//	if err != nil {
//		log.Println("error to post"+url)
//		return err
//	}
//	defer response.Body.Close()
//	resp_body, err := ioutil.ReadAll(response.Body)
//	log.Println(string(resp_body))
//	json.Unmarshal(resp_body,result)
//	return nil
//}

//封装文件内容

func mimeHeader(fieldname, filename,contenttype string) textproto.MIMEHeader {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contenttype)
	return h
}
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}