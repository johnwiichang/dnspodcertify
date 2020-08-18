# dnspodcertify
DNSPodCertify is a tiny project for Certify to send DNS record actions to DNSPod. Written in Golang.

## Installation

If you'd like to customise your own version or change default configurations, you ought to fork this repository and use

```bash
go build -ldflags "-s -w"
```

Or you can download artifacts from `Release` section.

> Tips:
> This application is a script tool for Certify, an ACME client for Microsoft Windows. So only Windows releases will be compiled and listed.

### Platforms

- Microsoft Windows (AMD64)
- Microsoft Windows (IA32)
- Microsoft Windows (ARM)

## How to...

### Get Credentials
You should apply your access token via from [DNSPod](https://console.dnspod.cn/account/token), and encode `{id}.{token}` UTF-8 bytes in HEX code with lower format.

### Rename Binary
You should make a copy of `dnspodcertify`, and rename they two with:

- The script used to create the record should be named with `add_{hex}.exe`
- The script used to delete the record should be named with `del_{hex}.exe`

### Update Settings

You need Certify client with `Customise DNS Authentication`:

1. Select `Authorization`.
2. Select `dns-01` in Authorization Settings.
3. Switch DNS Update Method to `(Use Custom Script)`
4. Fill `Create Script Path` and `Delete Script Path` with `dnspodcertify` binary paths.
5. Test and save configuration if succeed.