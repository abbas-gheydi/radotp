To configure the radius server, follow these steps:
•  Edit the config/radiusd.conf file and make the following changes:

•  Replace the ldap server address with your domain controller (ldap server) address.

•  Replace the radius secret with a strong and secure passphrase.

•  Set the Authentication_Mode in the radius section according to your needs.

•  Edit the docker-compose.yml file and make the following changes:

•  Replace the database username and password with new credentials.

•  Make sure the database username and password matches the one in the config/radiusd.conf file.

•  Save the changes and restart the docker containers.

```bash  
docker-compose up -d  

```
