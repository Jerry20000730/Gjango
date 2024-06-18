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
