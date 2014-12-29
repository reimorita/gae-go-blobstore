package pkg

import (
	"appengine"
	"appengine/blobstore"
	"html/template"
	"net/http"
	"fmt"
)

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/redirect", redirect)
}

func index(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	option := blobstore.UploadURLOptions{
		MaxUploadBytes: 1024 * 1024 * 1024,
		StorageBucket:  "tomorier001/subDir1/subDir2/subDir3",
	}

	uploadUrl, err := blobstore.UploadURL(c, "/redirect", &option)
	c.Infof("Up URL: %v", uploadUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var template = template.Must(template.ParseFiles("template/index.html"))
	if err := template.Execute(w, uploadUrl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file := blobs["file"]
	if len(file) == 0 {
		c.Errorf("no file uploaded")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	
	c.Infof("blobKey=> %v", file[0].BlobKey) // blobKey
	c.Infof("Filename=> %v", file[0].Filename) // ファイル名
	c.Infof("ContentType=> %v", file[0].ContentType) // ContentType
	c.Infof("Size=> %v", file[0].Size) // ファイルサイズ(byte)
	
	//http.Redirect(w, r, "/serve/?blobKey="+string(file[0].BlobKey), http.StatusFound)
	//blobstore.Send(w, appengine.BlobKey(string(file[0].BlobKey)))
	blobstore.Delete(c, appengine.BlobKey(string(file[0].BlobKey)))
	fmt.Fprintf(w, "%s", file)

}
