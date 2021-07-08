# Postee

![Postee](postee.jpeg)

![Docker Pulls][docker-pull]
[![Coverage Status][cov-img]][cov]
[![Go Report Card][report-card-img]][report-card]
![](https://github.com/aquasecurity/postee/workflows/Go/badge.svg)
[![License][license-img]][license]

[download]: https://img.shields.io/github/downloads/aquasecurity/postee/total?logo=github
[release-img]: https://img.shields.io/github/release/aquasecurity/postee.svg?logo=github
[release]: https://github.com/aquasecurity/postee/releases
[docker-pull]: https://img.shields.io/docker/pulls/aquasec/postee?logo=docker&label=docker%20pulls%20%2F%20postee
[go-doc-img]: https://godoc.org/github.com/aquasecurity/postee?status.svg
[cov-img]: https://codecov.io/github/aquasecurity/postee/branch/main/graph/badge.svg
[cov]: https://codecov.io/github/aquasecurity/postee
[report-card-img]: https://goreportcard.com/badge/github.com/aquasecurity/postee
[report-card]: https://goreportcard.com/report/github.com/aquasecurity/postee
[license-img]: https://img.shields.io/badge/License-mit-blue.svg
[license]: https://github.com/aquasecurity/postee/blob/master/LICENSE

# Table of Contents
- [Abstract](#abstract)
- [Features](#features)
- [Installation](#installation)
- [Configure the Aqua Server with Webhook Integration](#configure-the-aqua-server-with-webhook-integration)
- [Set up the Configuration File](#set-up-the-configuration-file)

## Abstract
Postee is a simple application that receives JSON messages from one hand, and delivers them (after reformatting) to different collaboration systems, like: JIRA, Email, Slack, Microsoft Teams, Generic WebHook, Splunk and ServiceNow.

Primary use of Postee is act as notification component for Aqua Security products. It's extremely useful for sending vulnerability scan results to collaboration systems
## Features

### New features in Postee V2
Main goal of V2 changes is to make every aspect of product customizable. It now can work with any incoming JSON messages (not only vulnerability scan results). Once message is received by application it's evaluated against app config to make a decision whether it needs to be forwarded or dropped. All features related to represent messages as html or slack markdown still available but to support custom formatting new feature is introduced: Rego Templates. It uses Rego Language to define message body. More details are in Rego templates section below

### Policy related features in Postee V2
Many options which were intended to limit sending messages for specific events only are redesigning in favor of using OPA rules.
Here is list of options which are not supported anymore:
- Policy-Min-Vulnerability
- Policy-Registry
- Policy-Image-Name
- Policy-Only-Fix-Available
- Policy-Non-Compliant
- Ignore-Registry
- Ignore-Image-Name
- Policy-OPA

#### OPA policies examples
Code snippet below can be used instead legacy Policy-Image-Name option
```
contains(input.image, "alpine")
```

#### Options still supported in V2

##### Policy-Show-All
Image rescans: When an image is rescanned, the integration will not send a message if the scan results are same as the previous scan results. If the scan results are different then a message with the diff between the results will be sent.

##### Aggregate-Issues-Number and Aggregate-Issues-Timeout
Aggregation policy: You can aggregate multiple scan results in a single ticket/message. This is useful if you would like to get a digest on daily/weekly basis.
See [Route plugins](#route-plugins-section) for more details

#### Input message workflow
Interaction of Postee v2 modules.
![Postee v2 scheme](/postee-v2-scheme.png)

#### Rego Templates
[Rego language](https://www.openpolicyagent.org/docs/latest/policy-language/) is used to define templates. Message payload is referenced as `input` when template is rendered. Result of rendering is output. Several properties are picked from output and send to configured outputs.

Key | Description |Type
--- | --- | ---
*result* | message body| Can be either string or json
*title* | message title| string
*aggregation_pkg*|reference to package used to aggregate messages (when Aggregate-Issues-Timeout or Aggregate-Issues-Number options are used). If it's missed then aggregation feature is not supported| valid rego package

So the simplest example of Rego template would look like:
```rego
package example.vuls.html

title:="Vulnerabilities are found"
result:=sprintf("Vulnerabilities are found while scanning of image: <i>%s</i>", [input.image])
```

Two examples are shipped with the app. One produces output for slack integration and another one builds html output which can be used across several integrations. These example can be used as starting point for message customization

### Postee UI
Another new feature of Postee is Admin application. It supports new application modules and provides set of forms for that.
![Config app](/postee-output-config.png)
Application has tabs for routes, outputs and templates. All field values are validated if required. Code editor is provided for the properties which can contain inline rego language statements.
Besides validation, output form provides an option to test output configuration. Test email/issue/message will be created in integrated application while testing.

See more details [here](#configure-and-run-postee-ui-application) 

### Data Persistency
The Postee container uses BoltDB to store information about previously scanned images.
This is used to prevent resending messages that were already sent before.
The size of the database can grow over time. Every image that is saved in the database uses 20K of storage.

If you would like to persist the database file between restarts of the Postee container, then you should
use a persistent storage option to mount the "/server/database" directory of the container.
The "deploy/kubernetes" directory in this project contains an example deployment that includes a basic Host Persistency.
## Installation

### From Source
Clone this project: 
```bash
git clone git@github.com:aquasecurity/postee.git
make build
./bin/postee
```

### Docker
Build the postee Docker image: 
```bash
docker build -t aquasec/postee:latest .
```

Run the Aqua Postee container with the configuration file:
```bash
docker run -d --name=postee -v /<path to configuration file>/cfg.yaml:/config/cfg.yaml \
    -e AQUAALERT_CFG=/config/cfg.yaml -e AQUAALERT_URL=0.0.0.0:8084 -e AQUAALERT_TLS=0.0.0.0:8444 \ 
    -p 8444:8444 -p 8084:8084 aquasec/postee:latest
```
### [Kubernetes](./deploy/kubernetes/README.md)

### [Helm](./deploy/helm/README.md)


## Configure the Aqua Server with Webhook Integration

Configure the Aqua Server to send a Webhook notification when a new vulnerability is found
![Screenshot](webhook-integration.png)

Validate that a ticket has been opened, or email was sent (depending on your configuration file).

You can configure the Aqua Server to send a Webhook notification whenever a new vulnerability is found.
Navigate to the **Settings** page in the System section, menu, under the "Image Scan Results Webhook" section.

Click "Enable sending image scan results to Postee server", and specify the URL of the Aqua Webhook server.

The URL is in the following formats:
**HTTPS**: https://<Postee IP or DNS>:8444
or
**HTTP**: http://<Postee IP or DNS>:8084

### Validate the Integration

To validate that the integration is working, you can scan a new image for security vulnerabilities from the Aqua Server UI (Images > Add Image > Specify Image Name > Add).

When vulnerabilities are found in an image, you will see that a JIRA ticket is created/ Email is received/ Slack message is posted to the channel.

###### *To troubleshoot the integration, you can look at both the Aqua Postee container logs and the Aqua Server logs. Use the "docker logs <container name>" command to view these logs.*

## Run the Aqua Postee Container

Build and run the Aqua Webhook Server container on the same host where the JIRA configuration file is located, as follows:

```bash
docker build -t postee:latest .

docker run -d --name=aqua-postee -v /<path to configuration file>/cfg.yaml:/config/cfg.yaml \
    -e AQUAALERT_CFG=/config/cfg.yaml -e AQUAALERT_URL=0.0.0.0:8084 -e AQUAALERT_TLS=0.0.0.0:8444 \
    -p 8444:8444 -p 8084:8084 postee:latest

```

###### *There is a volume mount that mounts the configuration file from the host to the container. There is also an environment variable, AQUAALERT_CFG, that specifies the location of the JIRA configuration file inside the container.*

## Configure and run Postee UI application
### Requirements
Postee Admin application shares location of `cfg.yaml` with main webhook app, also Bolt database needs to be in folder which is available for both apps.

**Important**: If application config is submitted by UI app then all yaml comments are removed. So if comments are important please make backup of config yaml.
### Docker Image for Postee UI application
Dockerfile to build image for UI app is [here](Dockerfile.ui)

### Orchestration example (Docker Compose)
There is an example of [docker-compose.yml](docker-compose.yml) that can be used to simplify deploying of both app. Notice that two shared volumes are used. One is for Bolt db and second to store app config. To start apps use: `docker-compose up`.

### Environment variables
Name | Description | Default value
--- | --- | ---
POSTEE_UI_CFG|Path to app config| required, no default value
POSTEE_UI_PORT|Port to use with UI app| 8090
POSTEE_UI_UPDATE_URL|Url of webhook application|required
POSTEE_ADMIN_USER|Admin account name|admin
POSTEE_ADMIN_PASSWORD|Admin account password|admin

## Outputs configuration

### Getting the JIRA connection details

Follow these steps to set up JIRA integration:

Login to Jira.
Go to the user profile API tokens (JIRA Cloud users can find it here: https://id.atlassian.com/manage-profile/security/api-tokens).
Click on the Create API Token. A new API token for the user is created.
Keep the token value, together with the JIRA URL and user name, for the next step.

### Getting the Slack connection details: [Slack Custom App](https://api.slack.com/)
1. Visit api.slack.com
2. Press "Create custom app"
3. Fill app name and select slack workspace
4. Open "Incoming webhooks" tab
5. Enable "Incoming webhooks"
6. Add webhook to workspace
7. On next screen pick slack channel and click allow
8. Copy webhook url to the Postee config

### Getting the MS Teams connection details
Open your Microsoft Teams client. Click on the "..." near the channel you would like to send notifications to.
Choose "Connectors". The connectors window will open.
Look for the "Incoming Webhook" connector (it is under the "All" category).
Click "Add" near the Incoming Webhook connector. Click "Add" again.
Provide a name and click "Create".
You will be provided with a URL address. Copy this URL and put it in the cfg.yaml.

### Configure the Splunk Integration
You will need to craate an HTTP Event Collector in Splunk Enterprise or Splunk Cloud.
This can usually be found in the Splunk console under "Settings -> Data Inputs -> HTTP Event Collector -> Add New".
Once you create an HTTP Event Collector you will receive a token. You should provide this token, together with the Splunk HTTP Collector
URL, as part of the cfg.yaml settings.

## Set up the Configuration File
To set up the integration, you will need to create a `cfg.yaml` file, which contains the connection settings. Edit the configuration file with the connection details of your JIRA, Slack, etc.

###### *IMPORTANT: Application config yaml is re-designed in V2 release and has no backward compatibility. Besides the structure changes all option names are now lowercase words separated by hyphens. So `UseMX` becomes `use-mx` for example*

### General options
General options are specified at the root level of config yaml. All these options were also available in Postee V1
Key | Description | Possible Values
--- | --- | ---
*aqua-server*|Aqua Console URL. This is used for some of the integrations to include link to scan results| any valid url
*delete-old-data*|delete data older than N day(s).  If empty then we do not delete.| any integer value
*db-verify-interval*|hours. an Interval between tests of DB. Default: 1 hour| any integer value
*max-db-size*|Max size of DB. MB. if empty then unlimited| any integer value

### Routes sections
Route is used to control messages flow. It must include references to one or more outputs and reference to the template used for message rendering.
Key | Description | Possible Values
--- | --- | ---
*name*|Unique name of route| string
*input*|Rego rule to filter message| Rego language statements
*outputs*|Outputs associated with route| Set of output names, like ["my-slack", "my-email"]. At least one element is required
*template*| Reference to template, required| any template name
#### Route plugins section
'Plugins' section contains configuration for useful Postee features. 
Key | Description | Possible Values
--- | --- | ---
*policy-show-all*|Optional. Open a ticket even if a ticket was opened for same image with same amount of vulnerabilities. Default is false.| boolean
*aggregate-issues-number*|Number of scans to aggregate into one ticket.| any integer value
*aggregate-issues-timeout*|number of seconds, minutes, hours to aggregate|Maximum is 24 hours Xs or Xm or Xh

### Templates
There are several options to configure templates. One required template property is `name` (to allow references to template within route configuration). For further configuration pick none option from the list below:
- Use buildin template. Postee loads bundle of templates from `rego-templates` folder. Root folder includes several templates shipped with application. Any subfolder can be used to append user templates. To specify particular rego rule use `rego-package` option.  Example is `postee.vuls.html` 
- Specify inline template. Relative small templates can be added to config directly. `body` option can be used for that 
- Load from url. Rego template can be loaded from url. There is an `url` option
- Legacy template. Legacy templates are introduced to support Postee V1 renderers. Option is `legacy-scan-renderer`. Available values are  "jira", "slack", "html". "jira" should be used for jira integration, "slack" is for slack and "html" is for everything else.

### Outputs
Outputs were known as plugins before. 
#### ServiceNow integration parameters
Key | Description | Possible Values
--- | --- | ---
*user* | ServiceNow user name | 
*password* | User API key / password |
*instance* | Name of ServiceNow Instance (usually the XXX at XXX.servicenow.com)|
*board* | ServiceNow board name to open tickets on. Default is "incident" |

#### Jira integration parameters
Key | Description | Possible Values
--- | --- | ---
*url* | Jira project url |
*user* | Jira user name | 
*password* | User's API key | 
*project-key* | The JIRA project key |
*board* |  Optional: JIRA board key |
*priority*|  Optional: ticket priority, e.g., High |
*assignee*| Optional: comma separated list of users (emails) that will be assigned to ticket, e.g., ["john@yahoo.com"]. To assign a ticket to the Application Owner email address (as defined in Aqua Application Scope, owner email field), specify ["<%application_scope_owner%>"] as the assignee value |
*issuetype*| Optional: issue type, e.g., Bug |
*labels*| Optional: comma separated list of labels that will be assigned to ticket, e.g., ["label1", "label2"]|
*sprint*| Optional: Sprint name, e.g., "3.5 Sprint 8" |

For Jira you can also specify custom fields that will be populated with values.
Use the `unknowns` parameter in cfg.yaml for custom fields.
Under the `unknowns` parameter, specify the list of fields names to provide value for.
You can add "-numeric-field", "-multiple-value", "multiple-line-text-field", "-date-time-picker" and "-field-url" as suffix to the custom field name, to specify what is the field type.

For example: 
```yaml
unknowns:
     mycustom: "this is a text custom field"
     mycustom-numeric-field: 123
     mycustom-multiple-value: 1,2,3 
     mycustom-multiple-line-text-field: "text \n moretext" 
     mycustom-date-time-picker: 2014-04-11T12:14:26.880+0400
     mycustom-url: https://tour.golang.org/moretypes/7
```

#### Email integration parameters
Key | Description | Possible Values
--- | --- | ---
*use-mx* | Whether to send the email as an SMTP server or a client. Specify 'true' if you would like to send email as an smtp server, in this case you don't need to provide user, password, host and port. | true, false
*user* | User name (usually email address) |
*password* | Password | 
*host* | SMTP host name | 
*port* | SMTP port |
*sender* |  Sender's email address |
*recipients*|  Recipients (array of comma separated emails), e.g. ["john@yahoo.com"]. To send the email to the Application Owner email address (as defined in Aqua Application Scope, owner email field), specify ["<%application_scope_owner%>"] as the recipients value |

#### Slack integration parameters
Key | Description | Possible Values
--- | --- | ---
*url* | Slack WebHook URL (includes the access key) |


#### MS Teams integration parameters
Key | Description | Possible Values
--- | --- | ---
*url* | MS Teams WebHook URL |

#### Splunk integration parameters
Key | Description | Possible Values
--- | --- | ---
*token* | The Splunk HTTP event collector token | 
*url* | URL to Splunk HTTP event collector (e.g. http://server:8088) |
*size-limit* | Optional. Maximum scan length, in bytes. Default: 10000 | 10000

#### Generic Webhook integration parameters
Key | Description | Possible Values
--- | --- | ---
*url* | Webhook URL |
