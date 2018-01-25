[![wercker status](https://app.wercker.com/status/620e72eed2fdf6afbc5a7f1e5ca943f4/s/master "wercker status")](https://app.wercker.com/project/byKey/620e72eed2fdf6afbc5a7f1e5ca943f4)
[![Go Report Card](https://goreportcard.com/badge/github.com/wayneashleyberry/dayzero-app)](https://goreportcard.com/report/github.com/wayneashleyberry/dayzero-app)
[![GoDoc](https://godoc.org/github.com/wayneashleyberry/dayzero-app?status.svg)](https://godoc.org/github.com/wayneashleyberry/dayzero-app)

> This project exposes the data available on the [City of Cape Town's water dashboard](http://coct.co/water-dashboard/) as a simple [JSON API](https://day-zero-cape-town.appspot.com/api/dashboard).

### Related

#### Command-line

Day Zero status for your command-line: https://github.com/wayneashleyberry/dayzero

#### BitBar Integration

If you're using [BitBar](https://getbitbar.com/) on macOS, then you could use the following script.

```sh
#!/bin/bash
time=`curl -s https://day-zero-cape-town.appspot.com/api/dashboard | /usr/local/bin/jq -r '.dayzero_humane'`
prefix='Day zero is '
echo $prefix$time
```

If you run into any issues, make sure you have [jq](https://stedolan.github.io/jq/) installed.

```sh
brew install jq
```
