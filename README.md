# terraform-provider-remotefile

Simple terraform provider that allows you to download a file from a remote source over HTTP(S). This may in future be extended to support other file protocols.

## Resources

### remotefile_download
The download resource makes an HTTP GET request to the given `source_url` and downloads the response to the `destination` file path. When refreshing the state, the remote file is checked to see if it has changed. If it has, it is re-downloaded. Similarly, if the previously downloaded file is missing or has been altered the remote file is re-downloaded.

#### Example Usage
```
resource "remotefile_download" "lambda_zip" {
  source_url  = "https://github.com/keattang/expire-iam-access-keys-lambda/archive/master.zip"
  destination = "/tmp/lambda.zip"
}
```

#### Argument Reference
- `source_url` - The URL to request data from. If this server supports E-Tags they will be used to check if the remote file has changed. If not, the remote file will have to be downloaded in order to check.
- `destination` - The file path to which to download the remote file.

#### Attributes Reference
- `file_base64sha256` - The base64-encoded SHA256 checksum of the downloaded file.
- `file_base64sha256` - The MD5 checksum of the downloaded file.
- `etag` - The value of the E-Tag header received when downloading the remote file.