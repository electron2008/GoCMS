// +----------------------------------------------------------------------
// | GoCMS 0.1
// +----------------------------------------------------------------------
// | Copyright (c) 2013-2014 http://www.6574.com.cn All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( http://www.apache.org/licenses/LICENSE-2.0 )
// +----------------------------------------------------------------------
// | Author: zzdboy <zzdboy1616@163.com>
// +----------------------------------------------------------------------

package controllers

//后台首页
import "os"
import "strconv"
import "runtime"
import "strings"
import "github.com/revel/revel"
import "github.com/revel/config"
import "admin/app/models"

type App struct {
	*revel.Controller
}

//首页
func (c App) Index(admin *models.Admin) revel.Result {
	title := "首页--GoCMS管理系统"

	if UserID, ok := c.Session["UserID"]; ok {
		UserID, err := strconv.ParseInt(UserID, 10, 64)
		if err != nil {
			revel.WARN.Println(err)
		}

		admin_info := admin.GetById(UserID)
		if admin_info.Id <= 0 {
			return c.Redirect("/User/Login")
		}

		//控制器
		c.RenderArgs["Controller"] = c.Name
		//动作
		c.RenderArgs["action"] = c.Action
		//模型
		c.RenderArgs["Model"] = c.MethodName

		//导航菜单
		menu := new(models.Menu)
		c.RenderArgs["menus"] = menu.GetAdminMenu(0, admin_info)

		//登陆用户信息
		c.RenderArgs["admin_info"] = admin_info

		//是否锁屏
		if c.Session["lock_screen"] == "" || c.Session["lock_screen"] == "0" {
			c.RenderArgs["lock_screen"] = "0"
		} else {
			c.RenderArgs["lock_screen"] = "1"
		}
	} else {
		return c.Redirect("/User/Login/")
	}

	c.Render(title)
	return c.RenderTemplate("App/Index.html")
}

func (c App) Main(admin *models.Admin) revel.Result {

	title := "首页--GoCMS管理系统"

	if UserID, ok := c.Session["UserID"]; ok {
		UserID, err := strconv.ParseInt(UserID, 10, 64)
		if err != nil {
			revel.WARN.Println(err)
		}

		admin_info := admin.GetById(UserID)

		//判断是否是系统的分隔符
		separator := "/"
		if os.IsPathSeparator('\\') {
			separator = "\\"
		} else {
			separator = "/"
		}

		config_file := (revel.BasePath + "/conf/config.conf")
		config_file = strings.Replace(config_file, "/", separator, -1)
		config_conf, _ := config.ReadDefault(config_file)

		system_info := make(map[string]string)

		//版本
		version, _ := config_conf.String("website", "website.version")
		system_info["version"] = version

		//前台网站地址
		sitedomain, _ := config_conf.String("website", "website.sitedomain")
		system_info["sitedomain"] = sitedomain

		//操作系统
		system_info["os"] = strings.ToUpper(runtime.GOOS + " " + runtime.GOARCH)

		//Go版本
		system_info["go_varsion"] = strings.ToUpper(runtime.Version())

		//MySQL版本
		system_info["mysql_varsion"] = admin.GetMysqlVer()

		//快捷面板
		admin_panel := new(models.Admin_Panel)
		panel_list := admin_panel.GetPanelList(admin_info)

		c.Render(title, admin_info, system_info, panel_list)
	} else {
		c.Render(title)
	}

	return c.RenderTemplate("App/Main.html")
}