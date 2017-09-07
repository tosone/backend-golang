package mvc

import (
	"reflect"
	"strings"

	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/memstore"
	"github.com/kataras/iris/mvc/activator"
)

// Controller is the base controller for the high level controllers instances.
//
// This base controller is used as an alternative way of building
// APIs, the controller can register all type of http methods.
//
// Keep note that controllers are bit slow
// because of the reflection use however it's as fast as possible because
// it does preparation before the serve-time handler but still
// remains slower than the low-level handlers
// such as `Handle, Get, Post, Put, Delete, Connect, Head, Trace, Patch`.
//
//
// All fields that are tagged with iris:"persistence"` or binded
// are being persistence and kept the same between the different requests.
//
// An Example Controller can be:
//
// type IndexController struct {
// 	Controller
// }
//
// func (c *IndexController) Get() {
// 	c.Tmpl = "index.html"
// 	c.Data["title"] = "Index page"
// 	c.Data["message"] = "Hello world!"
// }
//
// Usage: app.Controller("/", new(IndexController))
//
//
// Another example with bind:
//
// type UserController struct {
// 	mvc.Controller
//
// 	DB        *DB
// 	CreatedAt time.Time
// }
//
// // Get serves using the User controller when HTTP Method is "GET".
// func (c *UserController) Get() {
// 	c.Tmpl = "user/index.html"
// 	c.Data["title"] = "User Page"
// 	c.Data["username"] = "kataras " + c.Params.Get("userid")
// 	c.Data["connstring"] = c.DB.Connstring
// 	c.Data["uptime"] = time.Now().Sub(c.CreatedAt).Seconds()
// }
//
// Usage: app.Controller("/user/{id:int}", new(UserController), db, time.Now())
// Note: Binded values of context.Handler type are being recognised as middlewares by the router.
//
// Look `core/router/APIBuilder#Controller` method too.
type Controller struct {
	// Name contains the current controller's full name.
	Name string

	// currentRoute is the current request context's route.
	currentRoute context.RouteReadOnly

	// contains the `Name` as different words, all lowercase,
	// without the "Controller" suffix if exists.
	// we need this as field because the activator
	// we will not try to parse these if not needed
	// it's up to the end-developer to call `RelPath()` or `RelTmpl()`
	// which will result to fill them.
	nameAsWords []string

	// relPath the "as assume" relative request path.
	//
	// If UserController and request path is "/user/messages" then it's "/messages"
	// if UserPostController and request path is "/user/post" then it's "/"
	// if UserProfile and request path is "/user/profile/likes" then it's "/likes"
	relPath string

	// request path and its parameters, read-write.
	// Path is the current request path.
	Path string
	// Params are the request path's parameters, i.e
	// for route like "/user/{id}" and request to "/user/42"
	// it contains the "id" = 42.
	Params *context.RequestParams

	// some info read and write,
	// can be already set-ed by previous handlers as well.
	Status int
	Values *memstore.Store

	// relTmpl the "as assume" relative path to the view root folder.
	//
	// If UserController then it's "user/"
	// if UserPostController then it's "user/post/"
	// if UserProfile then it's "user/profile/".
	relTmpl string
	// view read and write,
	// can be already set-ed by previous handlers as well.
	Layout string
	Tmpl   string
	Data   map[string]interface{}

	// give access to the request context itself.
	Ctx context.Context
}

// SetName sets the controller's full name.
// It's called internally.
func (c *Controller) SetName(name string) {
	c.Name = name
}

func (c *Controller) getNameWords() []string {
	if len(c.nameAsWords) == 0 {
		c.nameAsWords = findCtrlWords(c.Name)
	}
	return c.nameAsWords
}

// Route returns the current request controller's context read-only access route.
func (c *Controller) Route() context.RouteReadOnly {
	if c.currentRoute == nil {
		c.currentRoute = c.Ctx.GetCurrentRoute()
	}
	return c.currentRoute
}

const slashStr = "/"

// RelPath tries to return the controller's name
// without the "Controller" prefix, all lowercase
// prefixed with slash and splited by slash appended
// with the rest of the request path.
// For example:
// If UserController and request path is "/user/messages" then it's "/messages"
// if UserPostController and request path is "/user/post" then it's "/"
// if UserProfile and request path is "/user/profile/likes" then it's "/likes"
//
// It's useful for things like path checking and redirect.
func (c *Controller) RelPath() string {
	if c.relPath == "" {
		w := c.getNameWords()
		rel := strings.Join(w, slashStr)

		reqPath := c.Ctx.Path()
		if len(reqPath) == 0 {
			// it never come here
			// but to protect ourselves just return an empty slash.
			return slashStr
		}
		// [1:]to ellimuate the prefixes like "//"
		// request path has always "/"
		rel = strings.Replace(reqPath[1:], rel, "", 1)
		if rel == "" {
			rel = slashStr
		}
		c.relPath = rel
		// this will return any dynamic path after the static one
		// or a a slash "/":
		//
		// reqPath := c.Ctx.Path()
		// if len(reqPath) == 0 {
		// 	// it never come here
		// 	// but to protect ourselves just return an empty slash.
		// 	return slashStr
		// }
		// var routeVParams []string
		// c.Params.Visit(func(key string, value string) {
		// 	routeVParams = append(routeVParams, value)
		// })

		// rel := c.Route().StaticPath()
		// println(rel)
		// // [1:]to ellimuate the prefixes like "//"
		// // request path has always "/"
		// rel = strings.Replace(reqPath, rel[1:], "", 1)
		// println(rel)
		// if rel == "" {
		// 	rel = slashStr
		// }
		// c.relPath = rel
	}

	return c.relPath
}

// RelTmpl tries to return the controller's name
// without the "Controller" prefix, all lowercase
// splited by slash and suffixed by slash.
// For example:
// If UserController then it's "user/"
// if UserPostController then it's "user/post/"
// if UserProfile then it's "user/profile/".
//
// It's useful to locate templates if the controller and views path have aligned names.
func (c *Controller) RelTmpl() string {
	if c.relTmpl == "" {
		c.relTmpl = strings.Join(c.getNameWords(), slashStr) + slashStr
	}
	return c.relTmpl
}

// BeginRequest starts the main controller
// it initialize the Ctx and other fields.
//
// It's called internally.
// End-Developer can ovverride it but it still MUST be called.
func (c *Controller) BeginRequest(ctx context.Context) {
	// path and path params
	c.Path = ctx.Path()
	c.Params = ctx.Params()
	// response status code
	c.Status = ctx.GetStatusCode()
	// share values
	c.Values = ctx.Values()
	// view
	c.Data = make(map[string]interface{}, 0)
	// context itself
	c.Ctx = ctx
}

// EndRequest is the final method which will be executed
// before response sent.
//
// It checks for the fields and calls the necessary context's
// methods to modify the response to the client.
//
// It's called internally.
// End-Developer can ovveride it but still should be called at the end.
func (c *Controller) EndRequest(ctx context.Context) {
	if path := c.Path; path != "" && path != ctx.Path() {
		// then redirect
		ctx.Redirect(path)
		return
	}

	if status := c.Status; status > 0 && status != ctx.GetStatusCode() {
		ctx.StatusCode(status)
	}

	if view := c.Tmpl; view != "" {
		if layout := c.Layout; layout != "" {
			ctx.ViewLayout(layout)
		}
		if data := c.Data; data != nil {
			for k, v := range data {
				ctx.ViewData(k, v)
			}
		}
		ctx.View(view)
	}
}

var ctrlSuffix = reflect.TypeOf(Controller{}).Name()

var _ activator.BaseController = &Controller{}
