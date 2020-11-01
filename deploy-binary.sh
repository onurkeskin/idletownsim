# Cross compile the app for linux/amd64
GOOS=linux GOARCH=amd64 go build -v -o $TMP/app ../app

# Add the app binary
tar -c -f $TMP/bundle.tar -C $TMP app

# Add static files.
tar -u -f $TMP/bundle.tar -C ../app templates

# BOOKSHELF_DEPLOY_LOCATION is something like "gs://my-bucket/bookshelf-VERSION.tar".
gsutil cp $TMP/bundle.tar $APP_DEPLOY_LOCATION