    package main

    import (
            "flag"
            "fmt"
            "io"
            "io/ioutil"
            "log"
            "net/http"
            "os"
            "path"
            "strconv"
            "strings"
            "time"
            "github.com/gin-gonic/gin"
    )

    var defaultPath string

    var baseURL string
    var upload_path string

    func Logger(req *http.Request, statusCode int) {
            const layout = "[ 2/Jan/2006 15:04:05 ]"
            fmt.Println(baseURL + " --- " + time.Now().Format(layout) + " " + req.Method + "  " + strconv.Itoa(statusCode) + "  " + req.URL.Path)
    }

    func Handler(w http.ResponseWriter, req *http.Request) {


            filename := defaultPath + req.URL.Path[1:]
            fmt.Println("=+=== file name: ",filename)
            if last := len(filename) - 1; last >= 0 && filename[last] == '/' && len(filename) != 1 {
                    filename = filename[:last]
            }

            if req.Method == "POST" {
                    file, head, err := req.FormFile("file")
                    if err != nil {
                            fmt.Println(err)
                            return
                    }

                    defer file.Close()


                    //创建文件
                    fW, err := os.Create(upload_path + head.Filename)
                    if err != nil {
                            fmt.Println("文件创建失败")
                            return
                    }
                    defer fW.Close()

                    _, err = io.Copy(fW, file)
                    if err != nil {
                            fmt.Println("文件保存失败")
                            return
                    }
            }

            // Empty request (Root)
            if filename == "" {
                    filename = "./"
            }

            file, err := os.Stat(filename)

            // 404 if file doesn't exist
            if os.IsNotExist(err) {
                    _, err = io.WriteString(w, "404 Not Found")
                    Logger(req, http.StatusNotFound)
                    return
            }

            // Serve directory
            if file.IsDir() {

                    slashCheck := ""

                    files, err := ioutil.ReadDir(filename)
                    // Catch the Error in reading from directory
                    if err != nil {
                            http.Redirect(w, req, "", http.StatusInternalServerError)
                            Logger(req, http.StatusInternalServerError)
                    }
                    // Checking for Root Directory
                    if filename != "./" {
                            if filename[len(filename)-1] != '/' {
                                    slashCheck = "/"
                            }
                    }

                    fmt.Println("=+=== slashCheck: ",slashCheck)
                    fmt.Println("=+=== req.URL.Path[0:]: ",req.URL.Path[0:])

                    responseString := "<html><body> <h3> Directory Listing for " + req.URL.Path[1:] + "/ </h3> <br/> <hr> <ul>"
                    for _, f := range files {
                            if f.Name()[0] != '.' {
                                    if f.IsDir() {
                                            responseString += "<li><a href=\"" + req.URL.Path[0:] + slashCheck + f.Name() + "\">" + f.Name() + "/" + "</a></li>"
                                    } else {
                                            responseString += "<li><a href=\"" + req.URL.Path[0:] + slashCheck + f.Name() + "\">" + f.Name() + "</a></li>"
                                    }
                            }
                    }

                    //Ending the list
                    responseString += "</ul><br/><hr/>"

                    p := req.URL.Path

                    // Display link to parent directory
                    if len(p) > 1 {
                            base := path.Base(p)

                            slice := len(p) - len(base) - 1

                            url := "/"

                            if slice > 1 {
                                    url = req.URL.Path[:slice]
                                    url = strings.TrimRight(url, "/") // Remove extra / at the end
                            }

                            responseString += "<br/><a href=\"" + url + "\">Parent directory</a>"
                    }

                    uploadStr := "<form action='#' method=\"post\" enctype=\"multipart/form-data\"> <label> </label><input type=\"file\" name='file'  /><br/><br/> <label><input type=\"submit\" value=\"上传文件\"/></label> </form>"

                    responseString = responseString + uploadStr + "</body></html>"
                    //fmt.Println("      =+=== responseString: ",responseString)
                    _, err = io.WriteString(w, responseString)
                    if err != nil {
                            // panic(err)
                            http.Redirect(w, req, "", http.StatusInternalServerError)
                            Logger(req, http.StatusInternalServerError)
                    } else {
                            Logger(req, http.StatusOK)
                    }

                    upload_path = "./" + req.URL.Path[0:] + "/"

                    return
            }

            // File exists and is no directory; Serve the file

            b, err := ioutil.ReadFile(filename)
            if err != nil {
                    http.Redirect(w, req, "", http.StatusInternalServerError)
                    Logger(req, http.StatusInternalServerError)
                    return
            }

            str := string(b)
            extension := path.Ext(filename)

            if extension == ".css" {
                    w.Header().Set("Content-Type", "text/css; charset=utf-8")
            } else if extension == ".js" {
                    w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
            }
            _, err = io.WriteString(w, str)
            if err != nil {
                    // panic(err)
                    http.Redirect(w, req, "", http.StatusInternalServerError)
            } else {
                    Logger(req, http.StatusOK)
            }

    }


    func FileServerGet(c *gin.Context) {
           
            filename := defaultPath + c.Request.URL.Path[1:]
            fmt.Println("=+=== file name: ",filename)
            if last := len(filename) - 1; last >= 0 && filename[last] == '/' && len(filename) != 1 {
                    filename = filename[:last]
            }

            // Empty request (Root)
            if filename == "" {
                    filename = "./"
            }

            file, err := os.Stat(filename)

            // 404 if file doesn't exist
            if os.IsNotExist(err) {
                    _, err = io.WriteString(c.Writer, "404 Not Found")
                    //Logger(c.Request, http.StatusNotFound)
                    return
            }

    // Serve directory
            if file.IsDir() {

                    slashCheck := ""

                    files, err := ioutil.ReadDir(filename)
                    // Catch the Error in reading from directory
                    if err != nil {
                            http.Redirect(c.Writer, c.Request, "", http.StatusInternalServerError)
                            //Logger(c.Request, http.StatusInternalServerError)
                    }
                    // Checking for Root Directory
                    if filename != "./" {
                            if filename[len(filename)-1] != '/' {
                                    slashCheck = "/"
                            }
                    }

                    fmt.Println("=+=== slashCheck: ",slashCheck)
                    fmt.Println("=+=== c.Request.URL.Path[0:]: ",c.Request.URL.Path[0:])

                    responseString := "<html><body> <h3> Directory Listing for " + c.Request.URL.Path[1:] + "/ </h3> <br/> <hr> <ul>"
                    for _, f := range files {
                            if f.Name()[0] != '.' {
                                    if f.IsDir() {
                                            responseString += "<li><a href=\"" + c.Request.URL.Path[0:] + slashCheck + f.Name() + "\">" + f.Name() + "/" + "</a></li>"
                                    } else {
                                            responseString += "<li><a href=\"" + c.Request.URL.Path[0:] + slashCheck + f.Name() + "\">" + f.Name() + "</a></li>"
                                    }
                            }
                    }

                    //Ending the list
                    responseString += "</ul><br/><hr/>"

                    p := c.Request.URL.Path

                    // Display link to parent directory
                    if len(p) > 1 {
                            base := path.Base(p)

                            slice := len(p) - len(base) - 1

                            url := "/"

                            if slice > 1 {
                                    url = c.Request.URL.Path[:slice]
                                    url = strings.TrimRight(url, "/") // Remove extra / at the end
                            }

                            responseString += "<br/><a href=\"" + url + "\">Parent directory</a>"
                    }

                    uploadStr := "<form action='#' method=\"post\" enctype=\"multipart/form-data\"> <label> </label><input type=\"file\" name='file'  /><br/><br/> <label><input type=\"submit\" value=\"上传文件\"/></label> </form>"

                    responseString = responseString + uploadStr + "</body></html>"
                    //fmt.Println("      =+=== responseString: ",responseString)
                    _, err = io.WriteString(c.Writer, responseString)
                    if err != nil {
                            // panic(err)
                            http.Redirect(c.Writer, c.Request, "", http.StatusInternalServerError)
                            //Logger(c.Request, http.StatusInternalServerError)
                    } else {
                            //Logger(c.Request, http.StatusOK)
                    }

                    upload_path = "./" + c.Request.URL.Path[0:] + "/"

                    return
            }

            // File exists and is no directory; Serve the file

            b, err := ioutil.ReadFile(filename)
            if err != nil {
                    http.Redirect(c.Writer, c.Request, "", http.StatusInternalServerError)
                    //Logger(c.Request, http.StatusInternalServerError)
                    return
            }

            str := string(b)
            extension := path.Ext(filename)

            if extension == ".css" {
                    c.Writer.Header().Set("Content-Type", "text/css; charset=utf-8")
            } else if extension == ".js" {
                    c.Writer.Header().Set("Content-Type", "application/javascript; charset=utf-8")
            }
            _, err = io.WriteString(c.Writer, str)
            if err != nil {
                    // panic(err)
                    http.Redirect(c.Writer, c.Request, "", http.StatusInternalServerError)
            } else {
                    //Logger(c.Request, http.StatusOK)
            }
    }

    func FileServerPost(c *gin.Context) {
            file, head, err := c.Request.FormFile ("file")
            if err != nil {
                    fmt.Println(err)
                    return
            }

            defer file.Close()


            //创建文件
            fW, err := os.Create(upload_path + head.Filename)
            if err != nil {
                    fmt.Println("文件创建失败")
                    return
            }
            defer fW.Close()

            _, err = io.Copy(fW, file)
            if err != nil {
                    fmt.Println("文件保存失败")
                    return
            }
    }

    func main() {

            router := gin.Default()
            router.GET("/",FileServerGet)
            router.POST("/",FileServerPost)


            defaultPortPtr := flag.String("p", "", "Port Number")
            defaultPathPtr := flag.String("d", "", "Root Directory")
            flag.Parse()

            portNum := "8080"

            // Handling the command line flags

            // Directory
            if *defaultPathPtr != "" {
                    defaultPath = "./" + *defaultPathPtr + "/"
            } else {
                    defaultPath = ""
            }
            // Port Number
            if *defaultPortPtr != "" {
                    portNum = *defaultPortPtr
            } else {
                    portNum = "8080"
            }

            baseURL = "http://localhost:" + portNum

            fmt.Println("Serving on ", baseURL, " subdirectory ", defaultPath)
    /*
            http.HandleFunc("/", Handler)
            err := http.ListenAndServe(":"+portNum, nil)
            if err != nil {
                    log.Fatal("ListenAndServe: ", err)
            }
            */
            err := http.ListenAndServe(":" + portNum, router)
            if err != nil {
                    log.Fatal("ListenAndServe: ", err)
            }
    }
