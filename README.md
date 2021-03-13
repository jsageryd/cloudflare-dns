# Cloudflare DNS record updater

Updates Cloudflare DNS A records with current public IP.

- Supports multiple hosts
- Only updates the records if the current IP has changed

## Usage
### Install
1. `go install github.com/jsageryd/cloudflare-dns@latest`

### Build for Raspberry 2
1. `GOOS=linux GOARCH=arm GOARM=7 go build`

### Run once
1. `CF_API_KEY=secret CF_API_EMAIL=me@example.com ./cloudflare-dns domain.tld [domain.tld ...]`

### Schedule
1. Put `cloudflare-dns` somewhere.
2. Edit crontab (`crontab -e`) to schedule the script to run for example every 15 minutes:

```
CF_API_KEY=secret
CF_API_EMAIL=me@example.com
*/15 * * * * /path/to/cloudflare-dns domain.tld [domain.tld ...]
```

## Licence
Copyright (c) 2016 Johan Sageryd <j@1616.se>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
