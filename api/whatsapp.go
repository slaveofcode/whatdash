package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"whatdash/wa"

	whatsapp "github.com/Rhymen/go-whatsapp"
	"github.com/globalsign/mgo/bson"
)

type WhatsApp struct {
	SessionHandler
}

func (c *WhatsApp) CreateSession(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
	}

	err := decoder.Decode(&params)
	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	wac, err := wa.Connect()

	if err != nil {
		ResponseJSON(w, 200, []byte(`{"status": "error", "err": "`+err.Error()+`"}`))
		return
	}

	stringQr := make(chan string)

	waMgr := wa.Manager{Conn: wac}
	go func(number string, waMgr *wa.Manager, wac *whatsapp.Conn, c *WhatsApp, stringQr chan string) {
		sess, _ := waMgr.LoginAccount(number, stringQr)
		c.Bucket.Save(number, wac, sess)
	}(params.Number, &waMgr, wac, c, stringQr)

	ResponseJSON(w, 200, []byte(`{"status": "created", "qr": "`+<-stringQr+`"}`))

	return
}

func (c *WhatsApp) CheckSession(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
	}

	err := decoder.Decode(&params)
	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	wrapper := c.Bucket.Get(params.Number)

	if wrapper == nil {
		ResponseJSON(w, 400, []byte(`{"status": "unregistered"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"status": "registered"}`))
	return
}

func (c *WhatsApp) Destroy(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
		Force  bool   `json:"force"`
	}

	err := decoder.Decode(&params)
	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	err = c.CloseManager(params.Number, params.Force)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"destroyed": false}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"destroyed": true}`))
	return
}

func (c *WhatsApp) TerminateConnection(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
	}

	err := decoder.Decode(&params)
	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	err = c.TerminateConn(params.Number)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"terminated": false, "error": "`+err.Error()+`"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"terminated": true}`))
	return
}

func (c *WhatsApp) SendText(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Message string `json:"message"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	waMgr, err := c.GetManager(params.From, false)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	err, msgID := waMgr.SendMessage(params.To, params.Message)

	if err != nil {
		ResponseJSON(w, 200, []byte(`{"status": "fail"}`))
		return
	}

	// insert text into DB
	(&wa.MessageKeeper{MgoSession: c.Bucket.MgoSession}).SaveText(&wa.MsgText{
		ID:          bson.NewObjectId(),
		OwnerNumber: params.From,
		WaMsg: wa.WaMsg{
			Type: "text",
			Info: wa.MsgInfo{
				ID:        msgID,
				RemoteJid: params.To,
				FromMe:    true,
				Timestamp: uint64(time.Now().Unix()),
			},
		},
		Text:      params.Message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	ResponseJSON(w, 200, []byte(`{"status": "sent"}`))
	return
}

const maxUploadSize = 5 * 1024 * 1024 // 5 MB
const uploadPath = "upload"

func (c *WhatsApp) SendMedia(w http.ResponseWriter, r *http.Request) {
	// validate filesize
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "Max file upload is 5MB"}`))
		return
	}

	file, _, err := r.FormFile("imageFile")
	number := r.FormValue("number")
	receipent := r.FormValue("receipentJid")
	caption := r.FormValue("caption")

	err = PrepareUploadDir(uploadPath)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "Unable to prepare directory: `+err.Error()+`"}`))
		return
	}

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "Invalid file"}`))
		return
	}

	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "Invalid file"}`))
		return
	}

	filetype := http.DetectContentType(fileBytes)
	var opType string
	var ext string
	switch filetype {
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
		opType = "image"
		break
	case "image/gif":
		ext = ".gif"
		opType = "image"
		break
	case "image/png":
		ext = ".png"
		opType = "image"
		break
	case "video/mp4":
		ext = ".mp4"
		opType = "video"
		break
	case "application/octet-stream", "video/3gpp":
		filetype = "video/3gpp"
		ext = ".3gp"
		opType = "video"
		break
	case "application/pdf":
		ext = ".pdf"
		opType = "document"
		break
	default:
		ResponseJSON(w, 400, []byte(`{"status": "File is not supported to send"}`))
		return
	}

	fileName := wa.CreateFileName(number, ext)
	baseDir, _ := GetBaseDir()
	newPath := filepath.Join(baseDir, uploadPath, fileName)

	newFile, err := os.Create(newPath)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "Cannot store file:`+err.Error()+`"}`))
		return
	}

	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		ResponseJSON(w, 400, []byte(`{"status": "Cannot store file:`+err.Error()+`"}`))
		return
	}

	waMgr, _ := c.GetManager(number, false)
	media, err := os.Open(newPath)
	var msgID string
	if opType == "video" {
		err, msgID = waMgr.SendVideo(receipent, media, filetype, caption)
	} else if opType == "image" {
		err, msgID = waMgr.SendImage(receipent, media, filetype, caption)
	} else if opType == "document" {
		err, msgID = waMgr.SendDocument(receipent, media, filetype, caption)
	}

	fileBytes, errFile := ioutil.ReadFile(newPath)
	if err == nil && errFile == nil {
		err = (&wa.MessageKeeper{MgoSession: c.Bucket.MgoSession}).SaveMedia(&wa.MsgMedia{
			ID:          bson.NewObjectId(),
			OwnerNumber: number,
			WaMsg: wa.WaMsg{
				Type: opType,
				Info: wa.MsgInfo{
					ID:        msgID,
					RemoteJid: receipent,
					FromMe:    true,
					Timestamp: uint64(time.Now().Unix()),
				},
			},
			Type:      filetype,
			Caption:   caption,
			Content:   fileBytes,
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		})
	}

	// store file into gridfs
	// mediaMgo, err := os.Open(newPath)
	// if err == nil {
	// 	(&wa.MessageKeeper{MgoSession: c.Bucket.MgoSession}).StoreFile(fileName, mediaMgo)
	// }

	// removing saved files
	defer func() {
		media.Close()
		// mediaMgo.Close()
		os.Remove(newPath)
	}()

	ResponseJSON(w, 200, []byte(`{"status": "OK"}`))
	return

}

func (c *WhatsApp) DownloadMedia(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
		MsgID  string `json:"messageId"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	mk := wa.MessageKeeper{MgoSession: c.Bucket.MgoSession}
	if !mk.IsMessageOwnerExist(params.Number, params.MsgID) {
		ShowError(w, "Not found")
		return
	}

	res, fileName := mk.GetMessageFile(params.Number, params.MsgID)
	if res == nil {
		ShowError(w, "Not found")
		return
	}

	storedFile, sessMgo, err := mk.GetStoredFile(fileName)
	defer sessMgo.Close()
	if err != nil {
		ShowError(w, "Not found")
		return
	}

	defer storedFile.Close()
	http.ServeContent(w, r, fileName, time.Now(), storedFile)
	return
}

func (c *WhatsApp) LoadContacts(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number       string `json:"number"`
		ReloadSocket bool   `json:"reloadSocket"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	waMgr, err := c.GetManager(params.Number, params.ReloadSocket)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "failed", "reason": "`+err.Error()+`"}`))
		return
	}

	go CollectContacts(&waMgr, c.Bucket.MgoSession)

	ResponseJSON(w, 200, []byte(`{"status": "requested"}`))
	return
}

func (c *WhatsApp) GetContacts(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	contactStorage := wa.ContactStorage{MgoSession: c.Bucket.MgoSession}
	err, contacts := contactStorage.FetchAll(params.Number)

	data, _ := json.Marshal(contacts)

	ResponseJSON(w, 200, data)
	return
}

func (c *WhatsApp) TriggerLoadMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number              string `json:"number"`
		Jid                 string `json:"jid"`
		ReferrenceMessageID string `json:"referenceMessageId"`
		MsgCount            int    `json:"messageCount"`
		ReloadSocket        bool   `json:"reloadSocket"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	waMgr, err := c.GetManager(params.Number, params.ReloadSocket)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "failed", "reason": "`+err.Error()+`"}`))
		return
	}

	err = waMgr.TriggerLoadMessage(params.Jid, params.ReferrenceMessageID, params.MsgCount)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "failed"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"status": "requested"}`))
	return
}

func (c *WhatsApp) TriggerLoadNewMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number        string `json:"number"`
		Jid           string `json:"jid"`
		LastMessageID string `json:"lastMessageId"`
		MsgCount      int    `json:"messageCount"`
		ReloadSocket  bool   `json:"reloadSocket"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	waMgr, err := c.GetManager(params.Number, params.ReloadSocket)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	err = waMgr.TriggerLoadNextMessage(params.Jid, params.LastMessageID, params.MsgCount)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "failed"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"status": "requested"}`))
	return
}

func (c *WhatsApp) TriggerLoadOldMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number       string `json:"number"`
		Jid          string `json:"jid"`
		MessageID    string `json:"messageId"`
		MsgCount     int    `json:"messageCount"`
		ReloadSocket bool   `json:"reloadSocket"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	waMgr, err := c.GetManager(params.Number, params.ReloadSocket)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	err = waMgr.TriggerLoadPrevMessage(params.Jid, params.MessageID, params.MsgCount)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "failed"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"status": "requested"}`))
	return
}
