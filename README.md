# Go Wiki

A simple HTTP server rendering Markdown styled documents on the fly and optionally shows its git history including diffs.

**NOTE** This is toy project to help me learn Go, so don't run this on anything publically available.

![Screenshot1](https://cloud.githubusercontent.com/assets/177685/5720761/2337178e-9b29-11e4-8a86-224f7905b3f6.png)

## Installation

```bash
$ go get github.com/renstrom/go-wiki
$ $GOPATH/bin/gowiki <path to wiki directory>
```

## Customize

It's only possible to customize the CSS. Put all your customizations in a file of your choosing and point to it using the `--custom-css` flag.

```bash
$ gowiki ~/www/wiki --custom-css=<path to custom css>
```

## Usage

Create git repository containing your Markdown formatted wiki pages.

### On the server

Create an empty repository.

``` bash
$ mkdir -p ~/www/wiki && cd $_
$ git init
$ git config core.worktree ~/www/wiki
$ git config receive.denycurrentbranch ignore
```

Setup a post-receive hook.

``` bash
$ cat > .git/hooks/post-receive <<EOF
#!/bin/sh
git checkout --force
EOF
$ chmod +x .git/hooks/post-receive
```

Start the server.

``` bash
$ gowiki ~/www/wiki
```

### On your local machine

Replace `<user>` and `<host>` with credentials for your specific machine.

``` bash
$ git init
$ git remote add origin \
    ssh://<user>@<host>/home/<user>/www/wiki
```

Now create some Markdown file and push.

``` bash
$ git add index.md
$ git commit -m 'Add index page'
$ git push origin master
```

## License

MIT
