# s3buster
S3 recon tool written in Go

It reads a wordlist of candidate bucket names, builds the regional S3
hostname for each, resolves it via DNS, and attempts to list the bucket's
contents — reporting the keys and sizes of any publicly listable bucket.

## Install

```sh
go install github.com/mikemadden42/s3buster@latest
```

Or build from a checkout:

```sh
go build -o s3buster .
```

## Usage

```sh
s3buster [flags]
```

| Flag           | Default      | Description                                   |
| -------------- | ------------ | --------------------------------------------- |
| `-wordlist`    | `words.txt`  | Path to the file of candidate bucket names.   |
| `-region`      | `us-west-1`  | AWS region of the S3 endpoint.                |
| `-concurrency` | `10`         | Number of hosts to scan in parallel.          |

The wordlist is one candidate bucket name per line; blank lines are skipped.

### Examples

Scan the default `words.txt` against `us-west-1`:

```sh
s3buster
```

Scan a custom wordlist in another region, 25 hosts at a time:

```sh
s3buster -wordlist names.txt -region eu-west-1 -concurrency 25
```

For each host, s3buster prints the hostname, any resolved IP addresses, and
the `key`/`size` of every object it can list. Because hosts are scanned
concurrently, their output blocks appear in completion order rather than
wordlist order.

## Disclaimers

1. All the code/content here is for demonstration purposes only.

2. Do not perform any activities against systems or applications without authorized consent.

3. It is your fault if you do not follow the aforementioned recommendations or use the code in an unethical manner.
