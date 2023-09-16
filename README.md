# RadOTP
## About The Project
RadOTP is a dedicated Radius server expertly tailored to serve as a two-factor authentication source for VPN users on firewalls like Fortigate and Cisco ASA.



```
  +----------------------------------+
  | internet                         |
  |                                  |
  |   +--------------------------+   |
  |   | User                     |   |
  |   |--------------------------|   |
  |   | Username: [YourUsername] |   |
  |   | Password: [YourPassword] |   |
  |   | OTP:      [YourOTP]      |   |
  |   +--------------------------+   |
  +----------------------------------+

                 |
                 v
 
  +----------------------------------+
  | Corporate Network                |
  |                                  |
  |   +--------------------------+   |
  |   | Fortigate SSL VPN        |   |
  |   |--------------------------|   |
  |   | Connects to RadOTP       |   |
  |   +--------------------------+   |
  |              |                   |
  |              v                   |
  |   +--------------------------+   |
  |   | RadOTP Radius Server     |   |
  |   |--------------------------|   |
  |   | 1. AD Password Check     |   |
  |   | 2. OTP Validation        |   |
  |   +--------------------------+   |
  +----------------------------------+

```

 **It has the following features**:

-  It supports two-factor authentication (2FA) using one-time passwords (OTP) that are generated and stored in a PostgreSQL database.

-  It can connect to Microsoft Active Directory using an LDAP client and verify users' credentials.

-  It has a web interface that allows administrators to manage users, view logs, and configure settings.

-  It exposes users' activity to Prometheus and Grafana for monitoring and alerting purposes.

-  It uses radius Access-Challenge to interact with users and request additional information.

-  It offers high availability and performance by using PostgreSQL replication and docker-compose deployment.

-  It has a REST API that enables external applications to manage users programmatically.


**The radius server supports four modes of authentication:**
  
-  only_password: This mode authenticates users against an Active Directory LDAP/LDAPS server. Users only need to enter their AD password to log in.

-  only_otp: This mode authenticates users with an OTP database only. Users only need to enter a one-time password (OTP) code to log in.

-  two_fa: This mode enables two-factor authentication (2FA). Users need to enter both their AD password and an OTP code to log in.

-  two_fa_optional_otp: This mode is similar to two_fa, but it only applies 2FA to users who have an OTP in the database. Users who do not have an OTP can log in with their AD password only.

![RadOTP](https://github.com/Abbas-gheydi/radotp/blob/main/assets/radotp.jpg)

### How to Install
first install docker and docker-compose then, Download [RadOTP](https://github.com/Abbas-gheydi/radotp/releases) and deploy it by:  
```bash
docker-compose up -d  
```

### How to Use it:
  
To configure the radius server, follow these steps:
•  Edit the config/radiusd.conf file and make the following changes:

•  Replace the ldap server address with your domain controller (ldap server) address.

•  Replace the radius secret with a strong and secure passphrase.

•  Set the Authentication_Mode in the radius section according to your needs.

•  Edit the docker-compose.yml file and make the following changes:

•  Replace the database username and password with new credentials.

•  Make sure the database username and password matches the one in the config/radiusd.conf file.

•  Save the changes and restart the docker containers.   

• in your browser, type http**s**://IP_ADDRESS and use the admin/admin password to login and manage users.   

• Download [Google Authenticator](https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&hl=en_US&gl=US) and scan the QR code.   

• Config Cisco or Fortinet firewalls to use radOTP (Radius) as authentication source:   
[Fortigate](https://docs.fortinet.com/document/fortigate/6.0.0/cookbook/200757/connecting-the-fortigate-to-the-radius-server)   
[Cisco ASA](https://www.cisco.com/c/en/us/support/docs/security/asa-5500-x-series-next-generation-firewalls/98594-configure-radius-authentication.html)   

## License

MIT
