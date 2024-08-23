# Dynamic Dns Cli GO app
## Overview

The Dynamic DNS CLI (DdnsCli) is a command-line tool built in Go for managing Dynamic DNS updates using the Name.com API. The tool helps to ensure that your domain's DNS records are kept up-to-date with your current IP address.

## Features

- Regularly checks the current IP address.
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

## Status

**Under Development**

The CLI tool is currently under development. The `main.go` file is operational, but additional features and improvements are planned for future updates.