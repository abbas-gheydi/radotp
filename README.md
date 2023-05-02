# RadOTP
## About The Project
RadOTP is a Radius server that is designed for "SSL VPN" authentication with two-factor authentication mechanisms.   
* It has a built-in "LDAP client," and it can connect to Microsoft Active Directory.   
* main reason to use radOTP is its fantastic feature: two-factor authentication using One Time password "OTP."   
* users activity is exposed to Prometheus and Grafana for alerting and monitoring purposes.   
* It has a web interface to manage users.   
* interactive mode using radius Access-Challenge.   
* high availability, data saved in Postgres SQL. if you want HA, then make a Postgres cluster.   
* it works in three modes:   
    only_password: authenticate users against Active Directory or any LDAP/LDAPS server.   

    only_otp: authenticate users with OTP database only.   

    two_fa: two factors authenticate mode. AD password + OTP code.   

### How to Install
first install docker and docker-compose then, Download [RadOTP](https://github.com/Abbas-gheydi/radotp/releases) and install it by:  
```bash
docker-compose build . 
docker-compose up -d  
```
if you want to run it as a service, you must make a new [systemd serivce](https://www.suse.com/support/kb/doc/?id=000019672).   

### How to Use it:
  
- edit radiusd.conf (radotp.conf can be in current directory or in /etc/radotp/) then start docker compose or radotp service.       

- in your browser, type IP_ADDRESS:8080 and use the admin/admin password to log in and manage users.   

- Download [Google Authenticator](https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&hl=en_US&gl=US) and scan the QR code.   

- Config Cisco or Fortinet firewalls to use radOTP (Radius) as authentication source:   
[Fortigate](https://docs.fortinet.com/document/fortigate/6.0.0/cookbook/200757/connecting-the-fortigate-to-the-radius-server)   
[Cisco ASA](https://www.cisco.com/c/en/us/support/docs/security/asa-5500-x-series-next-generation-firewalls/98594-configure-radius-authentication.html)   

## License

MIT
