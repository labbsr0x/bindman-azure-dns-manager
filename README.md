# Bindman Azure DNS Manager
![Build Status](https://travis-ci.com/labbsr0x/bindman-azure-dns-manager.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/labbsr0x/bindman-azure-dns-manager)](https://goreportcard.com/report/github.com/labbsr0x/bindman-azure-dns-manager)
[![Docker Pulls](https://img.shields.io/docker/pulls/labbsr0x/bindman-azure-dns-manager.svg)](https://hub.docker.com/r/labbsr0x/bindman-azure-dns-manager/)

This repository defines the component that manages Azure DNS Zone.

Azure DNS Zone commands get dispatched from REST API calls defined in the Bindman webhook project [Bindman DNS Webhook](https://github.com/labbsr0x/bindman-dns-webhook).

# Configuration

The Bindman is setup with the help of environment variables and volume mapping in the following way: 

## Volume Mapping

A store of records being managed is needed. Hence, a `/data` volume must be mapped to the host.

## Environment variables

1. `mandatory` **BINDMAN_AZURE_RESOURCE_GROUP**: specifies the app Resource Group to use.

2. `mandatory` **BINDMAN_AZURE_SUBSCRIPTION_ID**: specifies the subscription to use.

3. `mandatory` **BINDMAN_AZURE_CLIENT_ID**: specifies the app client ID to use.

4. `mandatory` **BINDMAN_AZURE_CLIENT_SECRET**: specifies the app secret to use.

5. `mandatory` **BINDMAN_AZURE_TENANT_ID**: specifies the Tenant to which to authenticate.

6. `mandatory` **BINDMAN_ZONE**: the zone that the bindman instance is responsible for managing.

7. `optional` **BINDMAN_DNS_TTL**: the dns recording rule expiration time (or time-to-live). By default, the TTL is **3600 seconds**.

8. `optional` **BINDMAN_DNS_REMOVAL_DELAY**: the delay in minutes to be applied to the removal of an DNS entry. The default is 10 minutes. This is to guarantee that in fact the removal should be processed.

9. `optional` **BINDMAN_MODE**: let the runtime know if the DEBUG mode is activated; useful for debugging the intermediary files created for sending `nsupdate` commands. Possible values: `DEBUG|PROD`. Empty defaults to `PROD`.