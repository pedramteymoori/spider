# Spider
Spider is a concurrent web page analyzer

## How to use

1. Download
    ```Bash
    go mod download
    ```

2. Start Spider
    ```Bash
    go run main.go
    ```

3. Check webpage by send a `GET` request to the app
    ```bash
    curl "http://localhost:8080?url=http://www.columbia.edu/~fdc/sample.html"
    ```

4. You'll see the analysis result:
    ```json
    {"HTMLVersion":"HTML 5","Title":"Sample Web Page","Headings":{"h2":1,"h3":13},"InternalLinks":["http://www.columbia.edu/~fdc/","http://www.columbia.edu/cu/computinghistory","http://www.columbia.edu/~fdc/family/dcmall.html","http://www.columbia.edu/~fdc/family/hallshill.html","http://www.columbia.edu/~fdc/family/frankfurt.html","http://www.columbia.edu/~fdc/sample.html#basics","http://www.columbia.edu/~fdc/sample.html#syntax","http://www.columbia.edu/~fdc/sample.html#chars","http://www.columbia.edu/~fdc/sample.html#convert","http://www.columbia.edu/~fdc/sample.html#effects","http://www.columbia.edu/~fdc/sample.html#lists","http://www.columbia.edu/~fdc/sample.html#links","http://www.columbia.edu/~fdc/sample.html#tables","http://www.columbia.edu/~fdc/sample.html#viewing","http://www.columbia.edu/~fdc/sample.html#install","http://www.columbia.edu/~fdc/sample.html#more","http://www.columbia.edu/~fdc/sample.html#fluid","http://www.columbia.edu/~fdc/sample.html#install","http://www.columbia.edu/~fdc/sample.html#tables","http://www.columbia.edu/entities.html","http://www.columbia.edu/~fdc/sample.html#contents","http://www.columbia.edu/~fdc/sample.html#lists","http://www.columbia.edu/rabbit.jpg","http://www.columbia.edu/~fdc/index.html"],"ExternalLinks":["https://kermitproject.org/newdeal/","http://kermitproject.org/unix.html","http://panix.com/~fdc/sample.html","https://en.wikipedia.org/wiki/HTML","https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes","https://en.wikipedia.org/wiki/UTF-8","https://en.wikipedia.org/wiki/UTF-8","http://kermitproject.org/about.html","http://kermitproject.org/html.html","http://www.kermitproject.org/","http://kermitproject.org/ckuins.html#x4.0","http://kermitproject.org/ckuins.html","https://www.w3schools.com/html/html_tables.asp","http://kermitproject.org/","https://www.cnn.com/~fdc/sample.html","https://www.w3.org/TR/html52/","https://developer.mozilla.org/en-US/docs/Web/CSS","http://validator.w3.org/"],"InAccessibleLinks":15,"HasLogin":false}
    ```
