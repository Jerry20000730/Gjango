# Gjango
A Self-designed Web Framework using Golang

# Components of a Web Framework
## 1. Router
A _route_ is a mapping between the **URL** and the **backend logic code**. For example, when we want to access a resource, we might write something like this:
```
GET http://localhost:8000/home
```

In this logic, we might want to get the html file of the homepage, so that the browser can render and display the frontend page to the user. Here:

- HTTP is a standardized protocol, other protocol should also be supported, such as HTTPS
- GET is a HTTP request method. Other methods include `POST`, `PUT`, `DELETE`
- If we change the URL, says we might want to access another resource

```
GET http://localhost:8000/about
```

Then these two URL should return two different results, one might be the html/css/js for the homepage, and the other might be the html/css/js for the about page, depending on the detailed logic of the backend code.

In this part, the web framework should support:

- [x] In-URL request parameters: when we do routing pattern matching, we usually include the request parameter in the request URL.
- [] Form submission: the web framework should also support form submission. In html, `<form>` is a basic component for submitting or changing the info that needed to be processed in the backend.

## Router Grouping
In most cases, we want the router to be able to register a group of name into one router group (e.g., in Django, you can create the URLPattern under another URLPattern, e.g., if your first URLPattern has a pattern named "/user", and in that user module, you create another URLPattern named "/getUser", and "/createUser").

In this part, the web framework should support:

- [x] 1. Register a router group name

## Request Method
In http, not only should we support the URL pattern matching, but also we must support different URL request methods under the same URL pattern. For example, even in "user/userInfo" pattern, we should support both `GET` abd `POST`. 

In this part, the web framework should support:

- [x] Support `GET`, `POST`, `DELETE`, `PUT`

## Middleware
Middleware is software that provide common services and capabilities to applications and help developers and operators build and deploy applications more efficiently. Middleware makes it easier to implement communication and input/output, so that the developer can focus on specific purpose of their application, while the middleware application will not change the original way of how the application is formed.

The goal of the middleware should be:

1. The middleware should not be coupled in the user code
2. The middleware should be independent but can get the context and make the decision upon the context

In this part, the web framework should support

- [x] Support registering a general middleware (pre-middleware and post-middleware)
- [x] Support registering a router-group-specific middleware

## Page Rendering
During the response, the interface should support returning

- HTML (support template)
- JSON
- XML

### HTML
In order to render the html, we must clarify on several element in the HTTP response so that the browser can identify the content and render it in the frontend.

1. `content-type = text/html; charset-utf-8`
2. Template elements
3. Data/context that need to display on the HTML

#### Usage
```go
func main() {
    engine := web.NewEngine()
    g := engine.Router.NewGroup("user")
    g.Get("/get/html", func(ctx *context.Context) {
        err := ctx.HTML(http.StatusOK, "<h1>HTML Template</h1><p>This is a template for html and test if the html is successfully returned and rendered</p>")
        if err != nil {
            log.Println(err)
        }   
    })
}
```

### JSON & XML
In order to support JSON & XML, we only need to make sure that the content-type is changed
1. `content-type = application/json; charset=utf-8` for json and `content-type = application/xml charset=utf-8`

#### Usage
```go
func main() {
    engine := web.NewEngine()
    g := engine.Router.NewGroup("user")
    g.Get("/json", func(ctx *context.Context) {
        user := &User{
            Name: "jerry",
        }
	err := ctx.JSON(http.StatusOK, user)
        if err != nil {
            log.Println(err)
        }
    })
}
```

### String
In order to only display some string on the web, we allow user to pass a specific string with extra arguments specified (like printf), so that the argument can be passed and displayed dynamically on the website.
1. `content-type = text/plain; charset=utf-8` for string
2. support dynamic parameter passing

#### Usage
```go
func main() {
    engine := web.NewEngine()
    g := engine.Router.NewGroup("user")
    g.Get("/string", func(ctx *context.Context) {
	err := ctx.String(http.StatusOK, "Test %s gjango web framework, int can also be passed: %d", "self-designed", 1)
    })
}
```

## Parameter Processing
Parameters are essential in passing the information, enabling the transfer of data between different parts of a web application. This can happen in various contexts, such as between the client and server or within different components of the application. 

Here are some examples this web framework can support:

### URL parameters
also known as query strings, are commonly used to pass information in a URL. They follow the ? character in a URL and consist of key-value pairs, for example:

```text
http://xxx.com/user/add?id=1&age=20
```

#### map parameters
In this case, the parameters are passed in the form of a map, so that the user can easily access the parameters by the key.
This approach allows for a structured and easily parseable method of passing complex and nested data through URLs. 
It is particularly useful for filtering, searching, or specifying multiple attributes of a resource in a single request.

```text
http://xxx.com/user/queryMap?user[id]=1&user[age]=20
```

##### Usage
```go
// http://xxx.com/user/queryMap?user[id]=1&user[age]=29
g.Get("/queryMap", func(ctx *context.Context) {
    m, ok := ctx.QueryMap("user")
    ctx.JSON(http.StatusOK, m)
})
```

### Post form parameters
Post form parameters are used to send data to the server in the body of the HTTP request. They are commonly used in forms and are sent as key-value pairs.
Web application use form to perform various tasks, such as user authentication, data submission, and more.

#### Usage
```go
g.Post("/formPost", func(ctx *context.Context) {
    m, _ := ctx.GetPostFormArray("name")
    ctx.JSON(http.StatusOK, m)
})
```

### Post form file
Post form file is used to upload files to the server. It is commonly used in web applications to allow users to upload images, videos, documents, and other types of files.

#### Usage
You can define the folder for uploaded file in the `SaveUploadedFile` function
```go
g.Post("/file", func(ctx *context.Context) {
    file, err := ctx.FormFile("file")
    if err != nil {
        log.Println(err)
    }
    err = ctx.SaveUploadedFile(file, "./upload/"+file.Filename)
    if err != nil {
        log.Println(err)
    }
})
```

The framework also allows for multiple file uploads at once. The user can specify the file in the body of the request, and the framework will save the file to the specified folder.
```go
g.Post("/multiFile", func(ctx *context.Context) {
    m, _ := ctx.GetPostFormMap("user")
    // here, you specify the key of the file 
    // inside the request body 
    // e.g., if the key is "file", then here you should enter "file"
	headers := ctx.FormFiles("file")
    for _, file := range headers {
        ctx.SaveUploadedFile(file, "./upload/"+file.Filename)
    }
    ctx.JSON(http.StatusOK, m)
})
```

### JSON Parameters
JSON parameters are used to send data to the server in the body of the HTTP request. They are commonly used in APIs and are sent as JSON objects.
When sending the request:
- header: "content-type: application/json"
- POST with parameters

The framework should support parsing the JSON parameters and returning the JSON response.

Apart from the above, the framework should also support:

- validate the parameter (mostly if the struct specify the parameter, while the JSON object parsed has not include such parameter)