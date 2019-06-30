# Teams & PD Plugin for PCF Event Alerts 

[PCF Event alert](https://docs.pivotal.io/event-alerts/1-2/index.html) supports Email, Slack and Webhook integration. However, the webhook integration is not easily integrated with Microsoft teams or Pagerduty and needs custom event transformation. This App enables the transformation and easy integration with these tools. 

## Getting Started
### Running on [Pivotal Web Services](https://run.pivotal.io/)


### Running locally
The following assumes you have a working, recent version of Go installed, and you have a properly set-up Go workspace.
```
|-<Go-workspace-name>
	|-src
	|-bin
	|-pkg
```
1. Install and start MySQL
```
brew install mysql
mysql.server start
mysql -u root
```

2. Create a database user and table 
```
MariaDB [(none)]> CREATE USER 'mapper'@'localhost' IDENTIFIED BY 'mapper';
Query OK, 0 rows affected (0.008 sec)

MariaDB [(none)]> CREATE DATABASE event_router_mapping;
Query OK, 1 row affected (0.002 sec)

MariaDB [(none)]> GRANT ALL ON event_router_mapping.* TO 'mapper'@'localhost';
Query OK, 0 rows affected (0.004 sec)
```

3. Install and start the application server
```
go get github.com/tushardag/webhook-handler
cd $GOPATH/src/github.com/webhook-handler
go install
$GOPATH/bin/webhook-handler
```

4. Export the test host in another shell, where you can then run the interactions:
```
export APPLINK=http://localhost:3000
```

Now follow the [interaction instructions](#interaction-instructions).

## Interaction instructions
Start by creating the route mapping either for MS Teams or for PagerDuty (HTTP 200 response code is expected)
`curl -v -H "Content-Type: application/json" -X PUT $APPLINK/teams/testIdentifier -d '{"URL": "https://outlook.office.com/webhook/9876-xyz/IncomingWebhook/1234/abc","Description": "Sample Teams Incoming webhook link"}'`
OR
`curl -v -H "Content-Type: application/json" -X PUT $APPLINK/pagerduty/testIdentifier -d '{"URL": "c576hhj7a88d99b0b23dc3htr0v","Description": "Sample Pagerduty Event API V2 integration key"}'`

List out the existing routes and respective Teams or Pagerduty mapping information 
`curl -v -X GET $APPLINK/routes`

Remove/Delete the existing route mapping (HTTP 200 response code is expected)
`curl -v -X DELETE $APPLINK/teams/testIdentifier`
OR
`curl -v -X DELETE $APPLINK/pagerduty/testIdentifier`

Post call to either open incident in PagerDuty or post message in Teams. This would be the webhook added in PCF Event Alert and called by EventAlert (HTTP 200 response code is expected)
```
curl -v -H "Content-Type: application/json" -X POST $APPLINK/pagerduty/testIdentifier -d \
'{
    "publisher": "healthwatch",
    "topic": "gorouter.latency.uaa",
    "metadata": {
        "status": "Critical",
        "statusColor": "#DD545B",
        "value": "200.00 ms",
        "job": "router",
        "index": "a8ffa403-dc5f-4d35-82cf-9dbed10b0f0f",
        "ip": "10.100.80.20",
        "deployment": "cf-abc123def321c0",
        "foundation": "sys.myfoundation.mydomain.com",
        "eventType": "Performance/Health Event",
        "eventDescription": "The UAA Request Latency measurement has crossed a critical threshold.",
        "url": "https://healthwatch.sys.myfoundation.mydomain.com/router/details",
        "docsUrl": "https://docs.pivotal.io/pivotalcf/2-5/monitoring/kpi.html#uaa_latency"
    }
}'
```

### Break down into end to end tests

Explain what these tests test and why

```
Give an example
```

### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Deployment

Add additional notes about how to deploy this on a live system

## Built With

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
* [Maven](https://maven.apache.org/) - Dependency Management
* [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Billie Thompson** - *Initial work* - [PurpleBooth](https://github.com/PurpleBooth)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc




















Local setup for mysql 

xx51309@L-SFO5BBG8WN-M  ~/github.gapinc.com/pcf-platform-automation   master  mysql -uroot
Welcome to the MariaDB monitor.  Commands end with ; or \g.
Your MariaDB connection id is 10
Server version: 10.3.15-MariaDB Homebrew

Copyright (c) 2000, 2018, Oracle, MariaDB Corporation Ab and others.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

MariaDB [(none)]> CREATE USER 'mapper'@'localhost' IDENTIFIED BY 'mapper';
Query OK, 0 rows affected (0.008 sec)

MariaDB [(none)]> CREATE DATABASE event_router_mapping;
Query OK, 1 row affected (0.002 sec)

MariaDB [(none)]> GRANT ALL ON event_router_mapping.* TO 'mapper'@'localhost';
Query OK, 0 rows affected (0.004 sec)

MariaDB [(none)]> exit
Bye