package controllers

import (
	"github.com/astaxie/beego"
	"mime/multipart"
	"time"
	"path"
	"github.com/astaxie/beego/orm"
	"testNews/models"
	"errors"
	"math"
	"bytes"
	"encoding/gob"
	"github.com/gomodule/redigo/redis"
)

type ArticleController struct {
	beego.Controller
}

func UploadImg(this *ArticleController, header *multipart.FileHeader, uploadname string) (filePath string, err error) {
	// 判断是否为图片
	suffix := path.Ext(header.Filename)
	if suffix != ".jpg" && suffix != ".png" && suffix != ".jpeg" {
		err = errors.New("上传图片格式不正确，请重新上传")
		return
	}

	if header.Size > 51200 {
		err = errors.New("上传图片过大,不能大于50KB")
		return
	}

	timeStamp := time.Now().Format("2006-01-02-15-04-15")
	this.SaveToFile("uploadname", "./static/img/"+timeStamp+suffix)
	filePath = "/static/img/" + timeStamp + suffix
	err = nil
	return
}

func (this *ArticleController) ShowIndex() {
	o := orm.NewOrm()
	qs := o.QueryTable("Article")
	var articles []models.Article

	pageSize := 2
	pageIndex, _ := this.GetInt("pageIndex")
	typeName := this.GetString("select")
	// 判断是否为空pageIndex
	var count int64
	start := pageSize * (pageIndex - 1)
	if typeName == "" {
		count, _ = qs.RelatedSel("ArticleType").Count()
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
	} else {
		count, _ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	}

	// 计算页数
	pageCount := math.Ceil(float64(count) / float64(pageSize))
	this.Data["pageIndex"] = pageIndex
	this.Data["typeName"] = typeName
	this.Data["count"] = count
	this.Data["pageCount"] = int(pageCount)
	this.Data["articles"] = articles

	// 获取所有tpyeName
	var articleTypes []models.ArticleType
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		beego.Error("redis 数据库连接失败")
	}

	resp, err := redis.Bytes(conn.Do("get", "articleTypes"))
	dec := gob.NewDecoder(bytes.NewReader(resp))
	dec.Decode(&articleTypes)

	if len(articleTypes) == 0 {
		o.QueryTable("ArticleType").All(&articleTypes)
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		enc.Encode(&articleTypes)
		conn.Do("set", "articleTypes", buffer.Bytes())
		beego.Info("从mysql中获取数据")
	}

	this.Data["articleTypes"] = articleTypes
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

func (this *ArticleController) ShowAdd() {
	o := orm.NewOrm()
	var articleType []models.ArticleType
	o.QueryTable("ArticleType").All(&articleType)
	this.Data["articleType"] = articleType
	this.Layout = "layout.html"
	this.TplName = "add.html"
}

func (this *ArticleController) HandleAdd() {
	// 获取信息
	uploadname := "uploadname"
	articleTitle := this.GetString("articleName")
	articleContent := this.GetString("content")
	file, header, err := this.GetFile(uploadname)
	if err != nil {
		this.Data["errmsg"] = "Get pic failed"
		this.TplName = "add.html"
		return
	}
	defer file.Close()

	// 判断 标题 和 内容 不能为空
	if articleTitle == "" || articleContent == "" {
		this.Data["errmsg"] = "Title, content is null"
		this.TplName = "add.html"
		return
	}

	// 上传文件，获取路径
	filePath, err := UploadImg(this, header, uploadname)
	if err != nil {
		beego.Info("上传图片失败", err)
		this.Redirect("/article/add", 302)
		return
	}
	beego.Info(filePath, articleTitle, articleContent)

	// 插入数据库
	o := orm.NewOrm()
	var article models.Article
	article.Img = filePath
	article.Title = articleTitle
	article.Content = articleContent
	article.Time = time.Now()

	// 获取标题
	typeName := this.GetString("select")
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType, "typeName")
	article.ArticleType = &articleType

	//article.ArticleType = acticleType
	id, err := o.Insert(&article)
	if err != nil {
		this.Data["errmsg"] = "Failed upload file"
		beego.Info(err)
		this.TplName = "add.html"
		return
	}
	beego.Info("UploadImg primad key = ", id)
	this.Redirect("/article/index", 302)
}

func (this *ArticleController) ShowAddType() {
	o := orm.NewOrm()
	var artilceTypes []models.ArticleType
	o.QueryTable("ArticleType").OrderBy("id").All(&artilceTypes)
	this.Data["artilceTypes"] = artilceTypes
	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

func (this *ArticleController) HandleShowType() {
	typeName := this.GetString("typeName")

	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	id, err := o.Insert(&articleType)
	if err != nil {
		beego.Info("Add type err")
		this.Redirect("/article/addType", 302)
		return
	}
	beego.Info(id)
	this.Redirect("/article/addType", 302)
}

func (this *ArticleController) DeleteType() {
	Id, _ := this.GetInt("Id")
	o := orm.NewOrm()
	var typeId models.ArticleType
	typeId.Id = Id
	o.Delete(&typeId)
	this.Redirect("/article/addType", 302)
}

func (this *ArticleController) ShowContent() {
	articleId, _ := this.GetInt("articleId")

	// 读取文章表
	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId
	o.Read(&article, "Id")

	// 添加访问次数
	article.Count += 1
	o.Update(&article)

	/*	m2m := o.QueryM2M(&article, "User")
		var user models.User
		userName := this.GetSession("userName")
		user.Name = userName.(string)
		o.Read(&user, "UserName")
		m2m.Add(user)*/

	this.Data["article"] = article
	this.TplName = "content.html"
}

func (this *ArticleController) ShowUpdate() {
	articleId, _ := this.GetInt("articleId")
	//beego.Info("articlesId = ", articleId)
	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId
	o.Read(&article)
	//beego.Info("showUpdate read", err)
	this.Data["article"] = article
	this.TplName = "update.html"
}

func (this *ArticleController) HandleUpdate() {
	Title := this.GetString("articleName")
	Content := this.GetString("content")
	file, header, err := this.GetFile("uploadname")
	articleId, _ := this.GetInt("articleId")
	beego.Info("get file", err)
	defer file.Close()

	filePath, err := UploadImg(this, header, "uploadname")
	if err != nil {
		beego.Info("上传图片失败", err)
		this.Redirect("/article/index", 302)
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId
	err = o.Read(&article)
	if err != nil {
		beego.Info("记录无法找到", err)
		this.Redirect("/article/index", 302)
		return
	}
	article.Title = Title
	article.Content = Content
	article.Img = filePath
	id, err := o.Update(&article)
	beego.Info(id, err)
	beego.Info(filePath, Title, Content)
	this.Redirect("/article/index", 302)
}

func (this *ArticleController) HandleDelete() {
	articleId, err := this.GetInt("articleId")
	if err != nil {
		beego.Error("删除文章连接错误")
		this.Redirect("/article/index", 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = articleId
	id, err := o.Delete(&article)
	if err != nil {
		beego.Error("删除失败")
		this.Redirect("/article/index", 302)
		return
	}
	beego.Info("id", id)
	this.Redirect("/article/index", 302)
}
