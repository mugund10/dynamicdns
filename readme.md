# Dynamic Dns Cli GO app

## Overview

The Dynamic DNS CLI (DdnsCli) is a command-line tool built in Go for managing Dynamic DNS updates using the Name.com API. The tool helps to ensure that your domain's DNS records are kept up-to-date with your current IP address.

## Workflow

- Regularly checks the current IP address from my [webservice](https://github.com/mugund10/whatsmyip).
- Compares it with the existing A record in DNS.
- Updates the record if there is a mismatch.
- Waits 5 minutes before checking again if the IP addresses match.

## Setup

Before running the application, ensure that you create a `.env` file in the root directory of your project with the following content:

```dotenv
APIUSERNAME="yourname.com_api_username"
APITOKEN="yourname.com_api_token"
DOMAIN="your domain address"
ID="id_of_your_dns_field"

```

## Installation

```bash
git clone github.com/mugund10/dynamicdns
cd dynamicdns
go build cmd/ddnsclient/main.go 
./main
```

## Status

- Cli is yet to be developed 


