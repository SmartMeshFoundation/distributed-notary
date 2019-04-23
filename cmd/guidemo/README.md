 

# Step 1: install the demo

Run the following commands:

    $ mv bind.go ../

# Step 2: install the bundler

Run the following command:

    $ go get -u github.com/asticode/go-astilectron-bundler/...
    $ go install github.com/asticode/go-astilectron-bundler/astilectron-bundler

Go get and build the binary.
And don't forget to add `$GOPATH/bin` to your `$PATH`.

# Step 3: bundle the app for your current environment

Run the following commands:

    $ cd $GOPATH/src/github.com/asticode/go-astilectron-demo
    $ astilectron-bundler -v
    $ rm bind_darwin_amd64.go  #删除/bind_darwin_amd64.go
    $ mv ../bind.go .
    
# Step 4: 
检查vendor 目录结构,应该如下:
```bash
astilectron
electron-darwin-amd642
```

# step 5
$ go build
$ ./guidemo 