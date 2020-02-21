package controllers

import (
	"bee_blog/models"
	"bee_blog/util"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
)

type AdminController struct {
	baseController
}

func (a *AdminController) Login() {
	if a.Ctx.Request.Method == "POST" {
		username := a.GetString("username")
		password := a.GetString("password")
		user := models.User{Username: username}

		user.Username = username
		//fmt.Println(user)
		//password:=a.GetString("password")
		err := a.o.Read(&user, "username")
		if err != nil {
			a.History("账号不存在", "")
		}
		if util.Md5(password) != strings.Trim(user.Password, " ") {
			a.History("密码错误", "")
		} else {
			a.History("登录成功", "/admin/main")
		}
		a.SetSession("user", user)
	}
	a.TplName = a.controllerName + "/login.html"
}

//主页
func (c *AdminController) Main() {
	c.TplName = c.controllerName + "/main.tpl"
}

//后台系统设置
func (a *AdminController) Config() {
	var result []*models.Config
	_, err := a.o.QueryTable("tb_config").All(&result)
	if err != nil {
		fmt.Println("没有查到数据")
	}
	mp := make(map[string]string)
	opt := make(map[string]*models.Config)
	for _, v := range result {
		mp[v.Name] = v.Value
		opt[v.Name] = v
	}
	if a.Ctx.Request.Method == "POST" {
		keys := []string{"url", "title", "keywords", "description", "email", "start", "qq"}
		for _, key := range keys {
			val := a.GetString(key)
			if _, ok := mp[key]; ok {
				s := opt[key]
				if _, err := a.o.Update(&models.Config{Id: s.Id, Name: s.Name, Value: val}); err != nil {
					continue
				}
			}
		}
		a.History("修改成功", "")
	}
	a.Data["config"] = mp
	a.TplName = a.controllerName + "/config.html"
}

//后台分类管理
func (a *AdminController) Category() {
	var category []*models.Category
	_, err := a.o.QueryTable("tb_category").All(&category)
	if err != nil {
		fmt.Println("没有查到数据")
	}
	a.Data["categorys"] = category
	a.TplName = a.controllerName + "/category.tpl"
}
func (a *AdminController) Categoryadd() {
	id := a.GetString("id")
	intId, _ := strconv.Atoi(id)
	cat := models.Category{Id: intId}
	err := a.o.Read(&cat)
	if err != nil {
		fmt.Println("没有查到数据！")
	}
	a.Data["cate"] = cat
	a.TplName = a.controllerName + "/category_add.tpl"
}
func (a *AdminController) CategorySave() {
	id := a.GetString("id")
	name := a.GetString("name")
	categroy := models.Category{}
	categroy.Name = name
	if id == "0" {
		if _, err := a.o.Insert(&categroy); err != nil {
			a.History("插入数据错误", "")
		} else {
			a.History("插入数据成功！", "/admin/category")
		}
	} else {
		categroy.Id, _ = strconv.Atoi(id)
		if _, err := a.o.Update(&categroy); err != nil {
			a.History("更新数据失败", "")
		} else {
			a.History("更新数据成功", "/admin/category")
		}
	}
}

//后台首页
func (a *AdminController) Index() {
	category := []*models.Category{}
	a.o.QueryTable("tb_category").All(&category)
	a.Data["categorys"] = category
	var (
		page     int
		pagesize int = 8
		offset   int
		list     []*models.Post
		keyword  string
		cateId   int
	)
	keyword = a.GetString("title")
	cateId, _ = a.GetInt("cate_id")
	if page, _ = a.GetInt("page"); page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize
	query := a.o.QueryTable(new(models.Post))
	if keyword != "" {
		query = query.Filter("title__contains", keyword)
	}
	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-is_top", "-created").Limit(pagesize, offset).All(&list)
	}
	a.Data["keyword"] = keyword
	a.Data["count"] = count
	a.Data["list"] = list
	a.Data["cate_id"] = cateId
	a.Data["pagebar"] = util.NewPager(page, int(count), pagesize,
		fmt.Sprintf("/admin/index.html?keyword=%s", keyword), true).ToString()
	a.TplName = a.controllerName + "/list.tpl"
}

//后台删除blog
func (a *AdminController) Delete() {
	id := a.GetString("id")
	stringId, _ := strconv.Atoi(id)
	_, err := a.o.Delete(&models.Post{Id: stringId})
	if err != nil {
		a.History(" 删除失败", "")
	} else {
		a.History("删除成功", "/admin/index.html")
	}
}
func (a *AdminController) Article() {
	category := []*models.Category{}
	a.o.QueryTable("tb_category").All(&category)
	a.Data["categorys"] = category
	id := a.GetString("id")
	intId, _ := strconv.Atoi(id)
	post := models.Post{Id: intId}
	a.o.Read(&post)
	a.Data["post"] = post
	a.TplName = a.controllerName + "/_form.tpl"

}

//上传文件
func (a *AdminController) Upload() {
	f, h, err := a.GetFile("uploadname")
	result := make(map[string]interface{})
	if err == nil {
		defer f.Close()
		exStrArr := strings.Split(h.Filename, ".")
		exStr := strings.ToLower(exStrArr[len(exStrArr)-1])
		if exStr != "jpg" && exStr != "png" && exStr != "gif" {
			result["code"] = 1
			result["message"] = "上传只能.jpg 或者png格式"
		} else {
			uploadDir := "static/upload/"
			if err != nil {
				beego.Error("创建文件失败")
			} else {
				filePath := uploadDir + util.UniqueId() + "." + exStr
				a.SaveToFile("uploadname", filePath)
				result["code"] = 0
				result["message"] = filePath
			}
		}
	} else {
		result["code"] = 2
		result["message"] = "上传异常" + err.Error()
	}
	a.Data["json"] = result
	a.ServeJSON()
}

//保存
func (c *AdminController) Save() {
	post := models.Post{}
	post.UserId = 1
	post.Title = c.Input().Get("title")
	post.Content = c.Input().Get("content")
	post.IsTop, _ = c.GetInt8("is_top")
	post.Types, _ = c.GetInt8("types")
	post.Tags = c.Input().Get("tags")
	post.Url = c.Input().Get("url")
	post.CategoryId, _ = c.GetInt("cate_id")
	post.Info = c.Input().Get("info")
	post.Image = c.Input().Get("image") //此处注意Git上原本的代码有bug get方法里的字段要和 页面name属性字段一模一样
	//fmt.Println("111111111111111" + c.Input().Get("image"))
	post.Created = time.Now()
	post.Updated = time.Now()

	id, _ := c.GetInt("id")
	if id == 0 {
		if _, err := c.o.Insert(&post); err != nil {
			c.History("插入数据错误"+err.Error(), "")
		} else {
			c.History("插入数据成功", "/admin/index.html")
		}
	} else {
		post.Id = id
		if _, err := c.o.Update(&post); err != nil {
			c.History("更新数据出错"+err.Error(), "")
		} else {
			c.History("插入数据成功", "/admin/index.html")
		}
	}
}
