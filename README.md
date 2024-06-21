# go-reload
SSE reload server

## Notes

Very opinionated way of handling "hot reload" via an installable server, templ component and a command to trigger the reload. Needs the url passed into the server + reload command.


Also - we call it hot reload - but it's not actually hot reload in a "React" sense. We detect a change, an SSE message is triggered to the front end and we do a `window.location.reload();`. Not magic, but nicer than pressing F5.

## Usage

For your app, you'll general have a `templ` file that serves as the base template for the app. It should do your styles, common js + page structure. 

1. Add `import "github.com/idreaminteractive/goreload"` to the top.
2. At the bottom, before the closing `</body>` tag, put the following in:


```
if isLocalDevelopment {
    @goreload.ReloadComponent(hotReloadUrl)
}
```

I generally attach both the hot reload url and the local dev flag into the context via middleware to allow `templ` to grab them easily.

3. To detect changes in your app + signal the reload, you'll need to call `err := goreload.SendReloadSignal(hotReloadUrl)`. I usually put this in the server itself, and use `air` for restarting the server process on changes. That way, the hot reload is triggered on each server start up.

4. Along side your app (or `air`), run the hot reload server via `goreload --url <hotReloadUrl> server`. This will handle the SSE signalling. We use [Task](https://github.com/go-task/task) for setting this up when we run the project in Gitpod. For example:

```
dev:
    deps: [air, hotreload]

air:
    cmds:
      - air

hotreload:
    cmds:
        - goreload --url {{.HOT_RELOAD_URL}} server
    vars:
        HOT_RELOAD_URL:
            sh: gp url 8082
```

Running `task dev` starts the `air` process to reload changes on code updates, and the `goreload` server is fed the dev url on port `8082`. We do a similar passing of the url internally in the app to fill out the template + the signal calls.

## Todo 

- Need some tests, _probably_.